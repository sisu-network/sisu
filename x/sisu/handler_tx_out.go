package sisu

import (
	"encoding/base64"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/echovl/cardano-go"
	etypes "github.com/ethereum/go-ethereum/core/types"
	hTypes "github.com/sisu-network/dheart/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxOut struct {
	pmm          PostedMessageManager
	keeper       keeper.Keeper
	globalData   common.GlobalData
	partyManager PartyManager
	dheartClient tssclients.DheartClient
	txTracker    TxTracker
}

func NewHandlerTxOut(mc ManagerContainer) *HandlerTxOut {
	return &HandlerTxOut{
		keeper:       mc.Keeper(),
		pmm:          mc.PostedMessageManager(),
		globalData:   mc.GlobalData(),
		partyManager: mc.PartyManager(),
		dheartClient: mc.DheartClient(),
		txTracker:    mc.TxTracker(),
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

	// Do key signing if this node is not catching up.
	if !h.globalData.IsCatchingUp() {
		// Only Deliver TxOut if the chain has been up to date.
		if libchain.IsETHBasedChain(txOut.OutChain) {
			h.signTx(ctx, txOut)
		}

		if libchain.IsCardanoChain(txOut.OutChain) {
			h.signCardanoTx(ctx, txOut)
		}
	}

	return nil, nil
}

// signTx sends a TxOut to dheart for TSS signing.
func (h *HandlerTxOut) signTx(ctx sdk.Context, tx *types.TxOut) {
	// Update the txOut to be delivered.
	h.txTracker.UpdateStatus(tx.OutChain, tx.OutHash, types.TxStatusDelivered)

	log.Info("Delivering TXOUT for chain", tx.OutChain, " tx hash = ", tx.OutHash)
	if tx.TxType == types.TxOutType_CONTRACT_DEPLOYMENT {
		log.Info("This TxOut is a contract deployment")
	}

	ethTx := &etypes.Transaction{}
	if err := ethTx.UnmarshalBinary(tx.OutBytes); err != nil {
		log.Error("cannot unmarshal tx, err =", err)
	}

	signer := libchain.GetEthChainSigner(tx.OutChain)
	if signer == nil {
		err := fmt.Errorf("cannot find signer for chain %s", tx.OutChain)
		log.Error(err)
	}

	hash := signer.Hash(ethTx)

	// 4. Send it to Dheart for signing.
	keysignReq := &hTypes.KeysignRequest{
		KeyType: libchain.KEY_TYPE_ECDSA,
		KeysignMessages: []*hTypes.KeysignMessage{
			{
				Id:          h.getKeysignRequestId(tx.OutChain, ctx.BlockHeight(), tx.OutHash),
				InChain:     tx.InChain,
				OutChain:    tx.OutChain,
				OutHash:     tx.OutHash,
				BytesToSign: hash[:],
			},
		},
	}

	pubKeys := h.partyManager.GetActivePartyPubkeys()

	err := h.dheartClient.KeySign(keysignReq, pubKeys)
	if err != nil {
		log.Error("Keysign: err =", err)
	}
}

func (h *HandlerTxOut) signCardanoTx(ctx sdk.Context, txOut *types.TxOut) {
	// Update the txOut to be delivered.
	h.txTracker.UpdateStatus(txOut.OutChain, txOut.OutHash, types.TxStatusDelivered)

	tx := &cardano.Tx{}
	if err := tx.UnmarshalCBOR(txOut.OutBytes); err != nil {
		log.Error("error when unmarshalling cardano tx out: ", err)
		return
	}

	bz := base64.StdEncoding.EncodeToString(txOut.OutBytes)
	log.Debug("bz in signCardanoTx = ", bz)

	txHash, err := tx.Hash()
	if err != nil {
		log.Error("error when getting cardano tx hash: ", err)
		return
	}

	signRequest := &hTypes.KeysignRequest{
		KeyType: libchain.KEY_TYPE_EDDSA,
		KeysignMessages: []*hTypes.KeysignMessage{
			{
				Id:          h.getKeysignRequestId(txOut.OutChain, ctx.BlockHeight(), txOut.OutHash),
				InChain:     txOut.InChain,
				OutChain:    txOut.OutChain,
				OutHash:     txOut.OutHash,
				BytesToSign: txHash[:],
			},
		},
	}

	pubKeys := h.partyManager.GetActivePartyPubkeys()
	err = h.dheartClient.KeySign(signRequest, pubKeys)
	if err != nil {
		log.Error("Keysign: err =", err)
	}
}

func (h *HandlerTxOut) getKeysignRequestId(chain string, blockHeight int64, txHash string) string {
	return chain + "_" + strconv.Itoa(int(blockHeight)) + "_" + txHash
}
