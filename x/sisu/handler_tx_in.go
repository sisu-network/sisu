package sisu

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxIn struct {
	pmm    PostedMessageManager
	keeper keeper.Keeper
}

func NewHandlerTxIn(mc ManagerContainer) *HandlerTxIn {
	return &HandlerTxIn{
		keeper: mc.Keeper(),
		pmm:    mc.PostedMessageManager(),
	}
}

func (h *HandlerTxIn) DeliverMsg(ctx sdk.Context, signerMsg *types.TxsInMsg) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doTxIn(ctx, signerMsg.Data)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

// Delivers observed Txs.
func (h *HandlerTxIn) doTxIn(ctx sdk.Context, msg *types.TxsIn) ([]byte, error) {
	log.Infof("Deliverying TxIn on chain %s with request length = %d", msg.Chain, len(msg.Requests))

	allTransfers := make(map[string][]*types.Transfer)
	// Add the message to the queue for later processing.
	for _, request := range msg.Requests {
		transfer := &types.Transfer{
			Id:        types.GetTransferId(msg.Chain, request.Hash),
			Recipient: request.Recipient,
			Token:     request.Token,
			Amount:    request.Amount,
		}

		if allTransfers[request.ToChain] == nil {
			allTransfers[request.ToChain] = h.keeper.GetTransferQueue(ctx, request.ToChain)
			if allTransfers[request.ToChain] == nil {
				allTransfers[request.ToChain] = make([]*types.Transfer, 0)
			}
		}

		fmt.Println("Adding transfer to the queue, transfer = ", *transfer)

		allTransfers[request.ToChain] = append(allTransfers[request.ToChain], transfer)
	}

	// Save all of transfer to the transfer queue
	for _, request := range msg.Requests {
		if allTransfers[request.ToChain] == nil {
			continue
		}

		fmt.Println("setting transfer queue for chain ", request.ToChain, " allTransfers length = ",
			len(allTransfers[request.ToChain]))

		h.keeper.SetTransferQueue(ctx, request.ToChain, allTransfers[request.ToChain])
		allTransfers[request.ToChain] = nil
	}

	return nil, nil
}
