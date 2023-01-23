package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/background"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxOutResult struct {
	pmm       components.PostedMessageManager
	keeper    keeper.Keeper
	transferQ background.TransferQueue
	privateDb keeper.PrivateDb
}

func NewHandlerTxOutResult(mc background.ManagerContainer) *HandlerTxOutResult {
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

	result := msg.Data
	txOut := h.keeper.GetTxOut(ctx, result.OutChain, result.OutHash)

	defer func(result *types.TxOutResult) {
		removeTxOut(ctx, h.privateDb, h.keeper, txOut)
	}(msg.Data)

	if txOut == nil {
		log.Errorf("Critical: cannot find txout from txOutConfirm message, chain = %s & hash = %s",
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

func removeTxOut(ctx sdk.Context, privateDb keeper.PrivateDb, k keeper.Keeper,
	txOut *types.TxOut) {
	// Remove the TxOut from the TxOutQueue
	q := k.GetTxOutQueue(ctx, txOut.Content.OutChain)
	if len(q) == 0 || q[0].GetId() != txOut.GetId() {
		// This is a rare case but it is possible to happen. The tx out is removed because it passes
		// expiration block. When this TxOut is confirmed, this queue does not have the txOut inside.
		if len(q) == 0 {
			log.Errorf("removeTransfer: The txout queue is empty")
		} else {
			log.Errorf("Id does not match. Id in the queue = %s, id in the message = %s", q[0].GetId(),
				txOut.GetId())
		}
		return
	}

	q = q[1:]
	k.SetTxOutQueue(ctx, txOut.Content.OutChain, q)

	// Unset the hold prcessing flag so that sisu can continue processing transfer/txout on this chain
	privateDb.SetHoldProcessing(types.TransferHoldKey, txOut.Content.OutChain, false)
	privateDb.SetHoldProcessing(types.TxOutHoldKey, txOut.Content.OutChain, false)

	// Remove the TxOut from the expired transactions.
	k.RemoveExpirationBlock(ctx, types.ExpirationBlock_TxOut, txOut.GetId())
}
