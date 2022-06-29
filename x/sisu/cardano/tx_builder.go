package cardano

import (
	"fmt"
	"math/big"

	"github.com/echovl/cardano-go"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/helper"
	"github.com/sisu-network/sisu/x/sisu/types"
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

// BuildTx constructs a cardano transaction that sends from sender address to receive address.
func BuildTx(node cardano.Node, network cardano.Network, sender, receiver cardano.Address,
	amount *cardano.Value, metadata cardano.Metadata, adaPrice int64, token *types.Token, destChain string, assetAmount uint64) (*cardano.Tx, error) {
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

	// Subtract required ADA (around 1,3 ADA) from multi-asset amount
	requiredAdaFeeInToken := helper.GetCardanoTxFeeInToken(big.NewInt(adaPrice), big.NewInt(token.Price), new(big.Int).SetUint64(uint64(minUTXO)))
	if requiredAdaFeeInToken.Cmp(big.NewInt(0)) < 0 {
		log.Error("tx fee is negative")
		return nil, err
	}

	if assetAmount <= requiredAdaFeeInToken.Uint64() {
		err := fmt.Errorf("token amount can not cover transaction fee. Expect %d, got %d", requiredAdaFeeInToken.Uint64(), assetAmount)
		log.Error(err)
		return nil, err
	}

	fee, err := GetCardanoMultiAsset(destChain, token, requiredAdaFeeInToken.Uint64())
	if err != nil {
		log.Error("error when getting cardano multi-asset: ", err)
		return nil, err
	}
	amount = amount.Sub(cardano.NewValueWithAssets(0, fee))

	// Find utxos that cover the amount to transfer
	pickedUtxos := []cardano.UTxO{}
	utxos, err := findUtxos(node, sender)

	// Pick at least <MinUTXO * 2> lovelace because we will produce at least 2 new utxos which contains multi-asset:
	// 1. Transfer coin + multi-asset for user
	// 2. Transfer remain coin + multi-asset for Cardano gateway (change address)
	targetUtxoBalance := cardano.NewValueWithAssets(amount.Coin*2, amount.MultiAsset)
	pickedUtxosAmount := cardano.NewValue(0)
	for _, utxo := range utxos {
		if pickedUtxosAmount.Cmp(targetUtxoBalance) == 1 {
			break
		}
		pickedUtxos = append(pickedUtxos, utxo)
		pickedUtxosAmount = pickedUtxosAmount.Add(utxo.Amount)
	}

	for _, utxo := range pickedUtxos {
		builder.AddInputs(&cardano.TxInput{TxHash: utxo.TxHash, Index: utxo.Index, Amount: utxo.Amount})
	}
	builder.AddOutputs(&cardano.TxOutput{Address: receiver, Amount: amount})

	if len(metadata) > 0 {
		builder.AddAuxiliaryData(&cardano.AuxiliaryData{Metadata: metadata})
	}

	tip, err := node.Tip()
	if err != nil {
		return nil, err
	}
	builder.SetTTL(tip.Slot + 1200)
	changeAddress := pickedUtxos[0].Spender
	builder.AddChangeIfNeeded(changeAddress)
	builder.Sign([]byte{}) // Use zoombie private key as the key holder.

	// Tx fee is calculated based on tx body, so we have to call the first build to construct the body
	tx, err := builder.Build2()
	if err != nil {
		return nil, err
	}

	// Subtract transaction fee from multi-asset amount
	txFeeInToken := helper.GetCardanoTxFeeInToken(big.NewInt(adaPrice), big.NewInt(token.Price), new(big.Int).SetUint64(uint64(tx.Body.Fee)))
	log.Debug("tx fee (unit token) = ", txFeeInToken.Uint64())
	fee, err = GetCardanoMultiAsset(destChain, token, txFeeInToken.Uint64())
	if err != nil {
		log.Error("error when getting cardano multi-asset: ", err)
		return nil, err
	}
	amount = amount.Sub(cardano.NewValueWithAssets(0, fee))

	clonedBuilder := CloneTxBuilder(pparams, pickedUtxos, receiver, amount, changeAddress, metadata, tip)

	// Second build here
	tx, err = clonedBuilder.Build2()
	if err != nil {
		log.Error("error when build cardano tx: ", err)
		return nil, err
	}

	return tx, nil
}

func CloneTxBuilder(pparams *cardano.ProtocolParams, inputUtxo []cardano.UTxO, receiver cardano.Address,
	amount *cardano.Value, changeAddress cardano.Address, metadata cardano.Metadata, tip *cardano.NodeTip) *cardano.TxBuilderV2 {
	builder := cardano.NewTxBuilderV2(pparams)

	for _, utxo := range inputUtxo {
		builder.AddInputs(&cardano.TxInput{TxHash: utxo.TxHash, Index: utxo.Index, Amount: utxo.Amount})
	}
	builder.AddOutputs(&cardano.TxOutput{Address: receiver, Amount: amount})
	if len(metadata) > 0 {
		builder.AddAuxiliaryData(&cardano.AuxiliaryData{Metadata: metadata})
	}

	builder.SetTTL(tip.Slot + 1200)
	builder.AddChangeIfNeeded(changeAddress)
	builder.Sign([]byte{}) // Use zoombie private key as the key holder.

	return builder
}
