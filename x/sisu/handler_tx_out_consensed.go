package sisu

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxOutConsensed struct {
	pmm    PostedMessageManager
	keeper keeper.Keeper
}

func NewHandlerTxOutConfirm(
	pmm PostedMessageManager,
	keeper keeper.Keeper,
) *HandlerTxOutConsensed {
	return &HandlerTxOutConsensed{
		pmm:    pmm,
		keeper: keeper,
	}
}

func (h *HandlerTxOutConsensed) DeliverMsg(ctx sdk.Context, msg *types.TxOutConsensedMsg) (*sdk.Result, error) {
	prefix := fmt.Sprintf("tx_out_consensed__%s", msg.Data.TxOutId)
	h.keeper.AddVoteResult(ctx, prefix, msg.Signer, &msg.Data.Vote)

	return &sdk.Result{}, nil
}

func (h *HandlerTxOutConsensed) checkVoteResult(ctx sdk.Context, txOutId string) {
	results := h.keeper.GetVoteResults(ctx, txOutId)
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

	if count >= int(tssParams.MajorityThreshold) {
		index := strings.Index(txOutId, "__")
		if index < 0 {
			log.Errorf("Invalid txout id ", txOutId)
			return
		}

		chain := txOutId[:index]
		hash := txOutId[index+2:]

		txOut := h.keeper.GetTxOut(ctx, chain, hash)

		doTxOut(ctx, h.keeper, txOut)
	}
}
