package sisu

import (
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func (p *DefaultTxOutputProducer) processPauseGw(chain string) (*types.TxResponse, error) {
	contractName := ContractErc20Gateway
	gw := p.publicDb.GetLatestContractAddressByName(chain, contractName)
	if len(gw) == 0 {
		err := fmt.Errorf("cannot find gw address for type: %s", contractName)
		log.Error(err)
		return nil, err
	}

	gwAddr := ethcommon.HexToAddress(gw)
	gwContract := SupportedContracts[contractName]

	input, err := gwContract.Abi.Pack(MethodPauseGw)
	if err != nil {
		log.Error("error when pack input: ", err)
		return nil, err
	}

	nonce := p.worldState.UseAndIncreaseNonce(chain)
	if nonce < 0 {
		err := fmt.Errorf("cannot find nonce for chain %s", chain)
		log.Error(err)
		return nil, err
	}

	gasPrice := p.privateDb.GetNetworkGasPrice(chain)
	if gasPrice < 0 {
		gasPrice = p.getDefaultGasPrice(chain).Int64()
	}

	log.Debug("Network gas price got: ", gasPrice)
	rawTx := ethTypes.NewTransaction(
		uint64(nonce),
		gwAddr,
		big.NewInt(0),
		p.getGasLimit(chain),
		big.NewInt(gasPrice),
		input,
	)

	bz, err := rawTx.MarshalBinary()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &types.TxResponse{
		OutChain: chain,
		EthTx:    rawTx,
		RawBytes: bz,
	}, nil
}
