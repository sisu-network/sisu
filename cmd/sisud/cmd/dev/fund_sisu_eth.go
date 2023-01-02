package dev

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/contracts/eth/vault"
	"github.com/sisu-network/sisu/utils"
)

// transferEth transfers a specific ETH amount to an address.
func (c *fundAccountCmd) transferEth(client *ethclient.Client, sisuRpc, chain, mnemonic, recipient string) {
	ch, err := queryChain(context.Background(), sisuRpc, chain)
	if err != nil {
		panic(fmt.Errorf("failed to get chain, err = %v", err))
	}

	_, account := getPrivateKey(mnemonic)
	log.Info("from address = ", account.String(), " to Address = ", recipient)

	nonce, err := client.PendingNonceAt(context.Background(), account)
	if err != nil {
		log.Error("Failed to get nonce on chain ", chain)
		panic(err)
	}

	genesisGas := big.NewInt(ch.EthConfig.GasPrice)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}
	if gasPrice.Cmp(genesisGas) < 0 {
		gasPrice = genesisGas
	}
	// Add some 10% premimum to the gas price
	gasPrice = gasPrice.Mul(gasPrice, big.NewInt(110))
	gasPrice = gasPrice.Quo(gasPrice, big.NewInt(100))

	log.Info("Gas price = ", gasPrice, " on chain ", chain)

	// 0.01 ETH
	amount := new(big.Int).Div(utils.EthToWei, big.NewInt(100))

	gasLimit := uint64(22000) // in units

	amountFloat := new(big.Float).Quo(new(big.Float).SetInt(amount), new(big.Float).SetInt(utils.ONE_ETHER_IN_WEI))
	log.Info("Amount in ETH: ", amountFloat, " on chain ", chain)

	toAddress := common.HexToAddress(recipient)
	var data []byte
	tx := ethtypes.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, data)

	privateKey, _ := getPrivateKey(mnemonic)
	signedTx, err := ethtypes.SignTx(tx, getSigner(client), privateKey)

	log.Info("Tx hash = ", signedTx.Hash(), " on chain ", chain)

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(fmt.Errorf("Failed to transfer ETH on chain %s, err = %s", chain, err))
	}

	bind.WaitDeployed(context.Background(), client, signedTx)

	waitForTx(client, signedTx.Hash())
}

func (c *fundAccountCmd) addVaultSpender(client *ethclient.Client, mnemonic string, vaultAddr, spender common.Address) {
	log.Info("Add vault spender, vault = ", vaultAddr.String(), " spender = ", spender.String())
	vaultInstance, err := vault.NewVault(vaultAddr, client)
	if err != nil {
		panic(err)
	}

	auth, err := getAuthTransactor(client, mnemonic)
	if err != nil {
		panic(err)
	}

	tx, err := vaultInstance.AddSpender(auth, spender)
	if err != nil {
		log.Error("Failed to add vault spender, vaultAddr = ", vaultAddr)
		panic(err)
	}

	bind.WaitDeployed(context.Background(), client, tx)

	waitForTx(client, tx.Hash())
}
