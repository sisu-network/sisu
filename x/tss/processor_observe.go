package tss

import (
	"fmt"

	sdk "github.com/sisu-network/cosmos-sdk/types"
	eyesTypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"
)

// Processed list of transactions sent from deyes to Sisu api server.
func (p *Processor) OnObservedTxs(txs *eyesTypes.Txs) {
	// Create ObservedTx messages and broadcast to the Sisu chain.
	for _, tx := range txs.Arr {
		// 1. Check if this tx is from one of our key. If it is, update the status of TxOut to confirmed.
		if p.db.ChainKeyExisted(txs.Chain, tx.From) {
			utils.LogVerbose("This is a transaction from us. We need to confirm it.")

			txHash := utils.KeccakHash32(string(tx.Serialized))
			p.db.UpdateTxOutStatus(txs.Chain, txHash, types.TxOutStatusConfirmed)

			// If this is a contract deployment, mark the contract as deployed.
			fmt.Println("tx.To  = ", tx.To)

			if utils.IsETHBasedChain(txs.Chain) && tx.To == "" {
				txOut := p.db.GetTxOutWithHashedSig(txs.Chain, txHash)

				if txOut != nil {
					utils.LogInfo("Updating contract status. Contract hash = ", txOut.ContractHash)
					p.db.UpdateContractsStatus([]*tssTypes.ContractEntity{
						{
							Chain: txs.Chain,
							Hash:  txOut.ContractHash,
						},
					}, tssTypes.ContractStateDeployed)
				}
			}
		}

		// 2. This is a transaction to our key account or one of our contracts. Create a message to
		// indicate that we have observed this transction and broadcast it to cosmos chain.
		if tx.To != "" {
			hash := utils.GetObservedTxHash(txs.Block, txs.Chain, tx.Serialized)

			arr := make([]*tssTypes.ObservedTx, 1)
			arr[0] = &tssTypes.ObservedTx{
				Chain:       txs.Chain,
				TxHash:      hash,
				BlockHeight: txs.Block,
				Serialized:  tx.Serialized,
			}

			observedTxs := tssTypes.NewObservedTxs(p.appKeys.GetSignerAddress().String(), arr)
			if len(observedTxs.Txs) > 0 {
				// Send to TxSubmitter. For now, we only want to include 1 observed tx per 1 Cosmos tx.
				go p.txSubmit.SubmitMessage(observedTxs)
			}
		}
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
