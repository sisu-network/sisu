package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxOut struct {
	pmm    PostedMessageManager
	keeper keeper.Keeper
}

func NewHandlerTxOut(mc ManagerContainer) *HandlerTxOut {
	return &HandlerTxOut{
		keeper: mc.Keeper(),
		pmm:    mc.PostedMessageManager(),
	}
}

func (h *HandlerTxOut) DeliverMsg(ctx sdk.Context, signerMsg *types.TxOutMsg) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doTxOut(ctx, signerMsg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

// deliverTxOut executes a TxOut transaction after it's included in Sisu block. If this node is
// catching up with the network, we would not send the tx to TSS for signing.
func (h *HandlerTxOut) doTxOut(ctx sdk.Context, txOutMsg *types.TxOutMsg) ([]byte, error) {
	txOut := txOutMsg.Data

	log.Info("Delivering TxOut")

	// Save this to KVStore
	h.keeper.SaveTxOut(ctx, txOut)

	// If this is a txOut deployment, mark the contract as being deployed.
	switch txOut.TxType {
	case types.TxOutType_TRANSFER_OUT:
		h.handlerTransfer(ctx, txOut)
	}

	return nil, nil
}

func (h *HandlerTxOut) handlerTransfer(ctx sdk.Context, txOut *types.TxOut) {
	// 1. Update TxOut queue
	h.addTxOutToQueue(ctx, txOut)

	// 2. Remove the transfers in txOut from the queue
	queue := h.keeper.GetTransferQueue(ctx, txOut.Content.OutChain)
	ids := make(map[string]bool, 0)
	for _, inHash := range txOut.Content.InHashes {
		ids[inHash] = true
	}

	newQueue := make([]*types.Transfer, 0)
	for _, transfer := range queue {
		if !ids[transfer.Id] {
			newQueue = append(newQueue, transfer)
		}
	}

	h.keeper.SetTransferQueue(ctx, txOut.Content.OutChain, newQueue)
}

func (h *HandlerTxOut) addTxOutToQueue(ctx sdk.Context, txOut *types.TxOut) {
	// Move the the transfers associated with this tx_out to pending.
	queue := h.keeper.GetTxOutQueue(ctx, txOut.Content.OutChain)
	queue = append(queue, txOut)
	h.keeper.SetTxOutQueue(ctx, txOut.Content.OutChain, queue)
}
