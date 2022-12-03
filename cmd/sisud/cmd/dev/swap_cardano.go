package dev

import (
	"context"
	"fmt"
	"math/big"

	"github.com/echovl/cardano-go"
	cgblockfrost "github.com/echovl/cardano-go/blockfrost"
	cardanocrypto "github.com/echovl/cardano-go/crypto"
	"github.com/echovl/cardano-go/wallet"
	hutils "github.com/sisu-network/dheart/utils"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	scardano "github.com/sisu-network/sisu/x/sisu/chains/cardano"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func getCardanoVault(ctx context.Context, sisuRpc string) string {
	allPubKeys := queryPubKeys(ctx, sisuRpc)
	cardanoKey, ok := allPubKeys[libchain.KEY_TYPE_EDDSA]
	if !ok {
		panic("can not find cardano pub key")
	}

	return hutils.GetAddressFromCardanoPubkey(cardanoKey).String()
}

func swapFromCardano(srcChain string, destChain string, token *types.Token,
	destRecipient, cardanoVault string, value *big.Int, network string, secret, mnemonic string, deyesUrl string) {
	privateKey, senderAddress := getCardanoSenderAddress(secret, mnemonic)
	receiver, err := cardano.NewAddress(cardanoVault)
	if err != nil {
		panic(err)
	}

	var amount *cardano.Value
	if token.Id == "ADA" {
		amount = cardano.NewValue(cardano.Coin(value.Uint64()))
	} else {
		multiAsset, err := scardano.GetCardanoMultiAsset(srcChain, token, value.Uint64())
		if err != nil {
			panic(err)
		}
		amount = cardano.NewValueWithAssets(cardano.Coin(1_600_000), multiAsset)
	}

	var metadata cardano.Metadata
	var nativeAda int

	if token.Id == "ADA" {
		nativeAda = 1
	}

	metadata = cardano.Metadata{
		0: map[string]interface{}{
			"chain":      destChain,
			"recipient":  destRecipient,
			"native_ada": nativeAda,
		},
	}

	deyesClient, err := external.DialDeyes(deyesUrl)
	if err != nil {
		panic(err)
	}

	tip, err := deyesClient.CardanoTip(srcChain, 20_000_000)
	if err != nil {
		panic(err)
	}

	utxos, err := deyesClient.CardanoUtxos(srcChain, senderAddress.String(), tip.Block+1000)
	tx, err := scardano.BuildTx(deyesClient, srcChain, senderAddress, []cardano.Address{receiver},
		[]*cardano.Value{amount}, metadata, utxos, tip.Block)
	if err != nil {
		panic(err)
	}
	if len(tx.WitnessSet.VKeyWitnessSet) != 1 {
		panic(fmt.Errorf("VKeyWitnessSet is expected to have length 1 but has length %d", len(tx.WitnessSet.VKeyWitnessSet)))
	}

	txHash, err := tx.Hash()
	if err != nil {
		panic(err)
	}

	// Sign tx
	tx.WitnessSet.VKeyWitnessSet = make([]cardano.VKeyWitness, 1)
	tx.WitnessSet.VKeyWitnessSet[0] = cardano.VKeyWitness{
		VKey:      privateKey.PubKey(),
		Signature: privateKey.Sign(txHash),
	}

	submitedHash, err := deyesClient.CardanoSubmitTx(srcChain, tx)
	if err != nil {
		panic(err)
	}

	if submitedHash.String() != txHash.String() {
		panic(fmt.Errorf("TxHash and submitted hash do not match, txhash = %s, submitted hash = %s", txHash, submitedHash))
	}

	log.Info("Cardano tx hash = ", txHash.String())
}

func getCardanoSenderAddress(blockfrostSecret, cardanoMnemonic string) (cardanocrypto.PrvKey, cardano.Address) {
	node := cgblockfrost.NewNode(cardano.Testnet, blockfrostSecret)
	opts := &wallet.Options{Node: node}
	client := wallet.NewClient(opts)

	wallet, err := client.RestoreWallet(DefaultCardanoWalletName, DefaultCardanoPassword, cardanoMnemonic)
	if err != nil {
		panic(err)
	}

	walletAddrs, err := wallet.Addresses()
	if err != nil {
		panic(err)
	}
	log.Info("sender = ", walletAddrs[0])

	key, _ := wallet.Keys()

	return key, walletAddrs[0]
}
