package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxOutResult struct {
	pmm       PostedMessageManager
	keeper    keeper.Keeper
	transferQ TransferQueue
	privateDb keeper.PrivateDb
}

func NewHandlerTxOutResult(mc ManagerContainer) *HandlerTxOutResult {
	return &HandlerTxOutResult{
		keeper:    mc.Keeper(),
		pmm:       mc.PostedMessageManager(),
		transferQ: mc.TransferQueue(),
		privateDb: mc.PrivateDb(),
	}
}

func (h *HandlerTxOutResult) DeliverMsg(ctx sdk.Context, msg *types.TxOutResultMsg) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, msg); process {
		data, err := h.doTxOutResult(ctx, msg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

func (h *HandlerTxOutResult) doTxOutResult(ctx sdk.Context, msg *types.TxOutResultMsg) ([]byte, error) {
	log.Info("Delivering TxOutResult")

	defer func(result *types.TxOutResult) {
		// Remove the TxOut from the TxOutQueue
		q := h.keeper.GetTxOutQueue(ctx, msg.Data.OutChain)
		if q[0].GetId() != msg.Data.TxOutId {
			// Critical error. The TxOutId should be the same like in the Message
			log.Errorf("Id does not match. Id in the queue = %s, id in the message = %s", q[0].GetId(),
				msg.Data.TxOutId)
		} else {
			q = q[1:]
			h.keeper.SetTxOutQueue(ctx, msg.Data.OutChain, q)
		}

		// Reset the TransferHold & TxOutHold variable so that these 2 queues can continue processing
		// TxOut.
		h.privateDb.SetHoldProcessing(types.TransferHoldKey, result.OutChain, false)
		h.privateDb.SetHoldProcessing(types.TxOutHoldKey, result.OutChain, false)
	}(msg.Data)

	result := msg.Data
	txOut := h.keeper.GetTxOut(ctx, result.OutChain, result.OutHash)
	if txOut == nil {
		log.Errorf("cannot find txout from txOutConfirm message, chain = %s & hash = %s",
			result.OutChain, result.OutHash)
		return nil, nil
	}

	log.Verbose("msg.Result = ", result.Result)

	switch result.Result {
	case types.TxOutResultType_IN_BLOCK_SUCCESS:
		return h.doTxOutConfirm(ctx, result, txOut)
	default:
		return h.doTxOutFailure(ctx, result, txOut)
	}
}

func (h *HandlerTxOutResult) doTxOutConfirm(ctx sdk.Context, msg *types.TxOutResult, txOut *types.TxOut) ([]byte, error) {
	log.Verbose("Transaction is successfully included in a block, hash (no sig)= ", msg.OutHash, " chain = ", msg.OutChain)
	return nil, nil
}

func (h *HandlerTxOutResult) doTxOutFailure(ctx sdk.Context, msg *types.TxOutResult, txOut *types.TxOut) ([]byte, error) {
	log.Warn("Transaction failed!, txOut.TxType = ", txOut.TxType)

	// TODO: Add TxOut and its transfer to the failure queue.

	return nil, nil
}
