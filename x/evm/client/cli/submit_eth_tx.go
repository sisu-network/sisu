package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	cTx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/sisu-network/sisu/x/evm/types"
	"github.com/spf13/cobra"
)

func GetCmdSubmitEthTx() *cobra.Command {
	return &cobra.Command{
		Use:   "submit-tx [base64Bytes]",
		Short: "Submit a new ETH transaction",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgEthTx(clientCtx.GetFromAddress().String(), []byte{})

			return cTx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}
