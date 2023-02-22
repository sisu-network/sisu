package sisu

import (
	"fmt"

	"github.com/sisu-network/sisu/x/sisu/components"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

const (
	VoteKey = "TxOutVote"
)

type HandlerTxOutVote struct {
	pmm       components.PostedMessageManager
	keeper    keeper.Keeper
	privateDb keeper.PrivateDb
}

func NewHandlerTxOutConsensed(
	pmm components.PostedMessageManager,
	keeper keeper.Keeper,
	privateDb keeper.PrivateDb,
) *HandlerTxOutVote {
	return &HandlerTxOutVote{
		pmm:       pmm,
		keeper:    keeper,
		privateDb: privateDb,
	}
}

func (h *HandlerTxOutVote) DeliverMsg(
	ctx sdk.Context,
	msg *types.TxOutVoteMsg,
) (*sdk.Result, error) {
	txOut := h.keeper.GetProposedTxOut(ctx, msg.Data.TxOutId, msg.Data.AssignedValidator)
	if txOut == nil {
		log.Error("Cannot get proposed txout, txOutId = ", msg.Data.TxOutId)
		return &sdk.Result{}, nil
	}

	done := h.keeper.IsTxRecordProcessed(ctx,
		[]byte(fmt.Sprintf("%s__%s", txOut.GetId(), msg.Data.AssignedValidator)))
	if done {
		log.Info("Ignore the completed proposed txout, txOutId = ", msg.Data.TxOutId)
		return &sdk.Result{}, nil
	}

	counter := h.keeper.GetTransferCounter(ctx, txOut.Input.TransferIds[0])
	prefix := fmt.Sprintf("%s__%s__%d", VoteKey, msg.Data.TxOutId, counter)
	h.keeper.AddVoteResult(ctx, prefix, msg.Signer, msg.Data.Vote)

	h.checkVoteResult(ctx, txOut, counter, msg.Data.AssignedValidator)

	return &sdk.Result{}, nil
}

func (h *HandlerTxOutVote) checkVoteResult(
	ctx sdk.Context,
	txOut *types.TxOut,
	counter int,
	assignedValidator string,
) {
	txOutId := txOut.GetId()

	prefix := fmt.Sprintf("%s__%s__%d", VoteKey, txOutId, counter)
	results := h.keeper.GetVoteResults(ctx, prefix)
	params := h.keeper.GetParams(ctx)
	if params == nil {
		log.Warn("tssParams is nil")
		return
	}

	approveCount := 0
	rejectCount := 0
	for _, result := range results {
		if result == types.VoteResult_APPROVE {
			approveCount++
		} else {
			rejectCount++
		}
	}

	threshold := int(params.MajorityThreshold)

	if approveCount < threshold && rejectCount < threshold {
		// TODO: handler the case or do timeout in the module.go
		return
	}

	h.keeper.ProcessTxRecord(ctx, []byte(fmt.Sprintf("%s__%s", txOut.GetId(), assignedValidator)))

	if approveCount >= threshold {
		finalizedTxOut := h.keeper.GetFinalizedTxOut(ctx, txOutId)
		if finalizedTxOut == nil {
			doTxOut(ctx, h.keeper, h.privateDb, txOut)
		} else {
			log.Verbosef("Finalized TxOut has been processed for txOut with id %s", txOutId)
		}
	} else {
		log.Verbose("TxOut is rejected, txOutId = ", txOutId)
		h.keeper.IncTransferCounter(ctx, txOut.Input.TransferIds[0])
		h.privateDb.SetHoldProcessing(types.TransferHoldKey, txOut.Content.OutChain, false)
	}
}
