package sisu

import (
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func (p *Processor) deliverMsgPauseGw(msg *types.MsgPauseGw) ([]byte, error) {
	currentPauseGwRecord := p.publicDb.GetPauseGwRecord(msg.Chain)

	// check duplicate, if this validator processed this msg then ignore
	if currentPauseGwRecord.HasSigned(msg.Signer) {
		return nil, nil
	}

	p.publicDb.SavePauseGwMsg(msg)

	// Check reach consensus
	savedRecord := p.publicDb.GetPauseGwRecord(msg.Chain)
	totalValidator := len(p.globalData.GetValidatorSet())
	if savedRecord == nil || !savedRecord.ReachConsensus(totalValidator) {
		return nil, nil
	}

	txOut := p.txOutputProducer.GetPauseGwTxOut(msg.Chain)
	go func() {
		if err := p.txSubmit.SubmitMessageAsync(txOut); err != nil {
			log.Error("error when submit async tx out: ", err)
		}

		log.Debug("Submitted tx out of pausing gateway successfully")
	}()
	return nil, nil
}
