package sisu

import (
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"

	ecommon "github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/sisu-network/sisu/x/sisu/world"
)

// Check if we can deploy contract after seeing some ETH being sent to our ethereum address.
func (p *DefaultTxOutputProducer) getEthContractDeploymentTx(ctx sdk.Context, height int64, chain string, contracts []*types.Contract) []*ethTypes.Transaction {
	txs := make([]*ethTypes.Transaction, 0)

	for _, contract := range contracts {
		nonce := p.worldState.UseAndIncreaseNonce(ctx, chain)
		log.Verbose("nonce for deploying contract:", nonce, " on chain ", chain)
		if nonce < 0 {
			log.Error("cannot get nonce for contract")
			continue
		}

		rawTx := p.getContractTx(ctx, contract, nonce)
		if rawTx == nil {
			log.Warn("raw Tx is nil")
			continue
		}

		txs = append(txs, rawTx)
	}

	return txs
}

func (p *DefaultTxOutputProducer) getContractTx(ctx sdk.Context, contract *types.Contract, nonce int64) *ethTypes.Transaction {
	erc20 := SupportedContracts[ContractErc20Gateway]
	switch contract.Hash {
	case erc20.AbiHash:
		// This is erc20gw contract.
		parsedAbi, err := abi.JSON(strings.NewReader(erc20.AbiString))
		if err != nil {
			log.Error("cannot parse erc20 abi. abi = ", erc20.AbiString, "err =", err)
			return nil
		}

		// Get all allowed chains
		supportedChains := make([]string, 0)
		for chain := range p.tssConfig.SupportedChains {
			if chain != contract.Chain {
				supportedChains = append(supportedChains, chain)
			}
		}

		sort.Strings(supportedChains)

		log.Info("Allowed chains for chain ", contract.Chain, " are: ", supportedChains)

		lp := p.keeper.GetLiquidity(ctx, contract.Chain)
		if lp == nil {
			log.Warn("Lp is nil for chain ", contract.Chain)
			return nil
		}

		log.Infof("Liquidity pool addr for chain %s is %s", contract.Chain, lp.Address)
		input, err := parsedAbi.Pack("", supportedChains, ecommon.HexToAddress(lp.Address))
		if err != nil {
			log.Error("cannot pack supportedChains, err =", err)
			return nil
		}

		byteCode := ecommon.FromHex(erc20.Bin)
		input = append(byteCode, input...)
		chain := p.keeper.GetChain(ctx, contract.Chain)
		if chain == nil {
			log.Error("getContractTx: chain is nil with id ", contract.Chain)
			return nil
		}

		gasPrice := chain.GasPrice
		if gasPrice <= 0 {
			gasPrice = p.getDefaultGasPrice(contract.Chain).Int64()
		}
		gasLimit := p.getGasLimit(contract.Chain)

		log.Verbose("Gas price = ", gasPrice, " on chain ", contract.Chain)
		log.Verbose("gasLimit = ", gasLimit, " on chain ", contract.Chain)

		rawTx := ethTypes.NewContractCreation(
			uint64(nonce),
			big.NewInt(0),
			gasLimit,
			big.NewInt(gasPrice),
			input,
		)

		return rawTx
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

func parseEthTransferOut(ethTx *ethTypes.Transaction, srcChain string, worldState world.WorldState) (*transferOutData, error) {
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

	return &transferOutData{
		tokenAddr: tokenAddr,
		destChain: destChain,
		token:     token,
		recipient: recipient,
		amount:    amount,
	}, nil
}

func parseTransferInData(ethTx *ethTypes.Transaction) (*transferInData, error) {
	erc20gatewayContract := SupportedContracts[ContractErc20Gateway]
	gwAbi := erc20gatewayContract.Abi
	callData := ethTx.Data()
	txParams, err := decodeTxParams(gwAbi, callData)
	if err != nil {
		return nil, err
	}

	token, ok := txParams["_token"].(ethcommon.Address)
	if !ok {
		err := fmt.Errorf("parseTransferInData: cannot convert _token to type ethcommon.Address: %v", txParams)
		return nil, err
	}

	recipient, ok := txParams["_recipient"].(ethcommon.Address)
	if !ok {
		err := fmt.Errorf("parseTransferInData: cannot convert _recipient to type ethcommon.Address: %v", txParams)
		return nil, err
	}

	amount, ok := txParams["_amount"].(*big.Int)
	if !ok {
		err := fmt.Errorf("parseTransferInData: cannot convert _amount to type *big.Int: %v", txParams)
		return nil, err
	}

	return &transferInData{
		token:     token,
		recipient: recipient.String(),
		amount:    amount,
	}, nil
}

func (p *DefaultTxOutputProducer) buildERC20TransferIn(
	ctx sdk.Context,
	token *types.Token,
	tokenAddress,
	recipient ethcommon.Address,
	amountIn *big.Int,
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

	nonce := p.worldState.UseAndIncreaseNonce(ctx, destChain)
	if nonce < 0 {
		err := errors.New("cannot find nonce for chain " + destChain)
		log.Error(err)
		return nil, err
	}

	gasPrice, err := p.worldState.GetGasPrice(destChain)
	if err != nil {
		return nil, err
	}

	if gasPrice.Cmp(big.NewInt(0)) <= 0 {
		gasPrice = p.getDefaultGasPrice(destChain)
	}

	log.Debug("Gas price for swapping  = ", gasPrice)

	// Calculate the output amount
	amountOut := new(big.Int).Set(amountIn)

	// 1. TODO: Subtract the commission fee.

	if token.Price == 0 {
		return nil, fmt.Errorf("token %s has price 0", token.Id)
	}

	gasPriceInToken, err := p.worldState.GetGasCostInToken(token.Id, destChain)
	if err != nil {
		return nil, err
	}

	if gasPriceInToken < 0 {
		log.Errorf("Gas price in token is negative: token id = %s", token.Id)
		gasPriceInToken = 0
	}

	amountOut.Sub(amountOut, big.NewInt(gasPriceInToken))

	if amountOut.Cmp(big.NewInt(0)) < 0 {
		return nil, world.ErrInsufficientFund
	}

	log.Verbosef("destChain: %s, gateway address on destChain: %s, tokenAddr: %s, recipient: %s, gasPriceInToken: %d, amountIn: %s, amountOut: %s",
		destChain, gatewayAddress.String(), tokenAddress, recipient, gasPriceInToken, amountIn.String(), amountOut.String(),
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
