package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerUpdateSolanaRecentHash struct {
	keeper keeper.Keeper
}

func NewHandlerUpdateSolanaRecentHash(keeper keeper.Keeper) *HandlerUpdateSolanaRecentHash {
	return &HandlerUpdateSolanaRecentHash{
		keeper: keeper,
	}
}

func (h *HandlerUpdateSolanaRecentHash) DeliverMsg(ctx sdk.Context, msg *types.UpdateSolanaRecentHashMsg) (*sdk.Result, error) {
	data := msg.Data
	h.keeper.SetSolanaConfirmedBlock(ctx, msg.Signer, data.Chain, data.Hash, data.Height)
	return nil, nil
}
