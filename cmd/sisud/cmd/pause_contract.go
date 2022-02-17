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

func PauseContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "pause-contract",
		Long: `Pause an ERC20 contract.
Usage:
pause-contract --chain [Chain] --name [ContractName] --index [Pause Command Index]

Example:
./sisu pause-contract --chain ganache1 --name erc20gateway --index 0
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			chain, _ := cmd.Flags().GetString(flags.Chain)
			name, _ := cmd.Flags().GetString(flags.Name)
			index, _ := cmd.Flags().GetInt(flags.Index)

			if len(chain) == 0 {
				return fmt.Errorf("invalid chain %s", chain)
			}

			if len(name) == 0 {
				return fmt.Errorf("invalid name %s", name)
			}

			hash := sisu.SupportedContracts[name].AbiHash
			if len(hash) == 0 {
				return fmt.Errorf("contract with name %s not supported", name)
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewPauseContractMsg(clientCtx.GetFromAddress().String(), chain, hash, index)
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
	cmd.Flags().Int(flags.Index, 0, "index of the command. This index is used to differentiate calling this contract multiple times")

	return cmd
}
