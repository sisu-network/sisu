package sisu

import (
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

	if len(newValidators) == 0 && len(oldValidators) == 0 {
		return []abci.ValidatorUpdate{}
	}

	log.Debug("newValidators = ", newValidators)
	log.Debug("oldValidators = ", oldValidators)

	newValidatorKeys := make([][]byte, 0, len(newValidators))
	oldValidatorKeys := make([][]byte, 0, len(oldValidators))
	validators := make([]abci.ValidatorUpdate, 0, len(newValidators)+len(oldValidators))
	for _, val := range newValidators {
		p.validatorManager.UpdateNodeStatus(ctx, val.AccAddress, val.ConsensusKey.GetBytes(), types.NodeStatus_Validator)
		validators = append(validators, abci.Ed25519ValidatorUpdate(val.ConsensusKey.GetBytes(), 100))
		newValidatorKeys = append(newValidatorKeys, val.ConsensusKey.GetBytes())
	}

	for _, val := range oldValidators {
		p.validatorManager.UpdateNodeStatus(ctx, val.AccAddress, val.ConsensusKey.GetBytes(), types.NodeStatus_Candidate)
		validators = append(validators, abci.Ed25519ValidatorUpdate(val.ConsensusKey.GetBytes(), 0))
		oldValidatorKeys = append(oldValidatorKeys, val.ConsensusKey.GetBytes())
	}

	changeValSetMsg := types.NewChangeValidatorSetMsg(p.appKeys.GetSignerAddress().String(), oldValidatorKeys, newValidatorKeys)
	p.txSubmit.SubmitMessageAsync(changeValSetMsg)

	return validators
}

// detects candidate nodes will be promoted to active node and active nodes will be removed from validator set
func (p *Processor) getChangedNodes(ctx sdk.Context) ([]*types.Node, []*types.Node, error) {
	removedNodes, err := p.validatorManager.GetExceedSlashThresholdValidators(ctx)
	if err != nil {
		return nil, nil, err
	}

	if len(removedNodes) == 0 {
		log.Debug("removedNodes is empty")
		return nil, nil, nil
	}

	topCandidates := p.keeper.GetTopBalance(ctx, len(removedNodes))
	newNodes := make([]*types.Node, 0, len(removedNodes))
	for _, candidate := range topCandidates {
		vals := p.keeper.LoadNodesByStatus(ctx, types.NodeStatus_Unknown)
		for _, node := range vals {
			if node.AccAddress == candidate.String() {
				newNodes = append(newNodes, node)
				break
			}
		}
	}

	return newNodes, removedNodes, nil
}
