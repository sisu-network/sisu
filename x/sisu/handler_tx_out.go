package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/eth"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxOut struct {
	pmm           PostedMessageManager
	keeper        keeper.Keeper
	globalData    common.GlobalData
	transferQueue TransferQueue
	txOutQueue    TxOutQueue
}

func NewHandlerTxOut(mc ManagerContainer) *HandlerTxOut {
	return &HandlerTxOut{
		keeper:        mc.Keeper(),
		pmm:           mc.PostedMessageManager(),
		txOutQueue:    mc.TxOutQueue(),
		globalData:    mc.GlobalData(),
		transferQueue: mc.TxInQueue(),
	}
}

func (h *HandlerTxOut) DeliverMsg(ctx sdk.Context, signerMsg *types.TxOutMsg) (*sdk.Result, error) {
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
	if txOut.TxType == types.TxOutType_CONTRACT_DEPLOYMENT {
		h.keeper.UpdateContractsStatus(ctx, txOut.OutChain, txOut.ContractHash, string(types.TxOutStatusSigning))
	}

	// Move the the transfers associated with this tx_out to pending.
	queue := h.keeper.GetTransferQueue(ctx, txOut.OutChain)
	newQueue := make([]*types.Transfer, 0)
	pending := make([]*types.Transfer, 0)
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
			pending = append(pending, transfer)
		}
	}

	h.keeper.SetTransferQueue(ctx, txOut.OutChain, newQueue)
	h.keeper.SetPendingTransfers(ctx, txOut.OutChain, pending)

	if !h.globalData.IsCatchingUp() {
		h.txOutQueue.AddTxOut(txOut)
	}

	return nil, nil
}

func (h *HandlerTxOut) removeTransfers(ctx sdk.Context, txOut *types.TxOut) {
	var transfers []*types.TransferOutData

	if libchain.IsETHBasedChain(txOut.OutChain) {
		erc20gatewayContract := SupportedContracts[ContractErc20Gateway]
		gwAbi := erc20gatewayContract.Abi

		ethTx := &ethTypes.Transaction{}
		err := ethTx.UnmarshalBinary(txOut.OutBytes)
		if err != nil {
			log.Error("Failed to unmarshall eth tx. err =", err)
			return
		}

		transfers, err = eth.ParseEthTransferIn(ctx, ethTx, txOut.OutChain, gwAbi, h.keeper)
		if err != nil {
			log.Error("removeTransfers: cannot parse transfer out")
			return
		}

		if len(transfers) != len(txOut.InHashes) {
			log.Error("transfers size does not match in chains or in hashes size")
			return
		}
	}
}
