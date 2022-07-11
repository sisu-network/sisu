package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTransferBatch struct {
	mc     ManagerContainer
	keeper keeper.Keeper
}

func NewHandlerTransferBatch() *HandlerTransferBatch {
	return &HandlerTransferBatch{}
}

func (h *HandlerTransferBatch) DeliverMsg(ctx sdk.Context, msg *types.TransferBatchMsg) (*sdk.Result, error) {
	pmm := h.mc.PostedMessageManager()
	if process, hash := pmm.ShouldProcessMsg(ctx, msg); process {
		data, err := h.doBatchTransfer(ctx, msg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

func (h *HandlerTransferBatch) doBatchTransfer(ctx sdk.Context, msg *types.TransferBatchMsg) ([]byte, error) {
	data := msg.Data
	pendings := h.keeper.GetPendingTransfers(ctx, data.Chain)
	if len(pendings) > 0 {
		return nil, nil
	}

	h.keeper.SetPendingTransfers(ctx, data.String(), data.Transfers)

	return nil, nil
}
