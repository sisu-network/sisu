package sisu

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxIn struct {
	pmm         PostedMessageManager
	keeper      keeper.Keeper
	valsManager ValidatorManager
	globalData  common.GlobalData
}

func NewHandlerTxIn(
	pmm PostedMessageManager,
	keeper keeper.Keeper,
	valsManager ValidatorManager,
	globalData common.GlobalData,
) *HandlerTxIn {
	return &HandlerTxIn{
		pmm:         pmm,
		keeper:      keeper,
		valsManager: valsManager,
		globalData:  globalData,
	}
}

func (h *HandlerTxIn) DeliverMsg(ctx sdk.Context, msg *types.TxInMsg) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, msg); process {
		fmt.Println("TxIn Hash 1 = ", hex.EncodeToString(hash))
		data, err := h.doTxIn(ctx, msg)
		h.keeper.ProcessTxRecord(ctx, hash)

		fmt.Println("Is processed hash 1 = ", h.keeper.IsTxRecordProcessed(ctx, hash))

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

func (h *HandlerTxIn) doTxIn(ctx sdk.Context, msg *types.TxInMsg) ([]byte, error) {
	// Check if we have the TxIn details .
	txInDetails := h.keeper.GetTxInDetails(ctx, msg.Data.Id)
	fmt.Println("BBBBB 000000")
	if txInDetails != nil {
		fmt.Println("BBBBB 111111")
		// 1. TODO: Do verificaiton on the tx in details to make sure this data is correct (including
		// the transfers)

		fmt.Println("BBBBB txInDetails.Data = ", txInDetails.Data)

		// 2. Add all the new transfers to the transfer queue.
		saveTransfers(ctx, h.keeper, txInDetails.Data.Transfers)
	} else {
		// We have not received the txInDetails yet
	}

	return nil, nil
}

func saveTransfers(ctx sdk.Context, k keeper.Keeper, transfers []*types.TransferDetails) {
	k.AddTransfers(ctx, transfers)

	for _, transfer := range transfers {
		// TODO: Optimize this path. We can save single transfer instead of the entire queue.
		queue := k.GetTransferQueue(ctx, transfer.ToChain)
		queue = append(queue, transfer)
		k.SetTransferQueue(ctx, transfer.ToChain, queue)
	}
}
