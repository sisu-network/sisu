package sisu

import (
	"encoding/json"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mr-tron/base58"
	eyessolanatypes "github.com/sisu-network/deyes/chains/solana/types"
	etypes "github.com/sisu-network/deyes/types"
	eyesTypes "github.com/sisu-network/deyes/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	scardano "github.com/sisu-network/sisu/x/sisu/chains/cardano"
	"github.com/sisu-network/sisu/x/sisu/chains/eth"
	solanatypes "github.com/sisu-network/sisu/x/sisu/chains/solana/types"
	"github.com/sisu-network/sisu/x/sisu/types"
)

// Processed list of transactions sent from deyes to Sisu api server.
// TODO: handle error correctly
func (a *ApiHandler) OnTxIns(txs *eyesTypes.Txs) error {
	log.Verbose("There is a new list of txs from deyes, len =", len(txs.Arr))

	transferRequests := &types.Transfers{
		Transfers: make([]*types.Transfer, 0),
	}

	ctx := a.globalData.GetReadOnlyContext()

	// Create TxIn messages and broadcast to the Sisu chain.
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

		transfers, err := a.parseDeyesTx(ctx, txs.Chain, tx)
		if err != nil {
			log.Error("Faield to parse transfer, err = ", err)
			continue
		}

		// Assign the id for all transfers
		for _, transfer := range transfers {
			transfer.Id = types.GetTransferId(transfer.FromChain, transfer.FromHash)
		}

		log.Verbose("Len(transfers) = ", len(transfers), " on chain ", txs.Chain)
		if transfers != nil {
			transferRequests.Transfers = append(transferRequests.Transfers, transfers...)
		}
	}

	if len(transferRequests.Transfers) > 0 {
		msg := types.NewTransfersMsg(a.appKeys.GetSignerAddress().String(), transferRequests)
		a.txSubmit.SubmitMessageAsync(msg)
	}

	if libchain.IsCardanoChain(txs.Chain) {
		log.Verbose("Updating block height for cardano")
		// Broadcast blockheight update
		msg := types.NewBlockHeightMsg(a.appKeys.GetSignerAddress().String(), &types.BlockHeight{
			Chain:  txs.Chain,
			Height: txs.Block,
			Hash:   txs.BlockHash,
		})
		a.txSubmit.SubmitMessageAsync(msg)
	}

	return nil
}

func (a *ApiHandler) parseDeyesTx(ctx sdk.Context, chain string, tx *eyesTypes.Tx) ([]*types.Transfer, error) {
	if libchain.IsETHBasedChain(chain) {
		parseResult := eth.ParseVaultTx(ctx, a.keeper, chain, tx)
		if parseResult.Error != nil {
			return nil, parseResult.Error
		}

		if parseResult.TransferOuts != nil {
			return parseResult.TransferOuts, nil
		}

		return []*types.Transfer{}, nil
	}

	if libchain.IsCardanoChain(chain) {
		return a.parseCardanoTx(ctx, chain, tx)
	}

	return nil, fmt.Errorf("Unknown chain %s", chain)
}

func (a *ApiHandler) parseCardanoTx(ctx sdk.Context, chain string, tx *eyesTypes.Tx) ([]*types.Transfer, error) {
	ret := make([]*types.Transfer, 0)
	cardanoTx := &etypes.CardanoTransactionUtxo{}
	err := json.Unmarshal(tx.Serialized, cardanoTx)
	if err != nil {
		return nil, err
	}

	if cardanoTx.Metadata != nil {
		nativeTransfer := cardanoTx.Metadata.NativeAda != 0
		log.Verbose("cardanoTx.Amount = ", cardanoTx.Amount)

		// Convert from ADA unit (10^6) to our standard unit (10^18)
		for _, amount := range cardanoTx.Amount {
			quantity, ok := new(big.Int).SetString(amount.Quantity, 10)
			if !ok {
				log.Error("Failed to get amount quantity in cardano tx")
				continue
			}
			quantity = utils.LovelaceToWei(quantity)

			// Remove the word wrap
			tokenUnit := amount.Unit
			if tokenUnit != "lovelace" {
				token := scardano.GetTokenFromCardanoAsset(ctx, a.keeper, tokenUnit, chain)
				if token == nil {
					log.Error("Failed to find token with id: ", tokenUnit)
					continue
				}
				tokenUnit = token.Id
			} else {
				if !nativeTransfer {
					// This ADA is for transaction transfer fee. It is not meant to be transfered.
					continue
				}
				tokenUnit = "ADA"
			}

			log.Verbose("tokenUnit = ", tokenUnit, " quantity = ", quantity)
			log.Verbose("cardanoTx.Metadata = ", cardanoTx.Metadata)

			ret = append(ret, &types.Transfer{
				FromHash:    cardanoTx.Hash,
				Token:       tokenUnit,
				Amount:      quantity.String(),
				ToChain:     cardanoTx.Metadata.Chain,
				ToRecipient: cardanoTx.Metadata.Recipient,
			})
		}
	}

	return ret, nil
}

func (a *ApiHandler) parseSolanaTx(ctx sdk.Context, chain string, tx *eyesTypes.Tx) ([]*types.Transfer, error) {
	ret := make([]*types.Transfer, 0)

	outerTx := new(eyessolanatypes.Transaction)
	err := json.Unmarshal(tx.Serialized, outerTx)
	if err != nil {
		return nil, err
	}

	cfg := a.mc.Config()
	accounts := outerTx.TransactionInner.Message.AccountKeys

	// Check that there is at least one instruction sent to the program id
	for _, ix := range outerTx.TransactionInner.Message.Instructions {
		if accounts[ix.ProgramIdIndex] != cfg.Solana.BridgeProgramId {
			continue
		}

		// Decode the instruction
		bytesArr, err := base58.Decode(ix.Data)
		if err != nil {
			return nil, err
		}

		if len(bytesArr) == 0 {
			return nil, fmt.Errorf("Data is empty")
		}

		ix := new(solanatypes.TransferOutInstruction)
		err = ix.Deserialize(bytesArr)
		if err != nil {
			return nil, err
		}

		transferData := ix.Data

		switch ix.Instruction {
		case solanatypes.TranserOut:
			ret = append(ret, &types.Transfer{
				FromHash:    outerTx.TransactionInner.Signatures[0],
				Token:       transferData.TokenAddress,
				Amount:      transferData.Amount.String(),
				ToChain:     libchain.GetChainNameFromInt(big.NewInt(int64(transferData.ChainId))),
				ToRecipient: transferData.Recipient,
			})
		}
	}

	return ret, nil
}
