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

const gasUnitPerSwap = 80_000

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
	ethCfg := chainCfg.EthConfig
	gasCost, _, _ := b.getGasCost(gasInfo, ethCfg.UseEip_1559, gasUnitPerSwap)

	inHashes := make([]string, 0, len(transfers))
	finalTokens := make([]ethcommon.Address, 0, len(transfers))
	finalRecipients := make([]ethcommon.Address, 0, len(transfers))
	finalAmounts := make([]*big.Int, 0, len(transfers))
	allTokens := b.keeper.GetAllTokens(ctx)
	tokenPrices := make([]string, 0, len(transfers))

	nativeTokenPrice, err := b.deyesClient.GetTokenPrice(chainCfg.NativeToken)
	if err != nil {
		return nil, err
	}
	if nativeTokenPrice.Cmp(utils.ZeroBigInt) == 0 {
		return nil, fmt.Errorf("token %s has price 0", chainCfg.NativeToken)
	}

	for _, transfer := range transfers {
		token := allTokens[transfer.Token]
		if token == nil {
			return nil, fmt.Errorf("token %s is invalid	", transfer.Token)
		}

		tokenPrice, err := b.deyesClient.GetTokenPrice(token.Id)
		if err != nil {
			return nil, err
		}
		if tokenPrice.Cmp(utils.ZeroBigInt) == 0 {
			return nil, fmt.Errorf("token %s has price 0", token.Id)
		}

		dstToken, amountOut, err := b.getTransferIn(ctx, allTokens[transfer.Token],
			transfer, gasCost, tokenPrice, nativeTokenPrice)
		if err != nil {
			return nil, err
		}

		finalTokens = append(finalTokens, dstToken)
		finalRecipients = append(finalRecipients, ethcommon.HexToAddress(transfer.ToRecipient))
		finalAmounts = append(finalAmounts, amountOut)
		inHashes = append(inHashes, transfer.Id)
		tokenPrices = append(tokenPrices, tokenPrice.String())

		amount, ok := new(big.Int).SetString(transfer.Amount, 10)
		if !ok {
			log.Warn("Cannot create big.Int value from amout ", transfer.Amount)
		}
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

	outMsg := &types.TxOut{
		TxType: types.TxOutType_TRANSFER_OUT,
		Content: &types.TxOutContent{
			OutChain: b.chain,
			OutHash:  responseTx.EthTx.Hash().String(),
			OutBytes: responseTx.RawBytes,
		},
		Input: &types.TxOutInput{
			TransferIds:      inHashes,
			NativeTokenPrice: nativeTokenPrice.String(),
			TokenPrices:      tokenPrices,
		},
	}

	return []*types.TxOut{outMsg}, nil
}

func (b *bridge) getTransferIn(
	ctx sdk.Context,
	token *types.Token,
	transfer *types.TransferDetails,
	gasCost *big.Int,
	tokenPrice *big.Int,
	nativeTokenPrice *big.Int,
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

	gasPriceInToken, err := helper.GasCostInToken(gasCost, tokenPrice, nativeTokenPrice)
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
	log.Verbosef("Nonce for %s on chain %s = %d", mpcAddr, b.chain, nonce)

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

	txGasInfo := &deyesethtypes.GasInfo{}
	if ethCfg.UseEip_1559 {
		txGasInfo.Tip = tx.GasTipCap().Int64()
		txGasInfo.BaseFee = tx.GasFeeCap().Int64()
	} else {
		txGasInfo.GasPrice = tx.GasPrice().Int64()
	}

	txGasCost, _, _ := b.getGasCost(txGasInfo, ethCfg.UseEip_1559, gasUnitPerSwap)
	log.Info("Validating gas cost, CurrentGasCost = %s, GasCostInTransaction=%s", currentGasCost, txGasCost)
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
	txAmountOuts, err := GetAmountOutFromTransaction(vaultInfo.Abi, tx, len(transfers))
	if err != nil {
		return err
	}

	allTokens := b.keeper.GetAllTokens(ctx)
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
		_, amountOut, err := b.getTransferIn(ctx, allTokens[transfer.Token],
			transfer, txGasCost, txTokenPrice, txNativeTokenPrice)
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
