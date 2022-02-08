package gen

import (
	"bytes"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/spf13/cobra"
	tmos "github.com/tendermint/tendermint/libs/os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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
		GrpcPortMap int
	}
}

type HeartConfiguration struct {
	PeerString    string
	SisuServerUrl string
	UseOnMemory   string
	Sql           SqlConfig
}

type GanacheConfig struct {
	Ip    string
	Index int
}

type DeyesConfiguration struct {
	Ganaches      []GanacheConfig
	SisuServerUrl string

	// sql
	Sql SqlConfig
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
			for i := range ips {
				dir := filepath.Join(outputDir, fmt.Sprintf("node%d", i))

				if err := os.MkdirAll(dir, nodeDirPerm); err != nil {
					panic(err)
				}

				nodeConfig := getNodeSettings(chainID, keyringBackend, i, mysqlIp, ips)
				nodeConfigs[i] = nodeConfig

				generateEyesToml(i, dir)
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

			for i := range ips {
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

	return cmd
}

func cleanData(outputDir string) {
	if !utils.IsFileExisted(outputDir) {
		log.Info("Creating output dir...")
		// Create the folder
		os.MkdirAll(outputDir, 0755)
		return
	}

	files, err := os.ReadDir(outputDir)
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
		docker.NodeData[i].GrpcPortMap = 9090 + i
	}

	return docker
}

func getNodeSettings(chainID, keyringBackend string, index int, mysqlIp string, ips []string) config.Config {
	majority := (len(ips) + 1) * 2 / 3

	return config.Config{
		Mode: "dev",
		Sisu: config.SisuConfig{
			ChainId:        chainID,
			KeyringBackend: keyringBackend,
			ApiHost:        "0.0.0.0",
			ApiPort:        25456,
		},
		Tss: config.TssConfig{
			MajorityThreshold: majority,
			DheartHost:        fmt.Sprintf("dheart%d", index),
			DheartPort:        5678,
			DeyesUrl:          fmt.Sprintf("http://deyes%d:31001", index),
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
}

// Ip assignments:
// - 192.168.10.2, .3: ganache(s)
// - 192.168.10.4: mysql
// - 192.168.10.5 onward: for Sisu and others
func generateDockerCompose(outputPath string, ips []string, dockerConfig DockerNodeConfig) {
	const dockerComposeTemplate = `version: "3"
services:
  ganache1:
    image: ganache
    environment:
      - port=7545
      - networkId=189985
    ports:
      - 7545:7545
  ganache2:
    image: ganache
    environment:
      - port=7545
      - networkId=189986
    ports:
      - 8545:7545
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

func generateEyesToml(index int, dir string) {
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

		SisuServerUrl: fmt.Sprintf("http://sisu%d:25456", index),
	}

	writeDeyesConfig(deyesConfig, dir)
}

func writeDeyesConfig(deyesConfig DeyesConfiguration, dir string) {
	eyesToml := `db_host = "{{ .Sql.Host }}"
db_port = {{ .Sql.Port }}
db_username = "{{ .Sql.Username }}"
db_password = "{{ .Sql.Password }}"
db_schema = "{{ .Sql.Schema }}"

server_port = 31001
sisu_server_url = "{{ .SisuServerUrl }}"

[chains]{{ range $k, $ganache := .Ganaches }}
[chains.ganache{{ $ganache.Index }}]
  chain = "ganache{{ $ganache.Index }}"
  block_time = 1000
  starting_block = 0
  rpc_url = "http://{{ $ganache.Ip }}:7545"{{ end }}
`

	tmpl := template.New("eyesToml")

	configTemplate, err := tmpl.Parse(eyesToml)
	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer
	err = configTemplate.Execute(&buffer, deyesConfig)

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

	useOnMemory := "false"
	if len(peerIds) == 1 {
		useOnMemory = "true"
	}

	peerString := strings.Join(peers, ", ")

	sqlConfig := SqlConfig{
		Host:     "mysql",
		Port:     3306,
		Schema:   fmt.Sprintf("dheart%d", index),
		Username: "root",
		Password: "password",
	}

	writeHeartConfig(index, dir, peerString, useOnMemory, sqlConfig)

}

func writeHeartConfig(index int, dir string, peerString string, useOnMemory string, sqlConfig SqlConfig) {
	heartConfig := HeartConfiguration{
		PeerString:    peerString,
		SisuServerUrl: fmt.Sprintf("http://sisu%d:25456", index),
		Sql:           sqlConfig,
		UseOnMemory:   useOnMemory,
	}

	heartToml := `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

home-dir = "/root/"
use-on-memory = {{ .UseOnMemory }}
shortcut-preparams = true
sisu-server-url = "{{ .SisuServerUrl }}"
port = 5678

###############################################################################
###                        Database Configuration                           ###
###############################################################################
[db]
  host = "{{ .Sql.Host }}"
  port = {{ .Sql.Port }}
  username = "{{ .Sql.Username }}"
  password = "{{ .Sql.Password }}"
  schema = "{{ .Sql.Schema }}"
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
	err = configTemplate.Execute(&buffer, heartConfig)

	tmos.MustWriteFile(filepath.Join(dir, "dheart.toml"), buffer.Bytes(), 0644)
}
