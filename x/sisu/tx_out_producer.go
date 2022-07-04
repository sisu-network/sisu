package sisu

import (
	"fmt"
	"math/big"

	"github.com/sisu-network/sisu/utils"
	scardano "github.com/sisu-network/sisu/x/sisu/cardano"

	ecommon "github.com/ethereum/go-ethereum/common"
	ethcommon "github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/echovl/cardano-go"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/sisu-network/sisu/x/sisu/world"
)

// This structs produces transaction output based on input. For a given tx input, this struct
// produces a list (could contain only one element) of transaction output.
type TxOutputProducer interface {
	// GetTxOuts returns a list of TxOut message and a list of un-processed transfer out request that
	// needs to be processed next time.
	GetTxOuts(ctx sdk.Context, transfers []*transferOutData) ([]*types.TxOutWithSigner, []*transferOutData)

	PauseContract(ctx sdk.Context, chain string, hash string) (*types.TxOutWithSigner, error)

	ResumeContract(ctx sdk.Context, chain string, hash string) (*types.TxOutWithSigner, error)

	ContractChangeOwnership(ctx sdk.Context, chain, contractHash, newOwner string) (*types.TxOutWithSigner, error)

	ContractSetLiquidPoolAddress(ctx sdk.Context, chain, contractHash, newAddress string) (*types.TxOutWithSigner, error)

	ContractEmergencyWithdrawFund(ctx sdk.Context, chain, contractHash string, tokens []string, newOwner string) (*types.TxOutWithSigner, error)
}

type DefaultTxOutputProducer struct {
	worldState  world.WorldState
	appKeys     common.AppKeys
	keeper      keeper.Keeper
	tssConfig   config.TssConfig
	txTracker   TxTracker
	valsManager ValidatorManager
	privateDb   keeper.Storage

	// Only use for cardano chain
	cardanoConfig        config.CardanoConfig
	cardanoNetwork       cardano.Network
	cardanoClient        scardano.CardanoClient
	minCaranoBlockHeight int64
}

type transferOutData struct {
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
	valsManager ValidatorManager, tssConfig config.TssConfig, cardanoConfig config.CardanoConfig,
	privateDb keeper.Storage, cardanoClient scardano.CardanoClient,
	txTracker TxTracker) TxOutputProducer {
	return &DefaultTxOutputProducer{
		keeper:         keeper,
		worldState:     worldState,
		appKeys:        appKeys,
		valsManager:    valsManager,
		tssConfig:      tssConfig,
		privateDb:      privateDb,
		txTracker:      txTracker,
		cardanoNetwork: cardanoConfig.GetCardanoNetwork(),
		cardanoClient:  cardanoClient,
	}
}

func (p *DefaultTxOutputProducer) parseCardanoTxIn(ctx sdk.Context, tx *types.TxIn) (*transferOutData, error) {
	return nil, nil

	// txIn := &etypes.CardanoTxInItem{}
	// if err := json.Unmarshal(tx.Serialized, txIn); err != nil {
	// 	log.Error("error when marshaling cardano tx in item: ", err)
	// 	return nil, err
	// }

	// extraInfo := txIn.Metadata
	// tokenName := txIn.Asset
	// if tokenName != "ADA" {
	// 	tokenName = tokenName[5:] // Remove the WRAP_ prefix
	// }
	// if len(tokenName) == 0 {
	// 	return nil, fmt.Errorf("Invalid token: %s", tokenName)
	// }

	// tokens := p.keeper.GetTokens(ctx, []string{tokenName})
	// token := tokens[tokenName]
	// if token == nil {
	// 	return nil, fmt.Errorf("Cannot find token in the keeper")
	// }

	// amount := new(big.Int).SetUint64(txIn.Amount)

	// return &transferOutData{
	// 	blockHeight: tx.BlockHeight,
	// 	destChain:   extraInfo.Chain,
	// 	recipient:   extraInfo.Recipient,
	// 	token:       token,
	// 	amount:      utils.LovelaceToWei(amount),
	// }, nil
}

func (p *DefaultTxOutputProducer) GetTxOuts(ctx sdk.Context, transfers []*transferOutData) ([]*types.TxOutWithSigner, []*transferOutData) {
	outMsgs := make([]*types.TxOutWithSigner, 0)

	params := p.keeper.GetParams(ctx)

	// Categories txs by their destination chains.
	transfersByChains := p.categorizeTransfer(transfers)
	for _, transfersInOneChain := range transfersByChains {
		// TODO: Don't use fixed batch size. Let the batch size dependent on the current data on-chain.
		dstChain := transfersInOneChain[0].destChain
		batchSize := params.GetMaxTransferOutBatch(dstChain)
		batches := splitTransfers(transfersInOneChain, batchSize)

		if libchain.IsETHBasedChain(dstChain) {
			msgs := p.processEthBatches(ctx, batches)
			outMsgs = append(outMsgs, msgs...)
		}

		if libchain.IsCardanoChain(dstChain) {
			msgs := p.processCardanoBatches(ctx, dstChain, batches)
			outMsgs = append(outMsgs, msgs...)
		}
	}

	return outMsgs, nil
}

func splitTransfers(transfers []*transferOutData, batchSize int) [][]*transferOutData {
	allBatches := make([][]*transferOutData, 0)
	var batch []*transferOutData
	for i := range transfers {
		if i%batchSize == 0 {
			if i > 0 {
				allBatches = append(allBatches, batch)
			}
			batch = make([]*transferOutData, 0)
		}
		batch = append(batch, transfers[i])
	}

	allBatches = append(allBatches, batch)

	return allBatches
}

// categorizeTransfer divides all transfer request by their destination chains.
func (p *DefaultTxOutputProducer) categorizeTransfer(transfers []*transferOutData) [][]*transferOutData {
	m := make(map[string][]*transferOutData)
	// We need to use an ordered array because map iteration is not deterministic and can result in
	// inconsistent data between nodes.
	orders := make([]string, 0)

	for _, transfer := range transfers {
		if m[transfer.destChain] == nil {
			m[transfer.destChain] = make([]*transferOutData, 0)
			orders = append(orders, transfer.destChain)
		}

		arr := m[transfer.destChain]
		arr = append(arr, transfer)
		m[transfer.destChain] = arr
	}

	ret := make([][]*transferOutData, 0)
	for _, chain := range orders {
		ret = append(ret, m[chain])
	}

	return ret
}

func (p *DefaultTxOutputProducer) processEthBatches(ctx sdk.Context, batches [][]*transferOutData) []*types.TxOutWithSigner {
	fmt.Println("AAAAA processEthBatches")

	dstChain := batches[0][0].destChain
	blockHeight := batches[0][0].blockHeight
	outMsgs := make([]*types.TxOutWithSigner, 0)

	for _, batch := range batches {
		tokens := make([]*types.Token, len(batch))
		recipients := make([]ethcommon.Address, len(batch))
		amounts := make([]*big.Int, len(batch))

		for k, transfer := range batch {
			fmt.Println("transfer = ", *transfer)
			tokens[k] = transfer.token
			recipients[k] = ecommon.HexToAddress(transfer.recipient)
			amounts[k] = transfer.amount
		}

		fmt.Println("len(tokens) = ", len(tokens), len(recipients), len(amounts))

		responseTx, err := p.buildERC20TransferIn(ctx, tokens, recipients, amounts, dstChain)
		if err != nil {
			log.Error("Failed to build erc20 transfer in, err = ", err)
			continue
		}

		bz, err := responseTx.EthTx.MarshalBinary()
		if err != nil {
			log.Error("processEthBatches: Failed to unmarshal eth tx, err = ", err)
			continue
		}

		outMsg := types.NewMsgTxOutWithSigner(
			p.appKeys.GetSignerAddress().String(),
			types.TxOutType_TRANSFER_OUT,
			blockHeight,
			"",
			"",
			dstChain,
			utils.KeccakHash32Bytes(bz),
			responseTx.RawBytes,
			"",
		)
		outMsgs = append(outMsgs, outMsg)
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
