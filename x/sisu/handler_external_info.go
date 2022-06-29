package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerExternalInfo struct {
	keeper keeper.Keeper
}

func NewHandlerExternalInfo(mc ManagerContainer) *HandlerExternalInfo {
	return &HandlerExternalInfo{
		keeper: mc.Keeper(),
	}
}

func (h *HandlerExternalInfo) DeliverMsg(ctx sdk.Context, sdkMsg *types.ExternalInfoMsg) (*sdk.Result, error) {
	msg := sdkMsg.Data

	if msg.GasPrice != nil {
		h.updateGasPrice(ctx, sdkMsg.Signer, msg.GasPrice)
	}

	if msg.BlockHeights != nil {
		h.updateBlockHeight(ctx, sdkMsg.Signer, msg.BlockHeights)
	}

	return &sdk.Result{}, nil
}

func (h *HandlerExternalInfo) updateGasPrice(ctx sdk.Context, signer string, gasPrice *types.GasPrice) {
	// TODO: Move the logic from handler_gas_price to here.
}

func (h *HandlerExternalInfo) updateBlockHeight(ctx sdk.Context, signer string, blocks []*types.BlockHeight) {
	record := h.keeper.GetBlockHeightRecord(ctx, signer)
	if record == nil {
		record = &types.BlockHeightRecord{
			BlockHeights: make([]*types.BlockHeight, 0),
		}
	}

	for _, blockHeight := range blocks {
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
	}

	h.keeper.SaveBlockHeights(ctx, signer, record)
}
