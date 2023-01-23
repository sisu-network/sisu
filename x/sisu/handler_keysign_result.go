package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/background"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerKeysignResult struct {
}

func NewHandlerKeysignResult(mc background.ManagerContainer) *HandlerKeysignResult {
	return &HandlerKeysignResult{}
}

func (h *HandlerKeysignResult) DeliverMsg(ctx sdk.Context, signerMsg *types.KeysignResult) (*sdk.Result, error) {
	// TODO: Implement this.
	return &sdk.Result{}, nil
}
