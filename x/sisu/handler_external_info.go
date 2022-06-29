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
		h.updateGasPrice(sdkMsg.Signer, msg.GasPrice)
	}

	if msg.BlockHeights != nil {
		h.updateBlockHeight(sdkMsg.Signer, msg.BlockHeights)
	}

	return &sdk.Result{}, nil
}

func (h *HandlerExternalInfo) updateGasPrice(signer string, gasPrice *types.GasPrice) {
	// TODO: Move the logic from handler_gas_price to here.
}

func (h *HandlerExternalInfo) updateBlockHeight(signer string, blocks []*types.BlockHeight) {

}
