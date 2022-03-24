package sisu

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
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
	GetTxOuts(ctx sdk.Context, height int64, tx *types.TxIn) []*types.TxOutWithSigner

	PauseContract(ctx sdk.Context, chain string, hash string) (*types.TxOutWithSigner, error)

	ResumeContract(ctx sdk.Context, chain string, hash string) (*types.TxOutWithSigner, error)

	ContractChangeOwnership(ctx sdk.Context, chain, contractHash, newOwner string) (*types.TxOutWithSigner, error)

	ContractSetLiquidPoolAddress(ctx sdk.Context, chain, contractHash, newAddress string) (*types.TxOutWithSigner, error)

	ContractEmergencyWithdrawFund(ctx sdk.Context, chain, contractHash string, tokens []string, newOwner string) (*types.TxOutWithSigner, error)
}

type DefaultTxOutputProducer struct {
	// List of key addresses in all eth based chain.
	// Map from: chain -> address -> bool.
	ethKeyAddrs map[string]map[string]bool

	worldState world.WorldState
	appKeys    common.AppKeys
	keeper     keeper.Keeper
	tssConfig  config.TssConfig
}

func NewTxOutputProducer(worldState world.WorldState, appKeys common.AppKeys,
	keeper keeper.Keeper, tssConfig config.TssConfig) TxOutputProducer {
	return &DefaultTxOutputProducer{
		keeper:     keeper,
		worldState: worldState,
		appKeys:    appKeys,
		tssConfig:  tssConfig,
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

	log.Verbose("ethTx.To() = ", ethTx.To())

	// 1. Check if this is a transaction sent to our key address. If this is true, it's likely a tx
	// that funds our account.
	if ethTx.To() != nil && p.keeper.IsKeygenAddress(ctx, libchain.KEY_TYPE_ECDSA, ethTx.To().String()) {
		contracts := p.keeper.GetPendingContracts(ctx, tx.Chain)
		log.Verbose("len(contracts) = ", len(contracts))

		if len(contracts) > 0 {
			// TODO: Check balance required to deploy all these contracts.

			// Get the list of deploy transactions. Those txs need to posted and verified (by validators)
			// to the Sisu chain.
			outEthTxs := p.getEthContractDeploymentTx(ctx, height, tx.Chain, contracts)

			for i, outTx := range outEthTxs {
				bz, err := outTx.MarshalBinary()
				if err != nil {
					log.Error("Cannot marshall binary", err)
					continue
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

				outMsgs = append(outMsgs, outMsg)
			}

			if len(outMsgs) > 0 {
				return outMsgs, nil
			}
		}
	}

	log.Verbose("len(ethTx.Data()) ", len(ethTx.Data()))

	// 2. Check if this is a tx sent to one of our contracts.
	if ethTx.To() != nil &&
		p.keeper.IsContractExistedAtAddress(ctx, tx.Chain, ethTx.To().String()) && len(ethTx.Data()) >= 4 {

		// TODO: compare method name to trigger corresponding contract method
		responseTx, err := p.processERC20TransferOut(ctx, ethTx)
		if err != nil {
			log.Error("cannot get response for erc20 tx, err = ", err)
			return nil, err
		}

		outMsg := types.NewMsgTxOutWithSigner(
			p.appKeys.GetSignerAddress().String(),
			types.TxOutType_TRANSFER_OUT,
			tx.BlockHeight,
			tx.Chain,
			tx.TxHash,
			responseTx.OutChain, // Could be different chain
			responseTx.EthTx.Hash().String(),
			responseTx.RawBytes,
			"",
		)

		outMsgs = append(outMsgs, outMsg)
	}

	// 3. Check other types of transaction.

	return outMsgs, nil
}

// Check if we can deploy contract after seeing some ETH being sent to our ethereum address.
func (p *DefaultTxOutputProducer) getEthContractDeploymentTx(ctx sdk.Context, height int64, chain string, contracts []*types.Contract) []*ethTypes.Transaction {
	txs := make([]*ethTypes.Transaction, 0)

	for _, contract := range contracts {
		nonce := p.worldState.UseAndIncreaseNonce(ctx, chain)
		log.Verbose("nonce for deploying contract:", nonce)
		if nonce < 0 {
			log.Error("cannot get nonce for contract")
			continue
		}

		rawTx := p.getContractTx(ctx, contract, nonce)
		if rawTx == nil {
			log.Warn("raw Tx is nil")
			continue
		}

		txs = append(txs, rawTx)
	}

	return txs
}

func (p *DefaultTxOutputProducer) getContractTx(ctx sdk.Context, contract *types.Contract, nonce int64) *ethTypes.Transaction {
	erc20 := SupportedContracts[ContractErc20Gateway]
	switch contract.Hash {
	case erc20.AbiHash:
		// This is erc20gw contract.
		parsedAbi, err := abi.JSON(strings.NewReader(erc20.AbiString))
		if err != nil {
			log.Error("cannot parse erc20 abi. abi = ", erc20.AbiString, "err =", err)
			return nil
		}

		// Get all allowed chains
		supportedChains := make([]string, 0)
		for chain := range p.tssConfig.SupportedChains {
			if chain != contract.Chain {
				supportedChains = append(supportedChains, chain)
			}
		}

		log.Info("Allowed chains for chain ", contract.Chain, " are: ", supportedChains)

		lp := p.keeper.GetLiquidity(ctx, contract.Chain)
		if lp == nil {
			return nil
		}

		log.Infof("Liquidity pool addr for chain %s is %s", contract.Chain, lp.Address)
		input, err := parsedAbi.Pack("", supportedChains, ecommon.HexToAddress(lp.Address))
		if err != nil {
			log.Error("cannot pack supportedChains, err =", err)
			return nil
		}

		byteCode := ecommon.FromHex(erc20.Bin)
		input = append(byteCode, input...)
		chain := p.keeper.GetChain(ctx, contract.Chain)
		if chain == nil {
			log.Error("getContractTx: chain is nil with id ", contract.Chain)
			return nil
		}

		gasPrice := chain.GasPrice
		if gasPrice < 0 {
			gasPrice = p.getDefaultGasPrice(contract.Chain).Int64()
		}

		rawTx := ethTypes.NewContractCreation(
			uint64(nonce),
			big.NewInt(0),
			p.getGasLimit(contract.Chain),
			big.NewInt(gasPrice),
			input,
		)

		return rawTx
	}

	return nil
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
		return big.NewInt(2_000_000_000) // 1 Gwei
	case "ganache2":
		return big.NewInt(2_000_000_000) // 1 Gwei
	case "eth-ropsten":
		return big.NewInt(1_700_000_000)
	case "eth-binance-testnet":
		return big.NewInt(10_000_000_000) // 10 Gwei
	}
	return big.NewInt(400_000_000_000) // 400 Gwei
}
