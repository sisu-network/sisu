package sisu

import (
	"bytes"
	"encoding/base64"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (p *Processor) ChangeValidatorSet(ctx sdk.Context) {
	newValidators, removedValidators, err := p.getChangedNodes(ctx)
	if err != nil {
		return
	}

	if len(newValidators) == 0 && len(removedValidators) == 0 {
		return
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
	changeValSetMsg := types.NewChangeValidatorSetMsg(p.appKeys.GetSignerAddress().String(), oldValSet, newValSet, msgIndex)
	hash, _, err := keeper.GetTxRecordHash(changeValSetMsg)
	if err != nil {
		log.Error(err)
		return
	}

	if !p.keeper.IsTxRecordProcessed(ctx, hash) {
		p.txSubmit.SubmitMessageAsync(changeValSetMsg)
	}
}

func (p *Processor) GetPendingValidatorUpdates(ctx sdk.Context) []abci.ValidatorUpdate {
	incomingValUpdate := p.keeper.GetIncomingValidatorUpdates(ctx)
	log.Debug("len of incomingValUpdate is", len(incomingValUpdate))

	if len(incomingValUpdate) == 0 {
		return []abci.ValidatorUpdate{}
	}

	p.keeper.ClearValidatorUpdates(ctx)

	for i, vUp := range incomingValUpdate {
		log.Debugf("incomingValUpdate[%d] pubkey = %s, power = %d\n", i, base64.StdEncoding.EncodeToString(vUp.PubKey.GetEd25519()), vUp.Power)
		// 100 is default power for validator
		if vUp.Power == 100 {
			p.validatorManager.UpdateNodeStatus(ctx, vUp.PubKey.GetEd25519(), types.NodeStatus_Validator)
			continue
		}

		p.validatorManager.UpdateNodeStatus(ctx, vUp.PubKey.GetEd25519(), types.NodeStatus_Candidate)
	}

	return incomingValUpdate
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
