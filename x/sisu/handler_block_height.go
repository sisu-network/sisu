package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerBlockHeight struct {
	keeper keeper.Keeper
}

func NewHandlerBlockHeight(keeper keeper.Keeper) *HandlerBlockHeight {
	return &HandlerBlockHeight{
		keeper: keeper,
	}
}

func (h *HandlerBlockHeight) DeliverMsg(ctx sdk.Context, signerMsg *types.BlockHeightMsg) (*sdk.Result, error) {
	record := h.keeper.GetBlockHeightRecord(ctx, signerMsg.Signer)
	if record == nil {
		record = &types.BlockHeightRecord{
			BlockHeights: make([]*types.BlockHeight, 0),
		}
	}

	blockHeight := signerMsg.Data
	found := false
	for j, recordHeight := range record.BlockHeights {
		if blockHeight.Chain == recordHeight.Chain {
			record.BlockHeights[j] = blockHeight
			found = true
			break
		}
	}

	if !found {
		record.BlockHeights = append(record.BlockHeights, blockHeight)
	}

	h.keeper.SaveBlockHeights(ctx, signerMsg.Signer, record)

	return &sdk.Result{}, nil
}
