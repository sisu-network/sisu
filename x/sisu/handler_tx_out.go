package sisu

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxOut struct {
	pmm         PostedMessageManager
	keeper      keeper.Keeper
	valsManager ValidatorManager
	globalData  common.GlobalData
	txSubmit    common.TxSubmit
}

func NewHandlerTxOut(mc ManagerContainer) *HandlerTxOut {
	return &HandlerTxOut{
		keeper:      mc.Keeper(),
		pmm:         mc.PostedMessageManager(),
		valsManager: mc.ValidatorManager(),
		globalData:  mc.GlobalData(),
		txSubmit:    mc.TxSubmit(),
	}
}

func (h *HandlerTxOut) DeliverMsg(ctx sdk.Context, msg *types.TxOutMsg) (*sdk.Result, error) {
	shouldProcess, hash := h.pmm.ShouldProcessMsg(ctx, msg)
	if shouldProcess {
		data, err := doTxOut(ctx, h.keeper, msg.Data)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	fmt.Println("AAAAA HandlerTxOut DeliverMsg")
	if ok, assignedVal := h.checkAssignedValMessage(ctx, msg); ok {
		// Submit the TxOut confirm
		txOutConfirmMsg := types.NewTxOutConsensedMsg(msg.Signer, &types.TxOutConsensed{
			AssignedValidator: assignedVal,
			TxOutId:           msg.Data.GetId(),
		})
		h.txSubmit.SubmitMessageAsync(txOutConfirmMsg)
	}

	return &sdk.Result{}, nil
}

// checkAssignedValMessage checks if a TxOutMsg comes from the assigned validator. If this is true,
// we can submit the confirm TxOut message.
func (h *HandlerTxOut) checkAssignedValMessage(ctx sdk.Context, msg *types.TxOutMsg) (bool, string) {
	// Check if this is the message from assigned validator.
	// TODO: Do a validation to verify that the this TxOut is still within the allowed time interval
	// since confirmed transfers.
	// TODO: if this is a transfer, make sure that the first transfer matches the first transfer in
	// Transfer queue
	transferIds := msg.Data.Input.TransferIds
	if len(transferIds) > 0 {
		queue := h.keeper.GetTransferQueue(ctx, msg.Data.Content.OutChain)
		if len(queue) < len(transferIds) {
			log.Errorf("Transfers list in the message is longer than the saved transfer queue.")
			return false, ""
		}

		if len(queue) > 0 {
			// Make sure that all transfers Ids are the first ids in the queue
			for i, transfer := range queue {
				if i >= len(transferIds) {
					break
				}

				if transfer.Id != transferIds[i] {
					log.Errorf(
						"Transfer ids do not match for index %s, id in the mesage = %s, id in the queue = %s",
						i, transferIds[i], transfer.Id,
					)
					return false, ""
				}
			}

			assignedNode := h.valsManager.GetAssignedValidator(ctx, queue[0].Id)
			if assignedNode.AccAddress == msg.Signer {
				return true, assignedNode.AccAddress
			}
		}
	}

	return false, ""
}

// doTxOut saves a TxOut in the keeper and add it the TxOut Queue.
func doTxOut(ctx sdk.Context, k keeper.Keeper, txOut *types.TxOutOld) ([]byte, error) {
	log.Info("Delivering TxOut")

	// Save this to KVStore
	k.SaveTxOut(ctx, txOut)

	// If this is a txOut deployment, mark the contract as being deployed.
	switch txOut.TxType {
	case types.TxOutType_TRANSFER_OUT:
		handlerTransfer(ctx, k, txOut)
	}

	return nil, nil
}

func handlerTransfer(ctx sdk.Context, k keeper.Keeper, txOut *types.TxOutOld) {
	// 1. Update TxOut queue
	addTxOutToQueue(ctx, k, txOut)

	// 2. Remove the transfers in txOut from the queue
	queue := k.GetTransferQueue(ctx, txOut.Content.OutChain)
	ids := make(map[string]bool, 0)
	for _, inHash := range txOut.Input.TransferIds {
		ids[inHash] = true
	}

	newQueue := make([]*types.TransferDetails, 0)
	for _, transfer := range queue {
		if !ids[transfer.Id] {
			newQueue = append(newQueue, transfer)
		}
	}

	k.SetTransferQueue(ctx, txOut.Content.OutChain, newQueue)
}

func addTxOutToQueue(ctx sdk.Context, k keeper.Keeper, txOut *types.TxOutOld) {
	// Move the the transfers associated with this tx_out to pending.
	queue := k.GetTxOutQueue(ctx, txOut.Content.OutChain)
	queue = append(queue, txOut)
	k.SetTxOutQueue(ctx, txOut.Content.OutChain, queue)
}
