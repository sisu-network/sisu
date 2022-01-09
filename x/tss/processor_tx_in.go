package tss

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"
)

func (p *Processor) checkTxIn(ctx sdk.Context, msgWithSigner *types.TxInWithSigner) error {
	// Make sure we should have seen this TxIn in our table.
	if !p.privateDb.IsTxInExisted(msgWithSigner.Data) {
		return ErrCannotFindMessage
	}

	// Make sure this message has been processed.
	if p.keeper.IsTxInExisted(ctx, msgWithSigner.Data) {
		return ErrMessageHasBeenProcessed
	}

	return nil
}

// Delivers observed Txs.
func (p *Processor) deliverTxIn(ctx sdk.Context, msgWithSigner *types.TxInWithSigner) ([]byte, error) {
	msg := msgWithSigner.Data

	if p.keeper.IsTxInExisted(ctx, msg) {
		// The tx has been processed before.
		return nil, nil
	}

	log.Info("Deliverying TxIn....")

	// Save this to KVStore & private db.
	p.keeper.SaveTxIn(ctx, msg)
	p.privateDb.SaveTxIn(msg)

	// Creates and broadcast TxOuts. This has to be deterministic based on all the data that the
	// processor has.
	txOutWithSigners := p.txOutputProducer.GetTxOuts(ctx, ctx.BlockHeight(), msg)

	// Save this TxOut to database
	log.Verbose("len(txOut) = ", len(txOutWithSigners))
	if len(txOutWithSigners) > 0 {
		txOuts := make([]*types.TxOut, len(txOutWithSigners))
		for i, outWithSigner := range txOutWithSigners {
			txOut := outWithSigner.Data
			txOuts[i] = txOut

			// We only save txOut to privateDb instead of keeper since it's not confirmed by everyone yet
			p.privateDb.SaveTxOut(txOut)

			// If this is a txOut deployment, mark the contract as being deployed.
			if txOut.TxType == types.TxOutType_CONTRACT_DEPLOYMENT {
				p.keeper.UpdateContractsStatus(ctx, txOut.OutChain, txOut.ContractHash, string(types.TxOutStatusSigning))
			}
		}
	}

	// If this node is not catching up, broadcast the tx.
	if !p.globalData.IsCatchingUp() && len(txOutWithSigners) > 0 {
		log.Info("Broadcasting txout....")

		// Creates TxOut. TODO: Only do this for top validator nodes.
		for _, msg := range txOutWithSigners {
			go func(m *types.TxOutWithSigner) {
				if err := p.txSubmit.SubmitMessage(m); err != nil {
					return
				}
			}(msg)
		}
	}

	return nil, nil
}
