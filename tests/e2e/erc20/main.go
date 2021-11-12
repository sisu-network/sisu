package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/contracts/eth/erc20gateway"
	"github.com/sisu-network/sisu/x/tss"
)

const (
	defaultMnemonic = "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum"
)

var (
	privateKeyHexes = []string{
		"9f575b88940d452da46a6ceec06a108fcd5863885524aec7fb0bc4906eb63ab1",
		"3d08e671b7457aaeb1dc0514d72f0871aa80bdcd3b9a37fbd0f8943e0771b6be",
		"03849d7555e05eda85c424611ffff03ec5133fb91e5abbcb7dd9a73eb5a9c8c4",
		"5c58d8757715e92b4c10390aed5a36c58ca696d8f85fd1c992a4018563c2e9d2",
		"4d50fac6b086099a5419117e1625b600ba7e4e4eeeb3e6b85cf10d28182c765c",
		"0a8d54b4988ec1746f46cf7386130bdcce889143d1c8f9ca93e085ed99e57eb1",
		"ea11c49dec5d5d293ecaa88e0b07010f066fc3731c0aaab48f61cdb93b35893d",
		"3eaa1f5761a7011aa8087bdb7a5d79de9bb610ca47d95b36fd5221330f1a702e",
		"1317d3214a927473c281ede09a5229b1d62cd51a4ce0c1d67b2f12f3a836f380",
		"06293b37ea13e22ea103e1145dd9e06d31cc4495d61508566fb5f261d470d005",
	}
)

func getNonce(client *ethclient.Client, address ethcommon.Address) uint64 {
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		panic(err)
	}

	return nonce
}

func getPrivateKeyAndAddress(hexString string) (*ecdsa.PrivateKey, ethcommon.Address) {
	privKey, err := crypto.HexToECDSA(hexString)
	if err != nil {
		panic(err)
	}

	return privKey, crypto.PubkeyToAddress(privKey.PublicKey)
}

func deployErc20Gateway(client *ethclient.Client, chain string) ethcommon.Address {
	erc20 := tss.SupportedContracts[tss.ContractErc20]
	parsedAbi, err := abi.JSON(strings.NewReader(erc20.AbiString))
	if err != nil {
		panic(err)
	}

	input, err := parsedAbi.Pack("", chain, []string{})
	if err != nil {
		panic(err)
	}

	byteCode := ethcommon.FromHex(erc20.Bin)
	input = append(byteCode, input...)
	privKey, accountAddress := getPrivateKeyAndAddress(privateKeyHexes[0])
	nonce := getNonce(client, accountAddress)

	rawTx := etypes.NewContractCreation(
		nonce,
		big.NewInt(0),
		uint64(5000000),
		big.NewInt(50),
		input,
	)

	chainId, err := client.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	tx, err := etypes.SignTx(rawTx, etypes.NewEIP2930Signer(chainId), privKey)
	if err != nil {
		panic(err)
	}

	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 3)

	return crypto.CreateAddress(accountAddress, tx.Nonce())
}

func testTransferIn(client *ethclient.Client, contractAddress ethcommon.Address,
	assetId string, recipient ethcommon.Address, amount *big.Int) {
	privKey, accountAddress := getPrivateKeyAndAddress(privateKeyHexes[0])

	erc20Contract := tss.SupportedContracts[tss.ContractErc20]

	input, err := erc20Contract.Abi.Pack(tss.MethodTransferIn, assetId, recipient, amount)
	if err != nil {
		panic(err)
	}

	nonce := getNonce(client, accountAddress)
	rawTx := etypes.NewTransaction(
		uint64(nonce),
		contractAddress,
		big.NewInt(0),
		uint64(5000000),
		big.NewInt(50),
		input,
	)

	chainId, err := client.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	tx, err := etypes.SignTx(rawTx, etypes.NewEIP155Signer(chainId), privKey)
	if err != nil {
		panic(err)
	}

	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 3)
}

func query(client *ethclient.Client, assetId string, contractAddress, recipient ethcommon.Address, expectedAmount *big.Int) {
	gateway, err := erc20gateway.NewErc20gateway(contractAddress, client)
	if err != nil {
		panic(err)
	}

	balance, err := gateway.GetBalance(&bind.CallOpts{Pending: true}, assetId, recipient)
	if err != nil {
		panic(err)
	}

	if balance.Cmp(expectedAmount) != 0 {
		panic(fmt.Errorf("balance does not match: %s %s", expectedAmount, balance))
	}
}

// Tests deploying erc20 gateway contracts, execute transaction and query balances
func main() {
	url := "http://0.0.0.0:8545"
	client, err := ethclient.Dial(url)
	if err != nil {
		panic(err)
	}

	// 1. Create contract
	contractAddress := deployErc20Gateway(client, "ganache1")
	log.Info("contractAddress = ", contractAddress.Hex())

	// // 2. Transfer In
	assetId := "eth__0x0d73608E5D226eAf90EDF51dF82d4afC24d8B9AA"
	recipient := ethcommon.HexToAddress("0xE8382821BD8a0F9380D88e2c5c33bc89Df17E466")
	amount := big.NewInt(1)
	testTransferIn(client, contractAddress, assetId, recipient, amount)

	query(client, assetId, contractAddress, recipient, amount)
}
