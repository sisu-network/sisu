package sisu

import (
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func (p *DefaultTxOutputProducer) ContractEmergencyWithdrawFund(ctx sdk.Context, chain, contractHash string,
	tokens []string, newOwner string) (*types.TxOutWithSigner, error) {

	if !libchain.IsETHBasedChain(chain) {
		return nil, fmt.Errorf("unsupported chain %s", chain)
	}

	nonce := p.worldState.UseAndIncreaseNonce(ctx, chain)
	if nonce < 0 {
		err := errors.New("ContractEmergencyWithdrawFund: cannot find nonce for chain " + chain)
		log.Error(err)
		return nil, err
	}

	gasPrice, err := p.worldState.GetGasPrice(chain)
	if err != nil {
		return nil, err
	}

	liquidityContract := SupportedContracts[ContractLiquidityPool]
	tokenHashes := make([]common.Address, 0, len(tokens))
	for _, token := range tokens {
		tokenHashes = append(tokenHashes, common.HexToAddress(token))
	}

	input, err := liquidityContract.Abi.Pack(MethodEmergencyWithdrawFund, tokenHashes, common.HexToAddress(newOwner))
	if err != nil {
		log.Error("ContractEmergencyWithdrawFund: error when pack input ", err)
		return nil, err
	}

	rawTx := ethTypes.NewTransaction(
		uint64(nonce),
		common.HexToAddress(contractHash),
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
		types.TxOutType_LIQUIDITY_WITHDRAW_FUND,
		0,
		"",                    // in chain
		"",                    // in hash
		chain,                 // out chain
		rawTx.Hash().String(), // out hash
		bz,
		contractHash, // contract hash
	), nil
}
