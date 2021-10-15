package dev

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/contracts/eth/erc20"
	erc20Gateway "github.com/sisu-network/sisu/contracts/eth/erc20gateway"
	"github.com/sisu-network/sisu/db"
	"github.com/sisu-network/sisu/utils"
	hdwallet "github.com/sisu-network/sisu/utils/hdwallet"
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
			fmt.Println("args = ", args)
			fromChain := args[0]
			contractType := args[1]
			tokenAddressString := args[2]
			toChain := args[3]

			// hash := tss.SupportedContracts[contractType].AbiHash
			// contract := database.GetContractFromHash(fromChain, hash)
			// if contract == nil {
			// 	return fmt.Errorf("cannot find contract")
			// }

			switch contractType {
			case "erc20":
				client, err := getEthClient(fromChain)
				if err != nil {
					return err
				}

				auth, err := getAuthTransactor(client)
				if err != nil {
					return err
				}

				gatewayAddress, tx, instance, err := erc20Gateway.DeployErc20Gateway(auth, client, "eth")
				_ = instance
				_ = gatewayAddress
				if err != nil {
					return err
				}
				bind.WaitDeployed(context.Background(), client, tx)
				utils.LogInfo("Gateway was deployed!")

				tokenAddress := common.HexToAddress(tokenAddressString)
				erc20Contract, err := erc20.NewErc20(tokenAddress, client)
				if err != nil {
					return err
				}

				firstBalance, err := erc20Contract.BalanceOf(&bind.CallOpts{Pending: true}, account0.Address)
				if err != nil {
					return err
				}
				fmt.Println("firstBalance = ", firstBalance)

				utils.LogInfo("Transfering token out....")
				auth, err = getAuthTransactor(client)
				if err != nil {
					return err
				}

				tx, err = instance.TransferOutFromContract(auth, tokenAddress, toChain, account0.Address.String(), big.NewInt(1))
				bind.WaitDeployed(context.Background(), client, tx)

				secondBalance, err := erc20Contract.BalanceOf(&bind.CallOpts{Pending: true}, account0.Address)
				if err != nil {
					return err
				}
				fmt.Println("secondBalance = ", secondBalance)
			}

			return nil
		},
	}
	return cmd
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
