package eth

import (
	"fmt"
	"math/big"
	"strings"

	libchain "github.com/sisu-network/lib/chain"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/contracts/eth/vault"
	"github.com/sisu-network/sisu/utils"
	chainstypes "github.com/sisu-network/sisu/x/sisu/chains/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// ParseVaultTx parses a transaction that is sent to the vault.
func ParseVaultTx(ctx sdk.Context, keeper keeper.Keeper, chain string, serialized []byte) *chainstypes.ParseResult {
	vaultAbi, err := abi.JSON(strings.NewReader(vault.VaultABI))
	if err != nil {
		return &chainstypes.ParseResult{Error: err}
	}

	ethTx := &ethtypes.Transaction{}
	err = ethTx.UnmarshalBinary(serialized)
	if err != nil {
		log.Error("Failed to unmarshall eth tx. err =", err)
		return &chainstypes.ParseResult{Error: err}
	}

	callData := ethTx.Data()
	if len(callData) < 4 {
		// This is just a normal transfer
		return &chainstypes.ParseResult{
			Method: chainstypes.MethodNativeTransfer,
		}
	}

	result := &chainstypes.ParseResult{}
	methodName, txParams, err := DecodeTxParams(vaultAbi, callData)
	if err != nil {
		return &chainstypes.ParseResult{Error: err}
	}

	switch methodName {
	case "transferOut":
		result.Method = chainstypes.MethodTransferOut
		result.TransferOuts, result.Error = parseTransferOut(ctx, keeper, ethTx, chain, true, txParams)
	case "transferOutNonEvm":
		result.Method = chainstypes.MethodTransferOutNonEvm
		result.TransferOuts, result.Error = parseTransferOut(ctx, keeper, ethTx, chain, false, txParams)
	case "addSpender":
		result.Method = chainstypes.MethodAddSpender
	case "remoteCall":
		result.Method = chainstypes.MethodRemoteCall
		result.TransferOuts, result.Error = parseRemoteCall(ctx, ethTx, chain, true, txParams)
	case "createApp":
		result.Method = chainstypes.MethodCreateApp
	case "setAppAnyCaller":
		result.Method = chainstypes.MethodSetAppAnyCaller
	default:
		result.Method = chainstypes.MethodUnknown
		result.Error = fmt.Errorf("Unknown method %s", methodName)
	}

	return result
}

func GetAmountOutFromTransaction(ctx sdk.Context, k keeper.Keeper, abi abi.ABI,
	tx *ethtypes.Transaction, transfers []*types.TransferDetails) ([]*big.Int, error) {
	methodName, params, err := DecodeTxParams(abi, tx.Data())
	if err != nil {
		return nil, err
	}

	switch methodName {
	case "transferIn":
		amount, err := parseTransferIn(ctx, k, params, transfers[0])
		if err != nil {
			return nil, err
		}
		return []*big.Int{amount}, nil

	case "transferInMultiple":
		return parseTransferInMultiple(ctx, k, params, transfers)

	default:
		return nil, fmt.Errorf("Unsupported method %s", methodName)
	}
}

func parseTransferOut(ctx sdk.Context, keeper keeper.Keeper, ethTx *ethtypes.Transaction, chain string,
	isEvm bool, txParams map[string]interface{}) ([]*types.TransferDetails, error) {
	msg, err := ethTx.AsMessage(ethtypes.NewLondonSigner(ethTx.ChainId()), nil)
	if err != nil {
		return nil, err
	}

	tokenAddr, ok := txParams["token"].(ethcommon.Address)
	if !ok {
		err := fmt.Errorf("cannot convert token to type ethcommon.Address: %v", txParams)
		return nil, err
	}

	allTokens := keeper.GetAllTokens(ctx)
	token := utils.GetTokenOnChain(allTokens, tokenAddr.String(), chain)
	if token == nil {
		return nil, fmt.Errorf("Cannot find token on chain %s with address %s", chain, tokenAddr.String())
	}

	destChain, ok := txParams["dstChain"].(*big.Int)
	if !ok {
		err := fmt.Errorf("cannot convert destChain to type string: %v", txParams)
		return nil, err
	}
	to := libchain.GetChainNameFromInt(destChain)
	if to == "" {
		return nil, fmt.Errorf("Unknown destChain %s", destChain)
	}

	var recipient string
	if isEvm {
		recipientAddr, ok := txParams["to"].(ethcommon.Address)
		if !ok {
			err := fmt.Errorf("cannot convert recipient to type ethcommon.Address: %v", txParams)
			return nil, err
		}
		recipient = recipientAddr.String()
	} else {
		recipient, ok = txParams["to"].(string)
		if !ok {
			err := fmt.Errorf("cannot convert recipient to type ethcommon.Address: %v", txParams)
			return nil, err
		}
	}

	amount, ok := txParams["amount"].(*big.Int)
	if !ok {
		err := fmt.Errorf("cannot convert _amount to type *big.Int: %v", txParams)
		return nil, err
	}

	return []*types.TransferDetails{
		{
			Id:          types.GetTransferId(chain, ethTx.Hash().String()),
			TxType:      types.TxInType_TOKEN_TRANSFER,
			FromChain:   chain,
			FromSender:  msg.From().Hex(),
			FromHash:    ethTx.Hash().String(),
			Token:       token.Id,
			Amount:      amount.String(),
			ToChain:     to,
			ToRecipient: recipient,
		},
	}, nil
}

func parseRemoteCall(ctx sdk.Context, ethTx *ethtypes.Transaction, chain string,
	isEvm bool, txParams map[string]interface{}) ([]*types.TransferDetails, error) {
	msg, err := ethTx.AsMessage(ethtypes.NewLondonSigner(ethTx.ChainId()), nil)
	if err != nil {
		return nil, err
	}

	appChain, ok := txParams["appChain"].(*big.Int)
	if !ok {
		err := fmt.Errorf("cannot convert appChain to type *big.Int: %v", txParams)
		return nil, err
	}
	appChainName := libchain.GetChainNameFromInt(appChain)
	if appChainName == "" {
		return nil, fmt.Errorf("Unknown destChain %s", appChain)
	}

	var app string
	if isEvm {
		appAddr, ok := txParams["app"].(ethcommon.Address)
		if !ok {
			err := fmt.Errorf("cannot convert recipient to type ethcommon.Address: %v", txParams)
			return nil, err
		}
		app = appAddr.String()
	} else {
		app, ok = txParams["app"].(string)
		if !ok {
			err := fmt.Errorf("cannot convert recipient to type ethcommon.Address: %v", txParams)
			return nil, err
		}
	}

	message, ok := txParams["message"].([]byte)
	if !ok {
		err := fmt.Errorf("cannot convert message to type []byte: %v", txParams)
		return nil, err
	}

	callGasLimit, ok := txParams["callGasLimit"].(uint64)
	if !ok {
		err := fmt.Errorf("cannot convert callGasLimit to type uint64: %v", txParams)
		return nil, err
	}

	return []*types.TransferDetails{
		{
			Id:           types.GetTransferId(chain, ethTx.Hash().String()),
			TxType:       types.TxInType_REMOTE_CALL,
			FromChain:    chain,
			FromSender:   msg.From().Hex(),
			FromHash:     ethTx.Hash().String(),
			ToChain:      appChainName,
			ToRecipient:  app,
			Message:      message,
			CallGasLimit: callGasLimit,
		},
	}, nil
}

func parseTransferIn(ctx sdk.Context, k keeper.Keeper, params map[string]interface{},
	transfer *types.TransferDetails) (*big.Int, error) {
	// 1. Validate amount
	amountParam, ok := params["amount"]
	if !ok {
		return nil, fmt.Errorf("param amount not found, params = %v", params)
	}
	amountInt, ok := amountParam.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("Amount param is not an integer, amountParam = %s", amountParam)
	}

	// 2. Validate recipient address
	if _, ok := params["to"].(ethcommon.Address); !ok {
		return nil, fmt.Errorf("Invalid recipient address, to = %s", params["to"])
	}

	// 3. Validate token
	tokenAddr, ok := params["token"].(ethcommon.Address)
	if !ok {
		return nil, fmt.Errorf("Invalid token address, to = %s", params["token"])
	}

	if err := validateToken(ctx, k, transfer, tokenAddr.String()); err != nil {
		return nil, err
	}

	return amountInt, nil
}

func parseTransferInMultiple(ctx sdk.Context, keeper keeper.Keeper, params map[string]interface{},
	transfers []*types.TransferDetails) ([]*big.Int, error) {
	// 1. Validate amounts array
	amountsParam, ok := params["amounts"]
	if !ok {
		return nil, fmt.Errorf("param amounts not found, params = %v", params)
	}
	amounts, ok := amountsParam.([]*big.Int)
	if !ok {
		return nil, fmt.Errorf("amountsParam is not an instance of []*bigInt, params = %s", amountsParam)
	}
	if len(amounts) != len(transfers) {
		return nil, fmt.Errorf("Amounts and transfers length do not match, expected %d, actual %d",
			len(transfers), len(amounts))
	}

	// 2. Validate recipient addrs
	tosParam, ok := params["tos"]
	if !ok {
		return nil, fmt.Errorf("param tos not found, params = %v", params)
	}
	tos, ok := tosParam.([]ethcommon.Address)
	if len(tos) != len(transfers) {
		return nil, fmt.Errorf("tos and transfers length do not match, expected %d, actual %d",
			len(transfers), len(tos))
	}

	// 3. Validate tokens
	tokensParam, ok := params["tokens"]
	if !ok {
		return nil, fmt.Errorf("param tokens not found, params = %v", params)
	}
	tokens, ok := tokensParam.([]ethcommon.Address)
	for i, token := range tokens {
		if err := validateToken(ctx, keeper, transfers[i], token.String()); err != nil {
			return nil, err
		}
	}

	return amounts, nil
}

func validateToken(ctx sdk.Context, k keeper.Keeper, transfer *types.TransferDetails,
	actualAddr string) error {
	token := k.GetToken(ctx, transfer.Token)
	if token == nil {
		// This is not a proposal's fault. The token validation should be done in at the handler TxIn.
		return fmt.Errorf("cannot find token %s in the keeper", transfer.Token)
	}

	expectedAddr := token.GetAddressForChain(transfer.ToChain)
	if !strings.EqualFold(expectedAddr, actualAddr) {
		return fmt.Errorf("Token address does not match, expected %s, actual %s", expectedAddr,
			actualAddr)
	}

	return nil
}

func DecodeTxParams(abi abi.ABI, callData []byte) (string, map[string]interface{}, error) {
	txParams := map[string]interface{}{}
	m, err := abi.MethodById(callData[:4])
	if err != nil {
		log.Error(err)
		return "", nil, err
	}

	if err := m.Inputs.UnpackIntoMap(txParams, callData[4:]); err != nil {
		log.Error(err)
		return "", nil, err
	}

	return m.Name, txParams, nil
}

func buildTransferIn(
	ctx sdk.Context,
	finalTokenAddrs []ethcommon.Address,
	finalRecipients []ethcommon.Address,
	finalAmounts []*big.Int,
) ([]byte, error) {
	vaultAbi, err := abi.JSON(strings.NewReader(vault.VaultABI))
	if err != nil {
		return nil, err
	}

	var output []byte
	if len(finalTokenAddrs) == 1 {
		output, err = vaultAbi.Pack(
			MethodTransferIn,
			finalTokenAddrs[0],
			finalRecipients[0],
			finalAmounts[0],
		)
	} else {
		output, err = vaultAbi.Pack(
			MethodTransferInMultiple,
			finalTokenAddrs,
			finalRecipients,
			finalAmounts,
		)
	}
	if err != nil {
		return nil, err
	}

	return output, nil
}

func buildRemoteExecute(
	ctx sdk.Context,
	commission *big.Int,
	finalCallerChains []*big.Int,
	finalCallers []ethcommon.Address,
	finalApps []ethcommon.Address,
	finalGasLimits []uint64,
	finalMessages [][]byte,
) ([]byte, error) {
	vaultAbi, err := abi.JSON(strings.NewReader(vault.VaultABI))
	if err != nil {
		return nil, err
	}

	var data []byte
	if len(finalMessages) == 1 {
		data, err = vaultAbi.Pack(
			MethodRemoteExecute,
			finalCallerChains[0],
			finalCallers[0],
			finalApps[0],
			finalGasLimits[0],
			finalMessages[0],
			commission,
		)
	} else {
		data, err = vaultAbi.Pack(
			MethodRemoteExecuteMultiple,
			finalCallerChains,
			finalCallers,
			finalApps,
			finalGasLimits,
			finalMessages,
			commission,
		)
	}
	if err != nil {
		return nil, err
	}

	return data, nil
}
