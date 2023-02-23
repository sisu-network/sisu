package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/chains"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTransferRetry struct {
	pmm           components.PostedMessageManager
	keeper        keeper.Keeper
	globalData    components.GlobalData
	bridgeManager chains.BridgeManager
	valsManager   components.ValidatorManager
	privateDb     keeper.PrivateDb
}

func NewHandlerTransferRetry(
	pmm components.PostedMessageManager,
	keeper keeper.Keeper,
	globalData components.GlobalData,
	bridgeManager chains.BridgeManager,
	valsManager components.ValidatorManager,
	privateDb keeper.PrivateDb,
) *HandlerTransferRetry {
	return &HandlerTransferRetry{
		pmm:           pmm,
		keeper:        keeper,
		globalData:    globalData,
		bridgeManager: bridgeManager,
		valsManager:   valsManager,
		privateDb:     privateDb,
	}
}

func (h *HandlerTransferRetry) DeliverMsg(
	ctx sdk.Context,
	msg *types.TransferRetryMsg,
) (*sdk.Result, error) {
	if shouldProcess, hash := h.pmm.ShouldProcessMsg(ctx, msg); shouldProcess {
		h.doRetryTransfer(ctx, msg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{}, nil
	}

	return &sdk.Result{}, nil
}

func (h *HandlerTransferRetry) doRetryTransfer(ctx sdk.Context, msg *types.TransferRetryMsg) {
	retryNum := h.keeper.GetFailedTransferRetryNum(ctx, msg.Data.TransferId)
	if retryNum == -1 {
		log.Errorf("The transfer isn't failed, transferId = %s", msg.Data.TransferId)
		return
	}

	if retryNum != msg.Data.Index {
		log.Errorf(
			"Mismatch between retry number of failed transfer and index of retried transfer, "+
				"transferId = %s, retryNum = %d, index = %d",
			msg.Data.TransferId, retryNum, msg.Data.Index)
		return
	}

	transfer := h.keeper.GetTransfer(ctx, msg.Data.TransferId)
	if transfer == nil {
		log.Error("Cannot get the transfer, transferId = ", msg.Data.TransferId)
		return
	}

	h.saveTransfers(ctx, transfer)
}

func (h *HandlerTransferRetry) saveTransfers(ctx sdk.Context, transfer *types.TransferDetails) {
	queue := h.keeper.GetTransferQueue(ctx, transfer.ToChain)
	queue = append(queue, transfer)

	h.keeper.SetTransferCounter(ctx, transfer.Id, 0)
	h.keeper.SetTransferQueue(ctx, transfer.ToChain, queue)
}
