package gen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/config"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type TestnetGenerator struct {
	ropstenUrl string
}

type TestnetNodes []TestnetNode

type SqlConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Schema   string `json:"schema"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type TestnetNode struct {
	SisuIp  string    `json:"sisu_ip"`
	HeartIp string    `json:"dheart_ip"`
	EyesIp  string    `json:"eyes_ip"`
	Sql     SqlConfig `json:"sql"`
}

// get cmd to initialize all files for tendermint localnet and application
func TestnetCmd(mbm module.BasicManager, genBalIterator banktypes.GenesisBalancesIterator) *cobra.Command {
	cmd := &cobra.Command{
		Use: "testnet",

		Short: "Initialize files for a simapp localnet",
		Long: `privatenet creates configuration for a network with N validators.
Example:
	For multiple nodes (running with docker):
	  ./sisu testnet --v 2 --output-dir ./output --config-string '[{"sisu_ip":"192.168.0.1","dheart_ip":"192.168.0.2","deyes_ip":"192.168.0.3"},{"sisu_ip":"192.168.1.1","dheart_ip":"192.168.1.2","deyes_ip":"192.168.1.3"}]'
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
			tmConfig.Consensus.TimeoutCommit = 5 * time.Second

			outputDir, _ := cmd.Flags().GetString(flagOutputDir)
			minGasPrices, _ := cmd.Flags().GetString(server.FlagMinGasPrices)
			nodeDirPrefix, _ := cmd.Flags().GetString(flagNodeDirPrefix)
			nodeDaemonHome, _ := cmd.Flags().GetString(flagNodeDaemonHome)
			chainId, _ := cmd.Flags().GetString(flagChainId)
			numValidators, _ := cmd.Flags().GetInt(flagNumValidators)
			algo, _ := cmd.Flags().GetString(flags.FlagKeyAlgorithm)
			generator.ropstenUrl, _ = cmd.Flags().GetString(flagRopstenUrl)
			configString, _ := cmd.Flags().GetString(flagConfigString)

			fmt.Println("configString = ", configString)

			testnetNodeData := TestnetNodes{}
			err = json.Unmarshal([]byte(configString), &testnetNodeData)
			if err != nil {
				panic(err)
			}

			chainId = "testnet"

			log.Info("testnet gen: chainId = ", chainId)

			err = os.MkdirAll(outputDir, os.ModePerm)
			if err != nil {
				panic(err)
			}

			// Clean data
			cleanData(outputDir)

			// TODO: Use backend file for keyring
			// keyringBackend := keyring.BackendFile
			keyringBackend := keyring.BackendTest

			monikers := make([]string, numValidators)
			for i := 0; i < numValidators; i++ {
				monikers[i] = "node" + strconv.Itoa(i)
			}

			nodes := testnetNodeData

			sisuIps := make([]string, numValidators)
			heartIps := make([]string, numValidators)

			for i := 0; i < numValidators; i++ {
				sisuIps[i] = nodes[i].SisuIp
				heartIps[i] = nodes[i].HeartIp
			}
			log.Info("ips = ", sisuIps)

			// Create configuration
			nodeConfigs := make([]config.Config, numValidators)
			for i := range sisuIps {
				dir := filepath.Join(outputDir, fmt.Sprintf("node%d", i))

				if err := os.MkdirAll(dir, nodeDirPerm); err != nil {
					panic(err)
				}

				nodeConfig := generator.getNodeSettings(chainId, keyringBackend, nodes[i], len(nodes))
				nodeConfigs[i] = nodeConfig
			}

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
			peerIds, err := getPeerIds(len(sisuIps), valPubKeys)
			if err != nil {
				panic(err)
			}

			for i := range heartIps {
				dir := filepath.Join(outputDir, fmt.Sprintf("node%d", i))

				// Dheart configs
				generator.generateHeartToml(i, dir, heartIps, peerIds, nodes[i].Sql)

				// Deyes configs
				generator.generateEyesToml(i, dir, sisuIps[i])
			}

			return err
		},
	}

	cmd.Flags().Int(flagNumValidators, 4, "Number of validators to initialize the localnet with")
	cmd.Flags().StringP(flagOutputDir, "o", "./output", "Directory to store initialization data for the localnet")
	cmd.Flags().String(flagNodeDirPrefix, "node", "Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().String(flagNodeDaemonHome, "main", "Home directory of the node's daemon configuration")
	cmd.Flags().String(flagTmpDir, "tmp-dir", "Location of temporary directory that contains list of peers ips and other configs.")
	cmd.Flags().String(flagChainId, "sisu-talon-01", "Name of the chain")
	cmd.Flags().String(server.FlagMinGasPrices, fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom), "Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01photino,0.001stake)")
	cmd.Flags().String(flags.FlagKeyAlgorithm, string(hd.Secp256k1Type), "Key signing algorithm to generate keys for")
	cmd.Flags().String(flagRopstenUrl, "", "RPC url for ropsten network")
	cmd.Flags().String(flagConfigString, "", "configuration string for all nodes")

	return cmd
}

func (g *TestnetGenerator) getNodeSettings(chainID string, keyringBackend string, testnetConfig TestnetNode, n int) config.Config {
	majority := (n + 1) * 2 / 3

	return config.Config{
		Mode: "testnet",
		Sisu: config.SisuConfig{
			ChainId:        chainID,
			KeyringBackend: keyringBackend,
			ApiHost:        "0.0.0.0",
			ApiPort:        25456,
		},
		Tss: config.TssConfig{
			MajorityThreshold: majority,
			DheartHost:        testnetConfig.HeartIp,
			DheartPort:        5678,
			DeyesUrl:          fmt.Sprintf("http://%s:31001", testnetConfig.EyesIp),
			SupportedChains: map[string]config.TssChainConfig{
				"ganache1": {
					Symbol: "ganache1",
				},
				"ganache2": {
					Symbol: "ganache2",
				},
			},
		},
	}
}

func (g *TestnetGenerator) generateHeartToml(index int, outputDir string, heartIps []string, peerIds []string, sqlConfig SqlConfig) {
	peers := make([]string, 0, len(peerIds)-1)
	for i := range peerIds {
		if i == index {
			continue
		}

		peers = append(peers, fmt.Sprintf(`"/dns/%s/tcp/28300/p2p/%s"`, heartIps[i], peerIds[i]))
	}

	peerString := strings.Join(peers, ", ")

	sqlConfig.Schema = "dheart"

	writeHeartConfig(index, outputDir, peerString, "false", sqlConfig)
}

func (g *TestnetGenerator) generateEyesToml(index int, dir string, sisuIp string) {
	deyesConfig := DeyesConfiguration{
		Ganaches: []GanacheConfig{
			{
				Ip:    "ganache1",
				Index: 1,
			},
			{
				Ip:    "ganache2",
				Index: 2,
			},
		},

		Sql: SqlConfig{
			Host:     "mysql",
			Port:     3306,
			Schema:   fmt.Sprintf("deyes%d", index),
			Username: "root",
			Password: "password",
		},

		SisuServerUrl: sisuIp,
	}

	writeDeyesConfig(deyesConfig, dir)
}
