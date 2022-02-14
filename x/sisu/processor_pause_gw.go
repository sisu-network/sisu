package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func (p *Processor) deliverMsgPauseGw(ctx sdk.Context, msg *types.MsgPauseGw) ([]byte, error) {
	currentPauseGwRecord := p.publicDb.GetPauseGwRecord(msg.Chain, msg.Address)

	// check duplicate, if this validator processed this msg then ignore
	if currentPauseGwRecord.HasSigned(msg.Signer) {
		return nil, nil
	}

	p.publicDb.SavePauseGwMsg(msg)

	// Check reach consensus
	savedRecord := p.publicDb.GetPauseGwRecord(msg.Chain, msg.Address)
	totalValidator := len(p.globalData.GetValidatorSet())
	if savedRecord == nil || !savedRecord.ReachConsensus(totalValidator) {
		return nil, nil
	}
	return nil, nil
}
