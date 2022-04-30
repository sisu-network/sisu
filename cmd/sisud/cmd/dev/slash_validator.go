package dev

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	sdkflags "github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func SlashValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "slash-validator",
		Long: `Slash validator
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			slashPoint, _ := cmd.Flags().GetInt64(flags.Amount)
			index, _ := cmd.Flags().GetInt32(flags.Index)
			nodeAddress, _ := cmd.Flags().GetString(flags.NodeAddress)
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewSlashValidatorMsg(clientCtx.GetFromAddress().String(), nodeAddress, slashPoint, index)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	sdkflags.AddTxFlagsToCmd(cmd)

	cmd.Flags().Int64(flags.Amount, 0, "Slash point amount")
	cmd.Flags().String(sdkflags.FlagChainID, "", "name of the sisu chain")
	cmd.Flags().Int32(flags.Index, 0, "index of the command. This index is used to differentiate calling this contract multiple times")
	cmd.Flags().String(flags.NodeAddress, "", "node address")

	return cmd
}
