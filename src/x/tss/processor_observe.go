package tss

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	eTypes "github.com/sisu-network/deyes/types"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"
)

// Processed list of transactions sent from deyes to Sisu api server.
func (p *Processor) ProcessObservedTxs(txs *eTypes.Txs) {
	// Create ObservedTx messages and broadcast to the Sisu chain.
	// TODO: Avoid sending too many messages. Find a way we can batch all txts together since SubmitTx
	// has 1s delay.

	observedTxs := &tssTypes.ObservedTxs{}
	observedTxs.Txs = make([]*tssTypes.ObservedTx, len(txs.Arr))

	for index, tx := range txs.Arr {
		observedTxs.Txs[index] = &tssTypes.ObservedTx{
			Chain:       txs.Chain,
			TxHash:      tx.Hash,
			BlockHeight: txs.Block,
		}
	}

	// Send to TxSubmitter.
	p.txSubmit.SubmitMessage(observedTxs)
}

// Delivers observed Txs.
func (p *Processor) DeliverObservedTxs(ctx sdk.Context, msg *tssTypes.ObservedTxs) ([]byte, error) {
	// Update the obsevation count for each transaction.
	for _, tx := range msg.Txs {
		p.keeper.UpdateObservedTxCount(ctx, tx, msg.Signer)
	}

	return nil, nil
}
