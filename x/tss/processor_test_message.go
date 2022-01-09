package tss

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"
)

func (p *Processor) checkTxTestMessage(ctx sdk.Context, msg *types.TestMessage) error {
	log.Info("Checking test message ....")

	return nil
}

func (p *Processor) deliverTestMessage(ctx sdk.Context, msg *types.TestMessage) ([]byte, error) {
	log.Info("Delivering test message ....")

	return nil, nil
}
