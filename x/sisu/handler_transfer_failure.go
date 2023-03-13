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
		doTransferFailure(h.keeper, ctx, msg.Data.Chain, msg.Data.TransferRetryIds)
		h.keeper.ProcessTxRecord(ctx, hash)
	}

	return &sdk.Result{}, nil
}

func doTransferFailure(keeper keeper.Keeper, ctx sdk.Context, chain string, transferRetryIds []string) {
	ids := make(map[string]bool)
	transferIds := types.GetIdsFromRetryIds(transferRetryIds)
	for _, id := range transferIds {
		ids[id] = true
	}

	// Remove all txout from transfer queue
	queue := keeper.GetTransferQueue(ctx, chain)
	newQ := make([]*types.TransferDetails, 0)

	for _, t := range queue {
		if !ids[t.Id] {
			newQ = append(newQ, t)
		} else {
			newTransfer := keeper.IncTransferRetryNum(ctx, t.Id)
			keeper.AddFailedTransfer(ctx, newTransfer.Id)
			log.Verbosef("Removing failed transfer from queue, transferRetryId = %s, chain = %s",
				newTransfer.GetRetryId(), chain)
		}
	}

	keeper.SetTransferQueue(ctx, chain, newQ)
}
