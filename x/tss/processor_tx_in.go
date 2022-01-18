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

	return nil, nil
}
