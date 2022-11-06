package dev

import (
	"github.com/echovl/cardano-go"
	"github.com/echovl/cardano-go/blockfrost"
	"github.com/echovl/cardano-go/wallet"
	"github.com/sisu-network/lib/log"
)

func (c *fundAccountCmd) fundCardano(receiver cardano.Address, funderMnemonic string,
	cardanoNetwork, blockfrostSecret string, sisuRpc string, tokens []string) {
	node := blockfrost.NewNode(cardano.Preprod, blockfrostSecret)
	opts := &wallet.Options{
		Node: node,
	}
	client := wallet.NewClient(opts)
	funderWallet, err := c.getWalletFromMnemonic(client, DefaultCardanoWalletName, DefaultCardanoPassword, funderMnemonic)
	if err != nil {
		panic(err)
	}

	addrs, err := funderWallet.Addresses()
	if err != nil {
		panic(err)
	}
	funderAddr := addrs[0]
	log.Info("Cardano funder address = ", funderAddr.String())

	// fund 30 ADA and 1000 WRAP_ADA
	txHash, err := funderWallet.Transfer(receiver, cardano.NewValueWithAssets(30*CardanoDecimals,
		c.getMultiAsset(sisuRpc, cardanoNetwork, tokens, 1e3)), nil) // 30ADA
	if err != nil {
		panic(err)
	}

	log.Infof("Address funded = %s, txHash = %s, explorer: https://testnet.cardanoscan.io/transaction/%s\n",
		receiver, txHash.String(), txHash.String())
}

func (c *fundAccountCmd) getWalletFromMnemonic(client *wallet.Client, name, password, mnemonic string) (*wallet.Wallet, error) {
	w, err := client.RestoreWallet(name, password, mnemonic)
	if err != nil {
		return nil, err
	}

	return w, nil
}
