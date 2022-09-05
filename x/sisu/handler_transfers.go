package sisu

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTransfers struct {
	pmm    PostedMessageManager
	keeper keeper.Keeper
}

func NewHandlerTransfers(mc ManagerContainer) *HandlerTransfers {
	return &HandlerTransfers{
		keeper: mc.Keeper(),
		pmm:    mc.PostedMessageManager(),
	}
}

func (h *HandlerTransfers) DeliverMsg(ctx sdk.Context, signerMsg *types.TransfersMsg) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doTransfers(ctx, signerMsg.Data)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

// Delivers observed Txs.
func (h *HandlerTransfers) doTransfers(ctx sdk.Context, msg *types.Transfers) ([]byte, error) {
	if len(msg.Transfers) == 0 {
		return nil, fmt.Errorf("Empty transfers array")
	}

	log.Infof("Deliverying TransferOut on chain %s with request length = %d",
		msg.Transfers[0].FromChain, len(msg.Transfers))

	allTransfers := make(map[string][]*types.Transfer)
	// Add the message to the queue for later processing.
	for _, request := range msg.Transfers {
		if allTransfers[request.ToChain] == nil {
			allTransfers[request.ToChain] = h.keeper.GetTransferQueue(ctx, request.ToChain)

		}

		log.Debug("Adding transfer to the queue, transfer = ", *request)

		allTransfers[request.ToChain] = append(allTransfers[request.ToChain], request)
	}

	// Save all of transfer to the transfer queue
	for _, request := range msg.Transfers {
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
