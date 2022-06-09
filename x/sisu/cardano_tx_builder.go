package sisu

import (
	"fmt"

	"github.com/echovl/cardano-go"
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

func BuildTx(node cardano.Node, network cardano.Network, sender, receiver cardano.Address, amount *cardano.Value) (*cardano.Tx, error) {
	// Calculate if the account has enough balance
	balance, err := Balance(node, sender)
	if err != nil {
		return nil, err
	}

	if cmp := balance.Cmp(amount); cmp == -1 || cmp == 2 {
		return nil, fmt.Errorf("Not enough balance, %v > %v", amount, balance)
	}

	// Find utxos that cover the amount to transfer
	pickedUtxos := []cardano.UTxO{}
	utxos, err := findUtxos(node, sender)
	pickedUtxosAmount := cardano.NewValue(0)
	for _, utxo := range utxos {
		if pickedUtxosAmount.Cmp(amount) == 1 {
			break
		}
		pickedUtxos = append(pickedUtxos, utxo)
		pickedUtxosAmount = pickedUtxosAmount.Add(utxo.Amount)
	}

	pparams, err := node.ProtocolParams()
	if err != nil {
		return nil, err
	}

	builder := cardano.NewTxBuilderV2(pparams)

	inputAmount := cardano.NewValue(0)
	for _, utxo := range pickedUtxos {
		builder.AddInputs(&cardano.TxInput{TxHash: utxo.TxHash, Index: utxo.Index, Amount: utxo.Amount})
		inputAmount = inputAmount.Add(utxo.Amount)
	}
	builder.AddOutputs(&cardano.TxOutput{Address: receiver, Amount: amount})

	tip, err := node.Tip()
	if err != nil {
		return nil, err
	}
	builder.SetTTL(tip.Slot + 1200)
	changeAddress := pickedUtxos[0].Spender
	if err = builder.AddChangeIfNeeded(changeAddress); err != nil {
		return nil, err
	}

	tx, err := builder.Build2()
	if err != nil {
		return nil, err
	}

	return tx, nil
}
