package dev

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sisu-network/sisu/contracts/eth/erc20gateway"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss"
	"github.com/spf13/cobra"
)

func Query() *cobra.Command {
	cmd := &cobra.Command{
		Use: "query",
		Long: `Transfer an ERC20 or ERC721 asset.
Usage:
query [ContractType] [chain] [AssetId] [AccountAddress]

Example:
query erc20 sisu-eth eth__0xB369Be7F62cfb3F44965db83404997Fa6EC9Dd58 0xE8382821BD8a0F9380D88e2c5c33bc89Df17E466
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			database := getDatabase()
			defer database.Close()

			utils.LogInfo("args = ", args)
			contractType := args[0]
			chain := args[1]
			assetId := args[2]
			accountAddr := args[3]

			switch contractType {
			case "erc20":
				client, err := getEthClient(chain)
				if err != nil {
					panic(err)
				}

				hash := tss.SupportedContracts[contractType].AbiHash
				contract := database.GetContractFromHash(chain, hash)
				if contract == nil {
					return fmt.Errorf("cannot find contract")
				}
				gatewayAddress := common.HexToAddress(contract.Address)
				gateway, err := erc20gateway.NewErc20gateway(gatewayAddress, client)
				if err != nil {
					panic(err)
				}

				balance, err := gateway.GetBalance(&bind.CallOpts{Pending: true}, assetId, common.HexToAddress(accountAddr))
				if err != nil {
					panic(err)
				}

				utils.LogInfo("balance = ", balance)
			}

			return nil
		},
	}

	return cmd
}
