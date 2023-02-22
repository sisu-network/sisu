package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HanlderTransferFailure struct {
	keeper keeper.Keeper
	pmm    components.PostedMessageManager
}

func NewHanlderTransferFailure(
	k keeper.Keeper,
	pmm components.PostedMessageManager,
) *HanlderTransferFailure {
	return &HanlderTransferFailure{
		keeper: k,
		pmm:    pmm,
	}
}

func (h *HanlderTransferFailure) DeliverMsg(
	ctx sdk.Context,
	msg *types.TransferFailureMsg,
) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, msg); process {
		err := h.doTransferFailure(ctx, msg.Data)
		h.keeper.ProcessTxRecord(ctx, hash)
		return &sdk.Result{}, err
	}

	return &sdk.Result{}, nil
}

func (h *HanlderTransferFailure) doTransferFailure(
	ctx sdk.Context,
	data *types.TransferFailure,
) error {
	ids := make(map[string]bool)
	for _, id := range data.Ids {
		ids[id] = true
	}

	// Remove all txout from transfer queue
	queue := h.keeper.GetTransferQueue(ctx, data.Chain)
	newQ := make([]*types.TransferDetails, 0)

	for _, t := range queue {
		if ids[t.Id] == false {
			newQ = append(newQ, t)
		} else {
			h.keeper.AddFailedTransfer(ctx, t.Id)
			retryNum := h.keeper.GetFailedTransferRetryNum(ctx, t.Id)
			log.Verbosef(
				"Removing failed transfer from queue, transferId = %s, chain = %s, retryNum = %d",
				t.Id, data.Chain, retryNum)
		}
	}

	h.keeper.SetTransferQueue(ctx, data.Chain, newQ)
	return nil
}
