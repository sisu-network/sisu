package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxOutResult struct {
	pmm           PostedMessageManager
	keeper        keeper.Keeper
	transferQueue TransferQueue
}

func NewHandlerTxOutResult(mc ManagerContainer) *HandlerTxOutResult {
	return &HandlerTxOutResult{
		keeper:        mc.Keeper(),
		pmm:           mc.PostedMessageManager(),
		transferQueue: mc.TransferQueue(),
	}
}

func (h *HandlerTxOutResult) DeliverMsg(ctx sdk.Context, signerMsg *types.TxOutResultMsg) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doTxOutResult(ctx, signerMsg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

func (h *HandlerTxOutResult) doTxOutResult(ctx sdk.Context, msgWithSigner *types.TxOutResultMsg) ([]byte, error) {
	log.Info("Delivering TxOutResult")

	msg := msgWithSigner.Data
	txOut := h.keeper.GetTxOut(ctx, msg.OutChain, msg.OutHash)
	if txOut == nil {
		log.Errorf("cannot find txout from txOutConfirm message, chain = %s & hash = %s",
			msg.OutChain, msg.OutHash)
		return nil, nil
	}

	log.Verbose("msg.Result = ", msg.Result)

	switch msg.Result {
	case types.TxOutResultType_IN_BLOCK_SUCCESS:
		return h.doTxOutConfirm(ctx, msg, txOut)
	case types.TxOutResultType_GENERIC_ERROR, types.TxOutResultType_NOT_ENOUGH_NATIVE_BALANCE,
		types.TxOutResultType_IN_BLOCK_FAILURE:
		return h.doTxOutFailure(ctx, msg, txOut)
	}

	return nil, nil
}

func (h *HandlerTxOutResult) doTxOutConfirm(ctx sdk.Context, msg *types.TxOutResult, txOut *types.TxOut) ([]byte, error) {
	log.Verbose("Transaction is successfully included in a block, hash = ", msg.OutHash, " chain = ", msg.OutChain)

	savedCheckPoint := h.keeper.GetGatewayCheckPoint(ctx, msg.OutChain)
	if savedCheckPoint == nil || savedCheckPoint.BlockHeight < msg.BlockHeight {
		// Save checkpoint
		checkPoint := &types.GatewayCheckPoint{
			Chain:       msg.OutChain,
			BlockHeight: msg.BlockHeight,
		}

		if libchain.IsETHBasedChain(msg.OutChain) {
			checkPoint.Nonce = msg.Nonce
		}

		// Update observed block height and nonce.
		h.keeper.AddGatewayCheckPoint(ctx, checkPoint)
	}

	// Clear the pending TxOut
	log.Verbose("Clearing pending out for chain ", txOut.Content.OutChain)
	h.keeper.SetPendingTxOutInfo(ctx, txOut.Content.OutChain, nil)

	return nil, nil
}

func (h *HandlerTxOutResult) doTxOutFailure(ctx sdk.Context, msg *types.TxOutResult, txOut *types.TxOut) ([]byte, error) {
	log.Verbose("Transaction failed!, txOut.TxType = ", txOut.TxType)

	switch txOut.TxType {
	case types.TxOutType_TRANSFER_OUT:
		ids := txOut.Input.TransferIds
		transfers := h.keeper.GetTransfers(ctx, ids)

		// Update the retry number of these transfers.
		for _, transfer := range transfers {
			transfer.RetryNum++
			h.keeper.AddTransfers(ctx, []*types.Transfer{transfer})
		}

		// TODO: Figure out when we should process this transfer since we do not have enough funding.
		// transferQ := h.keeper.GetTransferQueue(ctx, msg.OutChain)
		// for _, transfer := range transfers {
		// 	h.keeper.SetTransferQueue(ctx, msg.OutChain, transferQ)
		// 	transferQ = append(transferQ, transfer)
		// }
	}

	// Clear the pending TxOut
	log.Verbose("Clearing pending out for chain ", txOut.Content.OutChain)
	h.keeper.SetPendingTxOutInfo(ctx, txOut.Content.OutChain, nil)

	return nil, nil
}
