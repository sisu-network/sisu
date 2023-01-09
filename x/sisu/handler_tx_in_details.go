package sisu

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
}

func NewHandlerTxInDetails(
	pmm PostedMessageManager,
	keeper keeper.Keeper,
	globalData common.GlobalData,
	bridgeManager chains.BridgeManager,
	valsManager ValidatorManager,
) *HandlerTxInDetails {
	return &HandlerTxInDetails{
		pmm:           pmm,
		keeper:        keeper,
		globalData:    globalData,
		bridgeManager: bridgeManager,
		valsManager:   valsManager,
	}
}

func (h *HandlerTxInDetails) DeliverMsg(ctx sdk.Context, msg *types.TxInDetailsMsg) (*sdk.Result, error) {
	txIn := msg.Data.TxIn
	processed, hash := h.pmm.ShouldProcessMsg(ctx, msg)

	if h.keeper.IsTxRecordProcessed(ctx, []byte(txIn.Id)) {
		fmt.Println("AAAAA 0000")

		// Case 1: the thin tx is confirmed but no tx details is saved yet.
		if !processed {
			h.keeper.ProcessTxRecord(ctx, hash)

			// Save the tx in details.
			h.keeper.SetTxInDetails(ctx, msg.Data.FromChain, msg.Data)

			// 2. Save the transfers
			chain := msg.Data.FromChain
			q := h.keeper.GetTransferQueue(ctx, chain)
			q = append(q, msg.Data.Transfers...)
			h.keeper.SetTransferQueue(ctx, chain, q)
		}
	} else {
		fmt.Println("AAAAA 111111")
		// Case 2: thin tx is not confirmed yet. We only want to save the TxInDetails of the assigned
		// node for this TxIn. Later on, when the thin TxIn is confirmed, we already have the details
		// for it to process.
		assignedNode := h.valsManager.GetAssignedValidator(ctx, msg.Data.TxIn.Id)
		if assignedNode.AccAddress == msg.Signer {
			h.keeper.SetTxInDetails(ctx, msg.Data.FromChain, msg.Data)
		}
	}

	return &sdk.Result{}, nil
}
