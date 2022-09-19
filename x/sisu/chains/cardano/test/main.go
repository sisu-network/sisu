package main

import (
	"fmt"
	"os"

	"github.com/echovl/cardano-go"
	"github.com/echovl/cardano-go/blockfrost"
	"github.com/echovl/cardano-go/wallet"
)

const (
	Decimals = 1000 * 1000
)

// Localhost address: addr_test1vpljp50hd27w9s8mekvgxdwfhzyrhw4ytvz6xeuhvdrr9vsts7wex
func getBlockfrostSecret() string {
	secret := os.Getenv("BLOCKFROST_SECRET")
	if len(secret) == 0 {
		panic("Blockfrost secret is empty. Please set its value in the .env file")
	}
	return secret
}

func getWalletPassoword() string {
	password := os.Getenv("WALLET_PASSWORD")
	if len(password) == 0 {
		password = "12345678910"
	}
	return password
}

func getMnemonic() string {
	mnemonic := os.Getenv("MNEMONIC")
	if len(mnemonic) == 0 {
		mnemonic = "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum"
	}
	return mnemonic
}

func transferAda(recipientString string) {
	node := blockfrost.NewNode(cardano.Preprod, getBlockfrostSecret())
	client := wallet.NewClient(&wallet.Options{Node: node})
	w, err := client.RestoreWallet("sisu", getWalletPassoword(), getMnemonic())
	if err != nil {
		panic(err)
	}

	addrs, err := w.Addresses()
	if err != nil {
		panic(err)
	}
	fmt.Println("addr = ", addrs[0])

	balance, err := w.Balance()
	if err != nil {
		panic(err)
	}

	fmt.Println("Balance = ", balance)

	recipient, err := cardano.NewAddress(recipientString)
	if err != nil {
		panic(err)
	}

	txHash, err := w.Transfer(recipient, cardano.NewValue(2*Decimals), nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("hash = ", txHash.String())
}

func main() {
	transferAda("addr_test1vr40ggmush8wg8tdnpjm2d3sn65fftcsxfvprwjku2ekjhcgmz8f5")
}
