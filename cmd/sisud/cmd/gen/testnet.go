package gen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/BurntSushi/toml"
	heartcfg "github.com/sisu-network/dheart/core/config"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/spf13/cobra"

	"github.com/sisu-network/cosmos-sdk/client"
	"github.com/sisu-network/cosmos-sdk/client/flags"
	"github.com/sisu-network/cosmos-sdk/crypto/hd"
	"github.com/sisu-network/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/sisu-network/cosmos-sdk/crypto/types"
	"github.com/sisu-network/cosmos-sdk/server"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/cosmos-sdk/types/module"
	banktypes "github.com/sisu-network/cosmos-sdk/x/bank/types"
)

type TestnetGenerator struct {
}

type TestnetNodeConfig struct {
	SisuIp  string `json:"sisu_ip"`
	HeartIp string `json:"dheart_ip"`
	EyesIp  string `json:"eyes_ip"`
}

// get cmd to initialize all files for tendermint localnet and application
func TestnetCmd(mbm module.BasicManager, genBalIterator banktypes.GenesisBalancesIterator) *cobra.Command {
	cmd := &cobra.Command{
		Use: "testnet",

		Short: "Initialize files for a simapp localnet",
		Long: `privatenet creates configuration for a network with N validators.
Example:
	For multiple nodes (running with docker):
	  ./sisu testnet --v 2 --output-dir ./output --chain-id testnet --tmp-dir tmp
	`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			generator := &TestnetGenerator{}

			serverCtx := server.GetServerContextFromCmd(cmd)
			tmConfig := serverCtx.Config
			tmConfig.LogLevel = "info"

			outputDir, _ := cmd.Flags().GetString(flagOutputDir)
			minGasPrices, _ := cmd.Flags().GetString(server.FlagMinGasPrices)
			nodeDirPrefix, _ := cmd.Flags().GetString(flagNodeDirPrefix)
			nodeDaemonHome, _ := cmd.Flags().GetString(flagNodeDaemonHome)
			tempDir, _ := cmd.Flags().GetString(flagTmpDir)
			chainId, _ := cmd.Flags().GetString(flagChainId)
			numValidators, _ := cmd.Flags().GetInt(flagNumValidators)
			algo, _ := cmd.Flags().GetString(flags.FlagKeyAlgorithm)

			err = os.MkdirAll(outputDir, os.ModePerm)
			if err != nil {
				panic(err)
			}

			cleanData(outputDir)

			// TODO: Use backend file for keyring
			// keyringBackend := keyring.BackendFile
			keyringBackend := keyring.BackendTest

			monikers := make([]string, numValidators)
			for i := 0; i < numValidators; i++ {
				monikers[i] = "node-talon-" + strconv.Itoa(i)
			}

			sisuIps := make([]string, numValidators)
			heartIps := make([]string, numValidators)
			nodes := generator.readTestnetNodes(tempDir, numValidators)

			for i := 0; i < numValidators; i++ {
				sisuIps[i] = nodes[i].SisuIp
				heartIps[i] = nodes[i].HeartIp
			}
			utils.LogInfo("ips = ", sisuIps)

			nodeConfigs := generator.getTestnetNodeSettings(tempDir, numValidators)

			settings := &Setting{
				clientCtx:      clientCtx,
				cmd:            cmd,
				tmConfig:       tmConfig,
				mbm:            mbm,
				genBalIterator: genBalIterator,
				outputDir:      outputDir,
				chainID:        chainId,
				minGasPrices:   minGasPrices,
				nodeDirPrefix:  nodeDirPrefix,
				nodeDaemonHome: nodeDaemonHome,
				keyringBackend: keyringBackend,
				algoStr:        algo,
				numValidators:  numValidators,

				ips:         sisuIps,
				nodeConfigs: nodeConfigs,
			}

			valPubKeys, err := InitNetwork(settings)

			for i := 0; i < numValidators; i++ {
				inputDir := filepath.Join(tempDir, fmt.Sprintf("node%d", i))
				outputDir := filepath.Join(outputDir, fmt.Sprintf("node%d", i))

				generator.generateHeartToml(
					i,
					inputDir,
					outputDir,
					*nodes[i],
					heartIps,
					valPubKeys,
				)

				generator.generateEyesToml(i, inputDir, outputDir)
			}

			return err
		},
	}

	cmd.Flags().Int(flagNumValidators, 4, "Number of validators to initialize the localnet with")
	cmd.Flags().StringP(flagOutputDir, "o", "./output", "Directory to store initialization data for the localnet")
	cmd.Flags().String(flagNodeDirPrefix, "node", "Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().String(flagNodeDaemonHome, "main", "Home directory of the node's daemon configuration")
	cmd.Flags().String(flagTmpDir, "tmp-dir", "Location of temporary directory that contains list of peers ips and other configs.")
	cmd.Flags().String(flagChainId, "talon-01", "Name of the chain")
	cmd.Flags().String(server.FlagMinGasPrices, fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom), "Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01photino,0.001stake)")
	cmd.Flags().String(flags.FlagKeyAlgorithm, string(hd.Secp256k1Type), "Key signing algorithm to generate keys for")

	return cmd
}

func (g *TestnetGenerator) readTestnetNodes(root string, numValidators int) []*TestnetNodeConfig {
	nodeConfigs := make([]*TestnetNodeConfig, numValidators)

	for i := 0; i < numValidators; i++ {
		path := filepath.Join(root, fmt.Sprintf("node%d", i), "ips.json")
		content, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}

		nodeConfig := new(TestnetNodeConfig)
		err = json.Unmarshal([]byte(content), nodeConfig)
		if err != nil {
			panic(err)
		}

		nodeConfigs[i] = nodeConfig
	}

	return nodeConfigs
}

func (g *TestnetGenerator) getTestnetNodeSettings(root string, numValidators int) []config.Config {
	nodeConfigs := make([]config.Config, numValidators)

	for i := 0; i < numValidators; i++ {
		cfg := config.Config{}
		path := filepath.Join(root, fmt.Sprintf("node%d", i), "sisu.toml")
		_, err := toml.DecodeFile(path, &cfg)
		if err != nil {
			panic(err)
		}

		nodeConfigs[i] = cfg
	}

	return nodeConfigs
}

func (g *TestnetGenerator) generateHeartToml(index int, inputDir string, outputDir string, testNodeConfig TestnetNodeConfig, heartIps []string, valPubKeys []cryptotypes.PubKey) {
	peerIds, err := getPeerIds(len(heartIps), valPubKeys)
	if err != nil {
		panic(err)
	}

	peers := make([]string, 0, len(peerIds)-1)
	for i := range peerIds {
		if i == index {
			continue
		}

		peers = append(peers, fmt.Sprintf(`"/ip4/%s/tcp/28300/p2p/%s"`, heartIps[i], peerIds[i]))
	}

	heartConfig, err := heartcfg.ReadConfig(filepath.Join(inputDir, "dheart.toml"))
	if err != nil {
		panic(err)
	}

	heartConfig.Connection.BootstrapPeers = peers
	heartcfg.WriteConfigFile(filepath.Join(outputDir, "dheart.toml"), heartConfig)
}

func (g *TestnetGenerator) generateEyesToml(index int, inputDir string, outputDir string) {
	// Simply copy the input file to the output file
	data, err := ioutil.ReadFile(filepath.Join(inputDir, "deyes.toml"))
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filepath.Join(outputDir, "deyes.toml"), data, 0644)
	if err != nil {
		panic(err)
	}
}
