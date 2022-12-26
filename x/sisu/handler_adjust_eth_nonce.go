package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerAdjustEthNonce struct {
	pmm       PostedMessageManager
	keeper    keeper.Keeper
	privateDb keeper.PrivateDb
}

func NewHandlerAdjustEthNonce(pmm PostedMessageManager, k keeper.Keeper,
	privateDb keeper.PrivateDb) *HandlerAdjustEthNonce {
	return &HandlerAdjustEthNonce{
		pmm:       pmm,
		keeper:    k,
		privateDb: privateDb,
	}
}

func (h *HandlerAdjustEthNonce) DeliverMsg(ctx sdk.Context, msg *types.AdjustEthNonceMsg) (*sdk.Result, error) {
	data := msg.Data
	h.keeper.SetSignerNonce(ctx, data.Chain, msg.Signer, uint64(data.Nonce))

	if process, hash := h.pmm.ShouldProcessMsg(ctx, msg); process {
		h.updateNonce(ctx, msg)
		h.keeper.ProcessTxRecord(ctx, hash)
	}

	return &sdk.Result{}, nil
}

func (h *HandlerAdjustEthNonce) updateNonce(ctx sdk.Context, msg *types.AdjustEthNonceMsg) {
	log.Infof("Update mpc nonce for chain %s, new nonce = %d", msg.Data.Chain, msg.Data.Nonce)
	h.keeper.SetMpcNonce(ctx, &types.MpcNonce{Chain: msg.Data.Chain, Nonce: msg.Data.Nonce})

	// Update the nonce index in the private db.
	key := keeper.GetEthNonceKey(msg.Data.Chain)
	index := h.privateDb.GetTxHashIndex(key)
	h.privateDb.SetTxHashIndex(key, index+1)
}
