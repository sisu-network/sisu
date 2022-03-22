package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func (p *DefaultTxOutputProducer) ContractEmergencyWithdrawFund(ctx sdk.Context, chain, contractHash string,
	tokens []string, newOwner string) (*types.TxOutWithSigner, error) {

}
