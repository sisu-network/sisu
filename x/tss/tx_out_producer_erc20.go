package tss

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"
)

func (p *DefaultTxOutputProducer) processERC20TransferIn(ethTx *ethTypes.Transaction) (*types.TxResponse, error) {
	log.Debug("Processing ERC20 transfer In")
	erc20GatewayContract := SupportedContracts[ContractErc20Gateway]
	gwAbi := erc20GatewayContract.Abi
	callData := ethTx.Data()
	txParams, err := decodeTxParams(gwAbi, callData)
	if err != nil {
		return nil, err
	}

	tokenAddr, ok := txParams["_tokenIn"].(ethcommon.Address)
	if !ok {
		err := fmt.Errorf("cannot convert _token to type ethcommon.Address: %v", txParams)
		log.Error(err)
		return nil, err
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

	destChain, ok := txParams["_destChain"].(string)
	if !ok {
		err := fmt.Errorf("cannot convert _destChain to type string: %v", txParams)
		log.Error(err)
		return nil, err
	}

	return p.callERC20TransferIn(tokenAddr, recipient, amount, destChain)
}

func (p *DefaultTxOutputProducer) callERC20TransferIn(tokenAddress, recipient ethcommon.Address, amount *big.Int, destChain string) (*types.TxResponse, error) {
	gw := p.privateDb.GetLatestContractAddressByName(destChain, ContractErc20Gateway)
	if len(gw) == 0 {
		return nil, nil
	}

	gatewayAddress := ethcommon.HexToAddress(gw)
	erc20GatewayContract := SupportedContracts[ContractErc20Gateway]

	input, err := erc20GatewayContract.Abi.Pack(MethodTransferIn, tokenAddress, recipient, amount)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	nonce := p.worldState.UseAndIncreaseNonce(destChain)
	if nonce < 0 {
		err := errors.New("cannot find nonce for chain " + destChain)
		log.Error(err)
		return nil, err
	}

	log.Debugf("destChain: %s, gateway address on destChain: %s, tokenAddr: %s, recipient: %s, amount: %d", destChain, gatewayAddress.String(), tokenAddress, recipient, amount.Int64())
	rawTx := ethTypes.NewTx(&ethTypes.AccessListTx{
		Nonce:    uint64(nonce),
		GasPrice: p.getGasPrice(destChain),
		Gas:      p.getGasLimit(destChain),
		To:       &gatewayAddress,
		Value:    big.NewInt(0),
		Data:     input,
	})

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
