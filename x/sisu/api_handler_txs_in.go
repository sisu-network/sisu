package sisu

import (
	"fmt"

	eyesTypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"
)

// Processed list of transactions sent from deyes to Sisu api server.
// TODO: handle error correctly
func (a *ApiHandler) OnTxIns(txs *eyesTypes.Txs) error {
	log.Verbose("There is a new list of txs from deyes, len =", len(txs.Arr))

	ctx := a.globalData.GetReadOnlyContext()
	// Make sure that this chain is supported by Sisu
	params := a.keeper.GetParams(ctx)
	if !utils.IsChainSupported(params.SupportedChains, txs.Chain) {
		return fmt.Errorf("Unsupported chain: %s", txs.Chain)
	}

	bridge := a.bridgeManager.GetBridge(ctx, txs.Chain)
	if bridge == nil {
		return fmt.Errorf("cannot find bridge for chain %s", txs.Chain)
	}

	for _, tx := range txs.Arr {
		if !tx.Success {
			log.Verbose("Failed incoming transaction (not our fault), hash = ", tx.Hash, ", chain = ", txs.Chain)
			continue
		}

		// Check if this is a transaction from our sisu. If true, ignore it.
		sisu := a.keeper.GetMpcAddress(ctx, txs.Chain)
		if sisu == tx.From {
			log.Verbosef("This is a transaction sent from our sisu account %s on chain %s, ignore",
				sisu, txs.Chain)
			continue
		}

		transfers, err := bridge.ParseIncomingTx(ctx, txs.Chain, tx.Serialized)
		if err != nil {
			log.Errorf("Failed to parse transfer on chain %s, hash = %s", txs.Chain, tx.Hash)
			continue
		}

		log.Verbosef("len(transfers) = %d", len(transfers))

		if len(transfers) == 0 {
			// There is no transfer request from this transaction. Just ignore. In the future, we have
			// to check other type of transaction sent to the vault/gateway.
			continue
		}

		// Just send a thin tx in.
		txInId := fmt.Sprintf("%s__%s", txs.Chain, tx.Hash)
		// Parse the transfers
		msg := types.NewTxInMsg(
			a.appKeys.GetSignerAddress().String(),
			&types.TxIn{
				Id:        txInId,
				FromChain: txs.Chain,
				Serialize: tx.Serialized,
				Transfers: transfers,
			},
		)
		a.txSubmit.SubmitMessageAsync(msg)
	}

	return nil
}
