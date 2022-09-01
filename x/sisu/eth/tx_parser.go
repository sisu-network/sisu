package eth

import (
	"fmt"
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/contracts/eth/vault"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

func ParseEthTransferOut(ctx sdk.Context, ethTx *ethTypes.Transaction, srcChain string, keeper keeper.Keeper) (*types.TxIn, error) {
	gwAbi, err := abi.JSON(strings.NewReader(vault.VaultABI))
	if err != nil {
		return nil, err
	}

	callData := ethTx.Data()
	methodName, txParams, err := DecodeTxParams(gwAbi, callData)
	if err != nil {
		return nil, err
	}

	if methodName != "transferOut" && methodName != "transferOutMultiple" &&
		methodName != "transferOutNonEvm" && methodName != "transferOutMultipleNonEvm" {
		// This is not a transfer In function
		return nil, nil
	}

	msg, err := ethTx.AsMessage(ethtypes.NewLondonSigner(ethTx.ChainId()), nil)

	if err != nil {
		return nil, err
	}

	tokenAddr, ok := txParams["token"].(ethcommon.Address)
	if !ok {
		err := fmt.Errorf("cannot convert token to type ethcommon.Address: %v", txParams)
		return nil, err
	}

	// TODO: Optimize getting tokens
	allTokens := keeper.GetAllTokens(ctx)
	token := getTokenOnChain(allTokens, strings.ToLower(tokenAddr.String()), srcChain)
	if token == nil {
		return nil, fmt.Errorf("Cannot find token on chain %s with address %s", srcChain, strings.ToLower(tokenAddr.String()))
	}

	destChain, ok := txParams["dstChain"].(string)
	if !ok {
		err := fmt.Errorf("cannot convert _destChain to type string: %v", txParams)
		return nil, err
	}

	recipient, ok := txParams["to"].(ethcommon.Address)
	if !ok {
		err := fmt.Errorf("cannot convert _recipient to type ethcommon.Address: %v", txParams)
		return nil, err
	}

	amount, ok := txParams["amount"].(*big.Int)
	if !ok {
		err := fmt.Errorf("cannot convert _amount to type *big.Int: %v", txParams)
		return nil, err
	}

	return &types.TxIn{
		Sender:    msg.From().Hex(),
		ToChain:   destChain,
		Token:     token.Id,
		Recipient: recipient.String(),
		Amount:    amount.String(),
		Hash:      ethTx.Hash().String(),
	}, nil
}

func getTokenOnChain(allTokens map[string]*types.Token, tokenAddr, targetChain string) *types.Token {
	for _, t := range allTokens {
		if len(t.Chains) != len(t.Addresses) {
			log.Error("Chains length is not the same as address length ")
			log.Error("t.Chains = ", t.Chains)
			log.Error("t.Addresses = ", t.Addresses)
			return nil
		}

		for j, chain := range t.Chains {
			if chain == targetChain && t.Addresses[j] == tokenAddr {
				return t
			}
		}
	}

	return nil
}

func DecodeTxParams(abi abi.ABI, callData []byte) (string, map[string]interface{}, error) {
	if len(callData) < 4 {
		fmt.Println("len(callData) = ", len(callData))
		return "", nil, fmt.Errorf("decodeTxParams: call data size is smaller than 4")
	}

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
