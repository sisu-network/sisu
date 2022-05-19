package sisu

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (p *Processor) EndBlockValidator(ctx sdk.Context) []abci.ValidatorUpdate {
	newValidators, removedValidators, err := p.getChangedNodes(ctx)
	if err != nil {
		return []abci.ValidatorUpdate{}
	}

	if len(newValidators) == 0 && len(removedValidators) == 0 {
		return []abci.ValidatorUpdate{}
	}

	log.Debug("len = ", len(newValidators), " newValidators = ", newValidators)
	log.Debug("len = ", len(removedValidators), " removedValidators = ", removedValidators)

	newValidatorKeys := make([][]byte, 0)
	for _, val := range newValidators {
		newValidatorKeys = append(newValidatorKeys, val.ConsensusKey.GetBytes())
	}

	// 1. newValSet = current validator set - removedValidators + newValidators
	newValSet := make([][]byte, 0)
	currentVals := p.validatorManager.GetNodesByStatus(types.NodeStatus_Validator)

	// 1a. excludes removed validators
	for _, val := range currentVals {
		// if this validator is in removed validators list, ignore
		isOld := false
		for _, old := range removedValidators {
			if bytes.Equal(val.ConsensusKey.GetBytes(), old.ConsensusKey.GetBytes()) {
				isOld = true
				break
			}
		}

		if isOld {
			continue
		}

		newValSet = append(newValSet, val.ConsensusKey.GetBytes())
	}

	// 1b. includes new validators
	newValSet = append(newValSet, newValidatorKeys...)

	oldValSet := make([][]byte, 0)
	for _, oldVal := range currentVals {
		oldValSet = append(oldValSet, oldVal.ConsensusKey.GetBytes())
	}

	log.Debug("len newValSet = ", len(newValidatorKeys))
	msgIndex := p.keeper.GetValidatorUpdateIndex(ctx)
	log.Debug("msg index = ", msgIndex)
	changeValSetMsg := types.NewChangeValidatorSetMsg(p.appKeys.GetSignerAddress().String(), oldValSet, newValSet, int32(msgIndex))
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
