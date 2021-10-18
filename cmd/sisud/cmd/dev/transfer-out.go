package dev

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/contracts/eth/erc20"
	erc20Gateway "github.com/sisu-network/sisu/contracts/eth/erc20gateway"
	"github.com/sisu-network/sisu/db"
	"github.com/sisu-network/sisu/utils"
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
transfer-out [FromChain] [ContractType] [TokenAddress] [ToChain]

Example:
transfer-out eth erc20 0xB369Be7F62cfb3F44965db83404997Fa6EC9Dd58 sisu-eth
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get db config
			cfg, err := config.ReadConfig()
			if err != nil {
				return err
			}

			database := db.NewDatabase(cfg.Sisu.Sql)
			err = database.Init()
			if err != nil {
				return err
			}
			defer database.Close()

			// Get the contract address of token
			utils.LogInfo("args = ", args)
			fromChain := args[0]
			contractType := args[1]
			tokenAddressString := args[2]
			toChain := args[3]

			switch contractType {
			case "erc20":
				client, err := getEthClient(fromChain)
				if err != nil {
					panic(err)
				}

				hash := tss.SupportedContracts[contractType].AbiHash
				contract := database.GetContractFromHash(fromChain, hash)
				if contract == nil {
					return fmt.Errorf("cannot find contract")
				}
				gatewayAddress := common.HexToAddress(contract.Address)
				gateway, err := erc20Gateway.NewErc20gateway(gatewayAddress, client)
				if err != nil {
					panic(err)
				}

				tokenAddress := common.HexToAddress(tokenAddressString)
				erc20Contract, err := erc20.NewErc20(tokenAddress, client)
				if err != nil {
					return err
				}

				utils.LogInfo("Transfering token out....")
				approveAddress(erc20Contract, gatewayAddress, big.NewInt(1), client)

				auth, err := getAuthTransactor(client, account0.Address)
				tx, err := gateway.TransferOutFromContract(auth, tokenAddress, toChain, account0.Address.String(), big.NewInt(1))
				if err != nil {
					panic(err)
				}
				bind.WaitDeployed(context.Background(), client, tx)

				time.Sleep(time.Second * 3)

				gatewayBalance := getBalance(erc20Contract, gatewayAddress)
				utils.LogInfo("gatewayBalance = ", gatewayBalance)
			}

			return nil
		},
	}
	return cmd
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

	utils.LogInfo("Gateway was deployed!")

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
