package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(GetCmdPauseGw())
	for _, subCmd := range cmd.Commands() {
		flags.AddTxFlagsToCmd(subCmd)
	}

	return cmd
}

func GetCmdPauseGw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause [chain]",
		Short: "Pause gateway for a single chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			chain := args[0]

			msg := types.NewMsgPauseGw(clientCtx.GetFromAddress(), chain)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			log.Debugf("Pausing gateway for chain: %s ", chain)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}
