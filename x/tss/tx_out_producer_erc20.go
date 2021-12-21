package tss

import (
	"errors"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"
)

func (p *DefaultTxOutputProducer) createERC20TransferIn(gatewayAddress, tokenAddress, recipient string, amount *big.Int, destChain string) (*types.TxResponse, error) {
	erc20GatewayContract := SupportedContracts[ContractErc20]

	tokenAddr := ethcommon.HexToAddress(tokenAddress)
	recipientAddr := ethcommon.HexToAddress(recipient)
	input, err := erc20GatewayContract.Abi.Pack(MethodTransferIn, tokenAddr, recipientAddr, amount)
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

	gwAddress := ethcommon.HexToAddress(gatewayAddress)
	rawTx := ethTypes.NewTx(&ethTypes.AccessListTx{
		Nonce:    uint64(nonce),
		GasPrice: p.getGasPrice(destChain),
		Gas:      p.getGasLimit(destChain),
		To:       &gwAddress,
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

