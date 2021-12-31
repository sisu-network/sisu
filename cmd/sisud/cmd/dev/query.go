package dev

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/contracts/eth/erc20gateway"
	"github.com/sisu-network/sisu/x/tss"
	"github.com/spf13/cobra"
)

func Query() *cobra.Command {
	cmd := &cobra.Command{
		Use: "query",
		Long: `Query ERC20 or ERC721 balance of an address on a particular chain. Please note that the asset id is a global id cross chain.
Usage:
query [ContractType] [chain] [Port] [AssetId] [AccountAddress]

Example:
query erc20 ganache2 8545 ganache1__0xf0D676183dD5ae6b370adDdbE770235F23546f9d 0xE8382821BD8a0F9380D88e2c5c33bc89Df17E466
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("args = ", args)
			contractType := args[0]
			chain := args[1]
			port, err := strconv.Atoi(args[2])
			if err != nil {
				panic(err)
			}
			assetId := args[3]
			accountAddr := args[4]

			switch contractType {
			case "erc20":
				client, err := getEthClient(port)
				if err != nil {
					panic(err)
				}

				hash := tss.SupportedContracts[contractType].AbiHash
				contract := queryContract(cmd, chain, hash)
				if contract == nil {
					return fmt.Errorf("cannot find contract")
				}

				log.Info("contract.Address = ", contract.Address)

				gatewayAddress := common.HexToAddress(contract.Address)
				gateway, err := erc20gateway.NewErc20gateway(gatewayAddress, client)
				if err != nil {
					panic(err)
				}

				balance, err := gateway.GetBalance(&bind.CallOpts{Pending: true}, assetId, common.HexToAddress(accountAddr))
				if err != nil {
					panic(err)
				}

				log.Info("balance = ", balance)
			}

			return nil
		},
	}

	return cmd
}
