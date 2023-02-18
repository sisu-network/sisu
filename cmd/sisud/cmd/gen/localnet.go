package gen

import (
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"

	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	econfig "github.com/sisu-network/deyes/config"
	"github.com/sisu-network/sisu/config"
)

var (
	flagNodeDirPrefix     = "node-dir-prefix"
	flagNumValidators     = "v"
	flagOutputDir         = "output-dir"
	flagNodeDaemonHome    = "node-daemon-home"
	flagStartingIPAddress = "starting-ip-address"
	flagChainId           = "chain-id"
	flagConfigString      = "config-string"
)

type localnetGenerator struct{}

// get cmd to initialize all files for tendermint localnet and application
func LocalnetCmd(mbm module.BasicManager, genBalIterator banktypes.GenesisBalancesIterator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "localnet",
		Short: "Initialize files for a simapp localnet",
		Long: `localnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).
Note, strict routability for addresses is turned off in the config file.
Example:
	For running single instance:
		./sisu localnet
	`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			tmConfig := serverCtx.Config
			tmConfig.LogLevel = ""
			tmConfig.Consensus.TimeoutCommit = time.Second * 3

			generator := &localnetGenerator{}
			startingIPAddress, _ := cmd.Flags().GetString(flagStartingIPAddress)
			genesisFolder, _ := cmd.Flags().GetString(flags.GenesisFolder)

			// Get Chain id and keyring backend from .env file.
			chainID := "sisu-local"
			keyringBackend := keyring.BackendTest

			deyesChains := getDeyesChains(cmd, genesisFolder)

			nodeConfig := config.Config{
				Mode: "dev",
				Sisu: config.SisuConfig{
					ChainId:        chainID,
					KeyringBackend: keyringBackend,
					ApiHost:        "0.0.0.0",
					ApiPort:        25456,
				},
				Tss: config.TssConfig{
					DheartHost: "0.0.0.0",
					DheartPort: 5678,
					DeyesUrl:   "http://0.0.0.0:31001",
				},
			}

			generator.generateEyesToml("../deyes", deyesChains)

			settings := buildBaseSettings(cmd, mbm, genBalIterator)
			settings.tmConfig = tmConfig
			settings.chainID = chainID
			settings.ips = []string{startingIPAddress}
			settings.heartIps = []string{"dheart0"}
			settings.keyringBackend = keyringBackend
			settings.nodeConfigs = []config.Config{nodeConfig}

			_, err := InitNetwork(settings)
			return err
		},
	}

	cmd.Flags().Int(flagNumValidators, 1, "Number of validators to initialize the localnet with")
	cmd.Flags().StringP(flagOutputDir, "o", "./output", "Directory to store initialization data for the localnet")
	cmd.Flags().String(flagNodeDirPrefix, "node", "Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().String(flagNodeDaemonHome, "main", "Home directory of the node's daemon configuration")
	cmd.Flags().String(flagStartingIPAddress, "127.0.0.1", "Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:46656, ID1@192.168.0.2:46656, ...)")
	cmd.Flags().String(server.FlagMinGasPrices, fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom), "Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01photino,0.001stake)")
	cmd.Flags().String(flags.Algo, string(hd.Secp256k1Type), "Key signing algorithm to generate keys for")
	cmd.Flags().String(flags.KeyringBackend, keyring.BackendTest, "Keyring backend. file|os|kwallet|pass|test|memory")
	cmd.Flags().String(flags.GenesisFolder, "./misc/dev", "Relative path to the folder that contains genesis configuration.")

	return cmd
}

func (g *localnetGenerator) getIP(i int, startingIPAddr string) (ip string, err error) {
	if len(startingIPAddr) == 0 {
		ip, err = server.ExternalIP()
		if err != nil {
			return "", err
		}
		return ip, nil
	}
	return g.calculateIP(startingIPAddr, i)
}

func (g *localnetGenerator) calculateIP(ip string, i int) (string, error) {
	ipv4 := net.ParseIP(ip).To4()
	if ipv4 == nil {
		return "", fmt.Errorf("%v: non ipv4 address", ip)
	}

	for j := 0; j < i; j++ {
		ipv4[3]++
	}

	return ipv4.String(), nil
}

func (g *localnetGenerator) generateEyesToml(outputDir string, deyesConfig econfig.Deyes) {
	deyesConfig = updateOracleSecret(deyesConfig)
	deyesConfig.SisuServerUrl = fmt.Sprintf("http://%s:25456", "0.0.0.0")

	writeDeyesConfig(deyesConfig, outputDir)
}
