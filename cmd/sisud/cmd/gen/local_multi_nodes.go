package gen

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/spf13/cobra"
)

const (
	StartingAPIPort         = 25456
	StartingRPCPort         = 36656
	StartingInternalAPIPort = 1317
)

type localMultiNodesGenerator struct{}

func MultiNodesCmd(mbm module.BasicManager, genBalIterator banktypes.GenesisBalancesIterator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multi-nodes",
		Short: "Initialize files for local multi nodes",
		Long:  ``,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			serverCtx := server.GetServerContextFromCmd(cmd)
			tmConfig := serverCtx.Config
			tmConfig.P2P.AddrBookStrict = false
			tmConfig.LogLevel = ""
			tmConfig.Consensus.TimeoutCommit = time.Second * 4

			outputDir, _ := cmd.Flags().GetString(flagOutputDir)
			minGasPrices, _ := cmd.Flags().GetString(server.FlagMinGasPrices)
			nodeDirPrefix, _ := cmd.Flags().GetString(flagNodeDirPrefix)
			nodeDaemonHome, _ := cmd.Flags().GetString(flagNodeDaemonHome)
			//startingIPAddress, _ := cmd.Flags().GetString(flagStartingIPAddress)
			numValidators, _ := cmd.Flags().GetInt(flagNumValidators)
			numCandidates, _ := cmd.Flags().GetInt(flagNumCandidates)
			algo, _ := cmd.Flags().GetString(flags.FlagKeyAlgorithm)

			g := localMultiNodesGenerator{}

			// Make dir folder for mysql docker
			err = os.MkdirAll(filepath.Join(outputDir, "db"), 0755)
			if err != nil {
				panic(err)
			}

			// Get Chain id and keyring backend from .env file.
			chainID := "eth-sisu-local"
			keyringBackend := keyring.BackendTest

			//ips := getLocalIps(startingIPAddress, numValidators)
			ips := make([]string, numValidators)
			for i := 0; i < numValidators; i++ {
				ips[i] = "127.0.0.1"
			}

			nodeConfigs := g.getNodeSettings(chainID, keyringBackend, numValidators)
			//for i := range ips {
			//	dir := filepath.Join(outputDir, fmt.Sprintf("validator%d", i))
			//
			//	if err := os.MkdirAll(dir, nodeDirPerm); err != nil {
			//		panic(err)
			//	}
			//
			//	nodeConfig := g.getNodeSettings(chainID, keyringBackend, i, mysqlIp, ips)
			//	nodeConfigs[i] = nodeConfig
			//
			//	g.generateEyesToml(i, dir)
			//}

			settings := &Setting{
				clientCtx:        clientCtx,
				cmd:              cmd,
				tmConfig:         tmConfig,
				mbm:              mbm,
				genBalIterator:   genBalIterator,
				outputDir:        outputDir,
				chainID:          chainID,
				minGasPrices:     minGasPrices,
				nodeDirPrefix:    nodeDirPrefix,
				nodeDaemonHome:   nodeDaemonHome,
				keyringBackend:   keyringBackend,
				algoStr:          algo,
				numValidators:    numValidators,
				numCandidates:    numCandidates,
				isLocalMultiNode: true,

				ips:         ips,
				nodeConfigs: nodeConfigs,
				tokens:      getTokens("./misc/dev/tokens.json"),
				chains:      getChains("./misc/dev/chains.json"),
				liquidities: getLiquidity("./misc/dev/liquid.json"),
				params:      &types.Params{MajorityThreshold: int32(math.Ceil(float64(numValidators) * 2 / 3))},
			}

			_, err = InitNetwork(settings)
			if err != nil {
				return err
			}

			//peerIds, err := getPeerIds(len(ips), valPubKeys)
			//if err != nil {
			//	panic(err)
			//}

			//for i := range ips {
			//	dir := filepath.Join(outputDir, fmt.Sprintf("node%d", i))
			//	g.generateHeartToml(i, dir, DockerNodeConfig{}, peerIds, valPubKeys)
			//}

			return err
		},
	}

	cmd.Flags().Int(flagNumValidators, 2, "Number of validators to initialize the local multi nodes with")
	cmd.Flags().Int(flagNumCandidates, 0, "Number of candidates to initialize the local multi nodes with")
	cmd.Flags().StringP(flagOutputDir, "o", "./output", "Directory to store initialization data for the local multi nodes")
	cmd.Flags().String(flagNodeDirPrefix, "validator", "Prefix the directory name for each node with (node results in validator0, validator1, ...)")
	cmd.Flags().String(flagNodeDaemonHome, "main", "Home directory of the node's daemon configuration")
	cmd.Flags().String(flagStartingIPAddress, "127.0.0.1", "Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:46656, ID1@192.168.0.2:46656, ...)")
	cmd.Flags().String(server.FlagMinGasPrices, fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom), "Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01photino,0.001stake)")
	cmd.Flags().String(flags.FlagKeyAlgorithm, string(hd.Secp256k1Type), "Key signing algorithm to generate keys for")

	return cmd

}

func (g *localMultiNodesGenerator) getNodeSettings(chainID, keyringBackend string, n int) []config.Config {
	configs := make([]config.Config, n)
	for i := 0; i < n; i++ {
		cfg := config.Config{
			Mode: "dev",
			Sisu: config.SisuConfig{
				ChainId:         chainID,
				KeyringBackend:  keyringBackend,
				ApiHost:         "0.0.0.0",
				ApiPort:         uint16(StartingAPIPort + i),
				RpcPort:         StartingRPCPort + i,
				InternalApiPort: StartingInternalAPIPort + i,
			},
			Tss: config.TssConfig{
				DheartHost: "0.0.0.0",
				DheartPort: 5678,
				DeyesUrl:   fmt.Sprintf("http://0.0.0.0:%d", 31001),
				SupportedChains: map[string]config.TssChainConfig{
					"ganache1": {
						Id: "ganache1",
					},
					"ganache2": {
						Id: "ganache2",
					},
				},
			},
		}

		configs[i] = cfg
	}

	return configs
}
