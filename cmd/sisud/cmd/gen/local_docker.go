package gen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	tmos "github.com/sisu-network/tendermint/libs/os"
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
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
)

type DockerNodeConfig struct {
	Ganaches []struct {
		Ip       string
		ChainId  *big.Int
		HostPort int
	}
	MysqlIp string

	NodeData []struct {
		Ip          string
		EthPortMap  int
		GrpcPortMap int
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
	  ./sisu local-docker --v 2 --output-dir ./output
	`,
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
			// startingIPAddress, _ := cmd.Flags().GetString(flagStartingIPAddress)
			numValidators, _ := cmd.Flags().GetInt(flagNumValidators)
			algo, _ := cmd.Flags().GetString(flags.FlagKeyAlgorithm)

			// Clean up
			cleanData(outputDir)

			// Make dir folder for mysql docker
			err = os.MkdirAll(filepath.Join(outputDir, "db"), 0755)
			if err != nil {
				panic(err)
			}

			// Get Chain id and keyring backend from .env file.
			chainID := "eth-sisu-local"
			keyringBackend := keyring.BackendTest

			// startingIPAddress := "192.168.10.6"
			// ips := getLocalIps(startingIPAddress, numValidators)
			ips := make([]string, numValidators)
			for i := range ips {
				ips[i] = fmt.Sprintf("sisu%d", i)
			}

			mysqlIp := "192.168.10.4"
			chainIds := []*big.Int{libchain.GetChainIntFromId("eth-sisu-local"), libchain.GetChainIntFromId("ganache1")}
			dockerConfig := getDockerConfig([]string{"192.168.10.2", "192.168.10.3"}, chainIds, "192.168.10.4", ips)

			nodeConfigs := make([]config.Config, numValidators)
			for i, _ := range ips {
				dir := filepath.Join(outputDir, fmt.Sprintf("node%d", i))

				if err := os.MkdirAll(dir, nodeDirPerm); err != nil {
					panic(err)
				}

				nodeConfig := getNodeSettings(chainID, keyringBackend, i, mysqlIp, ips)
				nodeConfigs[i] = nodeConfig

				generateEyesToml(i, dockerConfig, dir)
			}

			generateDockerCompose(filepath.Join(outputDir, "docker-compose.yml"), ips, dockerConfig)

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

			valPubKeys, err := InitNetwork(settings)
			if err != nil {
				return err
			}

			peerIds, err := getPeerIds(len(ips), valPubKeys)
			if err != nil {
				panic(err)
			}

			for i, _ := range ips {
				dir := filepath.Join(outputDir, fmt.Sprintf("node%d", i))
				generateHeartToml(i, dir, dockerConfig, peerIds)
			}

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

func cleanData(outputDir string) {
	if !utils.IsFileExisted(outputDir) {
		log.Info("Creating output dir...")
		// Create the folder
		os.MkdirAll(outputDir, 0755)
		return
	}

	files, err := ioutil.ReadDir(outputDir)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if f.IsDir() {
			path := filepath.Join(outputDir, f.Name())
			err := os.RemoveAll(path)
			if err != nil {
				panic(err)
			}
		}
	}
}

func getDockerConfig(ganacheIps []string, chainIds []*big.Int, mysqlIp string, ips []string) DockerNodeConfig {
	docker := DockerNodeConfig{
		MysqlIp: mysqlIp,
	}
	docker.Ganaches = make([]struct {
		Ip       string
		ChainId  *big.Int
		HostPort int
	}, len(ganacheIps))
	docker.NodeData = make([]struct {
		Ip          string
		EthPortMap  int
		GrpcPortMap int
	}, len(ips))

	// Ganache
	for i := range ganacheIps {
		docker.Ganaches[i].Ip = ganacheIps[i]
		docker.Ganaches[i].ChainId = chainIds[i]
		docker.Ganaches[i].HostPort = 7545 + i
	}

	// Node data
	for i, ip := range ips {
		docker.NodeData[i].Ip = ip
		docker.NodeData[i].EthPortMap = 1234 + i
		docker.NodeData[i].GrpcPortMap = 9090 + i
	}

	return docker
}

func getNodeSettings(chainID, keyringBackend string, index int, mysqlIp string, ips []string) config.Config {
	return config.Config{
		Mode: "dev",
		Sisu: config.SisuConfig{
			ChainId:        chainID,
			KeyringBackend: keyringBackend,
			ApiHost:        "0.0.0.0",
			ApiPort:        25456,
			Sql: config.SqlConfig{
				Host:     "mysql",
				Port:     3306,
				Username: "root",
				Password: "password",
				Schema:   fmt.Sprintf("sisu%d", index),
			},
		},
		Eth: config.ETHConfig{
			Host:          "0.0.0.0",
			Port:          1234,
			ImportAccount: true,
		},
		Tss: config.TssConfig{
			Enable: true,
			// Enable:     false,
			DheartHost: fmt.Sprintf("dheart%d", index),
			DheartPort: 5678,
			SupportedChains: map[string]config.TssChainConfig{
				"eth-sisu-local": {
					Symbol:   "eth-sisu-local",
					DeyesUrl: fmt.Sprintf("http://deyes%d:31001", index),
				},
				"ganache1": {
					Symbol:   "ganache1",
					DeyesUrl: fmt.Sprintf("http://deyes%d:31001", index),
				},
			},
		},
	}
}

// Ip assignments:
// - 192.168.10.2, .3: ganache(s)
// - 192.168.10.4: mysql
// - 192.168.10.5 onward: for Sisu and others
func generateDockerCompose(outputPath string, ips []string, dockerConfig DockerNodeConfig) {
	const dockerComposeTemplate = `version: "3"
services:
  ganache1:
    image: ganache-cli
    environment:
      - port=7545
      - networkId=189985
    ports:
      - 7545:7545
  mysql:
    image: mysql:8.0.19
    command: "--default-authentication-plugin=mysql_native_password"
    restart: always
    environment:
      - MYSQL_ROOT_PASSWORD=password
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "127.0.0.1", "--silent"]
      interval: 3s
      retries: 5
      start_period: 30s
    volumes:
      - ./db:/var/lib/mysql
    expose:
      - 3306
    ports:
      - 4000:3306
{{ range $k, $nodeData := .NodeData }}
  sisu{{ $k }}:
    container_name: sisu{{ $k }}
    image: "sisu"
    expose:
      - 1317
      - 25456
      - 26656
      - 26657
    ports:
      - {{ $nodeData.EthPortMap }}:1234
      - {{ $nodeData.GrpcPortMap }}:9090
    restart: on-failure
    depends_on:
      - mysql
      - deyes{{ $k }}
      - dheart{{ $k }}
    volumes:
      - ./node{{ $k }}:/root/.sisu
  dheart{{ $k }}:
    image: dheart
    expose:
      - 28300
      - 5678
    restart: on-failure
    depends_on:
      - mysql
    volumes:
      - ./node{{ $k }}/dheart.toml:/root/dheart.toml
  deyes{{ $k }}:
    image: deyes
    restart: on-failure
    depends_on:
      - mysql
      - ganache1
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
	err = configTemplate.Execute(&buffer, dockerConfig)

	tmos.MustWriteFile(outputPath, buffer.Bytes(), 0644)
}

func generateEyesToml(index int, dockerConfig DockerNodeConfig, dir string) {
	data := struct {
		Index         int
		Ganache1      string
		SisuServerUrl string
		SqlSchema     string
		SqlHost       string
	}{
		Index:    index,
		Ganache1: "ganache1",

		SqlHost:   "mysql",
		SqlSchema: fmt.Sprintf("deyes%d", index),

		SisuServerUrl: fmt.Sprintf("http://sisu%d:25456", index),
	}

	eyesToml := `db_host = "{{ .SqlHost }}"
db_port = 3306
db_username = "root"
db_password = "password"
db_schema = "{{ .SqlSchema }}"

server_port = 31001
sisu_server_url = "{{ .SisuServerUrl }}"

[chains]
[chains.eth-sisu-local]
  chain = "eth-sisu-local"
  block_time = 3000
  starting_block = 0
  rpc_url = "http://sisu0:1234"

[chains.ganache1]
  chain = "ganache1"
  block_time = 1000
  starting_block = 0
  rpc_url = "http://ganache1:7545"
`

	tmpl := template.New("eyesToml")

	configTemplate, err := tmpl.Parse(eyesToml)
	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer
	err = configTemplate.Execute(&buffer, data)

	tmos.MustWriteFile(filepath.Join(dir, "deyes.toml"), buffer.Bytes(), 0644)
}

func getPeerIds(n int, pubKeys []cryptotypes.PubKey) ([]string, error) {
	ids := make([]string, n)

	for i := 0; i < n; i++ {
		p2pPubKey, err := crypto.UnmarshalEd25519PublicKey(pubKeys[i].Bytes())
		if err != nil {
			panic(err)
		}

		id, err := peer.IDFromPublicKey(p2pPubKey)
		if err != nil {
			panic(err)
		}

		ids[i] = id.String()
	}

	return ids, nil
}

func generateHeartToml(index int, dir string, dockerConfig DockerNodeConfig, peerIds []string) {
	peers := make([]string, 0, len(peerIds)-1)
	for i := range peerIds {
		if i == index {
			continue
		}

		peers = append(peers, fmt.Sprintf(`"/dns/dheart%d/tcp/28300/p2p/%s"`, i, peerIds[i]))
	}

	peerString := strings.Join(peers, ", ")
	haertConfig := struct {
		PeerString    string
		SisuServerUrl string
		SqlHost       string
		Schema        string
	}{
		PeerString:    peerString,
		SisuServerUrl: fmt.Sprintf("http://sisu%d:25456", index),
		SqlHost:       "mysql",
		Schema:        fmt.Sprintf("dheart%d", index),
	}

	heartToml := `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

home-dir = "/root/"
use-on-memory = false
shortcut-preparams = true
sisu-server-url = "{{ .SisuServerUrl }}"
port = 5678

###############################################################################
###                        Database Configuration                           ###
###############################################################################
[db]
  host = "{{ .SqlHost }}"
  port = 3306
  username = "root"
  password = "password"
  schema = "{{ .Schema }}"
	migration-path = "file://db/migrations/"
[connection]
  host = "0.0.0.0"
  port = 28300
  rendezvous = "rendezvous"
  peers = [{{ .PeerString }}]
`

	tmpl := template.New("heartToml")

	configTemplate, err := tmpl.Parse(heartToml)
	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer
	err = configTemplate.Execute(&buffer, haertConfig)

	tmos.MustWriteFile(filepath.Join(dir, "dheart.toml"), buffer.Bytes(), 0644)
}
