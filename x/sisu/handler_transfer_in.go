package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxIn struct {
	pmm    PostedMessageManager
	keeper keeper.Keeper
}

func NewHandlerTxIn(
	pmm PostedMessageManager,
	keeper keeper.Keeper,
) *HandlerTxIn {
	return &HandlerTxIn{
		pmm:    pmm,
		keeper: keeper,
	}
}

func (h *HandlerTxIn) DeliverMsg(ctx sdk.Context, signerMsg *types.TxInMsg) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doTxIn(ctx, signerMsg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

func (h *HandlerTxIn) doTxIn(ctx sdk.Context, signerMsg *types.TxInMsg) ([]byte, error) {
	return nil, nil
}
