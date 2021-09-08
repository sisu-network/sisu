package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/sisu-network/sisu/x/tss"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	hdwallet "github.com/sisu-network/sisu/utils/hdwallet"
)

// This package contains small tests that need real deployment (often on localhost). We cannot
// simply run it with go testing.

func testPrepareEthContracts() {
	o := tss.NewCrossChainLogic()

	// Replace the mnemonic value from what you get from ganache.
	wallet, err := hdwallet.NewFromMnemonic("arm misery utility choose box pelican loop lawn beauty result asset treat")
	path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", 0))
	account, err := wallet.Derive(path, false)
	privateKey, err := wallet.PrivateKey(account)

	rawTx := o.PrepareEthContractDeployment("eth", 0)
	if err != nil {
		panic(err)
	}

	signedTx, err := types.SignTx(rawTx, types.NewEIP155Signer(big.NewInt(1)), privateKey)
	client, err := ethclient.Dial("http://localhost:7545") // Ganache local address
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(err)
	}
}

func main() {
	testPrepareEthContracts()
}
