package dev

import (
	"github.com/sisu-network/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

func DevCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "dev",
		Short:                      "High level dev command that should be only used for local development.",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(FundAccount())
	cmd.AddCommand(DeployErc20())
	cmd.AddCommand(TransferOut())
	cmd.AddCommand(Query())

	return cmd
}
