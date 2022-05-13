package gen

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
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
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/spf13/cobra"

	econfig "github.com/sisu-network/deyes/config"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type TestnetGenerator struct {
}

type LogDNAConfig struct {
	Secret   string `json:"secret,omitempty"`
	AppName  string `json:"app_name,omitempty"`
	HostName string `json:"host_name,omitempty"`
}

type TestnetConfig struct {
	GenesisFolder   string                  `json:"genesis_folder"`
	Nodes           []TestnetNode           `json:"nodes"`
	DeyesChainsPath string                  `json:"deyes_chains_path"`
	LogDNAConfig    log.LogDNAConfig        `json:"log_dna_config"`
	EmailAlert      config.EmailAlertConfig `json:"email_alert"`
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
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			generator := &TestnetGenerator{}

			serverCtx := server.GetServerContextFromCmd(cmd)
			tmConfig := serverCtx.Config
			// tmConfig.LogLevel = "info"
			tmConfig.Consensus.TimeoutCommit = 5 * time.Second

			outputDir, _ := cmd.Flags().GetString(flagOutputDir)
			minGasPrices, _ := cmd.Flags().GetString(server.FlagMinGasPrices)
			nodeDirPrefix, _ := cmd.Flags().GetString(flagNodeDirPrefix)
			nodeDaemonHome, _ := cmd.Flags().GetString(flagNodeDaemonHome)
			chainId, _ := cmd.Flags().GetString(flagChainId)
			numValidators, _ := cmd.Flags().GetInt(flagNumValidators)
			algo, _ := cmd.Flags().GetString(flags.FlagKeyAlgorithm)
			configString, _ := cmd.Flags().GetString(flagConfigString)
			keyringBackend, _ := cmd.Flags().GetString(flags.FlagKeyringBackend)
			keyringPassphrase, _ := cmd.Flags().GetString(flagKeyringPassphrase)

			if keyringPassphrase == keyring.BackendFile && len(keyringPassphrase) == 0 {
				panic(fmt.Sprintf("Please input the passphrase if you're using keyring backend file by flag %s", keyringPassphrase))
			}

			testnetConfig := TestnetConfig{}
			err = json.Unmarshal([]byte(configString), &testnetConfig)

			if err != nil {
				panic(err)
			}

			if len(chainId) == 0 {
				chainId = "testnet"
			}

			log.Info("testnet chainId = ", chainId)

			err = os.MkdirAll(outputDir, os.ModePerm)
			if err != nil {
				panic(err)
			}

			monikers := make([]string, numValidators)
			for i := 0; i < numValidators; i++ {
				monikers[i] = "node" + strconv.Itoa(i)
			}

			nodes := testnetConfig.Nodes
			dnaConfig := testnetConfig.LogDNAConfig

			sisuIps := make([]string, numValidators)
			eyesIps := make([]string, numValidators)
			heartIps := make([]string, numValidators)

			for i := 0; i < numValidators; i++ {
				sisuIps[i] = nodes[i].SisuIp
				heartIps[i] = nodes[i].HeartIp
				eyesIps[i] = nodes[i].EyesIp
			}
			log.Info("ips = ", sisuIps)

			chains := getChains(filepath.Join(testnetConfig.GenesisFolder, "chains.json"))

			// Create configuration
			nodeConfigs := make([]config.Config, numValidators)
			for i := range sisuIps {
				dir := filepath.Join(outputDir, fmt.Sprintf("node%d", i))

				if err := os.MkdirAll(dir, nodeDirPerm); err != nil {
					panic(err)
				}

				nodeConfig := generator.getNodeSettings(i, chainId, keyringBackend, keyringPassphrase, nodes[i], chains, dnaConfig)
				nodeConfigs[i] = nodeConfig
			}

			settings := &Setting{
				clientCtx:         clientCtx,
				cmd:               cmd,
				tmConfig:          tmConfig,
				mbm:               mbm,
				genBalIterator:    genBalIterator,
				outputDir:         outputDir,
				chainID:           chainId,
				minGasPrices:      minGasPrices,
				nodeDirPrefix:     nodeDirPrefix,
				nodeDaemonHome:    nodeDaemonHome,
				keyringBackend:    keyringBackend,
				keyringPassphrase: keyringPassphrase,
				algoStr:           algo,
				numValidators:     numValidators,

				ips:         sisuIps,
				nodeConfigs: nodeConfigs,
				tokens:      getTokens(filepath.Join(testnetConfig.GenesisFolder, "tokens.json")),
				chains:      chains,
				liquidities: getLiquidity(filepath.Join(testnetConfig.GenesisFolder, "liquid.json")),
				params:      &types.Params{MajorityThreshold: int32(math.Ceil(float64(numValidators) * 2 / 3))},
				emailAlert:  testnetConfig.EmailAlert,
			}

			valPubKeys, err := InitNetwork(settings)
			peerIds, err := getPeerIds(len(sisuIps), valPubKeys)
			if err != nil {
				panic(err)
			}

			deyesChains := generator.readDeyesChainConfigs(testnetConfig.DeyesChainsPath)

			// Create config files for dheart and deyes.
			for i := range heartIps {
				dir := filepath.Join(outputDir, fmt.Sprintf("node%d", i))

				// Dheart configs
				generator.generateHeartToml(i, dir, heartIps, peerIds, sisuIps[i], nodes[i].Sql, valPubKeys, dnaConfig)

				// Deyes configs
				generator.generateEyesToml(i, dir, sisuIps[i], eyesIps[i], nodes[i].Sql, deyesChains, dnaConfig)
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
	cmd.Flags().String(flagConfigString, "", "configuration string for all nodes")
	cmd.Flags().String(flags.FlagKeyringBackend, keyring.BackendTest, "Keyring backend. file|os|kwallet|pass|test|memory")
	cmd.Flags().String(flagKeyringPassphrase, "", "Passphrase for keyring backend if using backend file. Leave it empty if use backend test")
	cmd.Flags().String(flagGenesisFolder, "./misc/dev", "Relative path to the folder that contains genesis configuration.")

	return cmd
}

func (g *TestnetGenerator) getNodeSettings(nodeIndex int, chainID, keyringBackend, keyringPassphrase string,
	testnetConfig TestnetNode, chainConfigs []*types.Chain, dnaConfig log.LogDNAConfig) config.Config {
	supportedChains := make(map[string]config.TssChainConfig)
	for _, chainConfig := range chainConfigs {
		supportedChains[chainConfig.Id] = config.TssChainConfig{
			Id: chainConfig.Id,
		}
	}

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
			DheartHost:      testnetConfig.HeartIp,
			DheartPort:      5678,
			DeyesUrl:        fmt.Sprintf("http://%s:31001", testnetConfig.EyesIp),
			SupportedChains: supportedChains,
		},
		LogDNA: dnaConfig,
	}
}

func (g *TestnetGenerator) readDeyesChainConfigs(path string) []econfig.Chain {
	deyesChains := make([]econfig.Chain, 0)
	file, _ := ioutil.ReadFile(path)
	err := json.Unmarshal([]byte(file), &deyesChains)
	if err != nil {
		panic(err)
	}

	return deyesChains
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

func (g *TestnetGenerator) generateEyesToml(index int, dir string, sisuIp, deyesIp string, sqlConfig SqlConfig,
	chainConfigs []econfig.Chain, dnaConfig log.LogDNAConfig) {
	sqlConfig.Schema = "deyes"

	dnaConfig.HostName = deyesIp
	dnaConfig.AppName = fmt.Sprintf("deyes%d", index)

	deyesConfig := DeyesConfiguration{
		Chains: chainConfigs,

		Sql:           sqlConfig,
		SisuServerUrl: fmt.Sprintf("http://%s:25456", sisuIp),
		LogDNA:        dnaConfig,
	}

	writeDeyesConfig(deyesConfig, dir)
}
