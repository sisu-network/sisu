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

func (h *HandlerTxOutVote) DeliverMsg(ctx sdk.Context, msg *types.TxOutVoteMsg) (*sdk.Result, error) {
	prefix := fmt.Sprintf("%s__%s", VoteKey, msg.Data.TxOutId)
	h.keeper.AddVoteResult(ctx, prefix, msg.Signer, msg.Data.Vote)

	h.checkVoteResult(ctx, msg.Data.TxOutId, msg.Data.AssignedValidator)

	return &sdk.Result{}, nil
}

func (h *HandlerTxOutVote) checkVoteResult(ctx sdk.Context, txOutId, assignedVal string) {
	prefix := fmt.Sprintf("%s__%s", VoteKey, txOutId)
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

	txOut := h.keeper.GetProposedTxOut(ctx, txOutId, assignedVal)
	if txOut == nil {
		log.Errorf("checkVoteResult: TxOut is nil, txOutId = %s", txOutId)
		return
	}

	if approveCount >= threshold {
		finalizedTxOut := h.keeper.GetFinalizedTxOut(ctx, txOutId)
		if finalizedTxOut == nil {
			doTxOut(ctx, h.keeper, h.privateDb, txOut)
		} else {
			log.Verbosef("Finalized TxOut has been processed for txOut with id %s", txOutId)
		}
	} else {
		log.Verbose("TxOut is rejected, txOutId = ", txOutId)
		transfer := h.keeper.GetTransfer(ctx, txOut.Input.TransferIds[0])
		if transfer == nil {
			log.Error("Invalid transfer, id = ", txOut.Input.TransferIds[0])
		} else {
			h.keeper.RemoveVoteResults(ctx, prefix)
			h.keeper.IncTransferCounter(ctx, transfer.Id)
			h.privateDb.SetHoldProcessing(types.TransferHoldKey, transfer.ToChain, false)
		}
	}
}
