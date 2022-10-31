package sisu

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HanlderTransferFailure struct {
	keeper keeper.Keeper
	pmm    PostedMessageManager
}

func NewHanlderTransferFailure(k keeper.Keeper, pmm PostedMessageManager) *HanlderTransferFailure {
	return &HanlderTransferFailure{
		keeper: k,
		pmm:    pmm,
	}
}

func (h *HanlderTransferFailure) DeliverMsg(ctx sdk.Context, msg *types.TransferFailureMsg) (*sdk.Result, error) {
	fmt.Println("BBBBBBBB DeliverMsg")

	if process, hash := h.pmm.ShouldProcessMsg(ctx, msg); process {
		data, err := h.doTransferFailure(ctx, msg.Data)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

func (h *HanlderTransferFailure) doTransferFailure(ctx sdk.Context, data *types.TransferFailure) ([]byte, error) {
	log.Verbose("Removing failed transfer from the queue, ids = ", data.Ids, ", chain = ", data.Chain)

	ids := make(map[string]bool)
	for _, id := range data.Ids {
		ids[id] = true
	}

	// Remove all txout from transfer queue
	queue := h.keeper.GetTransferQueue(ctx, data.Chain)
	newQ := make([]*types.Transfer, 0)

	for _, t := range queue {
		if ids[t.Id] == false {
			newQ = append(newQ, t)
		}
	}

	h.keeper.SetTransferQueue(ctx, data.Chain, newQ)

	return nil, nil
}
