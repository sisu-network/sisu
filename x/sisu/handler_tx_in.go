package sisu

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxIn struct {
	pmm        PostedMessageManager
	keeper     keeper.Keeper
	globalData common.GlobalData
	txInQueue  TxInQueue
}

func NewHandlerTxIn(mc ManagerContainer) *HandlerTxIn {
	return &HandlerTxIn{
		keeper:     mc.Keeper(),
		pmm:        mc.PostedMessageManager(),
		txInQueue:  mc.TxInQueue(),
		globalData: mc.GlobalData(),
	}
}

func (h *HandlerTxIn) DeliverMsg(ctx sdk.Context, signerMsg *types.TxsInMsg) (*sdk.Result, error) {
	fmt.Println("AAAAAA DeliverMsg, data = ", *signerMsg.Data)
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doTxIn(ctx, signerMsg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

// Delivers observed Txs.
func (h *HandlerTxIn) doTxIn(ctx sdk.Context, signerMsg *types.TxsInMsg) ([]byte, error) {
	msg := signerMsg.Data

	log.Info("Deliverying TxIn on chain ", msg.Chain)

	if !h.globalData.IsCatchingUp() {
		// Add the message to the queue for later processing.
		h.txInQueue.AddTxIn(ctx, msg)
	}

	return nil, nil
}
