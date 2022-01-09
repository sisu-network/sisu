package dev

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	defaultMnemonic = "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum"
	Blocktime       = time.Second * 3
)

var (
	privateKey0 *ecdsa.PrivateKey
	nonceMap    map[string]*big.Int
)

func init() {
}

func getEthClient(port int) (*ethclient.Client, error) {
	return ethclient.Dial(fmt.Sprintf("http://0.0.0.0:%d", port))
}

func getAuthTransactor(client *ethclient.Client, address common.Address) (*bind.TransactOpts, error) {
	addrString := address.Hex()
	if nonceMap[addrString] == nil {
		nonce, err := client.PendingNonceAt(context.Background(), address)
		if err != nil {
			return nil, err
		}

		nonceMap[addrString] = big.NewInt(int64(nonce))
	} else {
		nonceMap[addrString] = new(big.Int).Add(nonceMap[addrString], big.NewInt(1))
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	auth := bind.NewKeyedTransactor(privateKey0)
	auth.Nonce = nonceMap[addrString]
	auth.Value = big.NewInt(0)
	auth.GasPrice = gasPrice

	// auth.GasLimit = uint64(30 * 1000000) // 30M gas
	auth.GasLimit = uint64(3000000)

	return auth, nil
}
