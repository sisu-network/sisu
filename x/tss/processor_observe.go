package tss

import (
	"fmt"

	sdk "github.com/sisu-network/cosmos-sdk/types"
	eyesTypes "github.com/sisu-network/deyes/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"
)

// Processed list of transactions sent from deyes to Sisu api server.
// TODO: handle error correctly
func (p *Processor) OnObservedTxs(txs *eyesTypes.Txs) error {
	log.Verbose("There is a new list of txs from deyes, len =", len(txs.Arr))

	// Create ObservedTx messages and broadcast to the Sisu chain.
	for _, tx := range txs.Arr {
		// 1. Check if this tx is from one of our key. If it is, update the status of TxOut to confirmed.
		if p.db.IsChainKeyAddress(libchain.KEY_TYPE_ECDSA, tx.From) {
			return p.confirmTx(tx, txs.Chain)
		} else if len(tx.To) > 0 {
			// 2. This is a transaction to our key account or one of our contracts. Create a message to
			// indicate that we have observed this transaction and broadcast it to cosmos chain.
			hash := utils.GetObservedTxHash(txs.Block, txs.Chain, tx.Serialized)
			observedTxs := tssTypes.NewObservedTxs(
				p.appKeys.GetSignerAddress().String(),
				txs.Chain,
				hash,
				txs.Block,
				tx.Serialized,
			)

			// TODO: handle error correctly
			go func() {
				if err := p.txSubmit.SubmitMessage(observedTxs); err != nil {
					return
				}
			}()
		}
	}

	return nil
}

func (p *Processor) checkObservedTxs(ctx sdk.Context, tx *tssTypes.ObservedTx) error {
	if p.keeper.IsObservedTxExisted(ctx, tx) {
		return ErrMessageHasBeenProcessed
	}
	return nil
}

// Delivers observed Txs.
func (p *Processor) deliverObservedTxs(ctx sdk.Context, tx *tssTypes.ObservedTx) ([]byte, error) {
	if p.keeper.IsObservedTxExisted(ctx, tx) {
		// The tx has been processed before.
		return nil, nil
	}

	// Save this to KVStore
	p.keeper.SaveObservedTx(ctx, tx)

	// Creates and broadcast TxOuts
	txOuts := p.createTxOuts(ctx, tx)

	if !p.globalData.IsCatchingUp() {
		fmt.Println("Broadcasting txout....")

		// Creates TxOut. TODO: Only do this for top validator nodes.
		for _, msg := range txOuts {
			go func(m *tssTypes.TxOutWithSigner) {
				if err := p.txSubmit.SubmitMessage(m); err != nil {
					return
				}

				// p.db.UpdateTxOutStatus(m.OutChain, m.GetHash(), tssTypes.TxOutStatusBroadcasted, false)
			}(msg)
		}
	}

	return nil, nil
}

func (p *Processor) confirmTx(tx *eyesTypes.Tx, chain string) error {
	log.Verbose("This is a transaction from us. We need to confirm it. Chain =", chain)

	txHash := utils.KeccakHash32(string(tx.Serialized))
	if err := p.db.UpdateTxOutStatus(chain, txHash, tssTypes.TxOutStatusConfirmed, true); err != nil {
		return err
	}

	// If this is a contract deployment, mark the contract as deployed.
	if libchain.IsETHBasedChain(chain) && len(tx.To) == 0 {
		log.Info("This is a tx deployment")
		txOut := p.db.GetTxOutWithHash(chain, txHash, true)

		if txOut == nil {
			log.Warn("txOut by txHash", txHash, "is not found")
			return nil
		}

		// log.Info("Updating contract status. Contract hash = ", txOut.ContractHash)
		// if err := p.db.UpdateContractsStatus([]*types.ContractEntity{
		// 	{
		// 		Chain: chain,
		// 		Hash:  txOut.ContractHash,
		// 	},
		// }, tssTypes.ContractStateDeployed); err != nil {
		// 	return err
		// }

		if err := p.db.UpdateTxOutStatus(
			txOut.OutChain,
			txOut.HashWithoutSig,
			tssTypes.TxOutStatusDeployedToBlockchain,
			false); err != nil {
			return err
		}
	}

	return nil
}
