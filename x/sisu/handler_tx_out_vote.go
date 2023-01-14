package sisu

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

const (
	VoteKey = "TxOutVote"
)

type HandlerTxOutVote struct {
	pmm    PostedMessageManager
	keeper keeper.Keeper
}

func NewHandlerTxOutConsensed(
	pmm PostedMessageManager,
	keeper keeper.Keeper,
) *HandlerTxOutVote {
	return &HandlerTxOutVote{
		pmm:    pmm,
		keeper: keeper,
	}
}

func (h *HandlerTxOutVote) DeliverMsg(ctx sdk.Context, msg *types.TxOutVoteMsg) (*sdk.Result, error) {
	fmt.Println("HandlerTxOutConsensed, signer = ", msg.Signer)
	prefix := fmt.Sprintf("%s__%s", VoteKey, msg.Data.TxOutId)
	h.keeper.AddVoteResult(ctx, prefix, msg.Signer, msg.Data.Vote)

	h.checkVoteResult(ctx, msg.Data.TxOutId, msg.Signer)

	return &sdk.Result{}, nil
}

func (h *HandlerTxOutVote) checkVoteResult(ctx sdk.Context, txOutId, signer string) {
	results := h.keeper.GetVoteResults(ctx, VoteKey)
	tssParams := h.keeper.GetParams(ctx)
	if tssParams == nil {
		log.Warn("tssParams is nil")
		return
	}

	count := 0
	for _, result := range results {
		if result == types.VoteResult_APPROVE {
			count++
		}
	}

	fmt.Println("checkVoteResult, count = ", count)

	if count >= int(tssParams.MajorityThreshold) {
		txOut := h.keeper.GetProposedTxOut(ctx, txOutId, signer)
		fmt.Println("txOut = ", txOut)

		doTxOut(ctx, h.keeper, txOut)
	} else {
		// TODO: handler the case or do timeout in the module.go
	}
}
