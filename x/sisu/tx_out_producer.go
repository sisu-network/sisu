package sisu

import (
	"encoding/json"
	"fmt"
	"math/big"

	ecommon "github.com/ethereum/go-ethereum/common"
	ethcommon "github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/echovl/cardano-go"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	etypes "github.com/sisu-network/deyes/types"
	hutils "github.com/sisu-network/dheart/utils"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	scardano "github.com/sisu-network/sisu/x/sisu/cardano"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/sisu-network/sisu/x/sisu/world"
)

// This structs produces transaction output based on input. For a given tx input, this struct
// produces a list (could contain only one element) of transaction output.
type TxOutputProducer interface {
	GetTxOuts(ctx sdk.Context, height int64, txIns []*types.TxIn) []*types.TxOutWithSigner

	PauseContract(ctx sdk.Context, chain string, hash string) (*types.TxOutWithSigner, error)

	ResumeContract(ctx sdk.Context, chain string, hash string) (*types.TxOutWithSigner, error)

	ContractChangeOwnership(ctx sdk.Context, chain, contractHash, newOwner string) (*types.TxOutWithSigner, error)

	ContractSetLiquidPoolAddress(ctx sdk.Context, chain, contractHash, newAddress string) (*types.TxOutWithSigner, error)

	ContractEmergencyWithdrawFund(ctx sdk.Context, chain, contractHash string, tokens []string, newOwner string) (*types.TxOutWithSigner, error)
}

type DefaultTxOutputProducer struct {
	worldState world.WorldState
	appKeys    common.AppKeys
	keeper     keeper.Keeper
	tssConfig  config.TssConfig
	txTracker  TxTracker

	// Only use for cardano chain
	cardanoConfig  config.CardanoConfig
	cardanoNetwork cardano.Network
	cardanoNode    cardano.Node
}

type transferOutData struct {
	tokenAddr   ethcommon.Address
	blockHeight int64
	destChain   string
	token       *types.Token
	recipient   string
	amount      *big.Int

	// For tx_tracker
	txIn *types.TxIn
}

type transferInData struct {
	token     ethcommon.Address
	recipient string
	amount    *big.Int
}

func NewTxOutputProducer(worldState world.WorldState, appKeys common.AppKeys, keeper keeper.Keeper,
	tssConfig config.TssConfig, cardanoConfig config.CardanoConfig, cardanoNode cardano.Node,
	txTracker TxTracker) TxOutputProducer {
	return &DefaultTxOutputProducer{
		keeper:         keeper,
		worldState:     worldState,
		appKeys:        appKeys,
		tssConfig:      tssConfig,
		txTracker:      txTracker,
		cardanoNetwork: cardanoConfig.GetCardanoNetwork(),
		cardanoNode:    cardanoNode,
	}
}

func (p *DefaultTxOutputProducer) GetTxOuts(ctx sdk.Context, height int64, txIns []*types.TxIn) []*types.TxOutWithSigner {
	outMsgs := make([]*types.TxOutWithSigner, 0)
	transferOuts := make([]*transferOutData, 0)

	// 1. Extracts all the transfers requests from the incoming transactions.
	for _, txIn := range txIns {
		// ETH chain
		if libchain.IsETHBasedChain(txIn.Chain) {
			ethTx := &ethTypes.Transaction{}

			err := ethTx.UnmarshalBinary(txIn.Serialized)
			if err != nil {
				log.Error("Failed to unmarshall eth tx. err =", err)
				continue
			}

			// Check if this is a transaction that fund our account.
			if ethTx.To() != nil && p.keeper.IsKeygenAddress(ctx, libchain.KEY_TYPE_ECDSA, ethTx.To().String()) {
				txOuts, err := p.getContractDeploymentTx(ctx, txIn.BlockHeight, txIn)
				if err != nil {
					log.Error("Failed to get contract deployment tx, err = ", err)
				} else {
					outMsgs = append(outMsgs, txOuts...)
				}

				continue
			}

			// Check if this is a transaction to transfer
			if p.keeper.IsContractExistedAtAddress(ctx, txIn.Chain, ethTx.To().String()) && len(ethTx.Data()) >= 4 {
				transfer, err := parseEthTransferOut(ethTx, txIn.Chain, p.worldState)
				if err != nil {
					log.Error("faield to parse parseEthTransferOut, err = ", parseEthTransferOut)
					continue
				}
				transfer.txIn = txIn

				transferOuts = append(transferOuts, transfer)
			}
		}

		// Cardano chain
		if libchain.IsCardanoChain(txIn.Chain) {
			transfer, err := p.parseCardanoTxIn(ctx, txIn)
			if err != nil {
				log.Error("Failed to parse cardano transaction, err = ", err)
				continue
			}
			transfer.txIn = txIn

			transferOuts = append(transferOuts, transfer)
		}
	}

	transferOutMsgs := p.processTransferOut(ctx, transferOuts)
	outMsgs = append(outMsgs, transferOutMsgs...)

	return outMsgs
}

func (p *DefaultTxOutputProducer) getContractDeploymentTx(ctx sdk.Context, height int64, tx *types.TxIn) ([]*types.TxOutWithSigner, error) {
	outMsgs := make([]*types.TxOutWithSigner, 0)

	contracts := p.keeper.GetPendingContracts(ctx, tx.Chain)
	log.Verbose("len(contracts) = ", len(contracts))

	if len(contracts) > 0 {
		// TODO: Check balance required to deploy all these contracts. Also check if we are deploying
		// a contract to avoid duplication.

		// Get the list of deploy transactions. Those txs need to posted and verified (by validators)
		// to the Sisu chain.
		outEthTxs := p.getEthContractDeploymentTx(ctx, height, tx.Chain, contracts)

		for i, outTx := range outEthTxs {
			bz, err := outTx.MarshalBinary()
			if err != nil {
				return nil, err
			}

			outMsg := types.NewMsgTxOutWithSigner(
				p.appKeys.GetSignerAddress().String(),
				types.TxOutType_CONTRACT_DEPLOYMENT,
				tx.BlockHeight,
				tx.Chain,
				tx.TxHash,
				tx.Chain,
				outTx.Hash().String(),
				bz,
				contracts[i].Hash,
			)

			log.Verbose("ETH Tx Out hash = ", outTx.Hash().String(), " on chain ", tx.Chain)

			outMsgs = append(outMsgs, outMsg)
		}
	}

	return outMsgs, nil
}

// In Cardano chain, transferring multi-asset required at least 1 ADA (10^6 lovelace)
func (p *DefaultTxOutputProducer) getCardanoTx(ctx sdk.Context, data *transferOutData, assetAmount uint64) (*cardano.Tx, error) {
	// Subtract commission fee
	commissionFeeRate := float32(0)
	params := p.keeper.GetParams(ctx)
	if params != nil {
		commissionFeeRate = params.CommissionFeeRate
	}
	assetAmount = assetAmount - getCommissionFee(big.NewInt(int64(assetAmount)), commissionFeeRate).Uint64()

	pubkey := p.keeper.GetKeygenPubkey(ctx, libchain.KEY_TYPE_EDDSA)
	senderAddr := hutils.GetAddressFromCardanoPubkey(pubkey)
	log.Debug("cardano sender address = ", senderAddr.String())

	receiverAddr, err := cardano.NewAddress(data.recipient)
	if err != nil {
		log.Error("error when parsing receiver addr: ", err)
		return nil, err
	}

	multiAsset, err := scardano.GetCardanoMultiAsset(data.destChain, data.token, assetAmount)
	if err != nil {
		return nil, err
	}

	adaPrice, err := p.worldState.GetNativeTokenPriceForChain(data.destChain)
	if err != nil {
		log.Error("error when getting ada price: ", err)
		return nil, err
	}

	// We need at least 1 ada to send multi assets.
	tx, err := scardano.BuildTx(p.cardanoNode, p.cardanoNetwork, senderAddr, receiverAddr,
		cardano.NewValueWithAssets(cardano.Coin(utils.ONE_ADA_IN_LOVELACE.Uint64()), multiAsset), nil,
		adaPrice, data.token, data.destChain, assetAmount)

	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (p *DefaultTxOutputProducer) parseCardanoTxIn(ctx sdk.Context, tx *types.TxIn) (*transferOutData, error) {
	txIn := &etypes.CardanoTxInItem{}
	if err := json.Unmarshal(tx.Serialized, txIn); err != nil {
		log.Error("error when marshaling cardano tx in item: ", err)
		return nil, err
	}

	extraInfo := txIn.Metadata
	tokenName := txIn.Asset
	if tokenName != "ADA" {
		tokenName = tokenName[5:] // Remove the WRAP_ prefix
	}
	if len(tokenName) == 0 {
		return nil, fmt.Errorf("Invalid token: %s", tokenName)
	}

	tokens := p.keeper.GetTokens(ctx, []string{tokenName})
	token := tokens[tokenName]
	if token == nil {
		return nil, fmt.Errorf("Cannot find token in the keeper")
	}

	amount := new(big.Int).SetUint64(txIn.Amount)

	return &transferOutData{
		blockHeight: tx.BlockHeight,
		destChain:   extraInfo.Chain,
		recipient:   extraInfo.Recipient,
		token:       token,
		amount:      utils.LovelaceToWei(amount),
	}, nil
}

func (p *DefaultTxOutputProducer) processTransferOut(ctx sdk.Context, transfers []*transferOutData) []*types.TxOutWithSigner {
	outMsgs := make([]*types.TxOutWithSigner, 0)

	for _, transfer := range transfers {
		var bz []byte
		var txHash string

		// ETH chain
		if libchain.IsETHBasedChain(transfer.destChain) {
			responseTx, err := p.buildERC20TransferIn(ctx, transfer.token, transfer.tokenAddr, ecommon.HexToAddress(transfer.recipient),
				transfer.amount, transfer.destChain)
			if err != nil {
				log.Error("Failed to build erc20 transfer in, err = ", err)
				continue
			}

			bz = responseTx.RawBytes
			txHash = responseTx.EthTx.Hash().String()
		}

		// Cardano chain
		if libchain.IsCardanoChain(transfer.destChain) {
			multiAssetAmt := utils.WeiToLovelace(transfer.amount)
			log.Verbosef("data.amount = %v, multiAssetAmt = %v", transfer.amount, multiAssetAmt)

			// In real, this transaction transfers at least <1 ADA + additional tx fee>
			cardanoTx, err := p.getCardanoTx(ctx, transfer, multiAssetAmt.Uint64())
			if err != nil {
				log.Error("Failed to get cardano tx, err  = ", err)
				continue
			}

			cardanoTxHash, err := cardanoTx.Hash()
			if err != nil {
				log.Error("Failed to get cardano hash, err = ", err)
				continue
			}

			bz, err = cardanoTx.MarshalCBOR()
			if err != nil {
				log.Error("Faield to marshalcbor cardano tx, err = ", err)
				continue
			}

			txHash = cardanoTxHash.String()
		}

		// Add to the output array
		if bz != nil {
			outMsg := types.NewMsgTxOutWithSigner(
				p.appKeys.GetSignerAddress().String(),
				types.TxOutType_TRANSFER_OUT,
				transfer.blockHeight,
				transfer.destChain,
				"",
				transfer.destChain,
				txHash,
				bz,
				"",
			)

			// Track the txout
			p.txTracker.AddTransaction(
				outMsg.Data,
				transfer.txIn,
			)

			outMsgs = append(outMsgs, outMsg)
		}
	}

	return outMsgs
}

func (p *DefaultTxOutputProducer) getGasLimit(chain string) uint64 {
	// TODO: Make this dependent on different chains.
	return uint64(8_000_000)
}

// @Deprecated
func (p *DefaultTxOutputProducer) getDefaultGasPrice(chain string) *big.Int {
	// TODO: Make this dependent on different chains.
	switch chain {
	case "ganache1":
		return big.NewInt(2_000_000_000)
	case "ganache2":
		return big.NewInt(2_000_000_000)
	case "ropsten-testnet":
		return big.NewInt(4_000_000_000)
	case "binance-testnet":
		return big.NewInt(18_000_000_000)
	case "polygon-testnet":
		return big.NewInt(7_000_000_000)
	case "xdai":
		return big.NewInt(2_000_000_000)
	case "goerli-testnet":
		return big.NewInt(1_500_000_000)
	case "eth":
		return big.NewInt(70_000_000_000)
	case "arbitrum-testnet":
		return big.NewInt(50_000_000)
	case "fantom-testnet":
		return big.NewInt(75_000_000_000)
	}
	return big.NewInt(100_000_000_000)
}
