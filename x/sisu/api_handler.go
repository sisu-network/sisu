package sisu

import (
	"fmt"

	"github.com/echovl/cardano-go"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	chainstypes "github.com/sisu-network/deyes/chains/types"
	etypes "github.com/sisu-network/deyes/types"
	eyesTypes "github.com/sisu-network/deyes/types"
	dhtypes "github.com/sisu-network/dheart/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

const (
	// The number of block interval that we should update all token prices.
	TokenPriceUpdateInterval = 600 // About 30 mins for 3s block.
)

var (
	ErrInvalidMessageType      = fmt.Errorf("Invalid Message Type")
	ErrMessageHasBeenProcessed = fmt.Errorf("Message has been processed")
	ErrCannotFindMessage       = fmt.Errorf("Cannot find the message in node's private db")
	ErrValueDoesNotMatch       = fmt.Errorf("Value does not match")
)

// ApiHandler handles API callback from dheart or deyes. There are few functions (BeginBlock & EndBlock)
// that are still present for historical reason. They should be moved out of this file.
type ApiHandler struct {
	keeper     keeper.Keeper
	txSubmit   common.TxSubmit
	appKeys    common.AppKeys
	globalData common.GlobalData
	txTracker  TxTracker
	mc         ManagerContainer

	// Dheart & Deyes client
	dheartClient external.DheartClient
	deyesClient  external.DeyesClient

	privateDb keeper.Storage
}

func NewApiHandler(
	privateDb keeper.Storage,
	mc ManagerContainer,
) *ApiHandler {
	a := &ApiHandler{
		keeper:       mc.Keeper(),
		privateDb:    privateDb,
		appKeys:      mc.AppKeys(),
		txSubmit:     mc.TxSubmit(),
		globalData:   mc.GlobalData(),
		dheartClient: mc.DheartClient(),
		deyesClient:  mc.DeyesClient(),
		txTracker:    mc.TxTracker(),
		mc:           mc,
	}

	return a
}

// TODO: Move this function to module.go
func (a *ApiHandler) EndBlock(ctx sdk.Context) {
	if !a.globalData.IsCatchingUp() {
		// Inform dheart that we have reached end of block so that dheart could run presign works.
		height := ctx.BlockHeight()
		a.dheartClient.BlockEnd(height)
	}
}

/**
Process for generating a new key:
- Wait for the app to catch up
- If there is no support for a particular chain, creates a proposal to include a chain
- When other nodes receive the proposal, top N validator nodes vote to see if it should accept that.
- After M blocks (M is a constant) since a proposal is sent, count the number of yes vote. If there
are enough validator supporting the new chain, send a message to TSS engine to do keygen.
*/
func (a *ApiHandler) CheckTssKeygen(ctx sdk.Context, blockHeight int64) {
	// TODO: We can replace this by sending command from client instead of running at the beginning
	// of each block.
	if a.globalData.IsCatchingUp() || ctx.BlockHeight()%50 != 2 {
		return
	}

	keyTypes := []string{libchain.KEY_TYPE_ECDSA, libchain.KEY_TYPE_EDDSA}
	for _, keyType := range keyTypes {
		if a.keeper.IsKeygenExisted(ctx, keyType, 0) {
			continue
		}

		// Broadcast a message.
		signer := a.appKeys.GetSignerAddress()
		proposal := types.NewMsgKeygenWithSigner(
			signer.String(),
			keyType,
			0,
		)

		log.Info("Submitting proposal message for ", keyType)
		a.txSubmit.SubmitMessageAsync(proposal)
	}
}

// Called after having key generation result from Sisu's api server.
func (a *ApiHandler) OnKeygenResult(result dhtypes.KeygenResult) {
	var resultEnum types.KeygenResult_Result
	switch result.Outcome {
	case dhtypes.OutcomeSuccess:
		resultEnum = types.KeygenResult_SUCCESS
	case dhtypes.OutcomeFailure:
		resultEnum = types.KeygenResult_FAILURE
	case dhtypes.OutcometNotSelected:
		resultEnum = types.KeygenResult_NOT_SELECTED
	}

	if resultEnum == types.KeygenResult_NOT_SELECTED {
		// No need to send result when this node is not selected.
		return
	}

	signerMsg := types.NewKeygenResultWithSigner(
		a.appKeys.GetSignerAddress().String(),
		result.KeyType,
		result.KeygenIndex,
		resultEnum,
		result.PubKeyBytes,
		result.Address,
	)

	log.Info("There is keygen result from dheart, resultEnum = ", resultEnum, " keyType = ", result.KeyType)

	a.txSubmit.SubmitMessageAsync(signerMsg)
	ctx := a.globalData.GetReadOnlyContext()
	params := a.keeper.GetParams(ctx)
	// Add list the public key address to watch.
	for _, chain := range params.SupportedChains {
		if libchain.GetKeyTypeForChain(chain) != result.KeyType {
			continue
		}

		log.Verbose("adding chain account address ", result.Address, " for chain ", chain)
		if libchain.IsCardanoChain(chain) {
			a.deyesClient.SetVaultAddress(chain, result.Address, "")
		}
	}
}

// OnTxDeploymentResult is a callback after there is a deployment result from deyes.
func (a *ApiHandler) OnTxDeploymentResult(result *etypes.DispatchedTxResult) {
	if !result.Success {
		log.Verbosef("Result from deyes: failed to deploy tx, chain = %s, signed hash = %s, error = %v",
			result.Chain, result.TxHash, result.Err)
		txOut := a.getTxOutFromSignedHash(result.Chain, result.TxHash)

		if txOut == nil {
			log.Errorf("Cannot find txOut for dispath result with signed hash = %s, chain = %s", result.TxHash, result.Chain)
			return
		}

		// Report this as failure. Submit to the Sisu chain
		txOutResult := &types.TxOutResult{
			OutChain: txOut.Content.OutChain,
			OutHash:  txOut.Content.OutHash,
		}

		switch result.Err {
		case etypes.ErrNotEnoughBalance:
			txOutResult.Result = types.TxOutResultType_NOT_ENOUGH_NATIVE_BALANCE
		default:
			txOutResult.Result = types.TxOutResultType_UNKNOWN
		}

		a.submitTxOutResult(txOutResult)

		return
	}

	log.Info("The transaction has been sent to blockchain (but not included in a block yet). chain = ",
		result.Chain)
	a.txTracker.UpdateStatus(result.Chain, result.TxHash, types.TxStatusDepoyed)
}

// getTxOutFromSignedHash fetches txout in the TxOut store from the hash of a signed transaction.
func (a *ApiHandler) getTxOutFromSignedHash(chain, signedHash string) *types.TxOut {
	// The txOutSig is in private db while txOut should come from common db.
	txOutSig := a.privateDb.GetTxOutSig(chain, signedHash)
	if txOutSig == nil {
		log.Error("cannot find txOutSig with full signature hash: ", signedHash)
		return nil
	}

	ctx := a.globalData.GetReadOnlyContext()
	txOut := a.keeper.GetTxOut(ctx, chain, txOutSig.HashNoSig)
	if txOut == nil {
		log.Verbose("cannot find txOut with hash (with no sig): ", txOutSig.HashNoSig)
	}

	return txOut
}

// This function is called after dheart sends Sisu keysign result.
func (a *ApiHandler) OnKeysignResult(result *dhtypes.KeysignResult) {
	if result.Outcome == dhtypes.OutcometNotSelected {
		for _, msg := range result.Request.KeysignMessages {
			a.txTracker.RemoveTransaction(msg.OutChain, msg.OutHash)
		}
		return
	}

	if result.Outcome == dhtypes.OutcomeFailure {
		// TODO: Report failure and culprits here.
		log.Warn("Dheart signing failed")
		return
	}

	ctx := a.globalData.GetReadOnlyContext()

	// Post the keysign result to cosmos chain.
	request := result.Request

	for i, keysignMsg := range request.KeysignMessages {
		msg := types.NewKeysignResult(
			a.appKeys.GetSignerAddress().String(),
			keysignMsg.OutChain,
			keysignMsg.OutHash,
			result.Outcome == dhtypes.OutcomeSuccess,
			result.Signatures[i],
		)
		a.txSubmit.SubmitMessageAsync(msg)

		// Sends it to deyes for deployment.
		if result.Outcome == dhtypes.OutcomeSuccess {
			// TODO: clean code here
			if libchain.IsETHBasedChain(keysignMsg.OutChain) {
				a.processETHSigningResult(ctx, result, keysignMsg, i)
			}

			// TODO: Submit signing failure here.
			if libchain.IsCardanoChain(keysignMsg.OutChain) {
				if err := a.processCardanoSigningResult(ctx, result, keysignMsg, i); err != nil {
					log.Error("Failed to process cardano signing result, err = ", err)
					return
				}
			}

			// Mark the tx as signed
			a.txTracker.UpdateStatus(keysignMsg.OutChain, keysignMsg.OutHash, types.TxStatusSigned)
			// TODO: Check if we have any pending confirm tx that is waiting for this tx.
		} else {
			// TODO: handle failure case here.
			log.Warnf("Signing failed, in chain = %s, out chain = %s, out hash = %s", keysignMsg.InChain,
				keysignMsg.OutChain, keysignMsg.OutHash)

			a.txTracker.OnTxFailed(keysignMsg.OutChain, keysignMsg.OutHash, types.TxStatusSignFailed)
		}
	}
}

func (a *ApiHandler) processETHSigningResult(ctx sdk.Context, result *dhtypes.KeysignResult,
	signMsg *dhtypes.KeysignMessage, index int) {
	// Find the tx in txout table
	txOut := a.keeper.GetTxOut(ctx, signMsg.OutChain, signMsg.OutHash)
	if txOut == nil {
		log.Error("Cannot find tx out with hash", signMsg.OutHash)
		return
	}

	tx := &ethtypes.Transaction{}
	if err := tx.UnmarshalBinary(result.Request.KeysignMessages[index].Bytes); err != nil {
		log.Error("cannot unmarshal tx, err =", err)
		return
	}

	// Create full tx with signature.
	chainId := libchain.GetChainIntFromId(signMsg.OutChain)
	if len(result.Signatures[index]) != 65 {
		log.Error("Signature length is not 65 for chain: ", chainId)
	}
	signedTx, err := tx.WithSignature(ethtypes.NewLondonSigner(chainId), result.Signatures[index])
	if err != nil {
		log.Error("cannot set signature for tx, err =", err)
		return
	}

	bz, err := signedTx.MarshalBinary()
	if err != nil {
		log.Error("cannot marshal tx")
		return
	}

	// // TODO: Broadcast the keysign result that includes this TxOutSig.
	// // Save this to TxOutSig
	log.Verbosef("ETH keysign result chain = %s, hash (no sig) = %s, hash (signed) = %s",
		signMsg.OutChain, signMsg.OutHash, signedTx.Hash().String())
	a.privateDb.SaveTxOutSig(&types.TxOutSig{
		Chain:       signMsg.OutChain,
		HashWithSig: signedTx.Hash().String(),
		HashNoSig:   signMsg.OutHash,
	})

	err = a.deploySignedTx(ctx, bz, signMsg.OutChain, signedTx.Hash().String())
	if err != nil {
		log.Error("deployment error: ", err)
		return
	}
}

func (a *ApiHandler) processCardanoSigningResult(ctx sdk.Context, result *dhtypes.KeysignResult, signMsg *dhtypes.KeysignMessage, index int) error {
	log.Info("Processing Cardano signing result ...")
	txOut := a.keeper.GetTxOut(ctx, signMsg.OutChain, signMsg.OutHash)
	if txOut == nil {
		err := fmt.Errorf("cannot find tx out with hash %s", signMsg.OutHash)
		log.Error(err)
		return err
	}

	tx := &cardano.Tx{}
	if err := tx.UnmarshalCBOR(txOut.Content.OutBytes); err != nil {
		log.Error("error when unmarshalling cardano tx: ", err)
		return err
	}

	pubkey := a.keeper.GetKeygenPubkey(ctx, libchain.GetKeyTypeForChain(signMsg.OutChain))
	if len(pubkey) == 0 {
		err := fmt.Errorf("cannot find pubkey for type %s", libchain.GetKeyTypeForChain(signMsg.OutChain))
		log.Error(err)
		return err
	}

	for i := range tx.WitnessSet.VKeyWitnessSet {
		tx.WitnessSet.VKeyWitnessSet[i] = cardano.VKeyWitness{VKey: pubkey, Signature: result.Signatures[index]}
	}

	hashWSig, err := tx.Hash()
	if err != nil {
		log.Error(err)
		return err
	}

	a.privateDb.SaveTxOutSig(&types.TxOutSig{
		Chain:       signMsg.OutChain,
		HashWithSig: hashWSig.String(),
		HashNoSig:   signMsg.OutHash,
	})

	txBytes, err := tx.MarshalCBOR()
	if err != nil {
		log.Error("error when marshal cardano tx: ", err)
		return err
	}
	hash, err := tx.Hash()
	if err != nil {
		return nil
	}

	err = a.deploySignedTx(ctx, txBytes, signMsg.OutChain, result.Request.KeysignMessages[index].OutHash)
	if err != nil {
		log.Error("deployment error: ", err)
		return err
	}

	log.Info("Sent signed cardano tx to deyes, tx hash = ", hash)

	return nil
}

// deploySignedTx creates a deployment request and sends it to deyes.
func (a *ApiHandler) deploySignedTx(ctx sdk.Context, bz []byte, outChain string, outHash string) error {
	log.Verbose("Sending final tx to the deyes for deployment for chain ", outChain)

	pubkey := a.keeper.GetKeygenPubkey(ctx, libchain.GetKeyTypeForChain(outChain))
	if pubkey == nil {
		return fmt.Errorf("Cannot get pubkey for chain %s", outChain)
	}

	request := &etypes.DispatchedTxRequest{
		Chain:  outChain,
		TxHash: outHash,
		Tx:     bz,
		PubKey: pubkey,
	}

	go func(request *eyesTypes.DispatchedTxRequest) {
		result, err := a.deyesClient.Dispatch(request)

		// Handle failure case.
		if err != nil || (result != nil && !result.Success) {
			log.Error("Deployment failed!, err = ", err)

			txOut := a.getTxOutFromSignedHash(outChain, outHash)

			if txOut == nil {
				log.Errorf("Cannot find txOut for dispath result with signed hash = %s, chain = %s", outHash, outChain)
				return
			}

			// Report this as failure. Submit to the Sisu chain
			txOutResult := &types.TxOutResult{
				OutChain: txOut.Content.OutChain,
				OutHash:  txOut.Content.OutHash,
			}
			txOutResult.Result = types.TxOutResultType_GENERIC_ERROR

			if result != nil {
				log.Verbose("Result error = ", result.Err)
				switch result.Err {
				case etypes.ErrNotEnoughBalance:
					txOutResult.Result = types.TxOutResultType_NOT_ENOUGH_NATIVE_BALANCE
				}
			}

			a.submitTxOutResult(txOutResult)
		} else {
			log.Verbose("Tx is sent to deyes!")
		}
	}(request)

	return nil
}

// OnTxIncludedInBlock implements AppLogicListener
func (a *ApiHandler) OnTxIncludedInBlock(txTrack *chainstypes.TrackUpdate) {
	log.Verbosef("Confirming tx height = %d, chain = %s, signed hash = %s, nonce = %d",
		txTrack.BlockHeight, txTrack.Chain, txTrack.Hash, txTrack.Nonce)

	txOut := a.getTxOutFromSignedHash(txTrack.Chain, txTrack.Hash)

	txOutResult := &types.TxOutResult{
		OutChain:    txOut.Content.OutChain,
		OutHash:     txOut.Content.OutHash,
		BlockHeight: txTrack.BlockHeight,
	}
	if libchain.IsETHBasedChain(txTrack.Chain) {
		ethTx := &ethtypes.Transaction{}
		err := ethTx.UnmarshalBinary(txTrack.Bytes)
		if err != nil {
			log.Error("cannot unmarshal eth transaction, err = ", err)
			return
		}

		txOutResult.Nonce = int64(ethTx.Nonce()) + 1
	}

	if txTrack.Result == chainstypes.TrackResultConfirmed {
		log.Infof("confirming tx: chain = %s, signed hash = %s, type = %v", txTrack.Chain, txTrack.Hash, txOut.TxType)
		txOutResult.Result = types.TxOutResultType_IN_BLOCK_SUCCESS
	} else {
		// Tx is included in the block but fails to execute.
		txOutResult.Result = types.TxOutResultType_IN_BLOCK_FAILURE
	}

	a.submitTxOutResult(txOutResult)
}

func (a *ApiHandler) submitTxOutResult(txOutResult *types.TxOutResult) {
	msg := types.NewTxOutResultMsg(
		a.appKeys.GetSignerAddress().String(),
		txOutResult,
	)
	a.txSubmit.SubmitMessageAsync(msg)
}

// OnUpdateTokenPrice is called when there is a token price update from deyes. Post to the network
// until we reach a consensus about token price. The token price is only used to calculate gas price
// fee and not used for actual swapping calculation.
func (a *ApiHandler) OnUpdateTokenPrice(tokenPrices []*etypes.TokenPrice) {
	prices := make([]*types.TokenPrice, 0, len(tokenPrices))

	// Convert from deyes type to msg type
	for _, token := range tokenPrices {
		prices = append(prices, &types.TokenPrice{
			Id:    token.Id,
			Price: token.Price.String(),
		})
	}

	msg := types.NewUpdateTokenPrice(a.appKeys.GetSignerAddress().String(), prices)
	a.txSubmit.SubmitMessageAsync(msg)
}
