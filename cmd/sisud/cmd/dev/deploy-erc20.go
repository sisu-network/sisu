package dev

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/sisu-network/sisu/contracts/eth/erc20"
	"github.com/sisu-network/sisu/utils"
	"github.com/spf13/cobra"
)

func DeployErc20() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deploy-erc20",
		Long: `Deploy an ERC20 contract.
Usage:
deploy-erc20 [Chain]

Example:
deploy-erc20 eth
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			chain := args[0]

			client, err := getEthClient(chain)
			if err != nil {
				panic(err)
			}

			auth, err := getAuthTransactor(client, account0.Address)
			if err != nil {
				panic(err)
			}

			address, tx, instance, err := erc20.DeployErc20(auth, client, "name", "sisu-token")
			_ = instance
			if err != nil {
				panic(err)
			}

			bind.WaitDeployed(context.Background(), client, tx)

			time.Sleep(time.Second * 2)

			utils.LogInfo("Contract address: ", address.String())

			// Check sender's balance
			balance, err := instance.BalanceOf(&bind.CallOpts{Pending: true}, account0.Address)
			if err != nil {
				panic(err)
			}

			utils.LogInfo("balance of account0 = ", balance)

			return nil
		},
	}

	return cmd
}
