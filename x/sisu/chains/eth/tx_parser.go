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
	default:
		result.Method = chainstypes.MethodUnknown
		result.Error = fmt.Errorf("Unknown method %s", methodName)
	}

	return result
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

func parseTransferIn(ctx sdk.Context, keeper keeper.Keeper, ethTx *ethtypes.Transaction) (map[string]any, error) {
	vaultAbi, _ := abi.JSON(strings.NewReader(vault.VaultABI))
	callData := ethTx.Data()
	if len(callData) < 4 {
		// This is just a normal transfer
		return nil, fmt.Errorf("Invalid transferIn data")
	}

	_, txParams, err := DecodeTxParams(vaultAbi, callData)
	if err != nil {
		return nil, err
	}

	return txParams, nil
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
