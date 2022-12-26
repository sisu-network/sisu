package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerAdjustEthNonce struct {
	pmm    PostedMessageManager
	keeper keeper.Keeper
}

func NewHandlerAdjustEthNonce(pmm PostedMessageManager, k keeper.Keeper) *HandlerAdjustEthNonce {
	return &HandlerAdjustEthNonce{
		pmm:    pmm,
		keeper: k,
	}
}

func (h *HandlerAdjustEthNonce) DeliverMsg(ctx sdk.Context, signerMsg *types.AdjustEthNonceMsg) (*sdk.Result, error) {
	data := signerMsg.Data
	h.keeper.SetSignerNonce(ctx, data.Chain, signerMsg.Signer, data.Nonce)

	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		h.updateNonce(ctx, signerMsg)
		h.keeper.ProcessTxRecord(ctx, hash)
	}

	return &sdk.Result{}, nil
}

func (h *HandlerAdjustEthNonce) updateNonce(ctx sdk.Context, signerMsg *types.AdjustEthNonceMsg) {
	h.keeper.SetMpcNonce(ctx, &types.MpcNonce{Nonce: signerMsg.Data.Nonce})
}
