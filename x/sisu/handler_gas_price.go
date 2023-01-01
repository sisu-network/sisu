package sisu

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerGasPrice struct {
	keeper     keeper.Keeper
	globalData common.GlobalData
}

func NewHandlerGasPrice(mc ManagerContainer) *HandlerGasPrice {
	return &HandlerGasPrice{
		keeper:     mc.Keeper(),
		globalData: mc.GlobalData(),
	}
}

func (h *HandlerGasPrice) DeliverMsg(ctx sdk.Context, msg *types.GasPriceMsg) (*sdk.Result, error) {
	h.keeper.SetGasPrice(ctx, msg)

	params := h.keeper.GetParams(ctx)
	if params == nil {
		return nil, fmt.Errorf("Cannot find tss params")
	}

	h.keeper.SetGasPrice(ctx, msg)

	return &sdk.Result{}, nil
}
