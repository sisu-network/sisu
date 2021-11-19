package gen

import (
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"

	"github.com/sisu-network/cosmos-sdk/client"
	"github.com/sisu-network/cosmos-sdk/client/flags"
	"github.com/sisu-network/cosmos-sdk/crypto/hd"
	"github.com/sisu-network/cosmos-sdk/crypto/keyring"
	"github.com/sisu-network/cosmos-sdk/server"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/cosmos-sdk/types/module"
	banktypes "github.com/sisu-network/cosmos-sdk/x/bank/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/sisu/config"
)

var (
	flagNodeDirPrefix     = "node-dir-prefix"
	flagNumValidators     = "v"
	flagOutputDir         = "output-dir"
	flagNodeDaemonHome    = "node-daemon-home"
	flagStartingIPAddress = "starting-ip-address"
	flagEnableTss         = "enable-tss"
	flagTmpDir            = "tmp-dir"
	flagChainId           = "chain-id"
	flagRopstenUrl        = "ropsten-url"
)

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
		sisu localnet --v 1 --output-dir ./output --starting-ip-address 127.0.0.1
	For multiple nodes (running with docker):
	  sisu localnet --v 4 --output-dir ./output --starting-ip-address 192.168.10.2
	`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			serverCtx := server.GetServerContextFromCmd(cmd)
			tmConfig := serverCtx.Config
			tmConfig.LogLevel = ""
			tmConfig.Consensus.TimeoutCommit = time.Second * 3

			outputDir, _ := cmd.Flags().GetString(flagOutputDir)
			minGasPrices, _ := cmd.Flags().GetString(server.FlagMinGasPrices)
			nodeDirPrefix, _ := cmd.Flags().GetString(flagNodeDirPrefix)
			nodeDaemonHome, _ := cmd.Flags().GetString(flagNodeDaemonHome)
			startingIPAddress, _ := cmd.Flags().GetString(flagStartingIPAddress)
			numValidators, _ := cmd.Flags().GetInt(flagNumValidators)
			algo, _ := cmd.Flags().GetString(flags.FlagKeyAlgorithm)
			enableTss, _ := cmd.Flags().GetBool(flagEnableTss)

			// Get Chain id and keyring backend from .env file.
			chainID := "eth-sisu-local"
			keyringBackend := keyring.BackendTest

			nodeConfig := config.Config{
				Mode: "dev",
				Sisu: config.SisuConfig{
					ChainId:        chainID,
					KeyringBackend: keyringBackend,
					ApiHost:        "0.0.0.0",
					ApiPort:        25456,
					Sql: config.SqlConfig{
						Host:     "0.0.0.0",
						Port:     3306,
						Username: "root",
						Password: "password",
						Schema:   "sisu",
					},
				},
				Eth: config.ETHConfig{
					Host:          "0.0.0.0",
					Port:          1234,
					ImportAccount: true,
				},
				Tss: config.TssConfig{
					Enable:     enableTss,
					DheartHost: "0.0.0.0",
					DheartPort: 5678,
					SupportedChains: map[string]config.TssChainConfig{
						"eth-sisu-local": {
							Symbol:   "eth-sisu-local",
							Id:       int(libchain.GetChainIntFromId("eth-sisu-local").Int64()),
							DeyesUrl: "http://0.0.0.0:31001",
						},
						"ganache1": {
							Symbol:   "ganache1",
							Id:       int(libchain.GetChainIntFromId("ganache1").Int64()),
							DeyesUrl: "http://0.0.0.0:31001",
						},
					},
				},
			}

			settings := &Setting{
				clientCtx:      clientCtx,
				cmd:            cmd,
				tmConfig:       tmConfig,
				mbm:            mbm,
				genBalIterator: genBalIterator,
				outputDir:      outputDir,
				chainID:        chainID,
				minGasPrices:   minGasPrices,
				nodeDirPrefix:  nodeDirPrefix,
				nodeDaemonHome: nodeDaemonHome,
				ips:            getLocalIps(startingIPAddress, numValidators),
				keyringBackend: keyringBackend,
				algoStr:        algo,
				numValidators:  numValidators,
				nodeConfigs:    []config.Config{nodeConfig},
			}

			_, err = InitNetwork(settings)
			return err
		},
	}

	cmd.Flags().Int(flagNumValidators, 1, "Number of validators to initialize the localnet with")
	cmd.Flags().StringP(flagOutputDir, "o", "./output", "Directory to store initialization data for the localnet")
	cmd.Flags().String(flagNodeDirPrefix, "node", "Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().String(flagNodeDaemonHome, "main", "Home directory of the node's daemon configuration")
	cmd.Flags().String(flagStartingIPAddress, "127.0.0.1", "Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:46656, ID1@192.168.0.2:46656, ...)")
	cmd.Flags().String(server.FlagMinGasPrices, fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom), "Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01photino,0.001stake)")
	cmd.Flags().String(flags.FlagKeyAlgorithm, string(hd.Secp256k1Type), "Key signing algorithm to generate keys for")
	cmd.Flags().Bool(flagEnableTss, false, "Enable Tss. By default, this value is set to false.")

	return cmd
}

func getLocalIps(startingIPAddress string, count int) []string {
	ips := make([]string, count)
	for i := 0; i < count; i++ {
		ip, err := getIP(i, startingIPAddress)
		if err != nil {
			panic(err)
		}
		ips[i] = ip
	}

	return ips
}

func getIP(i int, startingIPAddr string) (ip string, err error) {
	if len(startingIPAddr) == 0 {
		ip, err = server.ExternalIP()
		if err != nil {
			return "", err
		}
		return ip, nil
	}
	return calculateIP(startingIPAddr, i)
}

func calculateIP(ip string, i int) (string, error) {
	ipv4 := net.ParseIP(ip).To4()
	if ipv4 == nil {
		return "", fmt.Errorf("%v: non ipv4 address", ip)
	}

	for j := 0; j < i; j++ {
		ipv4[3]++
	}

	return ipv4.String(), nil
}
