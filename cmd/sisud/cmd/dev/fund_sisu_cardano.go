package dev

import (
	"fmt"
	"strings"

	"github.com/echovl/cardano-go"
	"github.com/echovl/cardano-go/blockfrost"
	"github.com/echovl/cardano-go/wallet"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/helper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func (c *fundAccountCmd) fundCardano(genesisFolder string, receiver cardano.Address, funderMnemonic string,
	blockfrostSecret string, sisuRpc string, tokens []*types.Token) {
	cardanoConfig := helper.ReadCardanoConfig(genesisFolder)

	network := cardanoConfig.GetCardanoNetwork()
	node := blockfrost.NewNode(network, blockfrostSecret)
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
	multiAsset := c.getMultiAsset(sisuRpc, cardanoConfig.Chain, tokens, 1e3)
	txHash, err := funderWallet.Transfer(receiver, cardano.NewValueWithAssets(30*CardanoDecimals,
		multiAsset), nil) // 30ADA
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

func (c *fundAccountCmd) getMultiAsset(sisuRpc string, chain string, tokens []*types.Token, amt uint64) *cardano.MultiAsset {
	m := make(map[string]*cardano.Assets)

	for _, token := range tokens {
		if token.Id == "ADA" || len(token.Addresses) == 0 {
			continue
		}

		var tokenAddr string
		for i, c := range token.Chains {
			if c == chain {
				tokenAddr = token.Addresses[i]
				break
			}
		}

		index := strings.Index(tokenAddr, ":")
		policyString := tokenAddr[:index]
		assetName := tokenAddr[index+1:]

		if m[policyString] == nil {
			m[policyString] = cardano.NewAssets()
		}

		asset := cardano.NewAssetName(assetName)
		m[policyString].Set(asset, cardano.BigNum(amt*CardanoDecimals))
	}

	multiAsset := cardano.NewMultiAsset()
	for policy, assets := range m {
		policyHash, err := cardano.NewHash28(policy)
		if err != nil {
			err := fmt.Errorf("error when parsing policyID hash: %v", err)
			panic(err)
		}
		policyID := cardano.NewPolicyIDFromHash(policyHash)
		multiAsset.Set(policyID, assets)
	}

	return multiAsset
}
