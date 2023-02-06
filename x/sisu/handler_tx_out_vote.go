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

	count := 0
	for _, result := range results {
		if result == types.VoteResult_APPROVE {
			count++
		}
	}

	if count >= int(params.MajorityThreshold) {
		txOut := h.keeper.GetProposedTxOut(ctx, txOutId, assignedVal)
		if txOut == nil {
			log.Errorf("checkVoteResult: TxOut is nil, txOutId = %s", txOutId)
		} else {
			finalizedTxOut := h.keeper.GetFinalizedTxOut(ctx, txOutId)
			if finalizedTxOut == nil {
				doTxOut(ctx, h.keeper, h.privateDb, txOut)
			} else {
				log.Verbosef("Finalized TxOut has been processed for txOut with id %s", txOutId)
			}
		}
	} else {
		// TODO: handler the case or do timeout in the module.go
	}
}
