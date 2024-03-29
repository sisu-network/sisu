package gen

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/spf13/cobra"
	tmos "github.com/tendermint/tendermint/libs/os"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	econfig "github.com/sisu-network/deyes/config"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	heartconfig "github.com/sisu-network/dheart/core/config"
	p2ptypes "github.com/sisu-network/dheart/p2p/types"
)

type localDockerGenerator struct{}

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
			serverCtx := server.GetServerContextFromCmd(cmd)
			tmConfig := serverCtx.Config
			tmConfig.P2P.AddrBookStrict = false
			tmConfig.LogLevel = ""
			tmConfig.Consensus.TimeoutCommit = time.Second * 3

			outputDir, _ := cmd.Flags().GetString(flagOutputDir)
			nodeDirPrefix, _ := cmd.Flags().GetString(flagNodeDirPrefix)
			numValidators, _ := cmd.Flags().GetInt(flagNumValidators)
			genesisFolder, _ := cmd.Flags().GetString(flags.GenesisFolder)

			g := &localDockerGenerator{}

			// Clean up
			g.cleanData(outputDir)

			// Make dir folder for mysql docker
			err := os.MkdirAll(filepath.Join(outputDir, "db"), 0755)
			if err != nil {
				panic(err)
			}

			// Get Chain id and keyring backend from .env file.
			chainID := "sisu"
			keyringBackend := keyring.BackendTest
			deyesChains := getDeyesChains(cmd, genesisFolder)

			ips := make([]string, numValidators)
			for i := range ips {
				ips[i] = fmt.Sprintf("sisu%d", i)
			}

			mysqlIp := "192.168.10.4"
			chainIds := []*big.Int{libchain.GetChainIntFromId("ganache1"), libchain.GetChainIntFromId("ganache2")}
			dockerConfig := g.getDockerConfig([]string{"192.168.10.2", "192.168.10.3"}, chainIds, "192.168.10.4", ips)

			nodeConfigs := make([]config.Config, numValidators)
			for i := range ips {
				dir := filepath.Join(outputDir, fmt.Sprintf("node%d", i))

				if err := os.MkdirAll(dir, nodeDirPerm); err != nil {
					panic(err)
				}

				nodeConfig := g.getNodeSettings(chainID, keyringBackend, i, mysqlIp, ips)
				nodeConfigs[i] = nodeConfig

				g.generateEyesToml(deyesChains, i, dir)
			}
			g.generateDockerCompose(filepath.Join(outputDir, "docker-compose.yml"), ips, dockerConfig)

			settings := buildBaseSettings(cmd, mbm, genBalIterator)
			settings.tmConfig = tmConfig
			settings.chainID = chainID
			settings.nodeDirPrefix = nodeDirPrefix
			settings.keyringBackend = keyringBackend
			settings.ips = ips
			settings.nodeConfigs = nodeConfigs

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
				g.generateHeartToml(i, dir, dockerConfig, peerIds, valPubKeys)
			}

			return err
		},
	}

	cmd.Flags().Int(flagNumValidators, 1, "Number of validators to initialize the localnet with")
	cmd.Flags().StringP(flagOutputDir, "o", "./output", "Directory to store initialization data for the localnet")
	cmd.Flags().String(flagNodeDirPrefix, "node", "Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().String(flagNodeDaemonHome, "main", "Home directory of the node's daemon configuration")
	cmd.Flags().String(flags.GenesisFolder, "./misc/dev", "Relative path to the folder that contains genesis configuration.")
	cmd.Flags().String(server.FlagMinGasPrices, fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom),
		"Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01photino,0.001stake)")
	cmd.Flags().String(flags.Algo, string(hd.Secp256k1Type), "Key signing algorithm to generate keys for")

	return cmd
}

func (g *localDockerGenerator) cleanData(outputDir string) {
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
			log.Info("Deleting folder ", path)
			err := os.RemoveAll(path)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (g *localDockerGenerator) getDockerConfig(ganacheIps []string, chainIds []*big.Int, mysqlIp string, ips []string) DockerNodeConfig {
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

func (g *localDockerGenerator) getNodeSettings(chainID, keyringBackend string, index int,
	mysqlIp string, ips []string) config.Config {
	return config.Config{
		Mode: "dev",
		Sisu: config.SisuConfig{
			ChainId:        chainID,
			KeyringBackend: keyringBackend,
			ApiHost:        "0.0.0.0",
			ApiPort:        25456,
		},
		Tss: config.TssConfig{
			DheartHost: fmt.Sprintf("dheart%d", index),
			DheartPort: 5678,
			DeyesUrl:   fmt.Sprintf("http://deyes%d:31001", index),
		},
	}
}

func (g *localDockerGenerator) generateDockerCompose(outputPath string, ips []string, dockerConfig DockerNodeConfig) {
	const dockerComposeTemplate = `version: "3"
services:
  ganache1:
    image: ganache
    container_name: ganache1
    environment:
      - port=7545
      - networkId=189985
    expose:
      - 7545
    ports:
      - 7545:7545
  ganache2:
    image: ganache
    container_name: ganache2
    environment:
      - port=8545
      - networkId=189986
    expose:
      - 8545
    ports:
      - 8545:8545
  mysql:
    image: mysql:8.0.19
    container_name: mysql
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
    container_name: dheart{{ $k }}
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
    container_name: deyes{{ $k }}
    expose:
    - 31001
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
	if err != nil {
		panic(err)
	}

	tmos.MustWriteFile(outputPath, buffer.Bytes(), 0644)
}

func (g *localDockerGenerator) generateEyesToml(deyesCfg econfig.Deyes, index int, dir string) {
	deyesCfg = updateOracleSecret(deyesCfg)

	deyesCfg.DbSchema = fmt.Sprintf("deyes%d", index)
	deyesCfg.DbHost = "mysql"
	deyesCfg.SisuServerUrl = fmt.Sprintf("http://sisu%d:25456", index)

	writeDeyesConfig(deyesCfg, dir)
}

func (g *localDockerGenerator) generateHeartToml(index int, dir string, dockerConfig DockerNodeConfig,
	peerIds []string, valPubKeys []cryptotypes.PubKey) {
	peers := make([]*p2ptypes.Peer, 0, len(peerIds)-1)
	for i := range peerIds {
		if i == index {
			continue
		}

		peer := &p2ptypes.Peer{
			Address:    fmt.Sprintf(`"/dns/dheart%d/tcp/28300/p2p/%s"`, i, peerIds[i]),
			PubKey:     hex.EncodeToString(valPubKeys[i].Bytes()),
			PubKeyType: valPubKeys[i].Type(),
		}

		peers = append(peers, peer)
	}

	sisuUrl := fmt.Sprintf("http://sisu%d:25456", index)

	hConfig := heartconfig.HeartConfig{
		UseOnMemory:       len(peerIds) == 1,
		ShortcutPreparams: true,
		SisuServerUrl:     sisuUrl,
		Db: heartconfig.DbConfig{
			Host:     "mysql",
			Port:     3306,
			Username: "root",
			Password: "password",
			Schema:   fmt.Sprintf("dheart%d", index),
		},
		Connection: p2ptypes.ConnectionsConfig{
			Peers:      peers,
			Host:       "0.0.0.0",
			Port:       28300,
			Rendezvous: "rendezvous",
		},
	}

	writeHeartConfig(dir, hConfig)
}
