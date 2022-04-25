package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (p *Processor) EndBlockValidator(ctx sdk.Context) []abci.ValidatorUpdate {
	newVals, oldVals, err := p.getChangedNodes(ctx)
	if err != nil {
		return []abci.ValidatorUpdate{}
	}

	validators := make([]abci.ValidatorUpdate, 0, len(newVals)+len(oldVals))
	for _, val := range newVals {
		validators = append(validators, abci.Ed25519ValidatorUpdate(val.ConsensusKey.GetBytes(), 100))
	}

	for _, val := range oldVals {
		validators = append(validators, abci.Ed25519ValidatorUpdate(val.ConsensusKey.GetBytes(), 0))
	}

	return validators
}

// detects candidate nodes will be promoted to active node and active nodes will be removed from validator set
func (p *Processor) getChangedNodes(ctx sdk.Context) ([]*types.Node, []types.Node, error) {
	return nil, nil, nil
}
