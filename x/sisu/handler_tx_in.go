package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTransferOut struct {
	pmm    PostedMessageManager
	keeper keeper.Keeper
}

func NewHandlerTransferOut(mc ManagerContainer) *HandlerTransferOut {
	return &HandlerTransferOut{
		keeper: mc.Keeper(),
		pmm:    mc.PostedMessageManager(),
	}
}

func (h *HandlerTransferOut) DeliverMsg(ctx sdk.Context, signerMsg *types.TransferOutsMsg) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doTransferOut(ctx, signerMsg.Data)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

// Delivers observed Txs.
func (h *HandlerTransferOut) doTransferOut(ctx sdk.Context, msg *types.TransferOuts) ([]byte, error) {
	log.Infof("Deliverying TransferOut on chain %s with request length = %d", msg.Chain, len(msg.Requests))

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

		log.Debug("Adding transfer to the queue, transfer = ", *transfer)

		allTransfers[request.ToChain] = append(allTransfers[request.ToChain], transfer)
	}

	// Save all of transfer to the transfer queue
	for _, request := range msg.Requests {
		if allTransfers[request.ToChain] == nil {
			continue
		}

		log.Debug("setting transfer queue for chain ", request.ToChain, " allTransfers length = ",
			len(allTransfers[request.ToChain]))

		h.keeper.SetTransferQueue(ctx, request.ToChain, allTransfers[request.ToChain])
		allTransfers[request.ToChain] = nil
	}

	return nil, nil
}
