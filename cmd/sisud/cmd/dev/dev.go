package dev

import (
	"github.com/cosmos/cosmos-sdk/client"
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

	cmd.AddCommand(FundSisu())
	cmd.AddCommand(DeployContract())
	cmd.AddCommand(AddLiquidity())
	cmd.AddCommand(DeployAndFund())
	cmd.AddCommand(Swap())
	cmd.AddCommand(StressSwap())
	cmd.AddCommand(Query())
	cmd.AddCommand(RemoteCall())
	cmd.AddCommand(DeployExampleContract())

	return cmd
}
