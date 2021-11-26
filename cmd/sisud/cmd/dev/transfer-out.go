package dev

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/contracts/eth/erc20"
	erc20Gateway "github.com/sisu-network/sisu/contracts/eth/erc20gateway"
	"github.com/sisu-network/sisu/db"
	hdwallet "github.com/sisu-network/sisu/utils/hdwallet"
	"github.com/sisu-network/sisu/x/tss"
	"github.com/spf13/cobra"
)

// WIP. TODO: Complete and clean up this.
func TransferOut() *cobra.Command {
	cmd := &cobra.Command{
		Use: "transfer-out",
		Long: `Transfer an ERC20 or ERC721 asset.
Usage:
transfer-out [ContractType] [FromChain] [Port] [TokenAddress] [ToChain] [RecipientAddress]

Example:
transfer-out erc20 eth-sisu-local 1234 0xB369Be7F62cfb3F44965db83404997Fa6EC9Dd58 ganache1 0xE8382821BD8a0F9380D88e2c5c33bc89Df17E466
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Use this when running with docker.
			//sqlConfig := getDockerSqlConfig()

			// // Use this when running with single node on command line.
			cfg, cfgErr := config.ReadConfig()
			if cfgErr != nil {
				panic(cfgErr)
			}
			sqlConfig := cfg.Sisu.Sql

			database := db.NewDatabase(sqlConfig)
			err := database.Init()
			if err != nil {
				panic(err)
			}
			defer database.Close()

			// Get the contract address of token
			log.Info("args = ", args)
			contractType := args[0]
			fromChain := args[1]
			port, err := strconv.Atoi(args[2])
			if err != nil {
				panic(err)
			}

			tokenAddressString := args[3]
			toChain := args[4]
			recipient := args[5]

			switch contractType {
			case "erc20":
				client, err := getEthClient(port)
				if err != nil {
					panic(err)
				}

				hash := tss.SupportedContracts[contractType].AbiHash
				contract := database.GetContractFromHash(fromChain, hash)
				if contract == nil {
					panic(fmt.Errorf("cannot find contract"))
				}

				gatewayAddress := common.HexToAddress(contract.Address)
				gateway, err := erc20Gateway.NewErc20gateway(gatewayAddress, client)
				if err != nil {
					panic(err)
				}

				log.Info("gatewayAddress = ", gatewayAddress.String())

				tokenAddress := common.HexToAddress(tokenAddressString)
				erc20Contract, err := erc20.NewErc20(tokenAddress, client)
				if err != nil {
					return err
				}

				log.Info("Approving gateway address...")
				amount := big.NewInt(1)
				approveAddress(erc20Contract, gatewayAddress, amount, client)

				// Check the allowance
				allowance, err := erc20Contract.Allowance(&bind.CallOpts{Pending: true}, account0.Address, gatewayAddress)
				if err != nil {
					panic(err)
				}
				if allowance.Cmp(amount) != 0 {
					panic(fmt.Errorf("Invalid balance: expected %s, actual %s", amount, allowance))
				}

				log.Info("Transfering token out....")
				auth, err := getAuthTransactor(client, account0.Address)
				tx, err := gateway.TransferOutFromContract(auth, tokenAddress, toChain, recipient, amount)
				if err != nil {
					panic(err)
				}
				bind.WaitDeployed(context.Background(), client, tx)

				time.Sleep(Blocktime)

				gatewayBalance := getBalance(erc20Contract, gatewayAddress)
				log.Info("gatewayBalance = ", gatewayBalance)
			}

			return nil
		},
	}
	return cmd
}

func getDockerSqlConfig() config.SqlConfig {
	return config.SqlConfig{
		Host:     "0.0.0.0",
		Port:     4000,
		Username: "root",
		Password: "password",
		Schema:   "sisu0", // TODO: make this schema configurable
	}
}

func deployGatewayContract(toChain string, client *ethclient.Client) (common.Address, *erc20Gateway.Erc20gateway) {
	auth, err := getAuthTransactor(client, account0.Address)
	if err != nil {
		panic(err)
	}

	gatewayAddress, tx, gateway, err := erc20Gateway.DeployErc20gateway(
		auth,
		client,
		"eth",
		[]string{toChain},
	)
	if err != nil {
		panic(err)
	}

	_, err = bind.WaitDeployed(context.Background(), client, tx)
	if err != nil {
		panic(err)
	}

	log.Info("Gateway was deployed!")

	return gatewayAddress, gateway
}

func approveAddress(erc20Contract *erc20.Erc20, recipient common.Address, amount *big.Int, client *ethclient.Client) {
	auth, err := getAuthTransactor(client, account0.Address)
	if err != nil {
		panic(err)
	}

	_, err = erc20Contract.Approve(auth, recipient, amount)
	if err != nil {
		panic(err)
	}

	time.Sleep(Blocktime)
}

func getBalance(erc20Contract *erc20.Erc20, address common.Address) *big.Int {
	tokenBalance, err := erc20Contract.BalanceOf(&bind.CallOpts{Pending: true}, address)
	if err != nil {
		panic(err)
	}

	return tokenBalance
}

func getTransasctionOpts(wallet *hdwallet.Wallet, chainId *big.Int) *bind.TransactOpts {
	path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", 0))
	fromAccount, err := wallet.Derive(path, true)
	if err != nil {
		panic(err)
	}

	privateKey, err := wallet.PrivateKey(fromAccount)
	if err != nil {
		return nil
	}

	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		panic(err)
	}

	return opts
}
