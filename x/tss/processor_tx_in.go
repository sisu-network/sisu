package tss

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	eyesTypes "github.com/sisu-network/deyes/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"
)

// Processed list of transactions sent from deyes to Sisu api server.
// TODO: handle error correctly
func (p *Processor) OnTxIns(txs *eyesTypes.Txs) error {
	log.Verbose("There is a new list of txs from deyes, len =", len(txs.Arr))

	// Create TxIn messages and broadcast to the Sisu chain.
	for _, tx := range txs.Arr {
		// 1. Check if this tx is from one of our key. If it is, update the status of TxOut to confirmed.
		if p.db.IsChainKeyAddress(libchain.KEY_TYPE_ECDSA, tx.From) {
			return p.confirmTx(tx, txs.Chain)
		} else if len(tx.To) > 0 {
			// 2. This is a transaction to our key account or one of our contracts. Create a message to
			// indicate that we have observed this transaction and broadcast it to cosmos chain.
			// TODO: handle error correctly
			hash := utils.GetTxInHash(txs.Block, txs.Chain, tx.Serialized)
			txIn := types.NewTxInWithSigner(
				p.appKeys.GetSignerAddress().String(),
				txs.Chain,
				hash,
				txs.Block,
				tx.Serialized,
			)

			// Save tx in into db
			p.db.InsertTxIn(txIn.Data)

			go func(tx *types.TxInWithSigner) {
				if err := p.txSubmit.SubmitMessage(tx); err != nil {
					return
				}
			}(txIn)
		}
	}

	return nil
}

func (p *Processor) checkTxIn(ctx sdk.Context, msgWithSigner *types.TxInWithSigner) error {
	// Make sure we should have seen this TxIn in our table.
	if p.db.IsTxInExisted(msgWithSigner.Data) {
		return nil
	}

	return ErrCannotFindMessage
}

// Delivers observed Txs.
func (p *Processor) deliverTxIn(ctx sdk.Context, msgWithSigner *types.TxInWithSigner) ([]byte, error) {
	msg := msgWithSigner.Data

	if p.keeper.IsTxInExisted(ctx, msg) {
		// The tx has been processed before.
		return nil, nil
	}

	// Insert this TxIn into db table.
	p.db.InsertTxIn(msg)

	// Save this to KVStore
	p.keeper.SaveTxIn(ctx, msg)

	// Creates and broadcast TxOuts. This has to be deterministic based on all the data that the
	// processor has.
	txOutWithSigners := p.txOutputProducer.GetTxOuts(ctx, ctx.BlockHeight(), msg)

	// Save this TxOut to database
	log.Verbose("len(txOut) = ", len(txOutWithSigners))
	if len(txOutWithSigners) > 0 {
		txOuts := make([]*types.TxOut, len(txOutWithSigners))
		for i, outWithSigner := range txOutWithSigners {
			txOuts[i] = outWithSigner.Data
		}
		p.db.InsertTxOuts(txOuts)
	}

	// If this node is not catching up, broadcast the tx.
	if !p.globalData.IsCatchingUp() {
		log.Info("Broadcasting txout....")

		// Creates TxOut. TODO: Only do this for top validator nodes.
		for _, msg := range txOutWithSigners {
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
			string(txOut.OutBytes),
			tssTypes.TxOutStatusDeployedToBlockchain,
			false); err != nil {
			return err
		}
	}

	return nil
}
