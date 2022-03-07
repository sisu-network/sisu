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

func (p *DefaultTxOutputProducer) ContractChangeOwnership(_ sdk.Context, chain, contractHash, newOwner string) (*types.TxOutWithSigner, error) {
	if !libchain.IsETHBasedChain(chain) {
		return nil, fmt.Errorf("unsupported chain %s", chain)
	}

	// TODO: Support more than gateway contract
	targetContractName := ContractErc20Gateway
	gw := p.publicDb.GetLatestContractAddressByName(chain, targetContractName)
	if len(gw) == 0 {
		err := fmt.Errorf("ContractChangeOwnership: cannot find gw address for type: %s", targetContractName)
		log.Error(err)
		return nil, err
	}

	gatewayAddress := ethcommon.HexToAddress(gw)
	erc20gatewayContract := SupportedContracts[targetContractName]

	nonce := p.worldState.UseAndIncreaseNonce(chain)
	if nonce < 0 {
		err := errors.New("PauseEthContract: cannot find nonce for chain " + chain)
		log.Error(err)
		return nil, err
	}

	gasPrice, err := p.worldState.GetGasPrice(chain)
	if err != nil {
		return nil, err
	}

	input, err := erc20gatewayContract.Abi.Pack(MethodChangeOwnership, ethcommon.HexToAddress(newOwner))
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

	return types.NewMsgTxOutWithSigner(
		p.appKeys.GetSignerAddress().String(),
		types.TxOutType_NORMAL,
		0,
		"",                    // in chain
		"",                    // in hash
		chain,                 // out chain
		rawTx.Hash().String(), // out hash
		bz,
		contractHash, // contract hash
	), nil

}
