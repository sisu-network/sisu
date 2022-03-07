package cmd

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	sdkflags "github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/sisu-network/sisu/x/sisu"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func ContractChangeLiquidityAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "contract-change-liquidity",
		Long: `Change liquidity of a gateway.
Usage:
contract-change-liquidity --chain [Chain] --name [ContractName] --newLiquidityAddress [New liquidity pool address] --index [Index of this message]

Example:
./sisu contract-change-liquidity --chain ganache1 --name erc20gateway --newLiquidityAddress 0x2d532C099CA476780c7703610D807948ae47856A --index=0 --from=node0 --keyring-backend test --chain-id=eth-sisu-local -y
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			chain, _ := cmd.Flags().GetString(flags.Chain)
			name, _ := cmd.Flags().GetString(flags.Name)
			newLiquidityAddress, _ := cmd.Flags().GetString(flags.NewLiquidityAddress)
			index, _ := cmd.Flags().GetInt32(flags.Index)

			if len(chain) == 0 {
				return fmt.Errorf("invalid chain %s", chain)
			}

			if len(name) == 0 {
				return fmt.Errorf("invalid name %s", name)
			}

			if len(newLiquidityAddress) == 0 {
				return fmt.Errorf("invalid newLiquidityAddress %s", name)

			}
			hash := sisu.SupportedContracts[name].AbiHash
			if len(hash) == 0 {
				return fmt.Errorf("contract with name %s not supported", name)
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewChangePoolAddressMsg(clientCtx.GetFromAddress().String(), chain, hash, newLiquidityAddress, index)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	sdkflags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String(sdkflags.FlagChainID, "", "name of the sisu chain")
	cmd.Flags().String(flags.Chain, "", "target chain of the command")
	cmd.Flags().String(flags.Name, "", "name of the contract that identifies the contract")
	cmd.Flags().String(flags.NewLiquidityAddress, "", "new liquidity pool address")
	cmd.Flags().Int32(flags.Index, 0, "index of the command. This index is used to differentiate calling this contract multiple times")

	return cmd
}
