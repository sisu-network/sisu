package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxInDetails struct {
	pmm    PostedMessageManager
	keeper keeper.Keeper
}

func NewHandlerTxInDetails(
	pmm PostedMessageManager,
	keeper keeper.Keeper,
) *HandlerTxInDetails {
	return &HandlerTxInDetails{
		pmm:    pmm,
		keeper: keeper,
	}
}

func (h *HandlerTxInDetails) DeliverMsg(ctx sdk.Context, signerMsg *types.TxInDetailsMsg) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doTxIn(ctx, signerMsg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

func (h *HandlerTxInDetails) doTxIn(ctx sdk.Context, signerMsg *types.TxInDetailsMsg) ([]byte, error) {
	return nil, nil
}
