package sisu

import (
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func (p *DefaultTxOutputProducer) ContractSetLiquidPoolAddress(ctx sdk.Context, chain, contractHash, newAddress string) (*types.TxOutMsg, error) {
	if !libchain.IsETHBasedChain(chain) {
		return nil, fmt.Errorf("unsupported chain %s", chain)
	}

	// TODO: Support more than gateway contract
	targetContractName := ContractErc20Gateway
	gw := p.keeper.GetLatestContractAddressByName(ctx, chain, targetContractName)
	if len(gw) == 0 {
		err := fmt.Errorf("ContractSetLiquidPoolAddress: cannot find gw address for type: %s", targetContractName)
		log.Error(err)
		return nil, err
	}

	gatewayAddress := ethcommon.HexToAddress(gw)
	erc20gatewayContract := SupportedContracts[targetContractName]

	nonce := p.worldState.UseAndIncreaseNonce(ctx, chain)
	if nonce < 0 {
		err := errors.New("PauseEthContract: cannot find nonce for chain " + chain)
		log.Error(err)
		return nil, err
	}

	gasPrice, err := p.worldState.GetGasPrice(chain)
	if err != nil {
		return nil, err
	}

	input, err := erc20gatewayContract.Abi.Pack(MethodSetLiquidAddress, ethcommon.HexToAddress(newAddress))
	rawTx := ethTypes.NewTransaction(
		uint64(nonce),
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
		p.appKeys.GetSignerAddress().String(),
		types.TxOutType_CHANGE_LIQUIDITY,
		[]string{""},          // in hash
		chain,                 // out chain
		rawTx.Hash().String(), // out hash
		bz,
		contractHash, // contract hash
	), nil
}
