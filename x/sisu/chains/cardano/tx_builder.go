package cardano

import (
	"github.com/echovl/cardano-go"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/external"
)

// BuildTx contructs a cardano transaction that sends from sender address to receive address.
func BuildTx(deyesClient external.DeyesClient, chain string, sender cardano.Address, receivers []cardano.Address,
	amounts []*cardano.Value, metadata cardano.Metadata, utxos []cardano.UTxO, maxBlock uint64) (*cardano.Tx, error) {
	// Calculate if the account has enough balance
	balance, err := deyesClient.CardanoBalance(chain, sender.String(), int64(maxBlock))
	if err != nil {
		return nil, err
	}

	total := cardano.NewValue(cardano.Coin(0))
	for _, amount := range amounts {
		total = total.Add(amount)
	}
	if cmp := balance.Cmp(total); cmp == -1 || cmp == 2 {
		return nil, NewNotEnoughBalanceErr(total, balance)
	}

	pparams, err := deyesClient.CardanoProtocolParams(chain)
	if err != nil {
		return nil, err
	}

	builder := cardano.NewTxBuilderV2(pparams)

	// For multi-asset utxo, required minimum <1 ADA + additional fee>
	// Details: https://github.com/input-output-hk/cardano-ledger/blob/master/doc/explanations/min-utxo-alonzo.rst#example-min-ada-value-calculations-and-current-constants
	// txOut := &cardano.TxOutput{Address: receiver, Amount: amount}
	// minUTXO := builder.MinCoinsForTxOut(txOut)
	total.Coin = cardano.Coin(0)
	for _, amount := range amounts {
		var minFee uint64
		if amount.MultiAsset == nil {
			minFee = utils.MaxUint64(1_000_000, uint64(amount.Coin))
		} else {
			minFee = utils.MaxUint64(1_600_000, uint64(amount.Coin))
		}
		total.Coin = total.Coin + cardano.Coin(minFee)
	}
	// Additional amount to transfer remaining asset back to Sisu's wallet
	total.Coin = total.Coin + cardano.Coin(1_600_000)

	// Find utxos that cover the amount to transfer
	pickedUtxos := []cardano.UTxO{}
	// Pick at least <MinUTXO * 2> lovelace because we will produce at least 2 new utxos which contains multi-asset:
	// 1. Transfer coin + multi-asset for user
	// 2. Transfer remain coin + multi-asset for Cardano gateway (change address)
	log.Debug("Target utxo balance = ", total.Coin, " ", total.MultiAsset.String())
	pickedUtxosAmount := cardano.NewValue(0)
	ok := false
	for _, utxo := range utxos {
		if pickedUtxosAmount.Cmp(total) == 1 {
			ok = true
			break
		}
		pickedUtxos = append(pickedUtxos, utxo)
		pickedUtxosAmount = pickedUtxosAmount.Add(utxo.Amount)
	}
	if pickedUtxosAmount.Cmp(total) == 1 {
		ok = true
	}

	if !ok {
		return nil, common.InsufficientFundErr
	}

	log.Debug("picked utxo: ")
	for _, p := range pickedUtxos {
		log.Debug("txHash = ", p.TxHash.String(), ", utxo index = ", p.Index, ", coin amount = ", p.Amount.Coin)
	}

	for _, utxo := range pickedUtxos {
		builder.AddInputs(&cardano.TxInput{TxHash: utxo.TxHash, Index: utxo.Index, Amount: utxo.Amount})
	}
	for i, amount := range amounts {
		builder.AddOutputs(&cardano.TxOutput{Address: receivers[i], Amount: amount})
	}

	if len(metadata) > 0 {
		builder.AddAuxiliaryData(&cardano.AuxiliaryData{Metadata: metadata})
	}

	// TODO: Use data from context to make sure that the slot is deterministic.
	tip, err := deyesClient.CardanoTip(chain, maxBlock)
	if err != nil {
		return nil, err
	}

	builder.SetTTL(tip.Slot + 1200)
	changeAddress := sender
	builder.AddChangeIfNeeded(changeAddress)
	builder.Sign([]byte{}) // Use zoombie private key as the key holder.

	tx, err := builder.Build2()
	if err != nil {
		return nil, err
	}
	return tx, nil
}
