package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTokenPrice struct {
	keeper keeper.Keeper
}

func NewHandlerTokenPrice(mc ManagerContainer) *HandlerTokenPrice {
	return &HandlerTokenPrice{
		keeper: mc.Keeper(),
	}
}

func (h *HandlerTokenPrice) DeliverMsg(ctx sdk.Context, msg *types.UpdateTokenPrice) (*sdk.Result, error) {
	h.keeper.SetTokenPrices(ctx, uint64(ctx.BlockHeight()), msg)

	return &sdk.Result{}, nil
}
