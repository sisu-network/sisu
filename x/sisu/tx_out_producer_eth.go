package sisu

import (
	"fmt"
	"math/big"

	ecommon "github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/sisu-network/sisu/x/sisu/world"
)

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

func parseEthTransferOut(ethTx *ethTypes.Transaction, srcChain string, worldState world.WorldState) (*types.TransferOutData, error) {
	erc20gatewayContract := SupportedContracts[ContractErc20Gateway]
	gwAbi := erc20gatewayContract.Abi
	callData := ethTx.Data()
	txParams, err := decodeTxParams(gwAbi, callData)
	if err != nil {
		return nil, err
	}

	tokenAddr, ok := txParams["_tokenOut"].(ethcommon.Address)
	if !ok {
		err := fmt.Errorf("cannot convert _tokenOut to type ethcommon.Address: %v", txParams)
		return nil, err
	}

	destChain, ok := txParams["_destChain"].(string)
	if !ok {
		err := fmt.Errorf("cannot convert _destChain to type string: %v", txParams)
		return nil, err
	}

	token := worldState.GetTokenFromAddress(srcChain, tokenAddr.String())
	if token == nil {
		return nil, fmt.Errorf("invalid address %s on chain %s", tokenAddr, srcChain)
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

	return &types.TransferOutData{
		DestChain: destChain,
		Token:     token,
		Recipient: recipient,
		Amount:    amount,
	}, nil
}

func parseTransferInData(ethTx *ethTypes.Transaction) ([]*transferInData, error) {
	erc20gatewayContract := SupportedContracts[ContractErc20Gateway]
	gwAbi := erc20gatewayContract.Abi
	callData := ethTx.Data()
	txParams, err := decodeTxParams(gwAbi, callData)
	if err != nil {
		return nil, err
	}

	tokens, ok := txParams["tokens"].([]ethcommon.Address)
	if !ok {
		err := fmt.Errorf("parseTransferInData: cannot convert _token to type ethcommon.Address: %v", txParams)
		return nil, err
	}

	recipients, ok := txParams["recipients"].([]ethcommon.Address)
	if !ok {
		err := fmt.Errorf("parseTransferInData: cannot convert _recipient to type ethcommon.Address: %v", txParams)
		return nil, err
	}

	amounts, ok := txParams["amounts"].([]*big.Int)
	if !ok {
		err := fmt.Errorf("parseTransferInData: cannot convert _amount to type *big.Int: %v", txParams)
		return nil, err
	}

	txIns := make([]*transferInData, len(tokens))
	for i := range tokens {
		txIns[i] = &transferInData{
			token:     tokens[i],
			recipient: recipients[i].String(),
			amount:    amounts[i],
		}
	}

	return txIns, nil
}

func (p *DefaultTxOutputProducer) buildERC20TransferIn(
	ctx sdk.Context,
	k keeper.Keeper,
	tokens []*types.Token,
	recipients []ethcommon.Address,
	amounts []*big.Int,
	destChain string,
) (*types.TxResponse, error) {
	targetContractName := ContractErc20Gateway
	gw := p.keeper.GetLatestContractAddressByName(ctx, destChain, targetContractName)
	if len(gw) == 0 {
		err := fmt.Errorf("cannot find gw address for type: %s", targetContractName)
		log.Error(err)
		return nil, err
	}

	gatewayAddress := ethcommon.HexToAddress(gw)
	erc20gatewayContract := SupportedContracts[targetContractName]

	checkPoint := k.GetGatewayCheckPoint(ctx, destChain)
	if checkPoint == nil {
		return nil, fmt.Errorf("cannot find gateway checkout for chain %s", destChain)
	}

	nonce := checkPoint.Nonce

	gasPrice, err := p.worldState.GetGasPrice(destChain)
	if err != nil {
		return nil, err
	}

	if gasPrice.Cmp(big.NewInt(0)) <= 0 {
		gasPrice = p.getDefaultGasPrice(destChain)
	}

	log.Debug("Gas price for swapping  = ", gasPrice)

	finalTokenAddrs := make([]ethcommon.Address, 0)
	finalRecipients := make([]ethcommon.Address, 0)
	finalAmounts := make([]*big.Int, 0)
	amountIns := make([]*big.Int, 0)
	gasPrices := make([]int64, 0)

	for i := range amounts {
		fmt.Println("tokens = ", tokens)
		fmt.Println("i = ", i)
		if tokens[i].Price == 0 {
			return nil, fmt.Errorf("token %s has price 0", tokens[i].Id)
		}

		// 1. TODO: Subtract the commission fee.
		gasPriceInToken, err := p.worldState.GetGasCostInToken(tokens[i].Id, destChain)
		if err != nil {
			return nil, err
		}

		if gasPriceInToken < 0 {
			log.Errorf("Gas price in token is negative: token id = %s", tokens[i].Id)
			gasPriceInToken = 0
		}

		// Calculate the output amount
		amountOut := new(big.Int).Set(amounts[i])
		amountOut.Sub(amountOut, big.NewInt(gasPriceInToken))

		if amountOut.Cmp(big.NewInt(0)) < 0 {
			log.Error("Insufficient fund for transfer ", i)
			continue
		}

		// Find the address of the token.
		var tokenAddr string
		for _, token := range tokens {
			for j, chain := range token.Chains {
				if chain == destChain {
					tokenAddr = token.Addresses[j]
					break
				}
			}
			if len(tokenAddr) > 0 {
				break
			}
		}
		if len(tokenAddr) == 0 {
			continue
		}

		finalTokenAddrs = append(finalTokenAddrs, ecommon.HexToAddress(tokenAddr))
		finalAmounts = append(finalAmounts, amountOut)
		finalRecipients = append(finalRecipients, recipients[i])
		amountIns = append(amountIns, amounts[i])
		gasPrices = append(gasPrices, gasPriceInToken)
	}

	if len(finalTokenAddrs) == 0 {
		return nil, fmt.Errorf("No txOut is produced (might be due to insufficient fund")
	}

	log.Verbosef("destChain: %s, gateway address on destChain: %s", destChain, gatewayAddress.String())
	for i := range finalTokenAddrs {
		log.Verbosef("tokenAddr: %s, recipient: %s, gasPriceInToken: %d, amountIn: %s, amountOut: %s",
			finalTokenAddrs[i], finalRecipients[i], gasPrices[i], amountIns[i].String(), finalAmounts[i].String(),
		)
	}

	input, err := erc20gatewayContract.Abi.Pack(
		MethodTransferIn,
		finalTokenAddrs,
		finalRecipients,
		finalAmounts,
	)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	rawTx := ethTypes.NewTransaction(
		uint64(nonce),
		gatewayAddress,
		big.NewInt(0),
		100_000, // 100k for swapping operation.
		gasPrice,
		input,
	)

	bz, err := rawTx.MarshalBinary()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &types.TxResponse{
		OutChain: destChain,
		EthTx:    rawTx,
		RawBytes: bz,
	}, nil
}