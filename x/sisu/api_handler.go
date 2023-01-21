package sisu

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	chainstypes "github.com/sisu-network/deyes/chains/types"
	etypes "github.com/sisu-network/deyes/types"
	dhtypes "github.com/sisu-network/dheart/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/chains"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/service"
	"github.com/sisu-network/sisu/x/sisu/types"
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
	keeper        keeper.Keeper
	txSubmit      components.TxSubmit
	appKeys       components.AppKeys
	globalData    components.GlobalData
	txTracker     TxTracker
	bridgeManager chains.BridgeManager
	chainPolling  service.ChainPolling
	valManager    ValidatorManager
	mc            ManagerContainer

	// Dheart & Deyes client
	dheartClient external.DheartClient
	deyesClient  external.DeyesClient

	privateDb keeper.PrivateDb
}

func NewApiHandler(
	privateDb keeper.PrivateDb,
	mc ManagerContainer,
) *ApiHandler {
	a := &ApiHandler{
		mc:            mc,
		keeper:        mc.Keeper(),
		privateDb:     privateDb,
		appKeys:       mc.AppKeys(),
		txSubmit:      mc.TxSubmit(),
		globalData:    mc.GlobalData(),
		dheartClient:  mc.DheartClient(),
		deyesClient:   mc.DeyesClient(),
		txTracker:     mc.TxTracker(),
		chainPolling:  mc.ChainPolling(),
		bridgeManager: mc.BridgeManager(),
		valManager:    mc.ValidatorManager(),
	}

	return a
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
	)

	log.Info("There is keygen result from dheart, resultEnum = ", resultEnum, " keyType = ", result.KeyType)

	a.txSubmit.SubmitMessageAsync(signerMsg)
}

// OnTxDeploymentResult is a callback after there is a deployment result from deyes.
func (a *ApiHandler) OnTxDeploymentResult(result *etypes.DispatchedTxResult) {
	if !result.Success {
		log.Verbosef("Result from deyes: failed to deploy tx, chain = %s, signed hash = %s, error = %s",
			result.Chain, result.TxHash, result.Err)
		txOut := a.getTxOutFromSignedHash(result.Chain, result.TxHash)

		if txOut == nil {
			log.Errorf("Cannot find txOut for dispath result with signed hash = %s, chain = %s", result.TxHash, result.Chain)
			return
		}

		txOutId := txOut.GetId()

		// Report this as failure. Submit to the Sisu chain
		txOutResult := &types.TxOutResult{
			TxOutId:  txOutId,
			OutChain: txOut.Content.OutChain,
			OutHash:  txOut.Content.OutHash,
		}

		switch result.Err {
		case etypes.ErrNotEnoughBalance:
			txOutResult.Result = types.TxOutResultType_NOT_ENOUGH_NATIVE_BALANCE
		case etypes.ErrSubmitTx:
			txOutResult.Result = types.TxOutResultType_SUBMIT_TX_ERROR
		case etypes.ErrNonceNotMatched:
			txOutResult.Result = types.TxOutResultType_SUBMIT_TX_ERROR
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

			switch {
			case libchain.IsETHBasedChain(keysignMsg.OutChain):
				if err := a.processETHSigningResult(ctx, result, keysignMsg, i); err != nil {
					// TODO: Handle failure here.
					log.Error("Failed to process ETH signing result, err = ", err)
					return
				}

			case libchain.IsCardanoChain(keysignMsg.OutChain):
				if err := a.processCardanoSigningResult(ctx, result, keysignMsg, i); err != nil {
					// TODO: Handle failure here.
					log.Error("Failed to process cardano signing result, err = ", err)
					return
				}

			case libchain.IsSolanaChain(keysignMsg.OutChain):
				if err := a.processSolanaKeysignResult(ctx, result, keysignMsg, i); err != nil {
					// TODO: Handle failure here.
					log.Error("Failed to process solana signing result, err = ", err)
					return
				}

			case libchain.IsLiskChain(keysignMsg.OutChain):
				if err := a.processLiskKeysignResult(ctx, result, keysignMsg); err != nil {
					// TODO: Handle failure here.
					log.Error("Failed to process lisk signing result, err = ", err)
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

// OnTxIncludedInBlock implements AppLogicListener
func (a *ApiHandler) OnTxIncludedInBlock(txTrack *chainstypes.TrackUpdate) {
	log.Verbosef("Confirming tx height = %d, chain = %s, signed hash = %s",
		txTrack.BlockHeight, txTrack.Chain, txTrack.Hash)

	txOut := a.getTxOutFromSignedHash(txTrack.Chain, txTrack.Hash)
	txOutId := txOut.GetId()
	txOutResult := &types.TxOutResult{
		TxOutId:     txOutId,
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
