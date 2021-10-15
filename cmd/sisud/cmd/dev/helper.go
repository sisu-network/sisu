package dev

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/dcore/accounts"
	hdwallet "github.com/sisu-network/sisu/utils/hdwallet"
)

const (
	default_mnemonic = "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum"
)

var (
	localWallet *hdwallet.Wallet
	account0    accounts.Account
	privateKey0 *ecdsa.PrivateKey
)

func init() {
	var err error
	localWallet, err = hdwallet.NewFromMnemonic(default_mnemonic)
	if err != nil {
		panic(err)
	}

	path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", 0))
	account0, err = localWallet.Derive(path, true)
	if err != nil {
		panic(err)
	}

	privateKey0, err = localWallet.PrivateKey(account0)
	if err != nil {
		panic(err)
	}
}

func getEthClient(fromChain string) (*ethclient.Client, error) {
	switch fromChain {
	case "eth":
		return ethclient.Dial("http://0.0.0.0:7545")
	}

	return nil, fmt.Errorf("cannot find client for chain", fromChain)
}

func getAuthTransactor(client *ethclient.Client) (*bind.TransactOpts, error) {
	nonce, err := client.PendingNonceAt(context.Background(), account0.Address)
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	auth := bind.NewKeyedTransactor(privateKey0)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(1000000)

	return auth, nil
}
