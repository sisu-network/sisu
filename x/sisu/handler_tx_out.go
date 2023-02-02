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
	bacground   background.Background
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
		bacground:   mc.Background(),
	}
}

// There are 2 cases where a TxOut can be finalized:
// 1) The assigned validator submits the TxOut and it's approved 2/3 of validators
// 2) The proposed txOut is rejected or it is not produced during a timeout period. At this time,
// every validator node submits its own txOut and everyone to come up with a consensused txOut.
func (h *HandlerTxOut) DeliverMsg(ctx sdk.Context, msg *types.TxOutMsg) (*sdk.Result, error) {
	// TODO: In most cases, different txOuts produced by different validator will have different hash.
	// We should get the average gas price of all txOuts to calculate the final gas.
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

	// Add this message to the validation queue.
	h.bacground.AddVoteTxOut(ctx.BlockHeight(), msg)

	return &sdk.Result{}, nil
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

	// 4. Set Expiration Block
	params := k.GetParams(ctx)
	k.SetExpirationBlock(ctx, types.ExpirationBlock_TxOut, txOut.GetId(),
		ctx.BlockHeight()+int64(params.ExpirationBlock))
}
