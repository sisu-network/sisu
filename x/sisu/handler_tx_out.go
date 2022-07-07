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
	transferQueue TxInQueue
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

func (h *HandlerTxOut) DeliverMsg(ctx sdk.Context, signerMsg *types.TxOutWithSigner) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doTxOut(ctx, signerMsg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

// deliverTxOut executes a TxOut transaction after it's included in Sisu block. If this node is
// catching up with the network, we would not send the tx to TSS for signing.
func (h *HandlerTxOut) doTxOut(ctx sdk.Context, msgWithSigner *types.TxOutWithSigner) ([]byte, error) {
	txOut := msgWithSigner.Data

	log.Info("Delivering TxOut")

	// Save this to KVStore
	h.keeper.SaveTxOut(ctx, txOut)

	if !h.globalData.IsCatchingUp() {
		h.txOutQueue.AddTxOut(txOut)
	}

	// Remove all the transfer request in the tx in queue.

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

		if len(transfers) != len(txOut.InChains) || len(transfers) != len(txOut.InHashes) {
			log.Error("transfers size does not match in chains or in hashes size")
			return
		}
	}

	h.transferQueue.RemoveTransfers(ctx, transfers)
}
