package sisu

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func decodeTxParams(abi abi.ABI, callData []byte) (map[string]interface{}, error) {
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

func (p *DefaultTxOutputProducer) processERC20TransferIn(ethTx *ethTypes.Transaction) (*types.TxResponse, error) {
	log.Debug("Processing ERC20 transfer In")
	erc20gatewayContract := SupportedContracts[ContractErc20Gateway]
	gwAbi := erc20gatewayContract.Abi
	callData := ethTx.Data()
	txParams, err := decodeTxParams(gwAbi, callData)
	if err != nil {
		return nil, err
	}

	tokenAddr, ok := txParams["_tokenIn"].(ethcommon.Address)
	if !ok {
		err := fmt.Errorf("cannot convert _tokenIn to type ethcommon.Address: %v", txParams)
		log.Error(err)
		return nil, err
	}

	destChain, ok := txParams["_destChain"].(string)
	if !ok {
		err := fmt.Errorf("cannot convert _destChain to type string: %v", txParams)
		log.Error(err)
		return nil, err
	}

	token := p.worldState.GetTokenFromAddress(destChain, tokenAddr.String())
	if token == nil {
		return nil, fmt.Errorf("invalid address %s on chain %s", tokenAddr, destChain)
	}

	recipient, ok := txParams["_recipient"].(ethcommon.Address)
	if !ok {
		err := fmt.Errorf("cannot convert _recipient to type ethcommon.Address: %v", txParams)
		log.Error(err)
		return nil, err
	}

	amount, ok := txParams["_amount"].(*big.Int)
	if !ok {
		err := fmt.Errorf("cannot convert _amount to type *big.Int: %v", txParams)
		log.Error(err)
		return nil, err
	}

	return p.callERC20TransferIn(token, tokenAddr, recipient, amount, destChain)
}

func (p *DefaultTxOutputProducer) callERC20TransferIn(
	token *types.Token,
	tokenAddress,
	recipient ethcommon.Address,
	amountIn *big.Int,
	destChain string,
) (*types.TxResponse, error) {
	targetContractName := ContractErc20Gateway
	gw := p.publicDb.GetLatestContractAddressByName(destChain, targetContractName)
	if len(gw) == 0 {
		err := fmt.Errorf("cannot find gw address for type: %s", targetContractName)
		log.Error(err)
		return nil, err
	}

	gatewayAddress := ethcommon.HexToAddress(gw)
	erc20gatewayContract := SupportedContracts[targetContractName]

	nonce := p.worldState.UseAndIncreaseNonce(destChain)
	if nonce < 0 {
		err := errors.New("cannot find nonce for chain " + destChain)
		log.Error(err)
		return nil, err
	}

	gasPrice, err := p.worldState.GetGasPrice(destChain)
	if err != nil {
		return nil, err
	}

	// Calculate the output amount
	amountOut := new(big.Int).Set(amountIn)

	// 1. TODO: Subtract the commission fee.

	if token.Price == 0 {
		return nil, fmt.Errorf("token %s has price 0", token.Id)
	}

	// 2. Subtract the network gas fee on destination chain.
	gas := big.NewInt(8_000_000) // TODO: Show the correct gas cost here.
	gasPriceInToken, err := p.getGasCostInToken(gas, gasPrice, destChain, token)
	if err != nil {
		return nil, err
	}

	amountOut.Sub(amountOut, gasPriceInToken)

	if amountOut.Cmp(big.NewInt(0)) < 0 {
		return nil, ErrInsufficientFund
	}

	log.Debugf("destChain: %s, gateway address on destChain: %s, tokenAddr: %s, recipient: %s, gasPriceInToken: %s, amountIn: %s, amountOut: %s",
		destChain, gatewayAddress.String(), tokenAddress, recipient, gasPriceInToken.String(), amountIn.String(), amountOut.String(),
	)

	input, err := erc20gatewayContract.Abi.Pack(MethodTransferIn, tokenAddress, recipient, amountOut)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	rawTx := ethTypes.NewTransaction(
		uint64(nonce),
		gatewayAddress,
		big.NewInt(0),
		p.getGasLimit(destChain),
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

func (p *DefaultTxOutputProducer) getGasCostInToken(gas *big.Int, gasPrice *big.Int, chain string, token *types.Token) (*big.Int, error) {
	// Get total gas cost
	gasCost := new(big.Int).Mul(gas, gasPrice)

	chainTokenPrice, err := p.worldState.GetNativeTokenPriceForChain(chain)
	if err != nil {
		return nil, err
	}

	// amount := gasCost * chainTokenPrice / tokenPrice
	gasInToken := new(big.Int).Mul(gasCost, big.NewInt(chainTokenPrice))
	gasInToken = new(big.Int).Div(gasInToken, big.NewInt(token.Price))

	return gasInToken, nil
}
