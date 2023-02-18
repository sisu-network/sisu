package cmd

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	sdkflags "github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func RetryTransferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "retry-transfer",
		Long: `Retry a failed transfer.
Usage:
retry-transfer --transferId [transferId] --index [index]

Example:
./sisu retry-transfer --index 0
--transferId ganache1__0xe36b3b53f67eea926a629963e1e74bf14eb3bd6cb8f9c01f03453496364db8b4
--keyring-backend test --from node0 --chain-id=sisu-local
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			transferId, _ := cmd.Flags().GetString(flags.TransferId)
			index, _ := cmd.Flags().GetInt64(flags.Index)

			if len(transferId) == 0 {
				return fmt.Errorf("invalid transfer id")
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewTransferRetryMsg(clientCtx.GetFromAddress().String(), transferId, index)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	sdkflags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(sdkflags.FlagChainID, "", "name of the sisu chain")
	cmd.Flags().String(flags.TransferId, "", "the failed transfer id")
	cmd.Flags().Int64(flags.Index, 0, "index of the command. This index is used to differentiate "+
		"calling this contract multiple times")
	return cmd
}
