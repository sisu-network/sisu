package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxIn struct {
	pmm       PostedMessageManager
	keeper    keeper.Keeper
	txInQueue TxInQueue
}

func NewHandlerTxIn(mc ManagerContainer) *HandlerTxIn {
	return &HandlerTxIn{
		keeper:    mc.Keeper(),
		pmm:       mc.PostedMessageManager(),
		txInQueue: mc.TxInQueue(),
	}
}

func (h *HandlerTxIn) DeliverMsg(ctx sdk.Context, signerMsg *types.TxInWithSigner) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doTxIn(ctx, signerMsg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

// Delivers observed Txs.
func (h *HandlerTxIn) doTxIn(ctx sdk.Context, msgWithSigner *types.TxInWithSigner) ([]byte, error) {
	msg := msgWithSigner.Data

	log.Info("Deliverying TxIn, hash = ", msg.TxHash, " on chain ", msg.Chain)

	// Save this to db.
	h.keeper.SaveTxIn(ctx, msg)

	// Add the message to the queue for later processing.
	h.txInQueue.AddTxIn(msg)

	return nil, nil
}
