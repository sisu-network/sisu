package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerSetDheartIPAddress struct {
	keeper keeper.Keeper
}

func NewHandlerSetDheartIPAddress(mc ManagerContainer) *HandlerSetDheartIPAddress {
	return &HandlerSetDheartIPAddress{
		keeper: mc.Keeper(),
	}
}

func (h *HandlerSetDheartIPAddress) DeliverMsg(ctx sdk.Context, msg *types.SetDheartIpAddressMsg) (*sdk.Result, error) {
	if err := h.keeper.SaveDheartIPAddress(ctx, msg.GetSender(), msg.Data.IPAddress); err != nil {
		return &sdk.Result{}, err
	}

	return &sdk.Result{}, nil
}
