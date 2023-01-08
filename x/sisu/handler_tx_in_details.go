package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxInDetails struct {
	pmm        PostedMessageManager
	keeper     keeper.Keeper
	globalData common.GlobalData
}

func NewHandlerTxInDetails(
	pmm PostedMessageManager,
	keeper keeper.Keeper,
	globalData common.GlobalData,
) *HandlerTxInDetails {
	return &HandlerTxInDetails{
		pmm:        pmm,
		keeper:     keeper,
		globalData: globalData,
	}
}

func (h *HandlerTxInDetails) DeliverMsg(ctx sdk.Context, msg *types.TxInDetailsMsg) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, msg); process {
		h.keeper.ProcessTxRecord(ctx, hash)

		confirmedTxIn := h.keeper.GetConfirmedTxIn(ctx, msg.Data.TxIn.Id)
		if confirmedTxIn == nil {
			// We have the TxIn details, make this TxIn as confirmed.
			h.keeper.SetConfirmedTxIn(ctx, &types.ConfirmedTxIn{
				TxInId: msg.Data.TxIn.Id,
				Signer: msg.Signer,
			})

			h.globalData.ConfirmTxIn(msg.Data.TxIn)
		}
	}

	return &sdk.Result{}, nil
}
