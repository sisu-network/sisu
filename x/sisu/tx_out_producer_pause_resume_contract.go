package sisu

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func (p *DefaultTxOutputProducer) PauseContract(ctx sdk.Context, chain string, hash string) (*types.TxOutMsg, error) {
	if libchain.IsETHBasedChain(chain) {
		return p.PauseOrResumeEthContract(ctx, chain, hash, true)
	}

	return nil, fmt.Errorf("unsupported chain %s", chain)
}

func (p *DefaultTxOutputProducer) ResumeContract(ctx sdk.Context, chain string, hash string) (*types.TxOutMsg, error) {
	if libchain.IsETHBasedChain(chain) {
		return p.PauseOrResumeEthContract(ctx, chain, hash, false)
	}

	return nil, fmt.Errorf("unsupported chain %s", chain)
}

func (p *DefaultTxOutputProducer) PauseOrResumeEthContract(ctx sdk.Context, chain string, hash string, isPause bool) (*types.TxOutMsg, error) {
	if !libchain.IsETHBasedChain(chain) {
		return nil, fmt.Errorf("unsupported chain %s", chain)
	}

	// TODO: Support more than gateway contract
	targetContractName := ContractErc20Gateway
	gw := p.keeper.GetLatestContractAddressByName(ctx, chain, targetContractName)
	if len(gw) == 0 {
		err := fmt.Errorf("PauseEthContract: cannot find gw address for type: %s", targetContractName)
		log.Error(err)
		return nil, err
	}

	gatewayAddress := ethcommon.HexToAddress(gw)
	erc20gatewayContract := SupportedContracts[targetContractName]

	c := p.keeper.GetChain(ctx, chain)
	if c == nil {
		return nil, fmt.Errorf("PauseOrResumeEthContract: invalid chain %s", c)
	}
	gasPrice := big.NewInt(c.GasPrice)

	var input []byte
	var err error
	if isPause {
		input, err = erc20gatewayContract.Abi.Pack(MethodPauseGateway)
	} else {
		input, err = erc20gatewayContract.Abi.Pack(MethodResumeGateway)
	}

	rawTx := ethTypes.NewTransaction(
		0,
		gatewayAddress,
		big.NewInt(0),
		p.getGasLimit(chain),
		gasPrice,
		input,
	)

	bz, err := rawTx.MarshalBinary()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return types.NewTxOutMsg(
		p.signer,
		types.TxOutType_TRANSFER_OUT,
		[]string{""},          // in hash
		chain,                 // out chain
		rawTx.Hash().String(), // out hash
		bz,
		hash, // contract hash
	), nil

}
