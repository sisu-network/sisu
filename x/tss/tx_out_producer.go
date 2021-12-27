package tss

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/db"
	"github.com/sisu-network/sisu/x/tss/keeper"
	"github.com/sisu-network/sisu/x/tss/types"
)

// This structs produces transaction output based on input. For a given tx input, this struct
// produces a list (could contain only one element) of transaction output.
type TxOutputProducer interface {
	// AddKeyAddress(ctx sdk.Context, chain, addr string) error
	GetTxOuts(ctx sdk.Context, height int64, tx *types.TxIn) []*types.TxOutWithSigner
}

type DefaultTxOutputProducer struct {
	// List of key addresses in all eth based chain.
	// Map from: chain -> address -> bool.
	ethKeyAddrs map[string]map[string]bool

	worldState WorldState
	keeper     keeper.Keeper
	appKeys    common.AppKeys
	db         db.Database
	tssConfig  config.TssConfig
}

func NewTxOutputProducer(worldState WorldState, keeper keeper.Keeper, appKeys common.AppKeys, db db.Database, tssConfig config.TssConfig) TxOutputProducer {
	return &DefaultTxOutputProducer{
		worldState: worldState,
		keeper:     keeper,
		appKeys:    appKeys,
		tssConfig:  tssConfig,
		db:         db,
	}
}

func (p *DefaultTxOutputProducer) GetTxOuts(ctx sdk.Context, height int64, tx *types.TxIn) []*types.TxOutWithSigner {
	outMsgs := make([]*types.TxOutWithSigner, 0)
	var err error

	if libchain.IsETHBasedChain(tx.Chain) {
		log.Info("Getting tx out for chain ", tx.Chain)
		outMsgs, err = p.getEthResponse(ctx, height, tx)

		if err != nil {
			log.Error("Cannot get response for an eth tx, err = ", err)
		}
	}

	return outMsgs
}

// Get ETH out from an observed tx. Only do this if this is a validator node.
func (p *DefaultTxOutputProducer) getEthResponse(ctx sdk.Context, height int64, tx *types.TxIn) ([]*types.TxOutWithSigner, error) {
	ethTx := &ethTypes.Transaction{}

	err := ethTx.UnmarshalBinary(tx.Serialized)
	if err != nil {
		log.Error("Failed to unmarshall eth tx. err =", err)
		return nil, err
	}

	outMsgs := make([]*types.TxOutWithSigner, 0)

	// 1. Check if this is a transaction sent to our key address. If this is true, it's likely a tx
	// that funds our account.
	if ethTx.To() != nil && p.db.IsChainKeyAddress(libchain.KEY_TYPE_ECDSA, ethTx.To().String()) {
		contracts := p.keeper.GetPendingContracts(ctx, tx.Chain)
		log.Verbose("len(contracts) = ", len(contracts))

		if len(contracts) > 0 {
			// TODO: Check balance required to deploy all these contracts.

			// Get the list of deploy transactions. Those txs need to posted and verified (by validators)
			// to the Sisu chain.
			outEthTxs := p.getEthContractDeploymentTx(ctx, height, tx.Chain, contracts)

			for _, outTx := range outEthTxs {
				bz, err := outTx.MarshalBinary()
				if err != nil {
					log.Error("Cannot marshall binary", err)
					continue
				}

				outMsg := types.NewMsgTxOutWithSigner(
					p.appKeys.GetSignerAddress().String(),
					tx.BlockHeight,
					tx.Chain,
					tx.TxHash,
					tx.Chain,
					bz,
				)

				outMsgs = append(outMsgs, outMsg)
			}

			if len(outMsgs) > 0 {
				return outMsgs, nil
			}
		}
	}

	// 2. Check if this is a tx sent to one of our contracts.
	if ethTx.To() != nil && len(ethTx.Data()) >= 4 {
		log.Verbose("ethTx.To() = ", ethTx.To())

		responseTx, err := p.createErc20ContractResponse(ctx, ethTx, tx.Chain)
		if err == nil {
			outMsg := types.NewMsgTxOutWithSigner(
				p.appKeys.GetSignerAddress().String(),
				tx.BlockHeight,
				tx.Chain,
				tx.TxHash,
				responseTx.OutChain, // Could be different chain
				responseTx.RawBytes,
			)

			outMsgs = append(outMsgs, outMsg)
		} else {
			log.Error("cannot get response for erc20 tx, err =", err)
		}
	}

	// 3. Check other types of transaction.

	return outMsgs, nil
}

// Check if we can deploy contract after seeing some ETH being sent to our ethereum address.
func (p *DefaultTxOutputProducer) getEthContractDeploymentTx(ctx sdk.Context, height int64, chain string, contracts []*types.Contract) []*ethTypes.Transaction {
	txs := make([]*ethTypes.Transaction, 0)

	for _, contract := range contracts {
		nonce := p.worldState.UseAndIncreaseNonce(chain)
		log.Verbose("nonce for deploying contract:", nonce)
		if nonce < 0 {
			log.Error("cannot get nonce for contract")
			continue
		}
		rawTx := p.getContractTx(contract, nonce)
		txs = append(txs, rawTx)
	}

	// Update all contracts to "deploying" state. This is not the cleanest code to do this but passing
	// the contract back would also complicated.
	p.keeper.UpdateContractsStatus(ctx, contracts, types.ContractStateDeploying)

	return txs
}

func (p *DefaultTxOutputProducer) getContractTx(contract *types.Contract, nonce int64) *ethTypes.Transaction {
	erc20 := SupportedContracts[ContractErc20]
	switch contract.Hash {
	case erc20.AbiHash:
		// This is erc20 contract.
		parsedAbi, err := abi.JSON(strings.NewReader(erc20.AbiString))
		if err != nil {
			log.Error("cannot parse erc20 abi. abi = ", erc20.AbiString, "err =", err)
			return nil
		}

		// Get all allowed chains
		allowedChains := make([]string, 0)
		for chain := range p.tssConfig.SupportedChains {
			if chain != contract.Chain {
				allowedChains = append(allowedChains, chain)
			}
		}

		log.Info("Allowed chains for chain", contract.Chain, "are: ", allowedChains)

		input, err := parsedAbi.Pack("", contract.Chain, allowedChains)
		if err != nil {
			log.Error("cannot pack allowedChains, err =", err)
			return nil
		}

		byteCode := ecommon.FromHex(erc20.Bin)
		input = append(byteCode, input...)

		rawTx := ethTypes.NewContractCreation(
			uint64(nonce),
			big.NewInt(0),
			p.getGasLimit(contract.Chain),
			p.getGasPrice(contract.Chain),
			input)
		return rawTx
	}

	return nil
}

func (p *DefaultTxOutputProducer) getGasLimit(chain string) uint64 {
	// TODO: Make this dependent on different chains.
	return uint64(5000000)
}

func (p *DefaultTxOutputProducer) getGasPrice(chain string) *big.Int {
	// TODO: Make this dependent on different chains.
	switch chain {
	case "eth-ropsten":
		return big.NewInt(1700000000)
	case "eth-binance-testnet":
		return big.NewInt(10000000000) // 10 Gwei
	}
	return big.NewInt(10000000000) // 10 Gwei
}
