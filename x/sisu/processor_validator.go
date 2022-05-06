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

	log.Debug("len = ", len(newValidators), " newValidators = ", newValidators)
	log.Debug("len = ", len(oldValidators), " oldValidators = ", oldValidators)

	newValidatorKeys := make([][]byte, 0)
	oldValidatorKeys := make([][]byte, 0)
	for _, val := range newValidators {
		newValidatorKeys = append(newValidatorKeys, val.ConsensusKey.GetBytes())
	}

	for _, val := range oldValidators {
		oldValidatorKeys = append(oldValidatorKeys, val.ConsensusKey.GetBytes())
	}

	log.Debug("len oldValidatorKeys = ", len(oldValidatorKeys))
	log.Debug("len newValidatorKeys = ", len(newValidatorKeys))
	changeValSetMsg := types.NewChangeValidatorSetMsg(p.appKeys.GetSignerAddress().String(), oldValidatorKeys, newValidatorKeys)
	p.txSubmit.SubmitMessageAsync(changeValSetMsg)

	incomingValUpdate := p.keeper.GetIncomingValidatorUpdates(ctx)
	for i, vUp := range incomingValUpdate {
		log.Debugf("incomingValUpdate[%d] pubkey = %s, power = %d\n", i, vUp.PubKey.String(), vUp.Power)
		// 100 is default power for validator
		if vUp.Power == 100 {
			p.validatorManager.UpdateNodeStatus(ctx, vUp.PubKey.GetEd25519(), types.NodeStatus_Validator)
			continue
		}

		p.validatorManager.UpdateNodeStatus(ctx, vUp.PubKey.GetEd25519(), types.NodeStatus_Candidate)
	}
	return p.keeper.GetIncomingValidatorUpdates(ctx)
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

	newNodes := p.validatorManager.GetPotentialCandidates(ctx, len(removedNodes))
	return newNodes, removedNodes, nil
}
