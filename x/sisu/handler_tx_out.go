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
	pmm        PostedMessageManager
	keeper     keeper.Keeper
	globalData common.GlobalData
	txOutQueue TxOutQueue
}

func NewHandlerTxOut(mc ManagerContainer) *HandlerTxOut {
	return &HandlerTxOut{
		keeper:     mc.Keeper(),
		pmm:        mc.PostedMessageManager(),
		txOutQueue: mc.TxOutQueue(),
		globalData: mc.GlobalData(),
	}
}

func (h *HandlerTxOut) DeliverMsg(ctx sdk.Context, signerMsg *types.TxOutMsg) (*sdk.Result, error) {
	fmt.Println("BBB 000")
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doTxOut(ctx, signerMsg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

// deliverTxOut executes a TxOut transaction after it's included in Sisu block. If this node is
// catching up with the network, we would not send the tx to TSS for signing.
func (h *HandlerTxOut) doTxOut(ctx sdk.Context, txOutMsg *types.TxOutMsg) ([]byte, error) {
	txOut := txOutMsg.Data

	log.Info("Delivering TxOut")

	// Save this to KVStore
	h.keeper.SaveTxOut(ctx, txOut)

	// If this is a txOut deployment, mark the contract as being deployed.
	switch txOut.TxType {
	case types.TxOutType_CONTRACT_DEPLOYMENT:
		h.keeper.UpdateContractsStatus(ctx, txOut.OutChain, txOut.ContractHash, string(types.TxOutStatusSigning))
	case types.TxOutType_TRANSFER_OUT:
		h.handlerTransferOut(ctx, txOut)
	}

	return nil, nil
}

func (h *HandlerTxOut) handlerTransferOut(ctx sdk.Context, txOut *types.TxOut) {
	// If there are some transfer in the pendings queue, don't process this txOut.
	pendings := h.keeper.GetPendingTransfers(ctx, txOut.OutChain)
	if len(pendings) > 0 {
		log.Verbose("There are some pending transfers in the pending queue, don't process new transfers")
		return
	}

	// Move the the transfers associated with this tx_out to pending.
	queue := h.keeper.GetTransferQueue(ctx, txOut.OutChain)
	newQueue := make([]*types.Transfer, 0)
	pendings = make([]*types.Transfer, 0)
	for _, transfer := range queue {
		found := false
		for _, inHash := range txOut.InHashes {
			if transfer.Id == inHash {
				found = true
				break
			}
		}

		if !found {
			newQueue = append(newQueue, transfer)
		} else {
			pendings = append(pendings, transfer)
		}
	}

	h.keeper.SetTransferQueue(ctx, txOut.OutChain, newQueue)
	h.keeper.SetPendingTransfers(ctx, txOut.OutChain, pendings)

	if !h.globalData.IsCatchingUp() {
		h.txOutQueue.AddTxOut(txOut)
	}
}
