package tss

import (
	"fmt"
	"strconv"

	hTypes "github.com/sisu-network/dheart/types"

	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/tss/types"

	etypes "github.com/ethereum/go-ethereum/core/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
)

// checkTxOut checks if a TxOut message is valid before it is added into Sisu block.
func (p *Processor) checkTxOut(ctx sdk.Context, msg *types.TxOutWithSigner) error {
	if !p.privateDb.IsTxOutExisted(msg.Data) {
		return ErrCannotFindMessage
	}

	if p.keeper.IsTxOutExisted(ctx, msg.Data) {
		return ErrMessageHasBeenProcessed
	}

	return nil
}

// deliverTxOut executes a TxOut transaction after it's included in Sisu block. If this node is
// catching up with the network, we would not send the tx to TSS for signing.
func (p *Processor) deliverTxOut(ctx sdk.Context, msgWithSigner *types.TxOutWithSigner) ([]byte, error) {
	txOut := msgWithSigner.Data

	if p.keeper.IsTxOutExisted(ctx, txOut) {
		// The message has been processed
		return nil, nil
	}

	log.Info("Delivering TxOut")

	// Save this to KVStore
	p.keeper.SaveTxOut(ctx, txOut)
	p.privateDb.SaveTxOut(txOut)

	// If this is a txout deployment,

	// Do key signing if this node is not catching up.
	if !p.globalData.IsCatchingUp() {
		// Only Deliver TxOut if the chain has been up to date.
		if libchain.IsETHBasedChain(txOut.OutChain) {
			p.signTx(ctx, txOut)
		}
	}

	return nil, nil
}

// signTx sends a TxOut to dheart for TSS signing.
func (p *Processor) signTx(ctx sdk.Context, tx *types.TxOut) {
	outHash := tx.GetHash()

	log.Info("Delivering TXOUT for chain", tx.OutChain, " tx hash = ", tx.GetHash())
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
		Id:          p.getKeysignRequestId(tx.OutChain, ctx.BlockHeight(), outHash),
		InChain:     tx.InChain,
		OutChain:    tx.OutChain,
		OutHash:     outHash,
		BytesToSign: hash[:],
	}

	pubKeys := p.partyManager.GetActivePartyPubkeys()
	// if err := p.db.UpdateTxOutStatus(tx.OutChain, tx.GetHash(), tssTypes.TxOutStatusSigning, false); err != nil {
	// 	log.Error(err)
	// }

	err := p.dheartClient.KeySign(keysignReq, pubKeys)

	if err != nil {
		log.Error("Keysign: err =", err)
		// _ = p.db.UpdateTxOutStatus(tx.OutChain, tx.GetHash(), tssTypes.TxOutStatusSignFailed, false)
	}
	// _ = p.db.UpdateTxOutStatus(tx.OutChain, tx.GetHash(), tssTypes.TxOutStatusSigned, false)
}

func (p *Processor) getKeysignRequestId(chain string, blockHeight int64, txHash string) string {
	return chain + "_" + strconv.Itoa(int(blockHeight)) + "_" + txHash
}
