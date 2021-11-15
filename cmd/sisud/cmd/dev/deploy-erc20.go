package dev

import (
	"context"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/contracts/eth/erc20"
	"github.com/spf13/cobra"
)

func DeployErc20() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deploy-erc20",
		Long: `Deploy an ERC20 contract.
Usage:
deploy-erc20 [Port]

Example:
deploy-erc20 1234
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			port, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}

			client, err := getEthClient(port)
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

			log.Info("Deploying erc20....")

			bind.WaitDeployed(context.Background(), client, tx)

			time.Sleep(time.Second * 2)

			log.Info("Deployment done! Contract address: ", address.String())

			// Check sender's balance
			balance, err := instance.BalanceOf(&bind.CallOpts{Pending: true}, account0.Address)
			if err != nil {
				panic(err)
			}

			log.Info("balance of account0 = ", balance)

			return nil
		},
	}

	return cmd
}
