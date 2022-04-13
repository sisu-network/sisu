package dev

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	liquidity "github.com/sisu-network/sisu/contracts/eth/liquiditypool"
	"github.com/sisu-network/sisu/contracts/eth/sampleerc20"
	"github.com/spf13/cobra"
)

type DeployContractCmd struct {
	privateKey *ecdsa.PrivateKey
}

func DeployContract() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deploy",
		Long: `Deploy an ERC20 contract. You can list of empty string to expected addresses param.
Usage:
./sisu dev deploy --contract [contract-type] --chain-urls [list-of-urls] --token-name [TOKEN_NAME] --token-symbol [TOKEN_SYMBOL] --expected-addrs [List of Expected Addresses]

Example:
./sisu dev deploy --contract liquidity --chain-urls http://localhost:7545,http://localhost:8545
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			urlString, _ := cmd.Flags().GetString(flags.ChainUrls)
			contract, _ := cmd.Flags().GetString(flags.Contract)
			mnemonic, _ := cmd.Flags().GetString(flags.Mnemonic)
			expectedAddrString, _ := cmd.Flags().GetString(flags.ExpectedAddrs)
			tokenName, _ := cmd.Flags().GetString(flags.Erc20Name)
			tokenSymbol, _ := cmd.Flags().GetString(flags.Erc20Symbol)

			c := &DeployContractCmd{}
			c.doDeployment(urlString, contract, mnemonic, expectedAddrString, tokenName, tokenSymbol)

			return nil
		},
	}

	cmd.Flags().String(flags.Contract, "liquidity", "Contract name that we want to deploy.")
	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum", "Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.ChainUrls, "http://0.0.0.0:7545,http://0.0.0.0:8545", "RPCs of all the chains we want to fund.")
	cmd.Flags().String(flags.ExpectedAddrs, fmt.Sprintf("%s,%s", ExpectedLiquidPoolAddress, ExpectedLiquidPoolAddress), "Expected addressed of the contract after deployment. Empty string means do not check for address match.")
	cmd.Flags().String(flags.Erc20Name, "Sisu Token", "Token name")
	cmd.Flags().String(flags.Erc20Symbol, "SISU", "Token symbol")

	return cmd
}

func (c *DeployContractCmd) doDeployment(urlString, contract, mnemonic, expAddrString, tokenName, tokenSymbol string) []string {
	urls := strings.Split(urlString, ",")
	expectedAddrs := strings.Split(expAddrString, ",")

	clients := make([]*ethclient.Client, 0)

	if len(urls) != len(expectedAddrs) {
		panic("Expected addrs length does not match urls length")
	}

	// Get all urls from command arguments.
	for i := 0; i < len(urls); i++ {
		client, err := ethclient.Dial(urls[i])
		if err != nil {
			log.Error("please check chain is up and running, url = ", urls[i])
			panic(err)
		}
		clients = append(clients, client)
	}
	defer func() {
		for _, client := range clients {
			client.Close()
		}
	}()

	var owner common.Address
	c.privateKey, owner = c.getPrivateKey(mnemonic)

	log.Info("Key public addr = ", owner)
	deployedAddrs := make([]string, len(urls))

	for i, client := range clients {
		// If liquidity contract has been deployed, do nothing.
		if len(expectedAddrs[i]) > 0 && c.isContractDeployed(client, common.HexToAddress(expectedAddrs[i])) {
			log.Verbose("Contract ", i, " has been deployed")
			continue
		}

		var addr common.Address
		switch contract {
		case "erc20":
			addr = c.deployErc20(client, owner, expectedAddrs[i], tokenName, tokenSymbol)

		case "liquidity":
			addr = c.deployLiquidity(client, owner, expectedAddrs[i])

		default:
			panic(fmt.Sprintf("Unknown contract %s", contract))
		}

		deployedAddrs[i] = addr.String()
	}

	for _, addr := range deployedAddrs {
		log.Info("Deployed addr = ", addr)
	}

	return deployedAddrs
}

func (c *DeployContractCmd) getPrivateKey(mnemonic string) (*ecdsa.PrivateKey, common.Address) {
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

func (c *DeployContractCmd) queryNativeBalance(client *ethclient.Client, addr common.Address) *big.Int {
	balance, err := client.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		panic(err)
	}

	return balance
}

func (c *DeployContractCmd) isContractDeployed(client *ethclient.Client, tokenAddress common.Address) bool {
	bz, err := client.CodeAt(context.Background(), tokenAddress, nil)
	if err != nil {
		log.Error("Cannot get code at ", tokenAddress.String(), " err = ", err)
		return false
	}

	return len(bz) > 10
}

func (c *DeployContractCmd) deployErc20(client *ethclient.Client, accountAddr common.Address, expectedAddress string, tokenName, tokenSymbol string) common.Address {
	auth, err := c.getAuthTransactor(client, accountAddr)
	if err != nil {
		panic(err)
	}

	_, tx, _, err := sampleerc20.DeploySampleerc20(auth, client, tokenName, tokenSymbol)
	if err != nil {
		panic(err)
	}

	log.Info("Deploying erc20 contract ... ")
	contractAddr, err := bind.WaitDeployed(context.Background(), client, tx)
	if err != nil {
		panic(err)
	}

	if len(expectedAddress) > 0 && contractAddr.String() != expectedAddress {
		panic(fmt.Errorf(`Unmatched ERC20 address. We expect address %s but get %s.
You need to update the expected address (both in this file and the tokens_dev.json).`,
			expectedAddress, contractAddr.String()))
	}

	log.Info("Deployed contract successfully, addr: ", contractAddr.String())

	return contractAddr
}

func (c *DeployContractCmd) deployLiquidity(client *ethclient.Client, accountAddr common.Address, expectedAddress string) common.Address {
	auth, err := c.getAuthTransactor(client, accountAddr)
	if err != nil {
		panic(err)
	}

	_, tx, _, err := liquidity.DeployLiquiditypool(auth, client, []common.Address{}, []string{})
	if err != nil {
		panic(err)
	}

	log.Info("Deploying liquidity ... ")
	contractAddr, err := bind.WaitDeployed(context.Background(), client, tx)
	if err != nil {
		panic(err)
	}

	if len(expectedAddress) > 0 && contractAddr.String() != expectedAddress {
		panic(fmt.Errorf(`Unmatched Liquid pool address. We expect address %s but get %s.
You need to update the expected address (both in this file and the liquidity_dev.json).`,
			expectedAddress, contractAddr.String()))
	}

	log.Info("Deployed liquidity successfully, addr: ", contractAddr.String())

	return contractAddr
}

func (c *DeployContractCmd) getAuthTransactor(client *ethclient.Client, address common.Address) (*bind.TransactOpts, error) {
	nonce, err := client.PendingNonceAt(context.Background(), address)
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

	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, chainId)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasPrice = gasPrice

	auth.GasLimit = uint64(5_000_000)

	return auth, nil
}
