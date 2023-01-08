package sisu

import (
	"encoding/hex"
	"fmt"
	"strings"

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

	vals := a.valManager.GetValidators(ctx)
	for _, tx := range txs.Arr {
		// Just send a thin tx in.
		txInId := fmt.Sprintf("%s__%s", txs.Chain, tx.Hash)
		msg := types.NewTxInMsg(a.appKeys.GetSignerAddress().String(), &types.TxIn{Id: txInId})
		a.txSubmit.SubmitMessageAsync(msg)

		// Check if this node is assigned to confirm the next tx in.
		sortedVals := utils.GetSortedValidators(txInId, vals)
		if strings.EqualFold(a.appKeys.GetSignerAddress().String(), sortedVals[0].AccAddress) {
			// Parse the transfers
			transfers, err := bridge.ParseIncomginTx(ctx, txs.Chain, tx)
			if err != nil {
				log.Errorf("Failed to parse transfer on chain %s, hex of the tx's binary = %s",
					txs.Chain, hex.EncodeToString(tx.Serialized))
				continue
			}

			// Send a tx details instead
			msg := types.NewTxInDetailsMsg(
				a.appKeys.GetSignerAddress().String(),
				&types.TxInDetails{
					TxIn: &types.TxIn{
						Id: txInId,
					},
					FromChain: txs.Chain,
					Serialize: tx.Serialized,
					Transfers: transfers,
				},
			)
			a.txSubmit.SubmitMessageAsync(msg)
		}
	}

	// // Create TxIn messages and broadcast to the Sisu chain.
	// for _, tx := range txs.Arr {
	// 	if !tx.Success {
	// 		log.Verbose("Failed incoming transaction (not our fault), hash = ", tx.Hash, ", chain = ", txs.Chain)
	// 		continue
	// 	}

	// 	// Check if this is a transaction from our sisu. If true, ignore it.
	// 	sisu := a.keeper.GetMpcAddress(ctx, txs.Chain)
	// 	if sisu == tx.From {
	// 		log.Verbosef("This is a transaction sent from our sisu account %s on chain %s, ignore",
	// 			sisu, txs.Chain)
	// 		continue
	// 	}

	// 	transfers, err := a.parseDeyesTx(ctx, txs.Chain, tx)
	// 	if err != nil {
	// 		log.Error("Faield to parse transfer, err = ", err)
	// 		continue
	// 	}

	// 	// Assign the id for all transfers
	// 	for _, transfer := range transfers {
	// 		transfer.Id = types.GetTransferId(transfer.FromChain, transfer.FromHash)
	// 	}

	// 	log.Verbose("Len(transfers) = ", len(transfers), " on chain ", txs.Chain)
	// 	if transfers != nil {
	// 		transferRequests.Transfers = append(transferRequests.Transfers, transfers...)
	// 	}
	// }

	// if len(transferRequests.Transfers) > 0 {
	// 	msg := types.NewTransfersMsg(a.appKeys.GetSignerAddress().String(), transferRequests)
	// 	a.txSubmit.SubmitMessageAsync(msg)

	// 	if libchain.IsCardanoChain(txs.Chain) {
	// 		log.Verbose("Updating block height for cardano")
	// 		// Broadcast blockheight update
	// 		msg := types.NewBlockHeightMsg(a.appKeys.GetSignerAddress().String(), &types.BlockHeight{
	// 			Chain:  txs.Chain,
	// 			Height: txs.Block,
	// 			Hash:   txs.BlockHash,
	// 		})
	// 		a.txSubmit.SubmitMessageAsync(msg)
	// 	}

	// 	// Check to see if we need to update the gas price.
	// 	a.updateEthGasPrice(ctx, transferRequests.Transfers)
	// }

	return nil
}
