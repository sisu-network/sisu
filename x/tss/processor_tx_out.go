package tss

import (
	"fmt"
	"strconv"

	hTypes "github.com/sisu-network/dheart/types"

	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/tss/types"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"

	etypes "github.com/ethereum/go-ethereum/core/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
)

// checkTxOut checks if a TxOut message is valid before it is added into Sisu block.
func (p *Processor) checkTxOut(ctx sdk.Context, msg *types.TxOutWithSigner) error {
	if p.db.IsTxOutExisted(msg.Data) {
		return nil
	}

	return ErrCannotFindMessage
}

// deliverTxOut executes a TxOut transaction after it's included in Sisu block. If this node is
// catching up with the network, we would not send the tx to TSS for signing.
func (p *Processor) deliverTxOut(ctx sdk.Context, msgWithSigner *types.TxOutWithSigner) ([]byte, error) {
	txOut := msgWithSigner.Data

	if p.keeper.IsTxOutExisted(ctx, txOut) {
		return nil, nil
	}

	log.Info("Delivering TxOut")

	// Save this to KVStore
	p.keeper.SaveTxOut(ctx, txOut)

	// Save this to private db.
	txs := make([]*types.TxOut, 1)
	txs[0] = txOut
	p.db.InsertTxOuts(txs)

	// Do key signing if this node is not catching up.
	if !p.globalData.IsCatchingUp() {
		// Only Deliver TxOut if the chain has been up to date.
		if libchain.IsETHBasedChain(txOut.OutChain) {
			if err := p.db.UpdateTxOutStatus(txOut.OutChain, txOut.GetHash(), tssTypes.TxOutStatusPreSigning, false); err != nil {
				return nil, err
			}

			return p.signTx(ctx, txOut)
		}
	}

	return nil, nil
}

// signTx sends a TxOut to dheart for TSS signing.
func (p *Processor) signTx(ctx sdk.Context, tx *types.TxOut) ([]byte, error) {
	outHash := tx.GetHash()

	log.Verbose("Delivering TXOUT for chain", tx.OutChain, " tx hash = ", tx.GetHash())

	ethTx := &etypes.Transaction{}
	if err := ethTx.UnmarshalBinary(tx.OutBytes); err != nil {
		log.Error("cannot unmarshal tx, err =", err)
		return nil, err
	}

	signer := libchain.GetEthChainSigner(tx.OutChain)
	if signer == nil {
		err := fmt.Errorf("cannot find signer for chain %s", tx.OutChain)
		log.Error(err)
		return nil, err
	}

	hash := signer.Hash(ethTx)

	// 4. Send it to Dheart for signing.
	keysignReq := &hTypes.KeysignRequest{
		Id:             p.getKeysignRequestId(tx.OutChain, ctx.BlockHeight(), outHash),
		OutChain:       tx.OutChain,
		OutBlockHeight: p.currentHeight.Load().(int64),
		OutHash:        outHash,
		BytesToSign:    hash[:],
	}

	pubKeys := p.partyManager.GetActivePartyPubkeys()
	if err := p.db.UpdateTxOutStatus(tx.OutChain, tx.GetHash(), tssTypes.TxOutStatusSigning, false); err != nil {
		log.Error(err)
		return nil, err
	}

	err := p.dheartClient.KeySign(keysignReq, pubKeys)
	if err != nil {
		log.Error("Keysign: err =", err)
		_ = p.db.UpdateTxOutStatus(tx.OutChain, tx.GetHash(), tssTypes.TxOutStatusSignFailed, false)
		return nil, err
	}

	_ = p.db.UpdateTxOutStatus(tx.OutChain, tx.GetHash(), tssTypes.TxOutStatusSigned, false)

	return nil, nil
}

func (p *Processor) getKeysignRequestId(chain string, blockHeight int64, txHash string) string {
	return chain + "_" + strconv.Itoa(int(blockHeight)) + "_" + txHash
}
