package dev

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/contracts/eth/erc20"
	hdwallet "github.com/sisu-network/sisu/utils/hdwallet"
)

const (
	defaultMnemonic = "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum"
	Blocktime       = time.Second * 3
)

var (
	localWallet *hdwallet.Wallet
	account0    accounts.Account
	privateKey0 *ecdsa.PrivateKey
	nonceMap    map[string]*big.Int
)

func init() {
	var err error
	localWallet, err = hdwallet.NewFromMnemonic(defaultMnemonic)
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

	nonceMap = make(map[string]*big.Int)
}

func getEthClient(port int) (*ethclient.Client, error) {
	return ethclient.Dial(fmt.Sprintf("http://0.0.0.0:%d", port))
}

func getPrivateKey(mnemonic string) (*ecdsa.PrivateKey, common.Address) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		panic(err)
	}

	dpath, err := accounts.ParseDerivationPath("m/44'/60'/0'/0/0")
	if err != nil {
		panic(err)
	}

	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)

	key := masterKey
	for _, n := range dpath {
		key, err = key.Derive(n)
		if err != nil {
			panic(err)
		}
	}

	privateKey, err := key.ECPrivKey()
	if err != nil {
		panic(err)
	}

	privateKeyECDSA := privateKey.ToECDSA()
	publicKey := privateKeyECDSA.PublicKey
	addr := crypto.PubkeyToAddress(publicKey)

	return privateKeyECDSA, addr
}

func getAuthTransactor(client *ethclient.Client, mnemonic string) (*bind.TransactOpts, error) {
	// This is the private key of the accounts0
	privateKey, owner := getPrivateKey(mnemonic)
	nonce, err := client.PendingNonceAt(context.Background(), owner)
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	// This is the private key of the accounts0

	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasPrice = gasPrice

	auth.GasLimit = uint64(5_000_000)

	return auth, nil
}

func getEthClients(urlString string) []*ethclient.Client {
	urls := strings.Split(urlString, ",")
	clients := make([]*ethclient.Client, 0)

	// Get all urls from command arguments.
	for i := 0; i < len(urls); i++ {
		client, err := ethclient.Dial(urls[i])
		if err != nil {
			log.Error("please check chain is up and running, url = ", urls[i])
			panic(err)
		}
		clients = append(clients, client)
	}

	return clients
}

func queryErc20Balance(client *ethclient.Client, tokenAddr string, target string) (*big.Int, error) {
	store, err := erc20.NewErc20(common.HexToAddress(tokenAddr), client)
	if err != nil {
		return nil, err
	}

	balance, err := store.BalanceOf(nil, common.HexToAddress(target))

	return balance, err
}

func approveAddress(client *ethclient.Client, mnemonic string, erc20Addr string, target string) {
	contract, err := erc20.NewErc20(common.HexToAddress(erc20Addr), client)
	if err != nil {
		panic(err)
	}

	opts, err := getAuthTransactor(client, mnemonic)
	if err != nil {
		panic(err)
	}

	_, owner := getPrivateKey(mnemonic)
	ownerBalance, err := contract.BalanceOf(nil, owner)

	tx, err := contract.Approve(opts, common.HexToAddress(target), ownerBalance)
	bind.WaitDeployed(context.Background(), client, tx)
	time.Sleep(time.Second * 3)
}
