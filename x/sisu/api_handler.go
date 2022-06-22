package sisu

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/echovl/cardano-go"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	etypes "github.com/sisu-network/deyes/types"
	eyesTypes "github.com/sisu-network/deyes/types"
	dhtypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/lib/chain"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/sisu-network/sisu/x/sisu/world"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
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

type BlockSymbolPair struct {
	blockHeight int64
	chain       string
}

// ApiHandler handles API callback from dheart or deyes. There are few functions (BeginBlock & EndBlock)
// that are still present for historical reason. They should be moved out of this file.
type ApiHandler struct {
	keeper     keeper.Keeper
	config     config.TssConfig
	txSubmit   common.TxSubmit
	appKeys    common.AppKeys
	globalData common.GlobalData
	txTracker  TxTracker
	mc         ManagerContainer

	// Dheart & Deyes client
	dheartClient tssclients.DheartClient
	deyesClient  tssclients.DeyesClient

	worldState       world.WorldState
	keygenBlockPairs []BlockSymbolPair

	privateDb keeper.Storage
}

func NewApiHandler(
	privateDb keeper.Storage,
	mc ManagerContainer,
) *ApiHandler {
	a := &ApiHandler{
		keeper:     mc.Keeper(),
		privateDb:  privateDb,
		appKeys:    mc.AppKeys(),
		config:     mc.Config(),
		txSubmit:   mc.TxSubmit(),
		globalData: mc.GlobalData(),
		// And array that stores block numbers where we should do final vote count.
		keygenBlockPairs: make([]BlockSymbolPair, 0),
		dheartClient:     mc.DheartClient(),
		deyesClient:      mc.DeyesClient(),
		worldState:       mc.WorldState(),
		txTracker:        mc.TxTracker(),
		mc:               mc,
	}

	return a
}

// TODO: Move this function to module.go
func (a *ApiHandler) BeginBlock(ctx sdk.Context, blockHeight int64) {
	// Check keygen proposal
	if blockHeight > 1 {
		// We need to wait till block 2 for multistore of the app to be updated with latest account info
		// for signing.
		a.CheckTssKeygen(ctx, blockHeight)
	}

	oldValue := a.globalData.IsCatchingUp()
	a.globalData.UpdateCatchingUp()
	newValue := a.globalData.IsCatchingUp()

	if oldValue && !newValue {
		log.Info("Setting Sisu readiness for dheart.")
		// This node has fully catched up with the blockchain, we need to inform dheart about this.
		a.dheartClient.SetSisuReady(true)
		a.deyesClient.SetSisuReady(true)
	}

	// TODO: Make keygen to be command instead of embedding inside the code.
	// Check Vote result.
	for len(a.keygenBlockPairs) > 0 && !a.globalData.IsCatchingUp() {
		log.Verbose("blockHeight = ", blockHeight)
		if blockHeight < a.keygenBlockPairs[0].blockHeight {
			break
		}

		for len(a.keygenBlockPairs) > 0 && blockHeight >= a.keygenBlockPairs[0].blockHeight {
			// Remove the chain from processing queue.
			a.keygenBlockPairs = a.keygenBlockPairs[1:]
		}
	}

	// Calculate all token prices.
	a.calculateTokenPrices(ctx)
}

// calculateTokenPrices gets all token prices posted from all validators and calculate the median.
func (a *ApiHandler) calculateTokenPrices(ctx sdk.Context) {
	curBlock := ctx.BlockHeight()

	// We wait for 5 more blocks after we get prices from deyes so that any record can be posted
	// onto the blockchain.
	if curBlock%TokenPriceUpdateInterval != 5 {
		return
	}

	log.Info("Calcuating token prices....")

	// TODO: Fix the signer set.
	records := a.keeper.GetAllTokenPricesRecord(ctx)

	tokenPrices := make(map[string][]int64)
	for _, data := range records {
		for _, record := range data.Records {
			// Only calculate token prices that has been updated recently.
			if curBlock-int64(record.BlockHeight) > TokenPriceUpdateInterval {
				continue
			}

			m := tokenPrices[record.Token]
			if m == nil {
				m = make([]int64, 0)
			}

			m = append(m, record.Price)

			tokenPrices[record.Token] = m
		}
	}

	// Now sort all the array and get the median
	medians := make(map[string]int64)
	for token, list := range tokenPrices {
		if len(list) == 0 {
			log.Error("cannot find price list for token ", token)
			continue
		}

		sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
		median := list[len(list)/2]
		medians[token] = median
	}

	log.Verbose("Calculated prices = ", medians)

	// Update all the token data.
	arr := make([]string, 0, len(medians))
	for token, _ := range medians {
		arr = append(arr, token)
	}

	savedTokens := a.keeper.GetTokens(ctx, arr)

	for tokenId, price := range medians {
		savedTokens[tokenId].Price = price
	}

	a.keeper.SetTokens(ctx, savedTokens)

	// Update the world state
	a.worldState.SetTokens(savedTokens)
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

	// Add list the public key address to watch.
	for _, chainConfig := range a.config.SupportedChains {
		chain := chainConfig.Id

		if libchain.GetKeyTypeForChain(chain) != result.KeyType {
			continue
		}

		log.Verbose("adding watcher address ", result.Address, " for chain ", chain)

		a.addWatchAddress(chain, result.Address)
	}
}

func (a *ApiHandler) addWatchAddress(chain string, address string) {
	a.deyesClient.AddWatchAddresses(chain, []string{address})
}

// OnTxDeploymentResult is a callback after there is a deployment result from deyes.
func (a *ApiHandler) OnTxDeploymentResult(result *etypes.DispatchedTxResult) {
	log.Info("The transaction has been sent to blockchain (but not included in a block yet). chain = ",
		result.Chain, ", address = ", result.DeployedAddr)
	a.txTracker.UpdateStatus(result.Chain, result.TxHash, types.TxStatusDepoyed)
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

	ctx := a.mc.GetReadOnlyContext()

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

func (a *ApiHandler) processETHSigningResult(ctx sdk.Context, result *dhtypes.KeysignResult, signMsg *dhtypes.KeysignMessage, index int) {
	// Find the tx in txout table
	txOut := a.keeper.GetTxOut(ctx, signMsg.OutChain, signMsg.OutHash)
	if txOut == nil {
		log.Error("Cannot find tx out with hash", signMsg.OutHash)
		return
	}

	tx := &ethtypes.Transaction{}
	if err := tx.UnmarshalBinary(txOut.OutBytes); err != nil {
		log.Error("cannot unmarshal tx, err =", err)
		return
	}

	// Create full tx with signature.
	chainId := libchain.GetChainIntFromId(signMsg.OutChain)
	if len(result.Signatures[index]) != 65 {
		log.Error("Signature length is not 65 for chain: ", chainId)
	}
	signedTx, err := tx.WithSignature(ethtypes.NewEIP2930Signer(chainId), result.Signatures[index])
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
	a.privateDb.SaveTxOutSig(&types.TxOutSig{
		Chain:       signMsg.OutChain,
		HashWithSig: signedTx.Hash().String(),
		HashNoSig:   signMsg.OutHash,
	})

	log.Info("signedTx hash = ", signedTx.Hash().String())

	// If this is a contract deployment transaction, update the contract table with the hash of the
	// deployment tx bytes.
	isContractDeployment := chain.IsETHBasedChain(signMsg.OutChain) && txOut.TxType == types.TxOutType_CONTRACT_DEPLOYMENT
	err = a.deploySignedTx(ctx, bz, signMsg.OutChain, result.Request.KeysignMessages[index].OutHash, isContractDeployment)
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
	if err := tx.UnmarshalCBOR(txOut.OutBytes); err != nil {
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

	// TODO: temporary use json to encode/decode. In the future, should use cbor instead
	txBytes, err := json.Marshal(tx)
	if err != nil {
		log.Error("error when marshal cardano tx: ", err)
		return err
	}

	err = a.deploySignedTx(ctx, txBytes, signMsg.OutChain, result.Request.KeysignMessages[index].OutHash, false)
	if err != nil {
		log.Error("deployment error: ", err)
		return err
	}

	log.Info("Sent signed cardano tx to deyes")

	return nil
}

// deploySignedTx creates a deployment request and sends it to deyes.
func (a *ApiHandler) deploySignedTx(ctx sdk.Context, bz []byte, outChain string, outHash string, isContractDeployment bool) error {
	log.Verbose("Sending final tx to the deyes for deployment for chain ", outChain)

	pubkey := a.keeper.GetKeygenPubkey(ctx, libchain.GetKeyTypeForChain(outChain))
	if pubkey == nil {
		return fmt.Errorf("Cannot get pubkey for chain %s", outChain)
	}

	request := &etypes.DispatchedTxRequest{
		Chain:                   outChain,
		TxHash:                  outHash,
		Tx:                      bz,
		PubKey:                  pubkey,
		IsEthContractDeployment: isContractDeployment,
	}

	go a.deyesClient.Dispatch(request)

	return nil
}

// Processed list of transactions sent from deyes to Sisu api server.
// TODO: handle error correctly
func (a *ApiHandler) OnTxIns(txs *eyesTypes.Txs) error {
	log.Verbose("There is a new list of txs from deyes, len =", len(txs.Arr))

	ctx := a.mc.GetReadOnlyContext()

	// Create TxIn messages and broadcast to the Sisu chain.
	for _, tx := range txs.Arr {
		if !tx.Success {
			// TODO: Have a mechanism to handle failed transaction.
			a.txTracker.OnTxFailed(txs.Chain, tx.Hash, types.TxStatusReverted)
			continue
		}

		// 1. Check if this tx is from one of our key. If it is, update the status of TxOut to confirmed.
		if a.keeper.IsKeygenAddress(ctx, libchain.KEY_TYPE_ECDSA, tx.From) {
			return a.confirmTx(ctx, tx, txs.Chain, txs.Block)
		} else if len(tx.To) > 0 {
			// 2. This is a transaction to our key account or one of our contracts. Create a message to
			// indicate that we have observed this transaction and broadcast it to cosmos chain.
			// TODO: handle error correctly
			hash := utils.GetTxInHash(txs.Block, txs.Chain, tx.Serialized)
			signerMsg := types.NewTxInWithSigner(
				a.appKeys.GetSignerAddress().String(),
				txs.Chain,
				hash,
				txs.Block,
				tx.Serialized,
			)

			a.txSubmit.SubmitMessageAsync(signerMsg)
		}
	}

	return nil
}

// confirmTx confirms that a tx has been included in a block on the blockchain.
func (a *ApiHandler) confirmTx(ctx sdk.Context, tx *eyesTypes.Tx, chain string, blockHeight int64) error {
	log.Verbose("This is a transaction from us. We need to confirm it. Chain = ", chain)

	// The txOutSig is in private db while txOut should come from common db.
	txOutSig := a.privateDb.GetTxOutSig(chain, tx.Hash)
	if txOutSig == nil {
		// TODO: Add this to pending tx to confirm.
		log.Verbose("cannot find txOutSig with full signature hash: ", tx.Hash)
		return nil
	}

	txOut := a.keeper.GetTxOut(ctx, chain, txOutSig.HashNoSig)
	if txOut == nil {
		log.Verbose("cannot find txOut with hash (with no sig): ", txOutSig.HashNoSig)
		return nil
	}

	log.Info("confirming tx: chain, hash, type = ", chain, " ", tx.Hash, " ", txOut.TxType)

	a.txTracker.RemoveTransaction(chain, txOut.OutHash)

	contractAddress := ""
	if txOut.TxType == types.TxOutType_CONTRACT_DEPLOYMENT && libchain.IsETHBasedChain(chain) {
		ethTx := &ethTypes.Transaction{}
		err := ethTx.UnmarshalBinary(tx.Serialized)
		if err != nil {
			log.Error("cannot unmarshal eth transaction, err = ", err)
			return err
		}

		contractAddress = ethcrypto.CreateAddress(ethcommon.HexToAddress(tx.From), ethTx.Nonce()).String()
		log.Info("contractAddress = ", contractAddress)

		txConfirm := &types.TxOutContractConfirm{
			OutChain:        txOut.OutChain,
			OutHash:         txOut.OutHash,
			BlockHeight:     blockHeight,
			ContractAddress: contractAddress,
		}

		msg := types.NewTxOutContractConfirmWithSigner(
			a.appKeys.GetSignerAddress().String(),
			txConfirm,
		)
		a.txSubmit.SubmitMessageAsync(msg)
	}

	// We can assume that other tx transactions will succeed in majority of the time. Instead
	// broadcasting the tx confirmation to Sisu blockchain, we should only record missing or failed
	// transaction.
	// We only confirm if the tx out is a contract deployment to save the smart contract address.
	// TODO: Implement missing/ failed message and broadcast that to everyone after we have not seen
	// a tx for some blocks.

	return nil
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
			Price: int64(token.Price * utils.DecinmalUnit),
		})
	}

	msg := types.NewUpdateTokenPrice(a.appKeys.GetSignerAddress().String(), prices)
	a.txSubmit.SubmitMessageAsync(msg)
}
