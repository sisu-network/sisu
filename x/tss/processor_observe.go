package tss

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	eyesTypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/sisu/utils"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"
)

// Processed list of transactions sent from deyes to Sisu api server.
func (p *Processor) OnObservedTxs(txs *eyesTypes.Txs) {
	// Create ObservedTx messages and broadcast to the Sisu chain.
	for _, tx := range txs.Arr {
		hash := utils.GetObservedTxHash(txs.Block, txs.Chain, tx.Serialized)

		arr := make([]*tssTypes.ObservedTx, 1)
		arr[0] = &tssTypes.ObservedTx{
			Chain:       txs.Chain,
			TxHash:      hash,
			BlockHeight: txs.Block,
			Serialized:  tx.Serialized,
		}

		observedTxs := tssTypes.NewObservedTxs(p.appKeys.GetSignerAddress().String(), arr)

		// Send to TxSubmitter. For now, we only want to include 1 observed tx per 1 Cosmos tx.
		go p.txSubmit.SubmitMessage(observedTxs)
	}
}

func (p *Processor) CheckObservedTxs(ctx sdk.Context, msgs *tssTypes.ObservedTxs) error {
	// TODO: implement this. Compare this observed txs with what we have in database.
	return nil
}

// Delivers observed Txs.
func (p *Processor) DeliverObservedTxs(ctx sdk.Context, msg *tssTypes.ObservedTxs) ([]byte, error) {
	// Update the obsevation count for each transaction.
	for _, tx := range msg.Txs {
		if p.keeper.GetObservedTx(ctx, tx.Chain, tx.BlockHeight, tx.TxHash) != nil {
			utils.LogVerbose("This tx has been included in Sisu block: ", tx.Chain, tx.BlockHeight, tx.TxHash)
			continue
		}

		// Save this to KV store.
		p.keeper.SaveObservedTx(ctx, tx)

		// Save this to our local storage in case we have not seen it.
		p.createAndBroadcastTxOuts(ctx, tx)
	}

	return nil, nil
}
