package tss

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	eyesTypes "github.com/sisu-network/deyes/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
)

// Processed list of transactions sent from deyes to Sisu api server.
// TODO: handle error correctly
func (p *Processor) OnTxIns(txs *eyesTypes.Txs) error {
	log.Verbose("There is a new list of txs from deyes, len =", len(txs.Arr))

	// Create TxIn messages and broadcast to the Sisu chain.
	for _, tx := range txs.Arr {
		// 1. Check if this tx is from one of our key. If it is, update the status of TxOut to confirmed.
		if p.privateDb.IsKeygenAddress(libchain.KEY_TYPE_ECDSA, tx.From) {
			return p.confirmTx(tx, txs.Chain, txs.Block)
		} else if len(tx.To) > 0 {
			// 2. This is a transaction to our key account or one of our contracts. Create a message to
			// indicate that we have observed this transaction and broadcast it to cosmos chain.
			// TODO: handle error correctly
			hash := utils.GetTxInHash(txs.Block, txs.Chain, tx.Serialized)
			signerMsg := types.NewTxInWithSigner(
				p.appKeys.GetSignerAddress().String(),
				txs.Chain,
				hash,
				txs.Block,
				tx.Serialized,
			)

			// Save tx in into db
			p.privateDb.SaveTxIn(signerMsg.Data)

			go func(tx *types.TxInWithSigner) {
				if err := p.txSubmit.SubmitMessage(tx); err != nil {
					return
				}
			}(signerMsg)
		}
	}

	return nil
}

// confirmTx confirms that a tx has been included in a block on the blockchain.
func (p *Processor) confirmTx(tx *eyesTypes.Tx, chain string, blockHeight int64) error {
	log.Verbose("This is a transaction from us. We need to confirm it. Chain = ", chain)

	p.privateDb.PrintStoreKeys("txOut")
	txOut := p.privateDb.GetTxOutFromSigHash(chain, tx.Hash)
	if txOut == nil {
		// TODO: Add unconfirmed tx model
		log.Verbose("cannot find txOut with full signature hash: ", tx.Hash)
		return nil
	}

	log.Info("confirming tx: chain, hash, type = ", chain, tx.Hash, txOut.TxType)

	contractAddress := ""
	if txOut.TxType == types.TxOutType_CONTRACT_DEPLOYMENT && libchain.IsETHBasedChain(chain) {
		ethTx := &ethTypes.Transaction{}
		err := ethTx.UnmarshalBinary(tx.Serialized)
		if err != nil {
			log.Error("cannot unmarshal eth transaction, err = ", err)
			return err
		}

		contractAddress = crypto.CreateAddress(common.HexToAddress(tx.From), ethTx.Nonce()).String()
		log.Info("contractAddress = ", contractAddress)
	}

	confirmMsg := types.NewTxOutConfirmWithSigner(
		p.appKeys.GetSignerAddress().String(),
		txOut.TxType,
		txOut.OutChain,
		txOut.OutHash,
		blockHeight,
		contractAddress,
	)

	// Save this into db
	p.privateDb.SaveTxOutConfirm(confirmMsg.Data)

	go func() {
		p.txSubmit.SubmitMessage(confirmMsg)
	}()

	return nil
}

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
