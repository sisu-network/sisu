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
	utils.LogVerbose("There is a new list of txs from deyes, len =", len(txs.Arr))

	// Create ObservedTx messages and broadcast to the Sisu chain.
	for _, tx := range txs.Arr {
		// 1. Check if this tx is from one of our key. If it is, update the status of TxOut to confirmed.
		if p.db.IsChainKeyAddress(txs.Chain, tx.From) {
			p.confirmTx(tx, txs.Chain)
		} else if tx.To != "" {
			// 2. This is a transaction to our key account or one of our contracts. Create a message to
			// indicate that we have observed this transction and broadcast it to cosmos chain.
			hash := utils.GetObservedTxHash(txs.Block, txs.Chain, tx.Serialized)

			msgTx := &tssTypes.ObservedTx{
				Chain:       txs.Chain,
				TxHash:      hash,
				BlockHeight: txs.Block,
				Serialized:  tx.Serialized,
			}

			observedTxs := tssTypes.NewObservedTxs(p.appKeys.GetSignerAddress().String(), msgTx)
			go p.txSubmit.SubmitMessage(observedTxs)
		}
	}
}

func (p *Processor) CheckObservedTxs(ctx sdk.Context, msgs *tssTypes.ObservedTx) error {
	// TODO: implement this. Compare this observed txs with what we have in database.
	return nil
}

// Delivers observed Txs.
func (p *Processor) DeliverObservedTxs(ctx sdk.Context, tx *tssTypes.ObservedTx) ([]byte, error) {
	// Update the obsevation count for each transaction.
	if p.keeper.GetObservedTx(ctx, tx.Chain, tx.BlockHeight, tx.TxHash) != nil {
		utils.LogVerbose("This tx has been included in Sisu block: ", tx.Chain, tx.BlockHeight, tx.TxHash)
		return nil, fmt.Errorf("this observed tx has been added into our block")
	}

	// Save this to KV store.
	p.keeper.SaveObservedTx(ctx, tx)

	// Save this to our local storage in case we have not seen it.
	p.createAndBroadcastTxOuts(ctx, tx)

	return nil, nil
}

func (p *Processor) confirmTx(tx *eyesTypes.Tx, chain string) {
	utils.LogVerbose("This is a transaction from us. We need to confirm it. Chain =", chain)

	txHash := utils.KeccakHash32(string(tx.Serialized))
	p.db.UpdateTxOutStatus(chain, txHash, types.TxOutStatusConfirmed)

	// If this is a contract deployment, mark the contract as deployed.
	if utils.IsETHBasedChain(chain) && tx.To == "" {
		utils.LogInfo("This is a tx deployment")
		txOut := p.db.GetTxOutWithHashedSig(chain, txHash)

		if txOut != nil {
			utils.LogInfo("Updating contract status. Contract hash = ", txOut.ContractHash)
			p.db.UpdateContractsStatus([]*tssTypes.ContractEntity{
				{
					Chain: chain,
					Hash:  txOut.ContractHash,
				},
			}, tssTypes.ContractStateDeployed)
		}
	}
}
