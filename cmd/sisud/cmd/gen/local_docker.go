package gen

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	tmos "github.com/sisu-network/tendermint/libs/os"
	"github.com/spf13/cobra"

	"github.com/sisu-network/cosmos-sdk/client"
	"github.com/sisu-network/cosmos-sdk/client/flags"
	"github.com/sisu-network/cosmos-sdk/crypto/hd"
	"github.com/sisu-network/cosmos-sdk/crypto/keyring"
	"github.com/sisu-network/cosmos-sdk/server"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/cosmos-sdk/types/module"
	banktypes "github.com/sisu-network/cosmos-sdk/x/bank/types"
	"github.com/sisu-network/sisu/config"
)

type DockerNodeConfig struct {
	Ganaches []struct {
		chainName string
	}

	NodeData []struct {
		Ip     string
		Config config.Config
	}
}

// get cmd to initialize all files for tendermint localnet and application
func LocalDockerCmd(mbm module.BasicManager, genBalIterator banktypes.GenesisBalancesIterator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "local-docker",
		Short: "Initialize files for a simapp localnet",
		Long: `localnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).
Note, strict routability for addresses is turned off in the config file.
Example:
	For multiple nodes (running with docker):
	  sisu local-docker --v 4 --output-dir ./output --starting-ip-address 192.168.10.2
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
			startingIPAddress, _ := cmd.Flags().GetString(flagStartingIPAddress)
			numValidators, _ := cmd.Flags().GetInt(flagNumValidators)
			algo, _ := cmd.Flags().GetString(flags.FlagKeyAlgorithm)

			// Get Chain id and keyring backend from .env file.
			chainID := "sisu-dev"
			keyringBackend := keyring.BackendTest

			ips := getLocalIps(startingIPAddress, numValidators)

			nodeConfigs := make([]config.Config, numValidators)
			for i, _ := range ips {
				nodeConfig := getNodeSettings(chainID, keyringBackend, i)
				nodeConfigs[i] = nodeConfig

				generateDeyesToml(i, nodeConfig, filepath.Join(outputDir, fmt.Sprintf("node%d", i)))
			}

			generateDockerCompose(filepath.Join(outputDir, "docker-compose.yml"), ips, nodeConfigs)

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
				keyringBackend: keyringBackend,
				algoStr:        algo,
				numValidators:  numValidators,

				ips:         ips,
				nodeConfigs: nodeConfigs,
			}

			return InitNetwork(settings)
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

func getNodeSettings(chainID, keyringBackend string, index int) config.Config {
	return config.Config{
		Mode: "dev",
		Sisu: config.SisuConfig{
			ChainId:        chainID,
			KeyringBackend: keyringBackend,
			ApiHost:        "0.0.0.0",
			ApiPort:        25456,
			Sql: config.SqlConfig{
				Host:     fmt.Sprintf("sisu%d", index),
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
			Enable:     true,
			DheartHost: fmt.Sprintf("dheart%d", index),
			DheartPort: 5678,
			SupportedChains: map[string]config.TssChainConfig{
				"eth": {
					Symbol:   "eth",
					Id:       1,
					DeyesUrl: fmt.Sprintf("http://%s:31001", fmt.Sprintf("deyes%d", index)),
				},
				"sisu-eth": {
					Symbol:   "sisu-eth",
					Id:       36767,
					DeyesUrl: fmt.Sprintf("http://%s:31001", fmt.Sprintf("deyes%d", index)),
				},
			},
		},
	}
}

func generateDockerCompose(outputPath string, ips []string, nodeConfigs []config.Config) {
	docker := DockerNodeConfig{}
	docker.Ganaches = make([]struct{ chainName string }, 2)
	docker.NodeData = make([]struct {
		Ip     string
		Config config.Config
	}, len(ips))

	for i, ip := range ips {
		docker.NodeData[i].Ip = ip
		docker.NodeData[i].Config = nodeConfigs[i]
	}

	const dockerComposeTemplate = `version: "3"{{ $ganaches := .Ganaches }}
services:{{ range $k, $ganache := $ganaches }}
	ganache{{ $k }}:
		build: ganache-cli
		environment:
			- port=7545
			- networkId=1
{{ end }}
{{ range $k, $nodeData := .NodeData }}
	mysql{{ $k }}:
		image: mysql:8.0.19
		command: "--default-authentication-plugin=mysql_native_password"
		restart: always
		healthcheck:
			test: ["CMD", "mysqladmin", "ping", "-h", "127.0.0.1", "--silent"]
			interval: 3s
			retries: 5
			start_period: 30s
		volumes:
			- db-data:/var/lib/mysql
	sisu{{ $k }}:
		container_name: node{{ $k }}
		image: "sisu"
    ports:
      - "26656-26657:26656-26657"
      - "1317:1317"
      - "9090:9090"
    ports:
      - 26657:26657
      - 26656:26656
      - 9090:9090
      - 25456:25456
    restart: on-failure
    depends_on:
      - deyes{{ $k }}
      - dheart{{ $k }}
		volumes:
      - ./node{{ $k }}:/root/.sisu
	deyes{{ $k }}:
		build: deyes
		ports:
			- 31001:31001
		restart: on-failure
		depends_on:
			- mysql{{ range $j, $p := $ganaches }}
			- ganache{{ $j }}{{ end }}
		volumes:
			- ./node{{ $k }}/deyes.toml:/app/deyes.toml
{{ end }}
`

	tmpl := template.New("localDockerCompose")

	configTemplate, err := tmpl.Parse(dockerComposeTemplate)
	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer
	err = configTemplate.Execute(&buffer, docker)

	tmos.MustWriteFile(outputPath, buffer.Bytes(), 0644)
}

func generateDeyesToml(index int, nodeConfig config.Config, dir string) {
	data := struct {
		Index      int
		NodeConfig config.Config
	}{
		Index:      index,
		NodeConfig: nodeConfig,
	}

	eyesToml := `db_host = "mysql{{ .Index  }}"
db_port = 3306
db_username = "root"
db_password = "password"
db_schema = "deyes"

server_port = 31001
sisu_server_url = "http://sisu{{ .Index }}:25456"

[chains]
[chains.eth]
  chain = "eth"
  block_time = 1000
  starting_block = 0
  rpc_url = "http://ganache0:7545"

[chains.sisu-eth]
  chain = "sisu-eth"
  block_time = 1000
  starting_block = 0
  rpc_url = "http://ganache1:8545"
`

	tmpl := template.New("localDockerCompose")

	configTemplate, err := tmpl.Parse(eyesToml)
	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer
	err = configTemplate.Execute(&buffer, data)

	tmos.MustWriteFile(filepath.Join(dir, "deyes.toml"), buffer.Bytes(), 0644)
}
