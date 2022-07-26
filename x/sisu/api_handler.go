package sisu

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"

	"github.com/echovl/cardano-go"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	chainstypes "github.com/sisu-network/deyes/chains/types"
	etypes "github.com/sisu-network/deyes/types"
	eyesTypes "github.com/sisu-network/deyes/types"
	dhtypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/lib/chain"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	scardano "github.com/sisu-network/sisu/x/sisu/cardano"
	"github.com/sisu-network/sisu/x/sisu/eth"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"

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

	tokenPrices := make(map[string][]*big.Int)
	for _, data := range records {
		for _, record := range data.Records {
			// Only calculate token prices that has been updated recently.
			if curBlock-int64(record.BlockHeight) > TokenPriceUpdateInterval {
				continue
			}

			m := tokenPrices[record.Token]
			if m == nil {
				m = make([]*big.Int, 0)
			}

			value, _ := new(big.Int).SetString(record.Price, 10)
			m = append(m, value)

			tokenPrices[record.Token] = m
		}
	}

	// Now sort all the array and get the median
	medians := make(map[string]*big.Int)
	for token, list := range tokenPrices {
		if len(list) == 0 {
			log.Error("cannot find price list for token ", token)
			continue
		}

		sort.Slice(list, func(i, j int) bool {
			return list[i].Cmp(list[j]) < 0
		})
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
		savedTokens[tokenId].Price = price.String()
	}

	a.keeper.SetTokens(ctx, savedTokens)
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

		a.deyesClient.SetChainAccount(chain, result.Address)
		if libchain.IsCardanoChain(chain) {
			a.deyesClient.SetGatewayAddress(chain, result.Address)
		}
	}
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
	if err := tx.UnmarshalBinary(result.Request.KeysignMessages[index].Bytes); err != nil {
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
	log.Verbose("ETH keysign result signMsg.OutHash = ", signMsg.OutHash)
	a.privateDb.SaveTxOutSig(&types.TxOutSig{
		Chain:       signMsg.OutChain,
		HashWithSig: signedTx.Hash().String(),
		HashNoSig:   signMsg.OutHash,
	})

	log.Info("signedTx hash = ", signedTx.Hash().String())

	// If this is a contract deployment transaction, update the contract table with the hash of the
	// deployment tx bytes.
	isContractDeployment := chain.IsETHBasedChain(signMsg.OutChain) &&
		txOut.TxType == types.TxOutType_CONTRACT_DEPLOYMENT
	err = a.deploySignedTx(ctx, bz, signMsg.OutChain, signedTx.Hash().String(), isContractDeployment)
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
	txBytes, err := tx.MarshalCBOR()
	if err != nil {
		log.Error("error when marshal cardano tx: ", err)
		return err
	}
	hash, err := tx.Hash()
	if err != nil {
		return nil
	}

	err = a.deploySignedTx(ctx, txBytes, signMsg.OutChain, result.Request.KeysignMessages[index].OutHash, false)
	if err != nil {
		log.Error("deployment error: ", err)
		return err
	}

	log.Info("Sent signed cardano tx to deyes, tx hash = ", hash)

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

	go func(request *eyesTypes.DispatchedTxRequest) {
		result, err := a.deyesClient.Dispatch(request)
		if err != nil {
			log.Error("Failed to deploy, err = ", err)
			return
		}
		if result != nil && !result.Success {
			log.Error("Deployment failed!, err = ", result.Err)
		}
	}(request)

	return nil
}

// Processed list of transactions sent from deyes to Sisu api server.
// TODO: handle error correctly
func (a *ApiHandler) OnTxIns(txs *eyesTypes.Txs) error {
	log.Verbose("There is a new list of txs from deyes, len =", len(txs.Arr))

	transferRequests := &types.TxsIn{
		Chain:    txs.Chain,
		Hash:     txs.BlockHash,
		Height:   txs.Block,
		Requests: make([]*types.TxIn, 0),
	}

	ctx := a.globalData.GetReadOnlyContext()

	// Create TxIn messages and broadcast to the Sisu chain.
	// TODO: Prevent submitting fund gateway multiple times.
	gatewayFundTxSubmitted := false
	for _, tx := range txs.Arr {
		// Check if this is a transaciton that fund ETH gateway
		if libchain.IsETHBasedChain(txs.Chain) && !gatewayFundTxSubmitted {
			ethTx := &ethTypes.Transaction{}
			err := ethTx.UnmarshalBinary(tx.Serialized)
			if err != nil {
				log.Error("Failed to unmarshall eth tx. err =", err)
				continue
			}

			if ethTx.To() != nil && a.keeper.IsKeygenAddress(ctx, libchain.KEY_TYPE_ECDSA, ethTx.To().String()) {
				msg := types.NewFundGatewayMsg(a.appKeys.GetSignerAddress().String(), &types.FundGateway{
					Chain:  txs.Chain,
					TxHash: utils.KeccakHash32Bytes(tx.Serialized),
					Amount: ethTx.Value().Bytes(),
				})

				// For contract deployment
				a.txSubmit.SubmitMessageAsync(msg)
				gatewayFundTxSubmitted = true
				continue
			}
		}

		fmt.Println("Parsing transfer request....")
		txIns, err := a.parseTransferRequest(ctx, txs.Chain, tx)
		if err != nil {
			log.Error("Faield to parse transfer, err = ", err)
			continue
		}

		log.Verbose("Len(txIns) = ", len(txIns), " on chain ", txs.Chain)
		if txIns != nil {
			transferRequests.Requests = append(transferRequests.Requests, txIns...)
		}
	}

	if len(transferRequests.Requests) > 0 {
		msg := types.NewTxsInMsg(a.appKeys.GetSignerAddress().String(), transferRequests)
		a.txSubmit.SubmitMessageAsync(msg)
	}

	if libchain.IsCardanoChain(txs.Chain) {
		log.Verbose("Updating block height for cardano")
		// Broadcast blockheight update
		msg := types.NewBlockHeightMsg(a.appKeys.GetSignerAddress().String(), &types.BlockHeight{
			Chain:  txs.Chain,
			Height: txs.Block,
			Hash:   txs.BlockHash,
		})
		a.txSubmit.SubmitMessageAsync(msg)
	}

	return nil
}

func (a *ApiHandler) parseTransferRequest(ctx sdk.Context, chain string, tx *eyesTypes.Tx) ([]*types.TxIn, error) {
	if libchain.IsETHBasedChain(chain) {
		erc20gatewayContract := SupportedContracts[ContractErc20Gateway]
		gwAbi := erc20gatewayContract.Abi

		ethTx := &ethTypes.Transaction{}
		err := ethTx.UnmarshalBinary(tx.Serialized)
		if err != nil {
			log.Error("Failed to unmarshall eth tx. err =", err)
			return nil, err
		}

		transfer, err := eth.ParseEthTransferOut(ctx, ethTx, chain, gwAbi, a.keeper)
		if err != nil {
			return nil, err
		}

		return []*types.TxIn{transfer}, nil
	}

	if libchain.IsCardanoChain(chain) {
		ret := make([]*types.TxIn, 0)
		cardanoTx := &etypes.CardanoTransactionUtxo{}
		err := json.Unmarshal(tx.Serialized, cardanoTx)
		if err != nil {
			return nil, err
		}

		if cardanoTx != nil {
			fmt.Println("cardanoTx = ", *cardanoTx)
		} else {
			log.Error("cardanoTx is nil")
		}

		if cardanoTx.Metadata != nil {
			fmt.Println("cardanoTx.Amount = ", cardanoTx.Amount)
			// Convert from ADA unit (10^6) to our standard unit (10^18)
			for i, amount := range cardanoTx.Amount {
				if i == len(cardanoTx.Amount)-1 {
					// The last transaction returns the remaining token to the sender (or any other address).
					break
				}
				fmt.Println("AAAAAAA 0000000, amount = ", amount)
				quantity, ok := new(big.Int).SetString(amount.Quantity, 10)
				if !ok {
					log.Error("Failed to get amount quantity in cardano tx")
					continue
				}
				quantity = utils.LovelaceToWei(quantity)

				// Remove the word wrap
				tokenUnit := amount.Unit
				if tokenUnit != "lovelace" {
					token := scardano.GetTokenFromCardanoAsset(ctx, a.keeper, tokenUnit, chain)
					if token == nil {
						log.Error("Failed to find token with id: ", tokenUnit)
						continue
					}
					tokenUnit = token.Id
				} else {
					tokenUnit = "ADA"
				}

				fmt.Println("tokenUnit = ", tokenUnit, " quantity = ", quantity)
				fmt.Println("cardanoTx.Metadata = ", cardanoTx.Metadata)

				ret = append(ret, &types.TxIn{
					Hash:      cardanoTx.Hash,
					ToChain:   cardanoTx.Metadata.Chain,
					Token:     tokenUnit,
					Recipient: cardanoTx.Metadata.Recipient,
					Amount:    quantity.String(),
				})
			}
		}

		return ret, nil
	}

	return nil, fmt.Errorf("Unknown chain %s", chain)
}

// ConfirmTx implements AppLogicListener
func (a *ApiHandler) ConfirmTx(txTrack *chainstypes.TrackUpdate) {
	ctx := a.globalData.GetReadOnlyContext()

	log.Verbosef("Confirming tx height = %d, chain = %s, hash = %s, nonce = %d",
		txTrack.BlockHeight, txTrack.Chain, txTrack.Hash, txTrack.Nonce)

	// The txOutSig is in private db while txOut should come from common db.
	txOutSig := a.privateDb.GetTxOutSig(txTrack.Chain, txTrack.Hash)
	if txOutSig == nil {
		log.Error("cannot find txOutSig with full signature hash: ", txTrack.Hash)
		return
	}

	txOut := a.keeper.GetTxOut(ctx, txTrack.Chain, txOutSig.HashNoSig)
	if txOut == nil {
		log.Verbose("cannot find txOut with hash (with no sig): ", txOutSig.HashNoSig)
		return
	}
	log.Info("confirming tx: chain, hash, type = ", txTrack.Chain, " ", txTrack.Hash, " ", txOut.TxType)
	a.txTracker.RemoveTransaction(txTrack.Chain, txOut.OutHash)

	txConfirm := &types.TxOutConfirm{
		OutChain:    txOut.OutChain,
		OutHash:     txOut.OutHash,
		BlockHeight: txTrack.BlockHeight,
	}

	if libchain.IsETHBasedChain(txTrack.Chain) {
		ethTx := &ethTypes.Transaction{}
		err := ethTx.UnmarshalBinary(txTrack.Bytes)
		if err != nil {
			log.Error("cannot unmarshal eth transaction, err = ", err)
			return
		}

		txConfirm.Nonce = int64(ethTx.Nonce()) + 1
		if txOut.TxType == types.TxOutType_CONTRACT_DEPLOYMENT {
			sender, err := utils.GetEthSender(ethTx)
			if err != nil {
				log.Error("cannot get eth sender, err = ", err)
				return
			}

			contractAddress := ethcrypto.CreateAddress(sender, ethTx.Nonce()).String()
			log.Info("contractAddress = ", contractAddress)

			txConfirm.ContractAddress = contractAddress
		}
	}

	msg := types.NewTxOutConfirmMsg(
		a.appKeys.GetSignerAddress().String(),
		txConfirm,
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
