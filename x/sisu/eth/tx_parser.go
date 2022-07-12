package eth

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

func ParseEthTransferOut(ctx sdk.Context, ethTx *ethTypes.Transaction, srcChain string, gwAbi abi.ABI,
	keeper keeper.Keeper) (*types.TxIn, error) {
	callData := ethTx.Data()
	txParams, err := decodeTxParams(gwAbi, callData)
	if err != nil {
		return nil, err
	}

	msg, err := ethTx.AsMessage(ethtypes.NewEIP155Signer(ethTx.ChainId()), nil)
	if err != nil {
		return nil, err
	}

	tokenAddr, ok := txParams["_tokenOut"].(ethcommon.Address)
	if !ok {
		err := fmt.Errorf("cannot convert _tokenOut to type ethcommon.Address: %v", txParams)
		return nil, err
	}

	fmt.Println("tokenAddr = ", tokenAddr)

	// TODO: Optimize getting tokens
	allTokens := keeper.GetAllTokens(ctx)
	token := getTokenOnChain(allTokens, tokenAddr.String(), srcChain)
	if token == nil {
		return nil, fmt.Errorf("Cannot find token on chain %s with address %s", srcChain, tokenAddr)
	}

	destChain, ok := txParams["_destChain"].(string)
	if !ok {
		err := fmt.Errorf("cannot convert _destChain to type string: %v", txParams)
		return nil, err
	}

	recipient, ok := txParams["_recipient"].(string)
	if !ok {
		err := fmt.Errorf("cannot convert _recipient to type ethcommon.Address: %v", txParams)
		return nil, err
	}

	amount, ok := txParams["_amount"].(*big.Int)
	if !ok {
		err := fmt.Errorf("cannot convert _amount to type *big.Int: %v", txParams)
		return nil, err
	}

	return &types.TxIn{
		Sender:    msg.From().Hex(),
		ToChain:   destChain,
		Token:     token.Id,
		Recipient: recipient,
		Amount:    amount.String(),
		Hash:      ethTx.Hash().String(),
	}, nil
}

func ParseEthTransferIn(ctx sdk.Context, ethTx *ethTypes.Transaction, destChain string, gwAbi abi.ABI,
	keeper keeper.Keeper) ([]*types.TransferOutData, error) {
	callData := ethTx.Data()
	txParams, err := decodeTxParams(gwAbi, callData)
	if err != nil {
		return nil, err
	}

	tokensAddrs, ok := txParams["tokens"].([]ethcommon.Address)
	if !ok {
		err := fmt.Errorf("cannot convert tokens to type eth address]: %v", txParams)
		return nil, err
	}

	recipients, ok := txParams["recipients"].([]ethcommon.Address)
	if !ok {
		err := fmt.Errorf("cannot convert recipients to type eth address]: %v", txParams)
		return nil, err
	}

	amounts, ok := txParams["amounts"].([]*big.Int)
	if !ok {
		err := fmt.Errorf("cannot convert amounts to type []*big.Int: %v", txParams)
		return nil, err
	}

	allTokens := keeper.GetAllTokens(ctx)
	transfers := make([]*types.TransferOutData, 0)
	for i, tokenAddr := range tokensAddrs {
		token := getTokenOnChain(allTokens, tokenAddr.String(), destChain)
		if token == nil {
			return nil, fmt.Errorf("ParseEthTransferIn: Cannot find token on chain %s with address %s", destChain, tokenAddr)
		}

		transfers = append(transfers, &types.TransferOutData{
			DestChain: destChain,
			Token:     token,
			Recipient: recipients[i].String(),
			Amount:    amounts[i],
		})
	}

	return transfers, nil
}

func getTokenOnChain(allTokens map[string]*types.Token, tokenAddr, chain string) *types.Token {
	for _, t := range allTokens {
		for j, chain := range t.Chains {
			if chain == chain && t.Addresses[j] == tokenAddr {
				return t
			}
		}
	}

	return nil
}

func decodeTxParams(abi abi.ABI, callData []byte) (map[string]interface{}, error) {
	if len(callData) < 4 {
		return nil, fmt.Errorf("decodeTxParams: call data size is smaller than 4")
	}

	txParams := map[string]interface{}{}
	m, err := abi.MethodById(callData[:4])
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err := m.Inputs.UnpackIntoMap(txParams, callData[4:]); err != nil {
		log.Error(err)
		return nil, err
	}

	return txParams, nil
}
