package tss

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/tss/types"
)

func (p *Processor) checkKeysignResult(ctx sdk.Context, msg *types.KeysignResult) error {
	return nil
}

func (p *Processor) deliverKeysignResult(ctx sdk.Context, msg *types.KeysignResult) ([]byte, error) {
	// TODO: implements this to handle blame.

	return nil, nil
}
