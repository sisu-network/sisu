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

// Produces response for an observed tx. This has to be deterministic based on all the data that
// the processor has.
func (p *Processor) createAndBroadcastTxOuts(ctx sdk.Context, tx *types.ObservedTx) []*tssTypes.TxOut {
	outMsgs, outEntities := p.txOutputProducer.GetTxOuts(ctx, p.currentHeight, tx)

	// Save this to database
	log.Verbose("len(outEntities) = ", len(outEntities))
	if len(outEntities) > 0 {
		for _, outEntity := range outEntities {
			outEntity.Status = string(tssTypes.TxOutStatusPreBroadcast)
			log.Verbose("Inserting into db, tx hash = ", outEntity.HashWithoutSig)
		}
		p.db.InsertTxOuts(outEntities)
	}

	for _, msg := range outMsgs {
		go func(m *tssTypes.TxOut) {
			if err := p.txSubmit.SubmitMessage(m); err != nil {
				return
			}

			p.db.UpdateTxOutStatus(m.OutChain, m.GetHash(), tssTypes.TxOutStatusBroadcasted, false)
		}(msg)
	}

	return outMsgs
}

// CheckTxOut checks if a TxOut message is valid before it is added into Sisu block.
func (p *Processor) CheckTxOut(ctx sdk.Context, msg *types.TxOut) error {
	if p.keeper.IsTxOutExisted(ctx, msg) {
		return ErrMessageHasBeenProcessed
	}

	return nil
}

// DeliverTxOut executes a TxOut transaction after it's included in Sisu block. If this node is
// catching up with the network, we would not send the tx to TSS for signing.
func (p *Processor) DeliverTxOut(ctx sdk.Context, tx *types.TxOut) ([]byte, error) {
	if p.keeper.IsTxOutExisted(ctx, tx) {
		return nil, nil
	}

	p.keeper.SaveTxOut(ctx, tx)

	if !p.globalData.IsCatchingUp() {
		// Only Deliver TxOut if the chain has been up to date.
		if libchain.IsETHBasedChain(tx.OutChain) {
			if err := p.db.UpdateTxOutStatus(tx.OutChain, tx.GetHash(), tssTypes.TxOutStatusPreSigning, false); err != nil {
				return nil, err
			}

			return p.signTx(ctx, tx)
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
		OutBlockHeight: p.currentHeight,
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
