package cmd

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	sdkflags "github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func DepositSisuTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deposit-sisu-token",
		Long: `Deposit sisu token
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			amount, _ := cmd.Flags().GetInt64(flags.Amount)
			index, _ := cmd.Flags().GetInt32(flags.Index)

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewDepositSisuTokenMsg(clientCtx.GetFromAddress().String(), amount, index)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	sdkflags.AddTxFlagsToCmd(cmd)

	cmd.Flags().Int64(flags.Amount, 0, "Sisu token amount")
	cmd.Flags().String(sdkflags.FlagChainID, "", "name of the sisu chain")
	cmd.Flags().Int32(flags.Index, 0, "index of the command. This index is used to differentiate calling this contract multiple times")

	return cmd
}
