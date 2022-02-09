package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	eyesTypes "github.com/sisu-network/deyes/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"
)

// Processed list of transactions sent from deyes to Sisu api server.
// TODO: handle error correctly
func (p *Processor) OnTxIns(txs *eyesTypes.Txs) error {
	log.Verbose("There is a new list of txs from deyes, len =", len(txs.Arr))

	// Create TxIn messages and broadcast to the Sisu chain.
	for _, tx := range txs.Arr {
		// 1. Check if this tx is from one of our key. If it is, update the status of TxOut to confirmed.
		if p.publicDb.IsKeygenAddress(libchain.KEY_TYPE_ECDSA, tx.From) {
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

			p.txSubmit.SubmitMessageAsync(signerMsg)
		}
	}

	return nil
}

// confirmTx confirms that a tx has been included in a block on the blockchain.
func (p *Processor) confirmTx(tx *eyesTypes.Tx, chain string, blockHeight int64) error {
	log.Verbose("This is a transaction from us. We need to confirm it. Chain = ", chain)

	// The txOutSig is in private db while txOut should come from common db.
	txOutSig := p.privateDb.GetTxOutSig(chain, tx.Hash)
	if txOutSig == nil {
		// TODO: Add this to pending tx to confirm.
		log.Verbose("cannot find txOutSig with full signature hash: ", tx.Hash)
		return nil
	}

	txOut := p.publicDb.GetTxOut(chain, txOutSig.HashNoSig)
	if txOut == nil {
		log.Verbose("cannot find txOut with hash (with no sig): ", txOutSig.HashNoSig)
		return nil
	}

	log.Info("confirming tx: chain, hash, type = ", chain, " ", tx.Hash, " ", txOut.TxType)

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

		txConfirm := &types.TxOutContractConfirm{
			OutChain:        txOut.OutChain,
			OutHash:         txOut.OutHash,
			BlockHeight:     blockHeight,
			ContractAddress: contractAddress,
		}

		msg := types.NewTxOutContractConfirmWithSigner(
			p.appKeys.GetSignerAddress().String(),
			txConfirm,
		)
		p.txSubmit.SubmitMessageAsync(msg)
	}

	// We can assume that other tx transactions will succeed in majority of the time. Instead
	// broadcasting the tx confirmation to Sisu blockchain, we should only record missing or failed
	// transaction.
	// We only confirm if the tx out is a contract deployment to save the smart contract address.
	// TODO: Implement missing/ failed message and broadcast that to everyone after we have not seen
	// a tx for some blocks.

	return nil
}

func (p *Processor) deliverTxIn(ctx sdk.Context, signerMsg *types.TxInWithSigner) ([]byte, error) {
	if process, hash := p.shouldProcessMsg(ctx, signerMsg); process {
		p.doTxIn(ctx, signerMsg)
		p.publicDb.ProcessTxRecord(hash)
	}

	return nil, nil
}

// Delivers observed Txs.
func (p *Processor) doTxIn(ctx sdk.Context, msgWithSigner *types.TxInWithSigner) ([]byte, error) {
	msg := msgWithSigner.Data

	log.Info("Deliverying TxIn....")

	// Save this to KVStore & private db.
	p.publicDb.SaveTxIn(msg)

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

			// If this is a txOut deployment, mark the contract as being deployed.
			if txOut.TxType == types.TxOutType_CONTRACT_DEPLOYMENT {
				p.publicDb.UpdateContractsStatus(txOut.OutChain, txOut.ContractHash, string(types.TxOutStatusSigning))
			}
		}
	}

	// If this node is not catching up, broadcast the tx.
	if !p.globalData.IsCatchingUp() && len(txOutWithSigners) > 0 {
		log.Info("Broadcasting txout....")

		// Creates TxOut. TODO: Only do this for top validator nodes.
		for _, msg := range txOutWithSigners {
			p.txSubmit.SubmitMessageAsync(msg)
		}
	}

	return nil, nil
}
