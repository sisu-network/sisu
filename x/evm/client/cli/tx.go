package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sisu-network/cosmos-sdk/client"
	// "github.com/sisu-network/cosmos-sdk/client/flags"
	"github.com/sisu-network/sisu/x/evm/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(GetCmdSubmitEthTx())

	return cmd
}
