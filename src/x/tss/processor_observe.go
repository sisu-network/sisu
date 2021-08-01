package tss

import (
	"fmt"

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

	for _, tx := range txs.Arr {
		hash, err := utils.GetObservedTxHash(txs.Block, txs.Chain, tx.Serialized)
		if err != nil {
			utils.LogCritical(fmt.Sprintf("Cannot get hash for tx at block %d in chain %s", txs.Block, txs.Chain))
			continue
		}

		// Check local storage to see if this observed tx has been recorded in local storage. This is
		// very different from observed tx in kvstore as the kvstore only stores data from Sisu chain.
		if p.storage.GetObservedTx(txs.Chain, txs.Block, hash) != nil {
			// This tx has been recorded in the Sisu chain. We don't need to broadcast it anymore.
			utils.LogVerbose("This tx has been processed before.")
			continue
		}

		p.storage.SaveObservedTx(txs.Chain, txs.Block, hash, tx.Serialized)

		arr := make([]*tssTypes.ObservedTx, 1)
		arr[0] = &tssTypes.ObservedTx{
			Chain:       txs.Chain,
			TxHash:      hash,
			BlockHeight: txs.Block,
			Serialized:  tx.Serialized,
		}

		observedTxs := tssTypes.NewObservedTxs(p.appKeys.GetSignerAddress().String(), arr)

		// Send to TxSubmitter. For now, we only want to include 1 observed tx per 1 Cosmos tx.
		p.txSubmit.SubmitMessage(observedTxs)
	}
}

func (p *Processor) CheckObservedTxs(ctx sdk.Context, msgs *tssTypes.ObservedTxs) error {
	// Returns true if we have also observed the same tx.
	for _, msg := range msgs.Txs {
		bz := p.storage.GetObservedTx(msg.Chain, msg.BlockHeight, msg.TxHash)
		if bz == nil {
			return fmt.Errorf("We have not seen this transaction yet.")
		}
	}

	return nil
}

// Delivers observed Txs.
func (p *Processor) DeliverObservedTxs(ctx sdk.Context, msg *tssTypes.ObservedTxs) ([]byte, error) {
	// Update the obsevation count for each transaction.
	utils.LogVerbose("Deliver observed txs. Len = ", msg.Txs)

	for _, tx := range msg.Txs {
		if p.keeper.GetObservedTx(ctx, tx.Chain, tx.BlockHeight, tx.TxHash) != nil {
			utils.LogVerbose("This tx has been included in Sisu block: ", tx.Chain, tx.BlockHeight, tx.TxHash)
			continue
		}

		// Save this to KV store.
		p.keeper.SaveObservedTx(ctx, tx)

		// Save this to our local storage in case we have not seen it.
		p.storage.SaveObservedTx(tx.Chain, tx.BlockHeight, tx.TxHash, tx.Serialized)

		// Add this tx to pending list to process at the end of block.
		p.storage.AddPendingTx(tx)
	}

	return nil, nil
}
