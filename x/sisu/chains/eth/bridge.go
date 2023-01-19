package eth

import (
	"fmt"
	"math/big"

	libchain "github.com/sisu-network/lib/chain"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	deyesethtypes "github.com/sisu-network/deyes/chains/eth/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	ctypes "github.com/sisu-network/sisu/x/sisu/chains/types"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/helper"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type bridge struct {
	signer      string
	chain       string
	keeper      keeper.Keeper
	deyesClient external.DeyesClient
}

func NewBridge(chain string, signer string, k keeper.Keeper, deyesClient external.DeyesClient) ctypes.Bridge {
	return &bridge{
		chain:       chain,
		signer:      signer,
		keeper:      k,
		deyesClient: deyesClient,
	}
}

func (b *bridge) ProcessTransfers(ctx sdk.Context, transfers []*types.TransferDetails) ([]*types.TxOutMsg, error) {
	gasInfo, err := b.deyesClient.GetGasInfo(b.chain)
	if err != nil {
		return nil, err
	}
	chainCfg := b.keeper.GetChain(ctx, b.chain)
	ethCfg := chainCfg.EthConfig
	gasUnitPerSwap := 80_000
	gasCost, _, _ := b.getGasCost(gasInfo, ethCfg.UseEip_1559, gasUnitPerSwap)

	inHashes := make([]string, 0, len(transfers))
	finalTokens := make([]ethcommon.Address, 0, len(transfers))
	finalRecipients := make([]ethcommon.Address, 0, len(transfers))
	finalAmounts := make([]*big.Int, 0, len(transfers))
	allTokens := b.keeper.GetAllTokens(ctx)

	for _, transfer := range transfers {
		dstToken, amountOut, err := b.getTransferIn(ctx, allTokens, transfer, gasCost)
		if err != nil {
			log.Errorf("Failed to get transfer in, err = %s", err)
			break
		}

		token := allTokens[transfer.Token]
		if token == nil {
			log.Warn("cannot find token", transfer.Token)
			break
		}

		amount, ok := new(big.Int).SetString(transfer.Amount, 10)
		if !ok {
			log.Warn("Cannot create big.Int value from amout ", transfer.Amount)
			break
		}

		finalTokens = append(finalTokens, dstToken)
		finalRecipients = append(finalRecipients, ethcommon.HexToAddress(transfer.ToRecipient))
		finalAmounts = append(finalAmounts, amountOut)
		inHashes = append(inHashes, transfer.Id)

		log.Verbosef("Processing transfer in: id = %s, token = %s, recipient = %s, amount = %s, "+
			"inHash = %s, toChain = %s, toRecipient = %s",
			transfer.Id, token.Id, transfer.ToRecipient, amount, transfer.Id, transfer.ToChain,
			transfer.ToRecipient)
	}

	if len(finalTokens) == 0 {
		return nil, fmt.Errorf("Failed to get any transaction!")
	}

	responseTx, err := b.buildTransaction(ctx, finalTokens, finalRecipients, finalAmounts,
		gasUnitPerSwap, ethCfg.UseEip_1559, gasInfo)
	if err != nil {
		log.Error("Failed to build erc20 transfer in, err = ", err)
		return nil, err
	}

	outMsg := types.NewTxOutMsg(
		b.signer,
		types.TxOutType_TRANSFER_OUT,
		&types.TxOutContent{
			OutChain: b.chain,
			OutHash:  responseTx.EthTx.Hash().String(),
			OutBytes: responseTx.RawBytes,
		},
		&types.TxOutInput{
			TransferIds: inHashes,
		},
	)

	return []*types.TxOutMsg{outMsg}, nil
}

func (b *bridge) getTransferIn(
	ctx sdk.Context,
	allTokens map[string]*types.Token,
	transfer *types.TransferDetails,
	gasCost *big.Int,
) (ethcommon.Address, *big.Int, error) {
	targetContractName := ContractVault
	v := b.keeper.GetVault(ctx, b.chain, "")
	if v == nil {
		return ethcommon.Address{}, nil, fmt.Errorf("Cannot find vault for chain %s", b.chain)
	}
	gw := v.Address
	if len(gw) == 0 {
		err := fmt.Errorf("cannot find gw address for type: %s on chain %s", targetContractName, b.chain)
		log.Error(err)
		return ethcommon.Address{}, nil, err
	}

	chain := b.keeper.GetChain(ctx, b.chain)
	if chain == nil {
		return ethcommon.Address{}, nil, fmt.Errorf("Invalid chain: %s", chain)
	}

	commissionRate := b.keeper.GetParams(ctx).CommissionRate
	if commissionRate < 0 || commissionRate > 10_000 {
		return ethcommon.Address{}, nil, fmt.Errorf("Commission rate is invalid, rate = %d", commissionRate)
	}

	token := allTokens[transfer.Token]
	if token == nil {
		return ethcommon.Address{}, nil, fmt.Errorf("cannot find token %s", transfer.Token)
	}

	amountIn, ok := new(big.Int).SetString(transfer.Amount, 10)
	if !ok {
		return ethcommon.Address{}, nil, fmt.Errorf("Cannot create big.Int value from amout %s", transfer.Amount)
	}

	var tokenAddr string
	for j, chain := range token.Chains {
		if chain == b.chain {
			tokenAddr = token.Addresses[j]
			break
		}
	}

	if len(tokenAddr) == 0 {
		return ethcommon.Address{}, nil, fmt.Errorf("cannot find token address on chain %s", b.chain)
	}

	amountOut := new(big.Int).Set(amountIn)

	// Subtract commission rate
	amountOut = utils.SubtractCommissionRate(amountOut, commissionRate)

	price, err := b.deyesClient.GetTokenPrice(token.Id)
	if err != nil {
		return ethcommon.Address{}, nil, err
	}

	if price.Cmp(utils.ZeroBigInt) == 0 {
		return ethcommon.Address{}, nil, fmt.Errorf("token %s has price 0", token.Id)
	}

	gasPriceInToken, err := helper.GetChainGasCostInToken(ctx, b.keeper, b.deyesClient, token.Id,
		b.chain, gasCost)
	if err != nil {
		return ethcommon.Address{}, nil, fmt.Errorf("Cannot get gas cost in token, err = %s", err)
	}

	if gasPriceInToken.Cmp(utils.ZeroBigInt) < 0 {
		log.Errorf("Gas price in token is negative: token id = %s", token.Id)
		gasPriceInToken = utils.ZeroBigInt
	}

	// Subtract gas price in token.
	amountOut.Sub(amountOut, gasPriceInToken)

	// Check if the amountOut is smaller than 0 or not.
	if amountOut.Cmp(utils.ZeroBigInt) < 0 {
		return ethcommon.Address{}, nil,
			fmt.Errorf("Insufficient fund for transfer amountOut = %s, gasPriceInToken = %s",
				amountOut, gasPriceInToken)
	}

	log.Verbosef("tokenAddr: %s, recipient: %s, gasPriceInToken: %s, amountIn: %s, amountOut: %s",
		tokenAddr, transfer.ToRecipient, gasPriceInToken, amountIn.String(), amountOut,
	)

	return ethcommon.HexToAddress(tokenAddr), amountOut, nil
}

func (b *bridge) buildTransaction(
	ctx sdk.Context,
	finalTokenAddrs []ethcommon.Address,
	finalRecipients []ethcommon.Address,
	finalAmounts []*big.Int,
	gasUnitPerSwap int,
	useEip1559 bool,
	gasInfo *deyesethtypes.GasInfo,
) (*types.TxResponse, error) {
	targetContractName := ContractVault
	v := b.keeper.GetVault(ctx, b.chain, "")
	if v == nil {
		return nil, fmt.Errorf("Cannot find vault for chain %s", b.chain)
	}
	gw := v.Address
	if len(gw) == 0 {
		return nil, fmt.Errorf("cannot find gw address for type: %s on chain %s", targetContractName, b.chain)
	}

	gatewayAddress := ethcommon.HexToAddress(gw)
	vaultInfo := SupportedContracts[targetContractName]

	var input []byte
	var err error
	if len(finalTokenAddrs) == 1 {
		input, err = vaultInfo.Abi.Pack(
			MethodTransferIn,
			finalTokenAddrs[0],
			finalRecipients[0],
			finalAmounts[0],
		)
	} else {
		input, err = vaultInfo.Abi.Pack(
			MethodTransferInMultiple,
			finalTokenAddrs,
			finalRecipients,
			finalAmounts,
		)
	}
	if err != nil {
		return nil, err
	}

	mpcAddr := b.keeper.GetMpcAddress(ctx, b.chain)
	nonce, err := b.deyesClient.GetNonce(b.chain, mpcAddr)
	if err != nil {
		return nil, err
	}
	log.Verbosef("Nonce on chain %s = %d", b.chain, nonce)

	maxGas := uint64(gasUnitPerSwap * len(finalRecipients)) // max 80k per swapping operation.
	_, tipCap, feeCap := b.getGasCost(gasInfo, useEip1559, gasUnitPerSwap)

	var rawTx *ethtypes.Transaction
	if useEip1559 {
		dynamicFeeTx := &ethtypes.DynamicFeeTx{
			ChainID:   libchain.GetChainIntFromId(b.chain),
			Nonce:     uint64(nonce),
			GasTipCap: tipCap,
			GasFeeCap: feeCap,
			Gas:       maxGas,
			To:        &gatewayAddress,
			Value:     big.NewInt(0),
			Data:      input,
		}

		rawTx = ethtypes.NewTx(dynamicFeeTx)
	} else {
		rawTx = ethtypes.NewTransaction(
			uint64(nonce),
			gatewayAddress,
			big.NewInt(0),
			maxGas,
			big.NewInt(gasInfo.GasPrice),
			input,
		)
	}

	bz, err := rawTx.MarshalBinary()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &types.TxResponse{
		OutChain: b.chain,
		EthTx:    rawTx,
		RawBytes: bz,
	}, nil
}

// getGasCost returns total gas cost used for swapping transaction.
func (b *bridge) getGasCost(gasInfo *deyesethtypes.GasInfo, useEip1559 bool, maxGasUnit int) (*big.Int, *big.Int, *big.Int) {
	if useEip1559 {
		// Max fee = 2 * baseFee + Tip
		tipCap := big.NewInt(gasInfo.Tip)
		feeCap := big.NewInt(gasInfo.BaseFee*2 + gasInfo.Tip)

		return new(big.Int).Mul(feeCap, big.NewInt(int64(maxGasUnit))), tipCap, feeCap
	} else {
		return new(big.Int).Mul(big.NewInt(int64(maxGasUnit)), big.NewInt(gasInfo.GasPrice)), nil, nil
	}
}

func (b *bridge) ParseIncomingTx(ctx sdk.Context, chain string, serialized []byte) ([]*types.TransferDetails, error) {
	parseResult := ParseVaultTx(ctx, b.keeper, chain, serialized)
	if parseResult.Error != nil {
		return nil, parseResult.Error
	}

	if parseResult.TransferOuts != nil {
		return parseResult.TransferOuts, nil
	}

	return []*types.TransferDetails{}, nil
}
