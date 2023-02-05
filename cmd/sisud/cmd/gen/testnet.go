package gen

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	heartconfig "github.com/sisu-network/dheart/core/config"
	p2ptypes "github.com/sisu-network/dheart/p2p/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/config"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type TestnetGenerator struct {
}

type TestnetConfig struct {
	Nodes      []TestnetNode           `json:"nodes"`
	EmailAlert config.EmailAlertConfig `json:"email_alert"`
}

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
	EyesIp  string    `json:"deyes_ip"`
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
	  ./sisu testnet --v 2 --output-dir ./output --config-string '{"deyes_chains":[{"id":"ganache1","rpc":"http://ganache-0.ganache.ganache:7545","gas_price":5000000000,"block_time":3000},{"id":"ganache2","rpc":"http://ganache-1.ganache.ganache:7545","gas_price":10000000000,"block_time":3000}],"tokens":[{"id":"NATIVE_GANACHE1","price":2000000000},{"id":"NATIVE_GANACHE2","price":3000000000},{"id":"SISU","price":4000000000,"decimals":18,"chains":["ganache1","ganache2"],"addresses":["0x3A84fBbeFD21D6a5ce79D54d348344EE11EBd45C","0x3A84fBbeFD21D6a5ce79D54d348344EE11EBd45C"]}],"nodes":[{"sisu_ip":"sisud.sisu-0","dheart_ip":"dheart.sisu-0","deyes_ip":"deyes.sisu-0","sql":{"host":"mysql.mysql","port":3306,"username":"root","password":"password"}},{"sisu_ip":"sisud.sisu--1","dheart_ip":"dheart.sisu--1","deyes_ip":"deyes.sisu--1","sql":{"host":"mysql.mysql","port":3306,"username":"root","password":"password"}}]}'
	`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			generator := &TestnetGenerator{}

			serverCtx := server.GetServerContextFromCmd(cmd)
			tmConfig := serverCtx.Config
			// tmConfig.LogLevel = "info"
			tmConfig.Consensus.TimeoutCommit = 5 * time.Second

			outputDir, _ := cmd.Flags().GetString(flagOutputDir)
			nodeDirPrefix, _ := cmd.Flags().GetString(flagNodeDirPrefix)
			chainId, _ := cmd.Flags().GetString(flagChainId)
			numValidators, _ := cmd.Flags().GetInt(flagNumValidators)
			configString, _ := cmd.Flags().GetString(flagConfigString)

			testnetConfig := TestnetConfig{}
			err := json.Unmarshal([]byte(configString), &testnetConfig)
			if err != nil {
				panic(err)
			}

			if len(chainId) == 0 {
				chainId = "testnet"
			}

			err = os.MkdirAll(outputDir, os.ModePerm)
			if err != nil {
				panic(err)
			}

			monikers := make([]string, numValidators)
			for i := 0; i < numValidators; i++ {
				monikers[i] = "node" + strconv.Itoa(i)
			}

			nodes := testnetConfig.Nodes

			sisuIps := make([]string, numValidators)
			eyesIps := make([]string, numValidators)
			heartIps := make([]string, numValidators)

			for i := 0; i < numValidators; i++ {
				sisuIps[i] = nodes[i].SisuIp
				heartIps[i] = nodes[i].HeartIp
				eyesIps[i] = nodes[i].EyesIp
			}
			log.Info("ips = ", sisuIps)

			// Create configuration
			nodeConfigs := make([]config.Config, numValidators)
			keyringBackend := os.Getenv("KEYRING_BACKEND")
			keyringPassphrase := os.Getenv("KEYRING_PASSPHRASE")
			dnaConfig := generator.getLogDnaConfig()
			for i := range sisuIps {
				dir := filepath.Join(outputDir, fmt.Sprintf("node%d", i))
				if err := os.MkdirAll(dir, nodeDirPerm); err != nil {
					panic(err)
				}

				nodeConfig := generator.getNodeSettings(i, chainId, keyringBackend, keyringPassphrase,
					nodes[i], dnaConfig)
				nodeConfigs[i] = nodeConfig
			}

			settings := buildBaseSettings(cmd, mbm, genBalIterator)
			settings.tmConfig = tmConfig
			settings.outputDir = outputDir
			settings.chainID = chainId
			settings.nodeDirPrefix = nodeDirPrefix
			settings.keyringBackend = keyringBackend
			settings.keyringPassphrase = keyringPassphrase
			settings.ips = sisuIps
			settings.nodeConfigs = nodeConfigs
			settings.emailAlert = testnetConfig.EmailAlert

			valPubKeys, err := InitNetwork(settings)
			peerIds, err := getPeerIds(len(sisuIps), valPubKeys)
			if err != nil {
				panic(err)
			}

			// Create config files for dheart and deyes.
			for i := range heartIps {
				dir := filepath.Join(outputDir, fmt.Sprintf("node%d", i))

				// Dheart configs
				generator.generateHeartToml(i, dir, heartIps, peerIds, sisuIps[i], nodes[i].Sql, valPubKeys, dnaConfig)

				// Deyes configs
				generator.generateEyesToml(cmd, i, dir, sisuIps[i], eyesIps[i], nodes[i].Sql, dnaConfig)
			}

			return err
		},
	}

	cmd.Flags().Int(flagNumValidators, 4, "Number of validators to initialize the localnet with")
	cmd.Flags().StringP(flagOutputDir, "o", "./output", "Directory to store initialization data for the localnet")
	cmd.Flags().String(flagNodeDirPrefix, "node", "Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().String(flagNodeDaemonHome, "main", "Home directory of the node's daemon configuration")
	cmd.Flags().String(flagChainId, "sisu-talon-01", "Name of the chain")
	cmd.Flags().String(server.FlagMinGasPrices, fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom), "Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01photino,0.001stake)")
	cmd.Flags().String(flags.Algo, string(hd.Secp256k1Type), "Key signing algorithm to generate keys for")
	cmd.Flags().String(flagConfigString, "", "configuration string for all nodes")
	cmd.Flags().String(flags.KeyringBackend, keyring.BackendTest, "Keyring backend. file|os|kwallet|pass|test|memory")
	cmd.Flags().String(flags.GenesisFolder, "./misc/test", "Relative path to the folder that contains genesis configuration.")
	return cmd
}

func (g *TestnetGenerator) getNodeSettings(nodeIndex int, chainID, keyringBackend, keyringPassphrase string,
	testnetConfig TestnetNode, dnaConfig log.LogDNAConfig) config.Config {
	dnaConfig.HostName = testnetConfig.SisuIp
	dnaConfig.AppName = fmt.Sprintf("sisu%d", nodeIndex)

	return config.Config{
		Mode: "testnet",
		Sisu: config.SisuConfig{
			ChainId:           chainID,
			KeyringBackend:    keyringBackend,
			KeyringPassphrase: keyringPassphrase,
			ApiHost:           "0.0.0.0",
			ApiPort:           25456,
		},
		Tss: config.TssConfig{
			DheartHost: testnetConfig.HeartIp,
			DheartPort: 5678,
			DeyesUrl:   fmt.Sprintf("http://%s:31001", testnetConfig.EyesIp),
		},
		LogDNA: dnaConfig,
	}
}

func (g *TestnetGenerator) generateHeartToml(index int, outputDir string, heartIps []string,
	peerIds []string, sisuIp string, sqlConfig SqlConfig, valPubKeys []cryptotypes.PubKey, dnaConfig log.LogDNAConfig) {
	peers := make([]*p2ptypes.Peer, 0, len(peerIds)-1)
	for i := range peerIds {
		if i == index {
			continue
		}

		ipOrDns := "ip4"
		if net.ParseIP(heartIps[i]) == nil {
			ipOrDns = "dns"
		}

		address := fmt.Sprintf(`"/%s/%s/tcp/28300/p2p/%s"`, ipOrDns, heartIps[i], peerIds[i])
		peer := &p2ptypes.Peer{
			Address:    address,
			PubKey:     hex.EncodeToString(valPubKeys[i].Bytes()),
			PubKeyType: valPubKeys[i].Type(),
		}

		peers = append(peers, peer)
	}

	sisuUrl := fmt.Sprintf("http://%s:25456", sisuIp)
	dnaConfig.HostName = heartIps[index]
	dnaConfig.AppName = fmt.Sprintf("dheart%d", index)

	hConfig := heartconfig.HeartConfig{
		UseOnMemory:       false,
		ShortcutPreparams: false,
		SisuServerUrl:     sisuUrl,
		Db: heartconfig.DbConfig{
			Host:     sqlConfig.Host,
			Port:     sqlConfig.Port,
			Username: sqlConfig.Username,
			Password: sqlConfig.Password,
			Schema:   "dheart",
		},
		Connection: p2ptypes.ConnectionsConfig{
			Peers:      peers,
			Host:       "0.0.0.0",
			Port:       28300,
			Rendezvous: "rendezvous",
		},
		LogDNA: dnaConfig,
	}

	writeHeartConfig(outputDir, hConfig)
}

func (g *TestnetGenerator) getLogDnaConfig() log.LogDNAConfig {
	cfg := log.LogDNAConfig{}
	cfg.Secret = os.Getenv("LOG_DNA_SECRET")

	duration, err := time.ParseDuration("5s")
	if err != nil {
		panic(err)
	}
	cfg.FlushInterval.Duration = duration
	cfg.MaxBufferLen = 50
	cfg.LogLocal = true

	return cfg
}

func (g *TestnetGenerator) generateEyesToml(cmd *cobra.Command, index int, dir string, sisuIp,
	deyesIp string, sqlConfig SqlConfig, logDnaCfg log.LogDNAConfig) {
	sqlConfig.Schema = "deyes"

	genesisFolder, _ := cmd.Flags().GetString(flags.GenesisFolder)
	deyesCfg := getDeyesChains(cmd, genesisFolder)

	deyesCfg.DbHost = sqlConfig.Host
	deyesCfg.DbPort = sqlConfig.Port
	deyesCfg.DbUsername = sqlConfig.Username
	deyesCfg.DbPassword = sqlConfig.Password
	deyesCfg.DbSchema = sqlConfig.Schema
	deyesCfg.UseExternalRpcsInfo = true

	// deyesCfg.PriceOracleUrl = os.Getenv("ORACLE_URL")
	// deyesCfg.PriceOracleSecret = os.Getenv("ORACLE_SECRET")

	deyesCfg.LogDNA = logDnaCfg
	deyesCfg.LogDNA.HostName = deyesIp
	deyesCfg.LogDNA.AppName = fmt.Sprintf("deyes%d", index)

	deyesCfg.SisuServerUrl = fmt.Sprintf("http://%s:25456", sisuIp)

	writeDeyesConfig(deyesCfg, dir)
}
