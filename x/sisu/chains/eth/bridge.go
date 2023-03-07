package eth

import (
	"fmt"
	"math/big"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	deyesethtypes "github.com/sisu-network/deyes/chains/eth/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	ctypes "github.com/sisu-network/sisu/x/sisu/chains/types"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/helper"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

const gasUnitPerSwap uint64 = 80_000
const gasUnitPerExecution uint64 = 50_000

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

func (b *bridge) ProcessTransfers(ctx sdk.Context, transfers []*types.TransferDetails) ([]*types.TxOut, error) {
	gasInfo, err := b.deyesClient.GetGasInfo(b.chain)
	if err != nil {
		return nil, err
	}
	log.Verbosef("Gas info for chain %s %v", b.chain, *gasInfo)
	chainCfg := b.keeper.GetChain(ctx, b.chain)

	switch transfers[0].TxType {
	case types.TxInType_TOKEN_TRANSFER:
		return b.processTokenTransfer(ctx, gasInfo, chainCfg, transfers)
	case types.TxInType_REMOTE_CALL:
		return b.processRemoteCall(ctx, gasInfo, chainCfg, transfers)
	}

	return nil, fmt.Errorf("invalid transfer type %s", transfers[0].TxType)
}

func (b *bridge) processTokenTransfer(ctx sdk.Context, gasInfo *deyesethtypes.GasInfo,
	chainCfg *types.Chain, transfers []*types.TransferDetails,
) ([]*types.TxOut, error) {
	ethCfg := chainCfg.EthConfig
	gasCost, _, _ := b.getGasCost(gasInfo, ethCfg.UseEip_1559, gasUnitPerSwap)

	inHashes := make([]string, 0, len(transfers))
	maxGas := uint64(0)
	finalTokens := make([]ethcommon.Address, 0, len(transfers))
	finalRecipients := make([]ethcommon.Address, 0, len(transfers))
	finalAmounts := make([]*big.Int, 0, len(transfers))
	tokenPrices := make([]string, 0, len(transfers))

	nativeTokenPrice, err := b.deyesClient.GetTokenPrice(chainCfg.NativeToken)
	if err != nil {
		return nil, err
	}
	if nativeTokenPrice.Cmp(utils.ZeroBigInt) == 0 {
		return nil, fmt.Errorf("token %s has price 0", chainCfg.NativeToken)
	}

	for _, transfer := range transfers {
		dstToken, amountOut, tokenPrice, err := b.getTransferIn(ctx, transfer, gasCost, nativeTokenPrice)
		if err != nil {
			return nil, err
		}

		finalTokens = append(finalTokens, dstToken)
		finalRecipients = append(finalRecipients, ethcommon.HexToAddress(transfer.ToRecipient))
		finalAmounts = append(finalAmounts, amountOut)
		inHashes = append(inHashes, transfer.GetRetryId())
		tokenPrices = append(tokenPrices, tokenPrice.String())
		maxGas += gasUnitPerSwap

		amount, ok := new(big.Int).SetString(transfer.Amount, 10)
		if !ok {
			log.Warn("Cannot create big.Int value from amout ", transfer.Amount)
		}
		log.Verbosef("Processing transfer in: id = %s, token = %s, sender = %s, amount = %s, "+
			" toChain = %s, toRecipient = %s",
			transfer.Id, transfer.Token, transfer.FromSender, amount, transfer.ToChain,
			transfer.ToRecipient)
	}

	if len(finalTokens) == 0 {
		return nil, fmt.Errorf("failed to get any transaction!")
	}

	output, err := buildTransferIn(ctx, finalTokens, finalRecipients, finalAmounts)
	if err != nil {
		log.Error("Failed to build erc20 transfer in, err = ", err)
		return nil, err
	}

	responseTx, err := b.buildTransaction(ctx, output, maxGas, ethCfg.UseEip_1559, gasInfo)
	if err != nil {
		log.Error("Failed to build erc20 transaction, err = ", err)
		return nil, err
	}

	outMsg := &types.TxOut{
		TxType: types.TxOutType_TRANSFER,
		Content: &types.TxOutContent{
			OutChain: b.chain,
			OutHash:  responseTx.EthTx.Hash().String(),
			OutBytes: responseTx.RawBytes,
		},
		Input: &types.TxOutInput{
			TransferRetryIds: inHashes,
			NativeTokenPrice: nativeTokenPrice.String(),
			TokenPrices:      tokenPrices,
			EthData: &types.EthData{
				GasPrice: gasInfo.GasPrice,
				BaseFee:  gasInfo.BaseFee,
				Tip:      gasInfo.Tip,
			},
		},
	}

	return []*types.TxOut{outMsg}, nil
}

func (b *bridge) getTransferIn(
	ctx sdk.Context,
	transfer *types.TransferDetails,
	gasCost *big.Int,
	nativeTokenPrice *big.Int,
) (ethcommon.Address, *big.Int, *big.Int, error) {
	chain := b.keeper.GetChain(ctx, b.chain)
	if chain == nil {
		return ethcommon.Address{}, nil, nil, fmt.Errorf("Invalid chain: %s", chain)
	}

	token := b.keeper.GetToken(ctx, transfer.Token)
	if token == nil {
		return ethcommon.Address{}, nil, nil, fmt.Errorf("cannot find token %s", transfer.Token)
	}

	tokenPrice, err := b.deyesClient.GetTokenPrice(token.Id)
	if err != nil {
		return ethcommon.Address{}, nil, nil, err
	}
	if tokenPrice.Cmp(utils.ZeroBigInt) == 0 {
		return ethcommon.Address{}, nil, nil, fmt.Errorf("token %s has price 0", token.Id)
	}

	var tokenAddr string
	for i, chain := range token.Chains {
		if chain == b.chain {
			tokenAddr = token.Addresses[i]
			break
		}
	}

	if len(tokenAddr) == 0 {
		return ethcommon.Address{}, nil, nil, fmt.Errorf("cannot find token address on chain %s", b.chain)
	}

	amountOut, err := b.calcAmountOut(ctx, transfer, gasCost, tokenPrice, nativeTokenPrice)
	if err != nil {
		return ethcommon.Address{}, nil, nil, err
	}

	log.Verbosef("tokenAddr: %s, recipient: %s, amountOut: %s",
		tokenAddr, transfer.ToRecipient, amountOut,
	)

	return ethcommon.HexToAddress(tokenAddr), amountOut, tokenPrice, nil
}

func (b *bridge) processRemoteCall(ctx sdk.Context, gasInfo *deyesethtypes.GasInfo,
	chainCfg *types.Chain, transfers []*types.TransferDetails,
) ([]*types.TxOut, error) {
	ethCfg := chainCfg.EthConfig

	maxGas := uint64(0)
	inHashes := make([]string, 0, len(transfers))
	finalCallerChains := make([]*big.Int, 0, len(transfers))
	finalCallers := make([]ethcommon.Address, 0, len(transfers))
	finalApps := make([]ethcommon.Address, 0, len(transfers))
	finalGasLimits := make([]uint64, 0, len(transfers))
	finalMessages := make([][]byte, 0, len(transfers))

	for _, transfer := range transfers {
		inHashes = append(inHashes, transfer.GetRetryId())
		finalCallerChains = append(finalCallerChains, libchain.GetChainIntFromId(transfer.FromChain))
		finalCallers = append(finalCallers, ethcommon.HexToAddress(transfer.FromSender))
		finalApps = append(finalApps, ethcommon.HexToAddress(transfer.ToRecipient))
		finalGasLimits = append(finalGasLimits, transfer.CallGasLimit)
		finalMessages = append(finalMessages, transfer.Message)
		maxGas += gasUnitPerExecution + transfer.CallGasLimit

		log.Verbosef("Processing remote call: id = %s, caller = %s, appChain = %s, app = %s",
			transfer.Id, transfer.FromSender, transfer.ToChain, transfer.ToRecipient)
	}

	if len(finalMessages) == 0 {
		return nil, fmt.Errorf("Failed to get any transaction!")
	}

	commission := b.keeper.GetParams(ctx).RemoteCallCommission
	if commission < 0 {
		return nil, fmt.Errorf("invalid commission value %d", commission)
	}

	data, err := buildRemoteExecute(ctx, commission, finalCallerChains, finalCallers, finalApps,
		finalGasLimits, finalMessages)
	if err != nil {
		log.Error("Failed to build erc20 transfer in, err = ", err)
		return nil, err
	}

	responseTx, err := b.buildTransaction(ctx, data, maxGas, ethCfg.UseEip_1559, gasInfo)
	if err != nil {
		log.Error("Failed to build erc20 transaction, err = ", err)
		return nil, err
	}

	outMsg := &types.TxOut{
		TxType: types.TxOutType_TRANSFER,
		Content: &types.TxOutContent{
			OutChain: b.chain,
			OutHash:  responseTx.EthTx.Hash().String(),
			OutBytes: responseTx.RawBytes,
		},
		Input: &types.TxOutInput{
			TransferRetryIds: inHashes,
			EthData: &types.EthData{
				GasPrice: gasInfo.GasPrice,
				BaseFee:  gasInfo.BaseFee,
				Tip:      gasInfo.Tip,
			},
		},
	}

	return []*types.TxOut{outMsg}, nil
}

func (b *bridge) buildTransaction(
	ctx sdk.Context,
	data []byte,
	maxGas uint64,
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

	mpcAddr := b.keeper.GetMpcAddress(ctx, b.chain)
	nonce, err := b.deyesClient.GetNonce(b.chain, mpcAddr)
	if err != nil {
		return nil, err
	}
	log.Verbosef("Nonce for %s on chain %s = %d", mpcAddr, b.chain, nonce)

	_, tipCap, feeCap := b.getGasCost(gasInfo, useEip1559, maxGas)

	var inner ethtypes.TxData
	if useEip1559 {
		inner = &ethtypes.DynamicFeeTx{
			ChainID:   libchain.GetChainIntFromId(b.chain),
			Nonce:     nonce,
			GasTipCap: tipCap,
			GasFeeCap: feeCap,
			Gas:       maxGas,
			To:        &gatewayAddress,
			Value:     big.NewInt(0),
			Data:      data,
		}
	} else {
		inner = &ethtypes.LegacyTx{
			Nonce:    nonce,
			To:       &gatewayAddress,
			Value:    big.NewInt(0),
			Gas:      maxGas,
			GasPrice: big.NewInt(gasInfo.GasPrice),
			Data:     data,
		}
	}
	var rawTx = ethtypes.NewTx(inner)

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
func (b *bridge) getGasCost(
	gasInfo *deyesethtypes.GasInfo, useEip1559 bool, maxGasUnit uint64,
) (*big.Int, *big.Int, *big.Int) {
	if useEip1559 {
		// Max fee = 2 * baseFee + Tip
		tipCap := big.NewInt(gasInfo.Tip)
		feeCap := big.NewInt(gasInfo.BaseFee*2 + gasInfo.Tip)

		return new(big.Int).Mul(feeCap, big.NewInt(int64(maxGasUnit))), tipCap, feeCap
	} else {
		return new(big.Int).Mul(big.NewInt(int64(maxGasUnit)), big.NewInt(gasInfo.GasPrice)), nil, nil
	}
}

// ParseIncomingTx implements bridge interface
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

// ProcessCommand implements bridge interface
func (b *bridge) ProcessCommand(ctx sdk.Context, cmd *types.Command) (*types.TxOut, error) {
	return nil, fmt.Errorf("Invalid command")
}

func (b *bridge) ValidateTxOut(ctx sdk.Context, txOut *types.TxOut, transfers []*types.TransferDetails) error {
	if transfers[0].TxType != types.TxInType_TOKEN_TRANSFER {
		return nil
	}

	// 1. Validate gas cost.
	gasInfo, err := b.deyesClient.GetGasInfo(b.chain)
	if err != nil {
		return err
	}

	chainCfg := b.keeper.GetChain(ctx, b.chain)
	ethCfg := chainCfg.EthConfig

	currentGasCost, _, _ := b.getGasCost(gasInfo, ethCfg.UseEip_1559, gasUnitPerSwap)

	tx := &ethtypes.Transaction{}
	if err := tx.UnmarshalBinary(txOut.Content.OutBytes); err != nil {
		return err
	}

	txGasInfo := &deyesethtypes.GasInfo{
		GasPrice: txOut.Input.EthData.GasPrice,
		BaseFee:  txOut.Input.EthData.BaseFee,
		Tip:      txOut.Input.EthData.Tip,
	}

	txGasCost, _, _ := b.getGasCost(txGasInfo, ethCfg.UseEip_1559, gasUnitPerSwap)
	log.Infof("Validating gas cost, CurrentGasCost = %s, GasCostInTransaction=%s", currentGasCost, txGasCost)
	if ratio, ok := helper.CheckRatioThreshold(txGasCost, currentGasCost, 3.00); !ok {
		return fmt.Errorf(
			"cannot accept the transaction with too large difference in gas cost, ratio=%d%%",
			int(ratio*100))
	}

	// 2. Validate native token price.
	txNativeTokenPrice, ok := new(big.Int).SetString(txOut.Input.NativeTokenPrice, 10)
	if !ok {
		return fmt.Errorf("cannot convert nativeTokenPrice (%s) to big int", txOut.Input.NativeTokenPrice)
	}
	if err = b.checkDifferenceTokenPrice(chainCfg.NativeToken, txNativeTokenPrice); err != nil {
		return err
	}

	// 3. Validate token price and amountOut.
	targetContractName := ContractVault
	vaultInfo := SupportedContracts[targetContractName]
	txAmountOuts, err := GetAmountOutFromTransaction(ctx, b.keeper, vaultInfo.Abi, tx, transfers)
	if err != nil {
		return err
	}

	for i, transfer := range transfers {
		// 3a. Validate token price.
		txTokenPrice, ok := new(big.Int).SetString(txOut.Input.TokenPrices[i], 10)
		if !ok {
			return fmt.Errorf("cannot convert tokenPrice (%s) to big int", txOut.Input.TokenPrices[i])
		}
		if err = b.checkDifferenceTokenPrice(transfer.Token, txTokenPrice); err != nil {
			return err
		}

		// 3b. Validate amoutOut.
		amountOut, err := b.calcAmountOut(ctx, transfer, txGasCost, txTokenPrice, txNativeTokenPrice)
		if err != nil {
			return err
		}

		if txAmountOuts[i].Cmp(amountOut) != 0 {
			return fmt.Errorf(
				"cannot accept the transaction with incorrect amountOut, got %s, expected %s",
				txAmountOuts[i], amountOut)
		}
	}

	return nil
}

func (b *bridge) checkDifferenceTokenPrice(token string, txPrice *big.Int) error {
	currentPrice, err := b.deyesClient.GetTokenPrice(token)
	if err != nil {
		return err
	}

	if ratio, ok := helper.CheckRatioThreshold(currentPrice, txPrice, 1.1); !ok {
		return fmt.Errorf(
			"cannot accept the transaction with too large difference in token %s price, ratio=%d%%",
			token, int(ratio*100))
	}

	return nil
}

func (b *bridge) calcAmountOut(
	ctx sdk.Context, transfer *types.TransferDetails, gasCost, tokenPrice, nativeTokenPrice *big.Int,
) (*big.Int, error) {
	amountIn, ok := new(big.Int).SetString(transfer.Amount, 10)
	if !ok {
		return nil, fmt.Errorf("annot create big.Int value from amout %s", transfer.Amount)
	}

	commissionRate := b.keeper.GetParams(ctx).TransferCommissionRate
	if commissionRate < 0 || commissionRate > 10_000 {
		return nil, fmt.Errorf("Commission rate is invalid, rate = %d", commissionRate)
	}

	amountOut := utils.SubtractCommissionRate(amountIn, commissionRate)

	gasPriceInToken, err := helper.GasCostInToken(gasCost, tokenPrice, nativeTokenPrice)
	if err != nil {
		return nil, fmt.Errorf("Cannot get gas cost in token, err = %s", err)
	}

	if gasPriceInToken.Cmp(utils.ZeroBigInt) < 0 {
		log.Errorf("gas token in price is negative, gasCost = %s, tokenPrice = %s, nativeTokenPrice = %s",
			gasCost, tokenPrice, nativeTokenPrice)
		gasPriceInToken = utils.ZeroBigInt
	}

	finalAmount := new(big.Int).Sub(amountOut, gasPriceInToken)
	if err != nil {
		return nil, err
	}

	if finalAmount.Cmp(utils.ZeroBigInt) < 0 {
		return nil, fmt.Errorf("insufficient funds, amountIn = %s, amountOut = %s, gas = %s",
			amountIn, amountOut, gasPriceInToken)
	}

	return finalAmount, nil
}
