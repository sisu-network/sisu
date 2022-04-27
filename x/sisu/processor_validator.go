package sisu

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (p *Processor) EndBlockValidator(ctx sdk.Context) []abci.ValidatorUpdate {
	newValidators, oldValidators, err := p.getChangedNodes(ctx)
	if err != nil {
		return []abci.ValidatorUpdate{}
	}

	validators := make([]abci.ValidatorUpdate, 0, len(newValidators)+len(oldValidators))
	for _, val := range newValidators {
		p.validatorManager.UpdateNodeStatus(ctx, val.AccAddress, val.ConsensusKey.GetBytes(), types.NodeStatus_Validator)
		validators = append(validators, abci.Ed25519ValidatorUpdate(val.ConsensusKey.GetBytes(), 100))
	}

	for _, val := range oldValidators {
		p.validatorManager.UpdateNodeStatus(ctx, val.AccAddress, val.ConsensusKey.GetBytes(), types.NodeStatus_Candidate)
		validators = append(validators, abci.Ed25519ValidatorUpdate(val.ConsensusKey.GetBytes(), 0))
	}

	return validators
}

// detects candidate nodes will be promoted to active node and active nodes will be removed from validator set
func (p *Processor) getChangedNodes(ctx sdk.Context) ([]*types.Node, []*types.Node, error) {
	removedNodes, err := p.validatorManager.GetExceedSlashThresholdValidators(ctx)
	if err != nil {
		return nil, nil, err
	}

	// TODO: only get candidates node
	topCandidates := p.keeper.GetTopBalance(ctx, len(removedNodes))
	newNodes := make([]*types.Node, len(topCandidates))
	for _, candidate := range topCandidates {
		vals := p.validatorManager.GetVals(ctx)
		node, ok := vals[candidate.String()]
		if !ok {
			err = errors.New(fmt.Sprintf("can not find validator info. addr = %s", candidate.String()))
			log.Error(err)
			return nil, nil, err
		}

		newNodes = append(newNodes, node)
	}
	return newNodes, removedNodes, nil
}
