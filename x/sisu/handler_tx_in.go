package sisu

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxIn struct {
	pmm       PostedMessageManager
	keeper    keeper.Keeper
	txInQueue TransferQueue
}

func NewHandlerTxIn(mc ManagerContainer) *HandlerTxIn {
	return &HandlerTxIn{
		keeper:    mc.Keeper(),
		pmm:       mc.PostedMessageManager(),
		txInQueue: mc.TxInQueue(),
	}
}

func (h *HandlerTxIn) DeliverMsg(ctx sdk.Context, signerMsg *types.TxsInMsg) (*sdk.Result, error) {
	fmt.Println("AAAAAA DeliverMsg, data = ", *signerMsg.Data)
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doTxIn(ctx, signerMsg.Data)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

// Delivers observed Txs.
func (h *HandlerTxIn) doTxIn(ctx sdk.Context, msg *types.TxsIn) ([]byte, error) {
	log.Info("Deliverying TxIn on chain ", msg.Chain)

	allTransfers := make(map[string][]*types.Transfer)
	// Add the message to the queue for later processing.
	for _, request := range msg.Requests {
		transfer := &types.Transfer{
			Id:        fmt.Sprintf("%s__%s", msg.Chain, request.Hash),
			Recipient: request.Recipient,
			Token:     request.Token,
			Amount:    request.Amount,
		}

		if allTransfers[request.ToChain] == nil {
			allTransfers[request.ToChain] = h.keeper.GetTransferQueue(ctx, msg.Chain)
			if allTransfers[request.ToChain] == nil {
				allTransfers[request.ToChain] = make([]*types.Transfer, 0)
			}
		}

		allTransfers[request.ToChain] = append(allTransfers[request.ToChain], transfer)
	}

	// Save all of transfer to the transfer queue
	for _, request := range msg.Requests {
		if allTransfers[request.ToChain] == nil {
			continue
		}

		fmt.Println("transfer queue size: ", request.ToChain, len(allTransfers[request.ToChain]))

		h.keeper.SetTransferQueue(ctx, request.ToChain, allTransfers[request.ToChain])

		allTransfers[request.ToChain] = nil
	}

	return nil, nil
}
