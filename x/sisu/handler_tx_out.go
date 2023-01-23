package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/background"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxOut struct {
	pmm         components.PostedMessageManager
	keeper      keeper.Keeper
	valsManager components.ValidatorManager
	globalData  components.GlobalData
	txSubmit    components.TxSubmit
	appKeys     components.AppKeys
	privateDb   keeper.PrivateDb
}

func NewHandlerTxOut(mc background.ManagerContainer) *HandlerTxOut {
	return &HandlerTxOut{
		keeper:      mc.Keeper(),
		pmm:         mc.PostedMessageManager(),
		valsManager: mc.ValidatorManager(),
		globalData:  mc.GlobalData(),
		txSubmit:    mc.TxSubmit(),
		appKeys:     mc.AppKeys(),
		privateDb:   mc.PrivateDb(),
	}
}

func (h *HandlerTxOut) DeliverMsg(ctx sdk.Context, msg *types.TxOutMsg) (*sdk.Result, error) {
	shouldProcess, hash := h.pmm.ShouldProcessMsg(ctx, msg)
	if shouldProcess {
		doTxOut(ctx, h.keeper, h.privateDb, msg.Data)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{}, nil
	}

	if h.keeper.IsTxRecordProcessed(ctx, hash) {
		// This msg has been processed before.
		return &sdk.Result{}, nil
	}

	h.keeper.AddProposedTxOut(ctx, msg.Signer, msg.Data)
	// TODO: In case there every one does not approve the assigned node's transaction, we have to
	// use the proposed txOuts from everyone to calculate the final TxOut.

	// do message validation. This work can be done in the background.
	ok, assignedVal := h.validateTxOut(ctx, msg)
	vote := types.VoteResult_APPROVE
	if !ok {
		vote = types.VoteResult_REJECT
	}

	// Submit the TxOut confirm
	txOutConfirmMsg := types.NewTxOutVoteMsg(
		h.appKeys.GetSignerAddress().String(),
		&types.TxOutVote{
			AssignedValidator: assignedVal,
			TxOutId:           msg.Data.GetId(),
			Vote:              vote,
		},
	)

	h.txSubmit.SubmitMessageAsync(txOutConfirmMsg)

	return &sdk.Result{}, nil
}

// validateTxOut checks if a TxOutMsg comes from the assigned validator. If this is true,
// we can submit the confirm TxOut message.
func (h *HandlerTxOut) validateTxOut(ctx sdk.Context, msg *types.TxOutMsg) (bool, string) {
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
func doTxOut(ctx sdk.Context, k keeper.Keeper, privateDb keeper.PrivateDb,
	txOut *types.TxOut) ([]byte, error) {
	log.Info("Delivering TxOut")

	// Save this to KVStore
	k.SaveTxOut(ctx, txOut)

	// If this is a txOut deployment, mark the contract as being deployed.
	switch txOut.TxType {
	case types.TxOutType_TRANSFER_OUT:
		handlerTransfer(ctx, k, privateDb, txOut)
	}

	return nil, nil
}

func handlerTransfer(ctx sdk.Context, k keeper.Keeper, privateDb keeper.PrivateDb,
	txOut *types.TxOut) {
	// 1. Update TxOut txOutQ
	txOutQ := k.GetTxOutQueue(ctx, txOut.Content.OutChain)
	txOutQ = append(txOutQ, txOut)
	k.SetTxOutQueue(ctx, txOut.Content.OutChain, txOutQ)

	// 2. Remove the transfers in txOut from the queue
	transferQ := k.GetTransferQueue(ctx, txOut.Content.OutChain)
	ids := make(map[string]bool, 0)
	for _, inHash := range txOut.Input.TransferIds {
		ids[inHash] = true
	}

	newQueue := make([]*types.TransferDetails, 0)
	for _, transfer := range transferQ {
		if !ids[transfer.Id] {
			newQueue = append(newQueue, transfer)
		}
	}

	k.SetTransferQueue(ctx, txOut.Content.OutChain, newQueue)

	// 3. Update the HoldProcessing for transfer queue so that we do not process any more transfer.
	privateDb.SetHoldProcessing(types.TransferHoldKey, txOut.Content.OutChain, true)

	// Set Expiration Block
	params := k.GetParams(ctx)
	k.SetExpirationBlock(ctx, types.ExpirationBlock_TxOut, txOut.GetId(),
		ctx.BlockHeight()+int64(params.ExpirationBlock))
}
