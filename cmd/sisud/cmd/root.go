package cmd

import (
	"os"
	"strings"

	"github.com/sisu-network/sisu/app/params"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	tmcli "github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingcli "github.com/cosmos/cosmos-sdk/x/auth/vesting/client/cli"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/sisu-network/sisu/app"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/dev"
	gen "github.com/sisu-network/sisu/cmd/sisud/cmd/gen"
)

var ChainID string

// NewRootCmd creates a new root command for simd. It is called once in the
// main function.
func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	// Set config for prefixes
	app.SetConfig()

	encodingConfig := app.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithJSONMarshaler(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(app.MainAppHome)

	rootCmd := &cobra.Command{
		Use:   app.Name,
		Short: "Stargate CosmosHub App",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			return server.InterceptConfigsPreRunHandler(cmd)
		},
	}

	initRootCmd(rootCmd, encodingConfig)
	overwriteFlagDefaults(rootCmd, map[string]string{
		flags.FlagChainID: ChainID,
	})

	changeDescription(rootCmd)

	return rootCmd, encodingConfig
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig) {
	authclient.Codec = encodingConfig.Marshaler

	rootCmd.AddCommand(
		genutilcli.InitCmd(app.ModuleBasics, app.MainAppHome),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, app.MainAppHome),
		genutilcli.MigrateGenesisCmd(),
		genutilcli.GenTxCmd(app.ModuleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, app.MainAppHome),
		genutilcli.ValidateGenesisCmd(app.ModuleBasics),
		AddGenesisAccountCmd(app.MainAppHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		debug.Cmd(),
		dev.DevCmd(),
		// this line is used by starport scaffolding # stargate/root/commands
	)

	a := appCreator{
		encCfg: encodingConfig,
	}
	server.AddCommands(rootCmd, app.MainAppHome, a.newApp, a.appExport, addModuleInitFlags)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		queryCommand(),
		txCommand(),
		keys.Commands(app.MainAppHome),
		gen.LocalnetCmd(app.ModuleBasics, banktypes.GenesisBalancesIterator{}),
		gen.LocalDockerCmd(app.ModuleBasics, banktypes.GenesisBalancesIterator{}),
		gen.TestnetDockerCmd(app.ModuleBasics, banktypes.GenesisBalancesIterator{}),
		gen.TestnetCmd(app.ModuleBasics, banktypes.GenesisBalancesIterator{}),
		resetCmd(),
	)
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
	// this line is used by starport scaffolding # stargate/root/initFlags
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetAccountCmd(),
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
	)

	app.ModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetValidateSignaturesCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		flags.LineBreak,
		vestingcli.GetTxCmd(),
	)

	app.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func overwriteFlagDefaults(c *cobra.Command, defaults map[string]string) {
	set := func(s *pflag.FlagSet, key, val string) {
		if f := s.Lookup(key); f != nil {
			f.DefValue = val
			f.Value.Set(val)
		}
	}
	for key, val := range defaults {
		set(c.Flags(), key, val)
		set(c.PersistentFlags(), key, val)
	}
	for _, c := range c.Commands() {
		overwriteFlagDefaults(c, defaults)
	}
}

// change cosmos prefix to terra
func changeDescription(command *cobra.Command) {
	childCommands := command.Commands()
	if len(childCommands) == 0 {
		return
	}

	for _, childCommand := range childCommands {
		childCommand.Long = strings.ReplaceAll(childCommand.Long, "cosmos", "sisu")
		childCommand.Long = strings.ReplaceAll(childCommand.Long, "<appcli>", "sisucli")
		childCommand.Long = strings.ReplaceAll(childCommand.Long, "Atoms", "Sisu")

		changeDescription(childCommand)
	}
}
