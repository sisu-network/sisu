package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTokenPrice struct {
	publicDb keeper.Storage
}

func NewHandlerTokenPrice(mc ManagerContainer) *HandlerTokenPrice {
	return &HandlerTokenPrice{
		publicDb: mc.PublicDb(),
	}
}

func (h *HandlerTokenPrice) DeliverMsg(ctx sdk.Context, msg *types.UpdateTokenPrice) (*sdk.Result, error) {
	h.publicDb.SetTokenPrices(uint64(ctx.BlockHeight()), msg)

	return nil, nil
}
