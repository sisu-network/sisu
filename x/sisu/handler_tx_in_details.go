package sisu

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/chains"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxInDetails struct {
	pmm           PostedMessageManager
	keeper        keeper.Keeper
	globalData    common.GlobalData
	bridgeManager chains.BridgeManager
	valsManager   ValidatorManager
	privateDb     keeper.PrivateDb
}

func NewHandlerTxInDetails(
	pmm PostedMessageManager,
	keeper keeper.Keeper,
	globalData common.GlobalData,
	bridgeManager chains.BridgeManager,
	valsManager ValidatorManager,
	privateDb keeper.PrivateDb,
) *HandlerTxInDetails {
	return &HandlerTxInDetails{
		pmm:           pmm,
		keeper:        keeper,
		globalData:    globalData,
		bridgeManager: bridgeManager,
		valsManager:   valsManager,
		privateDb:     privateDb,
	}
}

func (h *HandlerTxInDetails) DeliverMsg(ctx sdk.Context, msg *types.TxInDetailsMsg) (*sdk.Result, error) {
	if shouldProcess, hash := h.pmm.ShouldProcessMsg(ctx, msg); shouldProcess {
		h.doTxInDetails(ctx, h.keeper, msg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{}, nil
	}

	return &sdk.Result{}, nil
}

// doTxInDetails should only be called when majority of nodes has submitted thin TxIn (to confirm)
// and either:
// 1) The assigned validator submitted the TxInDetails within a time frame
// 2) The assigned validator fails to submit TxInDetails but majority of nodes have submitted the
// TxIn details.
func (h *HandlerTxInDetails) doTxInDetails(ctx sdk.Context, k keeper.Keeper, msg *types.TxInDetailsMsg) {
	log.Verbosef("Process TxInDetails with TxIn id %s", msg.Data.TxIn.Id)

	hash, _, err := keeper.GetTxRecordHash(msg)
	if err != nil {
		log.Errorf("doTxInDetails: Failed to get tx record hash for TxInDetailsMsg")
		return
	}

	k.ProcessTxRecord(ctx, hash)

	// 1 .Save the tx in details.
	k.SetTxInDetails(ctx, msg.Data.FromChain, msg.Data)

	// 2. Save the transfers
	h.saveTransfers(ctx, k, msg.Data.Transfers)

	// 3. Save the transfer state
	fmt.Println("BBBBB 00")
	h.privateDb.SetTransferState(msg.Data.TxIn.GetId(), types.TransferState_Confirmed)
	fmt.Println("BBBBB 11")
}

func (h *HandlerTxInDetails) saveTransfers(ctx sdk.Context, k keeper.Keeper, transfers []*types.TransferDetails) {
	k.AddTransfers(ctx, transfers)

	for _, transfer := range transfers {
		// TODO: Optimize this path. We can save single transfer instead of the entire queue.
		queue := k.GetTransferQueue(ctx, transfer.ToChain)
		queue = append(queue, transfer)
		k.SetTransferQueue(ctx, transfer.ToChain, queue)
	}
}
