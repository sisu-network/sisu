package sisu

import (
	"math/big"

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
	GetTxOuts(ctx sdk.Context, chain string, transfers []*types.Transfer) ([]*types.TxOutMsg, error)

	PauseContract(ctx sdk.Context, chain string, hash string) (*types.TxOutMsg, error)

	ResumeContract(ctx sdk.Context, chain string, hash string) (*types.TxOutMsg, error)

	ContractChangeOwnership(ctx sdk.Context, chain, contractHash, newOwner string) (*types.TxOutMsg, error)

	ContractSetLiquidPoolAddress(ctx sdk.Context, chain, contractHash, newAddress string) (*types.TxOutMsg, error)

	ContractEmergencyWithdrawFund(ctx sdk.Context, chain, contractHash string, tokens []string, newOwner string) (*types.TxOutMsg, error)
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

func (p *DefaultTxOutputProducer) parseCardanoTxIn(ctx sdk.Context, tx *types.TxIn) (*types.TransferOutData, error) {
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

func (p *DefaultTxOutputProducer) GetTxOuts(ctx sdk.Context, chain string,
	transfers []*types.Transfer) ([]*types.TxOutMsg, error) {

	if libchain.IsETHBasedChain(chain) {
		msgs, err := p.processEthBatches(ctx, chain, transfers)
		if err != nil {
			return nil, err
		}

		return msgs, nil
	}

	if libchain.IsCardanoChain(chain) {
		msgs, err := p.processCardanoBatches(ctx, p.keeper, chain, transfers)
		if err != nil {
			return nil, err
		}

		return msgs, nil
	}

	log.Error("Unknown chain type: ", chain)

	return nil, nil
}

func (p *DefaultTxOutputProducer) processEthBatches(ctx sdk.Context,
	dstChain string, transfers []*types.Transfer) ([]*types.TxOutMsg, error) {
	inHashes := make([]string, 0, len(transfers))
	tokens := make([]*types.Token, 0, len(transfers))
	recipients := make([]ethcommon.Address, 0, len(transfers))
	amounts := make([]*big.Int, 0, len(transfers))

	allTokens := p.keeper.GetAllTokens(ctx)
	for _, transfer := range transfers {
		token := allTokens[transfer.Token]
		if token == nil {
			log.Warn("cannot find token", transfer.Token)
			continue
		}

		amount, ok := new(big.Int).SetString(transfer.Amount, 10)
		if !ok {
			log.Warn("Cannot create big.Int value from amout ", transfer.Amount)
			continue
		}

		tokens = append(tokens, token)
		recipients = append(recipients, ecommon.HexToAddress(transfer.Recipient))
		amounts = append(amounts, amount)
		inHashes = append(inHashes, transfer.Id)

		log.Verbosef("Processing transfer in: id = %s, recipient = %s, amount = %s, inHash = %s",
			token.Id, transfer.Recipient, amount, transfer.Id)
	}

	responseTx, err := p.buildERC20TransferIn(ctx, p.keeper, tokens, recipients, amounts, dstChain)
	if err != nil {
		log.Error("Failed to build erc20 transfer in, err = ", err)
		return nil, err
	}

	outMsg := types.NewTxOutMsg(
		p.appKeys.GetSignerAddress().String(),
		types.TxOutType_TRANSFER_OUT,
		inHashes,
		dstChain,
		responseTx.EthTx.Hash().String(),
		responseTx.RawBytes,
		"",
	)
	return []*types.TxOutMsg{outMsg}, nil
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
