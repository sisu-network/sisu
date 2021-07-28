package tss

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	deTypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/sisu/utils"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"
)

// Processed list of transactions sent from deyes to Sisu api server.
func (p *Processor) ProcessObservedTxs(txs *deTypes.Txs) {
	// Create ObservedTx messages and broadcast to the Sisu chain.
	// TODO: Avoid sending too many messages. Find a way we can batch all txts together since SubmitTx
	// has 1s delay.

	arr := make([]*tssTypes.ObservedTx, len(txs.Arr))

	for index, tx := range txs.Arr {
		arr[index] = &tssTypes.ObservedTx{
			Chain:       txs.Chain,
			TxHash:      tx.Hash,
			BlockHeight: txs.Block,
		}
	}

	observedTxs := tssTypes.NewObservedTxs(p.appKeys.GetSignerAddress().String(), arr)

	// Send to TxSubmitter.
	p.txSubmit.SubmitMessage(observedTxs)

	// Save all txs into database. We save this to local database instead of kvstore since this is a
	// set of txs that observed by this node only (not all the nodes). KVStore is used to store state
	// that have been agreed by all nodes in the network.
	p.storage.SaveTxs(txs)
}

// Delivers observed Txs.
func (p *Processor) DeliverObservedTxs(ctx sdk.Context, msg *tssTypes.ObservedTxs) ([]byte, error) {
	// Update the obsevation count for each transaction.
	utils.LogVerbose("Deliver observed txs. Len = ", msg.Txs)

	for _, tx := range msg.Txs {
		size, err := p.keeper.UpdateObservedTxCount(ctx, tx, msg.Signer)
		if err != nil {
			continue
		}

		if size >= (p.globalData.ValidatorSize()+2)/3 && !p.keeper.IsObservedTxPendingOrProcessed(ctx, tx) {
			utils.LogVerbose("Adding observed tx to pending")
			// Majority has been meet and the tx has not been processed yet. Put it in the pending queue.
			// They will be processed in the next block.
			p.keeper.AddObservedTxToPending(ctx, tx)
		}
	}

	return nil, nil
}
