package dev

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/contracts/eth/erc20"
)

func waitForTx(client *ethclient.Client, hash common.Hash) {
	start := time.Now()
	end := start.Add(time.Minute * 2)

	for {
		if time.Now().After(end) {
			log.Error("Time out for transaction with hash ", hash)
			panic("")
		}

		tx, isPending, err := client.TransactionByHash(context.Background(), hash)
		if err != nil && err != ethereum.NotFound {
			log.Error("Failed to get transaction with hash ", hash)
			panic(err)
		}

		if tx == nil || isPending {
			time.Sleep(time.Second * 3)
			continue
		}

		break
	}
}

func queryAllownace(client *ethclient.Client,
	tokenAddr, owner, target string) *big.Int {
	store, err := erc20.NewErc20(common.HexToAddress(tokenAddr), client)
	if err != nil {
		panic(err)
	}

	balance, err := store.Allowance(nil, common.HexToAddress(owner), common.HexToAddress(target))
	if err != nil {
		panic(err)
	}

	return balance
}
