package cardano

import (
	"fmt"

	"github.com/echovl/cardano-go"
	"github.com/sisu-network/lib/log"
)

func Balance(node cardano.Node, address cardano.Address) (*cardano.Value, error) {
	balance := cardano.NewValue(0)
	utxos, err := findUtxos(node, address)
	if err != nil {
		return nil, err
	}

	for _, utxo := range utxos {
		balance = balance.Add(utxo.Amount)
	}

	return balance, nil
}

func findUtxos(node cardano.Node, address cardano.Address) ([]cardano.UTxO, error) {
	walletUtxos := make([]cardano.UTxO, 0)
	addrUtxos, err := node.UTxOs(address)
	if err != nil {
		return nil, err
	}

	walletUtxos = append(walletUtxos, addrUtxos...)
	return walletUtxos, nil
}

func BuildTx(node cardano.Node, network cardano.Network, sender, receiver cardano.Address,
	amount *cardano.Value) (*cardano.Tx, error) {
	// Calculate if the account has enough balance
	balance, err := Balance(node, sender)
	if err != nil {
		return nil, err
	}

	if cmp := balance.Cmp(amount); cmp == -1 || cmp == 2 {
		return nil, fmt.Errorf("Not enough balance, %v > %v", amount, balance)
	}

	pparams, err := node.ProtocolParams()
	if err != nil {
		return nil, err
	}

	builder := cardano.NewTxBuilderV2(pparams)

	txOut := &cardano.TxOutput{Address: receiver, Amount: amount}
	// For multi-asset utxo, required minimum <1 ADA + additional fee>
	// Details: https://github.com/input-output-hk/cardano-ledger/blob/master/doc/explanations/min-utxo-alonzo.rst#example-min-ada-value-calculations-and-current-constants
	minUTXO := builder.MinCoinsForTxOut(txOut)
	amount.Coin = minUTXO
	log.Debug("amount.Coin = ", amount.Coin)

	// Find utxos that cover the amount to transfer
	pickedUtxos := []cardano.UTxO{}
	utxos, err := findUtxos(node, sender)
	log.Debug("all utxos: ")
	for _, utxo := range utxos {
		log.Debug("txHash = ", utxo.TxHash.String(), " coin amount = ", utxo.Amount.Coin)
	}
	log.Debug("--------------------------")

	// Pick at least <MinUTXO * 2> lovelace because we will produce at least 2 new utxos which contains multi-asset:
	// 1. Transfer coin + multi-asset for user
	// 2. Transfer remain coin + multi-asset for Cardano gateway (change address)
	targetUtxoBalance := cardano.NewValueWithAssets(amount.Coin*2, amount.MultiAsset)
	log.Debug("Target utxo balance = ", targetUtxoBalance.Coin, targetUtxoBalance.MultiAsset.String())
	pickedUtxosAmount := cardano.NewValue(0)
	for _, utxo := range utxos {
		if pickedUtxosAmount.Cmp(targetUtxoBalance) == 1 {
			break
		}
		pickedUtxos = append(pickedUtxos, utxo)
		pickedUtxosAmount = pickedUtxosAmount.Add(utxo.Amount)
	}

	log.Debug("picked utxo: ")
	for _, p := range pickedUtxos {
		log.Debug("txHash = ", p.TxHash.String(), " coin amount = ", p.Amount.Coin)
	}

	for _, utxo := range pickedUtxos {
		builder.AddInputs(&cardano.TxInput{TxHash: utxo.TxHash, Index: utxo.Index, Amount: utxo.Amount})
	}
	builder.AddOutputs(&cardano.TxOutput{Address: receiver, Amount: amount})

	tip, err := node.Tip()
	if err != nil {
		return nil, err
	}
	builder.SetTTL(tip.Slot + 1200)
	changeAddress := pickedUtxos[0].Spender
	builder.AddChangeIfNeeded(changeAddress)
	builder.Sign([]byte{}) // Use zoombie private key as the key holder.

	tx, err := builder.Build2()
	if err != nil {
		return nil, err
	}
	return tx, nil
}
