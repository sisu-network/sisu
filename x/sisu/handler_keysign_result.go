package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/background"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerKeysignResult struct {
	pmm        components.PostedMessageManager
	keeper     keeper.Keeper
	background background.Background
}

func NewHandlerKeysignResult(mc background.ManagerContainer) *HandlerKeysignResult {
	return &HandlerKeysignResult{
		pmm:        mc.PostedMessageManager(),
		keeper:     mc.Keeper(),
		background: mc.Background(),
	}
}

func (h *HandlerKeysignResult) DeliverMsg(ctx sdk.Context, msg *types.KeysignResultMsg) (*sdk.Result, error) {
	if shouldProcess, hash := h.pmm.ShouldProcessMsg(ctx, msg); shouldProcess {
		h.doKeysignResult(ctx, msg.Data)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{}, nil
	}

	return &sdk.Result{}, nil
}

func (h *HandlerKeysignResult) doKeysignResult(ctx sdk.Context, keysignResult *types.KeysignResult) {
	if keysignResult.Success {
		// TODO: award participants of the keysign with some tokesn.
	} else {
		// TODO: Find culprits and penalize them.
		// Do a retry
		txOutId := keysignResult.TxOutId
		retryCount := h.keeper.GetKeySignRetryCount(ctx, txOutId)
		log.Verbosef("Keysign failed, doing retry number %d for tx out %s", retryCount, txOutId)
		params := h.keeper.GetParams(ctx)
		if retryCount < int(params.MaxKeysignRetry) { // We can make 2 in the config?
			// Do retry signing txOut again.
			txOut := h.keeper.GetFinalizedTxOut(ctx, txOutId)
			if txOut == nil {
				log.Errorf("Cannot find TxOut to retry, txOut id = %s", txOutId)
				return
			}

			h.background.AddRetryTxOut(ctx.BlockHeight(), txOut)
			h.keeper.SetKeySignRetryCount(ctx, txOutId, retryCount+1)
		}
	}
}
