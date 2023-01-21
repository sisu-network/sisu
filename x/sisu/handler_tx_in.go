package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/chains"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxIn struct {
	pmm           PostedMessageManager
	keeper        keeper.Keeper
	globalData    common.GlobalData
	bridgeManager chains.BridgeManager
	valsManager   ValidatorManager
	privateDb     keeper.PrivateDb
}

func NewHandlerTxIn(
	pmm PostedMessageManager,
	keeper keeper.Keeper,
	globalData common.GlobalData,
	bridgeManager chains.BridgeManager,
	valsManager ValidatorManager,
	privateDb keeper.PrivateDb,
) *HandlerTxIn {
	return &HandlerTxIn{
		pmm:           pmm,
		keeper:        keeper,
		globalData:    globalData,
		bridgeManager: bridgeManager,
		valsManager:   valsManager,
		privateDb:     privateDb,
	}
}

func (h *HandlerTxIn) DeliverMsg(ctx sdk.Context, msg *types.TxInMsg) (*sdk.Result, error) {
	if shouldProcess, hash := h.pmm.ShouldProcessMsg(ctx, msg); shouldProcess {
		h.doTxIn(ctx, h.keeper, msg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{}, nil
	}

	return &sdk.Result{}, nil
}

func (h *HandlerTxIn) doTxIn(ctx sdk.Context, k keeper.Keeper, msg *types.TxInMsg) {
	log.Verbosef("Process doTxIn with TxIn id %s", msg.Data.Id)

	hash, _, err := keeper.GetTxRecordHash(msg)
	if err != nil {
		log.Errorf("doTxIn: Failed to get tx record hash for doTxInMsg")
		return
	}

	k.ProcessTxRecord(ctx, hash)

	// Save the transfers
	h.saveTransfers(ctx, k, msg.Data.Transfers)
}

func (h *HandlerTxIn) saveTransfers(ctx sdk.Context, k keeper.Keeper, transfers []*types.TransferDetails) {
	if len(transfers) == 0 {
		log.Warnf("There is no transfer in the TxIn message.")
		return
	}
	k.AddTransfers(ctx, transfers)

	chain := transfers[0].ToChain
	queue := k.GetTransferQueue(ctx, chain)
	for _, transfer := range transfers {
		// TODO: Optimize this path. We can save single transfer instead of the entire queue.
		queue = append(queue, transfer)
	}

	k.SetTransferQueue(ctx, chain, queue)
}
