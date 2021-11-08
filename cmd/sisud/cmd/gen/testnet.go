package gen

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/spf13/cobra"

	"github.com/sisu-network/cosmos-sdk/client"
	"github.com/sisu-network/cosmos-sdk/client/flags"
	"github.com/sisu-network/cosmos-sdk/crypto/hd"
	"github.com/sisu-network/cosmos-sdk/crypto/keyring"
	"github.com/sisu-network/cosmos-sdk/server"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/cosmos-sdk/types/module"
	banktypes "github.com/sisu-network/cosmos-sdk/x/bank/types"
)

type TestnetNode struct {
	SisuIp  string `json:"sisu_ip"`
	HeartIp string `json:"heart_ip"`
	EyesIp  string `json:"eyes_ip"`
}

// type TestnetNodes struct {
// 	Nodes []TestnetNode `json:"nodes"`
// }

const (
	flagNodeFile = "node-file"
	flagTomlFile = "toml-file"
	flagChainId  = "chain-id"
)

// get cmd to initialize all files for tendermint localnet and application
func TestnetCmd(mbm module.BasicManager, genBalIterator banktypes.GenesisBalancesIterator) *cobra.Command {
	cmd := &cobra.Command{
		Use: "testnet",

		Short: "Initialize files for a simapp localnet",
		Long: `privatenet creates configuration for a network with N validators.
Example:
	For multiple nodes (running with docker):
	  ./sisu testnet --v 2 --output-dir ./output --chain-id testnet --node-file tmp/ips.json --toml-file tmp/sisu.toml
	`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			serverCtx := server.GetServerContextFromCmd(cmd)
			tmConfig := serverCtx.Config

			outputDir, _ := cmd.Flags().GetString(flagOutputDir)
			minGasPrices, _ := cmd.Flags().GetString(server.FlagMinGasPrices)
			nodeDirPrefix, _ := cmd.Flags().GetString(flagNodeDirPrefix)
			nodeDaemonHome, _ := cmd.Flags().GetString(flagNodeDaemonHome)
			nodeFile, _ := cmd.Flags().GetString(flagNodeFile)
			chainId, _ := cmd.Flags().GetString(flagChainId)
			numValidators, _ := cmd.Flags().GetInt(flagNumValidators)
			flagTomlFile, _ := cmd.Flags().GetString(flagTomlFile)
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

			node := readNodeFile(nodeFile)
			ips := make([]string, numValidators)

			for i := 0; i < numValidators; i++ {
				ips[i] = node.SisuIp
			}
			utils.LogInfo("ips = ", ips)

			nodeConfigs := getTestnetNodeSettings(numValidators, flagTomlFile)

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

				ips:         ips,
				nodeConfigs: nodeConfigs,
			}

			_, err = InitNetwork(settings)
			return err
		},
	}

	cmd.Flags().Int(flagNumValidators, 4, "Number of validators to initialize the localnet with")
	cmd.Flags().StringP(flagOutputDir, "o", "./output", "Directory to store initialization data for the localnet")
	cmd.Flags().String(flagNodeDirPrefix, "node", "Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().String(flagNodeDaemonHome, "main", "Home directory of the node's daemon configuration")
	cmd.Flags().String(flagNodeFile, "tmp/ips.json", "Location of ips.json file")
	cmd.Flags().String(flagTomlFile, "tmp/sisu.toml", "Location of sisu.toml file")
	cmd.Flags().String(flagChainId, "talon-01", "Name of the chain")
	cmd.Flags().String(server.FlagMinGasPrices, fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom), "Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01photino,0.001stake)")
	cmd.Flags().String(flags.FlagKeyAlgorithm, string(hd.Secp256k1Type), "Key signing algorithm to generate keys for")

	return cmd
}

func readNodeFile(path string) *TestnetNode {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	nodeConfigs := new(TestnetNode)
	err = json.Unmarshal([]byte(content), nodeConfigs)
	if err != nil {
		panic(err)
	}

	return nodeConfigs
}

func getTestnetNodeSettings(numValidators int, flagTomlFile string) []config.Config {
	nodeConfigs := make([]config.Config, numValidators)

	for i := 0; i < numValidators; i++ {
		cfg := config.Config{}
		_, err := toml.DecodeFile(flagTomlFile, &cfg)
		if err != nil {
			panic(err)
		}

		nodeConfigs[i] = cfg
	}

	return nodeConfigs
}
