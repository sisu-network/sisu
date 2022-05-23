package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxIn struct {
	pmm              PostedMessageManager
	keeper           keeper.Keeper
	txOutputProducer TxOutputProducer
	globalData       common.GlobalData
	txSubmit         common.TxSubmit
	txTracker        TxTracker
}

func NewHandlerTxIn(mc ManagerContainer) *HandlerTxIn {
	return &HandlerTxIn{
		keeper:           mc.Keeper(),
		pmm:              mc.PostedMessageManager(),
		txOutputProducer: mc.TxOutProducer(),
		globalData:       mc.GlobalData(),
		txSubmit:         mc.TxSubmit(),
		txTracker:        mc.TxTracker(),
	}
}

func (h *HandlerTxIn) DeliverMsg(ctx sdk.Context, signerMsg *types.TxInWithSigner) (*sdk.Result, error) {
	rcHash, _, err := keeper.GetTxRecordHash(signerMsg)
	if err != nil {
		return &sdk.Result{}, err
	}

	if h.keeper.IsTxRecordProcessed(ctx, rcHash) {
		return &sdk.Result{}, nil
	}

	log.Debug("IncSlashToken in HandlerTxIn")
	if err := h.keeper.IncSlashToken(ctx, types.ObserveSlashPoint, signerMsg.GetSender()); err != nil {
		return &sdk.Result{}, nil
	}
	log.Debug("inc slash token successfully for sender: ", signerMsg.GetSender().String())

	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doTxIn(ctx, signerMsg)
		if err != nil {
			return &sdk.Result{}, err
		}

		h.keeper.ProcessTxRecord(ctx, hash)

		voters := h.keeper.GetVotersInAccAddress(ctx, hash)
		if err := h.keeper.DecSlashToken(ctx, types.ObserveSlashPoint, voters...); err != nil {
			return &sdk.Result{}, err
		}

		log.Debug("dec slash token successfully")
		for _, v := range voters {
			slash, err := h.keeper.GetSlashToken(ctx, v)
			if err != nil {
				log.Error(err)
				continue
			}

			log.Debugf("address %s has slash %d", v.String(), slash)
		}

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

// Delivers observed Txs.
func (h *HandlerTxIn) doTxIn(ctx sdk.Context, msgWithSigner *types.TxInWithSigner) ([]byte, error) {
	msg := msgWithSigner.Data

	log.Info("Deliverying TxIn, hash = ", msg.TxHash, " on chain ", msg.Chain)

	// Save this to KVStore & private db.
	h.keeper.SaveTxIn(ctx, msg)

	// Creates and broadcast TxOuts. This has to be deterministic based on all the data that the
	// processor has.
	txOutWithSigners := h.txOutputProducer.GetTxOuts(ctx, ctx.BlockHeight(), msg)

	// Save this TxOut to database
	log.Verbose("len(txOut) = ", len(txOutWithSigners))
	if len(txOutWithSigners) > 0 {
		txOuts := make([]*types.TxOut, len(txOutWithSigners))
		for i, outWithSigner := range txOutWithSigners {
			txOut := outWithSigner.Data
			txOuts[i] = txOut

			// If this is a txOut deployment, mark the contract as being deployed.
			if txOut.TxType == types.TxOutType_CONTRACT_DEPLOYMENT {
				h.keeper.UpdateContractsStatus(ctx, txOut.OutChain, txOut.ContractHash, string(types.TxOutStatusSigning))
			}
		}
	}

	// If this node is not catching up, broadcast the tx.
	if !h.globalData.IsCatchingUp() && len(txOutWithSigners) > 0 {
		log.Info("Broadcasting txout....")

		// Creates TxOut. TODO: Only do this for top validator nodes.
		for _, txOutWithSigner := range txOutWithSigners {
			h.txSubmit.SubmitMessageAsync(txOutWithSigner)

			// Track the txout
			h.txTracker.AddTransaction(
				txOutWithSigner.Data,
				msg,
			)
		}
	}

	return nil, nil
}
