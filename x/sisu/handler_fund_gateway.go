package sisu

import (
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerFundGateway struct {
	mc     ManagerContainer
	keeper keeper.Keeper
}

func NewHandlerFundGateway(mc ManagerContainer) *HandlerFundGateway {
	return &HandlerFundGateway{
		mc:     mc,
		keeper: mc.Keeper(),
	}
}

func (h *HandlerFundGateway) DeliverMsg(ctx sdk.Context, msg *types.FundGatewayMsg) (*sdk.Result, error) {
	pmm := h.mc.PostedMessageManager()
	if process, hash := pmm.ShouldProcessMsg(ctx, msg); process {
		h.doContractDeployment(ctx, msg.Data)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{}, nil
	}

	return &sdk.Result{}, nil
}

func (h *HandlerFundGateway) doContractDeployment(ctx sdk.Context, data *types.FundGateway) {
	log.Info("Create txout to deploy a contract")
	outMsgs := make([]*types.TxOutMsg, 0)
	contracts := h.keeper.GetPendingContracts(ctx, data.Chain)
	log.Verbose("len(contracts) = ", len(contracts))

	if len(contracts) > 0 {
		// TODO: Check balance required to deploy all these contracts. Also check if we are deploying
		// a contract to avoid duplication.

		// Get the list of deploy transactions. Those txs need to posted and verified (by validators)
		// to the Sisu chain.
		outEthTxs := h.getEthContractDeploymentTx(ctx, data.Chain, contracts)

		for i, outTx := range outEthTxs {
			bz, err := outTx.MarshalBinary()
			if err != nil {
				return
			}

			outMsg := types.NewTxOutMsg(
				h.mc.AppKeys().GetSignerAddress().String(),
				types.TxOutType_CONTRACT_DEPLOYMENT,
				[]string{data.TxHash},
				data.Chain,
				outTx.Hash().String(),
				bz,
				contracts[i].Hash,
			)

			log.Verbose("ETH Tx Out hash = ", outTx.Hash().String(), " on chain ", data.Chain)

			outMsgs = append(outMsgs, outMsg)
		}
	}

	for _, outMsg := range outMsgs {
		h.mc.TxSubmit().SubmitMessageAsync(outMsg)
	}
}

func (h *HandlerFundGateway) getEthContractDeploymentTx(ctx sdk.Context, chain string, contracts []*types.Contract) []*ethTypes.Transaction {
	txs := make([]*ethTypes.Transaction, 0)

	for _, contract := range contracts {
		rawTx := h.getContractTx(ctx, contract)
		if rawTx == nil {
			log.Warn("raw Tx is nil")
			continue
		}

		txs = append(txs, rawTx)
	}

	return txs
}

func (h *HandlerFundGateway) getContractTx(ctx sdk.Context, contract *types.Contract) *ethTypes.Transaction {
	erc20 := SupportedContracts[ContractErc20Gateway]
	switch contract.Hash {
	case erc20.AbiHash:
		// This is erc20gw contract.
		parsedAbi, err := abi.JSON(strings.NewReader(erc20.AbiString))
		if err != nil {
			log.Error("cannot parse erc20 abi. abi = ", erc20.AbiString, "err =", err)
			return nil
		}

		params := h.keeper.GetParams(ctx)
		log.Info("Allowed chains for chain ", contract.Chain, " are: ", params.SupportedChains)
		lp := h.keeper.GetLiquidity(ctx, contract.Chain)
		if lp == nil {
			log.Warn("Lp is nil for chain ", contract.Chain)
			return nil
		}

		log.Infof("Liquidity pool addr for chain %s is %s", contract.Chain, lp.Address)
		input, err := parsedAbi.Pack("", params.SupportedChains, ecommon.HexToAddress(lp.Address))
		if err != nil {
			log.Error("cannot pack supportedChains, err =", err)
			return nil
		}

		byteCode := ecommon.FromHex(erc20.Bin)
		input = append(byteCode, input...)
		chain := h.keeper.GetChain(ctx, contract.Chain)
		if chain == nil {
			log.Error("getContractTx: chain is nil with id ", contract.Chain)
			return nil
		}

		gasPrice := chain.GasPrice
		if gasPrice <= 0 {
			gasPrice = h.getDefaultGasPrice(contract.Chain).Int64()
		}
		gasLimit := h.getGasLimit(contract.Chain)
		rawTx := ethTypes.NewContractCreation(
			0,
			big.NewInt(0),
			gasLimit,
			big.NewInt(gasPrice),
			input,
		)

		return rawTx
	}

	return nil
}

func (h *HandlerFundGateway) getGasLimit(chain string) uint64 {
	// TODO: Make this dependent on different chains.
	return uint64(8_000_000)
}

// @Deprecated
func (h *HandlerFundGateway) getDefaultGasPrice(chain string) *big.Int {
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
