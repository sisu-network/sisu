package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxIn struct {
	pmm         PostedMessageManager
	keeper      keeper.Keeper
	valsManager ValidatorManager
	globalData  common.GlobalData
}

func NewHandlerTxIn(
	pmm PostedMessageManager,
	keeper keeper.Keeper,
	valsManager ValidatorManager,
	globalData common.GlobalData,
) *HandlerTxIn {
	return &HandlerTxIn{
		pmm:         pmm,
		keeper:      keeper,
		valsManager: valsManager,
		globalData:  globalData,
	}
}

func (h *HandlerTxIn) DeliverMsg(ctx sdk.Context, msg *types.TxInMsg) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, msg); process {
		data, err := h.doTxIn(ctx, msg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

func (h *HandlerTxIn) doTxIn(ctx sdk.Context, msg *types.TxInMsg) ([]byte, error) {
	// Check if we have the TxIn details .
	txInDetails := h.keeper.GetTxInDetails(ctx, msg.Data.Id)
	if txInDetails != nil {
		// 1. TODO: Do verificaiton on the tx in details to make sure this data is correct (including
		// the transfers)
		// 2. Add all the new transfers to the transfer queue.
		for _, transfer := range txInDetails.Data.Transfers {
			// TODO: Optimize this path. We can save single transfer instead of the entire queue.
			queue := h.keeper.GetTransferQueue(ctx, transfer.ToChain)
			queue = append(queue, transfer)
			h.keeper.SetTransferQueue(ctx, transfer.ToChain, queue)
		}
	}

	return nil, nil
}
