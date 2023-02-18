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
		data, err := h.doTransferFailure(ctx, msg.Data)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

func (h *HanlderTransferFailure) doTransferFailure(
	ctx sdk.Context,
	data *types.TransferFailure,
) ([]byte, error) {
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
			nonce := h.keeper.AddFailedTransfer(ctx, t.Id)
			log.Verbosef(
				"Removing failed transfer from the queue, transferId = %s, chain = %s, nonce = %d",
				t.Id, data.Chain, nonce)
		}
	}

	h.keeper.SetTransferQueue(ctx, data.Chain, newQ)

	return nil, nil
}
