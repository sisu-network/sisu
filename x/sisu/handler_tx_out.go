package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxOut struct {
	pmm        PostedMessageManager
	keeper     keeper.Keeper
	txOutQueue TxOutQueue
}

func NewHandlerTxOut(mc ManagerContainer) *HandlerTxOut {
	return &HandlerTxOut{
		keeper:     mc.Keeper(),
		pmm:        mc.PostedMessageManager(),
		txOutQueue: mc.TxOutQueue(),
	}
}

func (h *HandlerTxOut) DeliverMsg(ctx sdk.Context, signerMsg *types.TxOutWithSigner) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doTxOut(ctx, signerMsg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

// deliverTxOut executes a TxOut transaction after it's included in Sisu block. If this node is
// catching up with the network, we would not send the tx to TSS for signing.
func (h *HandlerTxOut) doTxOut(ctx sdk.Context, msgWithSigner *types.TxOutWithSigner) ([]byte, error) {
	txOut := msgWithSigner.Data

	log.Info("Delivering TxOut")

	// Save this to KVStore
	h.keeper.SaveTxOut(ctx, txOut)

	h.txOutQueue.AddTxOut(txOut)

	return nil, nil
}
