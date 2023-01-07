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

func (h *HandlerTxIn) DeliverMsg(ctx sdk.Context, signerMsg *types.TxInMsg) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doTxIn(ctx, signerMsg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

func (h *HandlerTxIn) doTxIn(ctx sdk.Context, signerMsg *types.TxInMsg) ([]byte, error) {
	h.globalData.ConfirmTxIn(signerMsg.Data)

	return nil, nil
}
