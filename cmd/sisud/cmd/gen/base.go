package gen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"sort"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	econfig "github.com/sisu-network/deyes/config"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/helper"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/spf13/cobra"
	tmconfig "github.com/tendermint/tendermint/config"
)

type Setting struct {
	clientCtx         client.Context
	cmd               *cobra.Command
	tmConfig          *tmconfig.Config
	mbm               module.BasicManager
	genBalIterator    banktypes.GenesisBalancesIterator
	outputDir         string
	chainID           string
	minGasPrices      string
	nodeDirPrefix     string
	nodeDaemonHome    string
	ips               []string
	monikers          []string
	keyringBackend    string
	keyringPassphrase string
	algoStr           string
	numValidators     int
	cardanoSecret     string
	vaults            []*types.Vault

	nodeConfigs  []config.Config
	tokens       []*types.Token // tokens in the genesis data
	chains       []*types.Chain // chains in the genesis data
	params       *types.Params
	solanaConfig *config.SolanaConfig

	emailAlert config.EmailAlertConfig
}

func buildBaseSettings(cmd *cobra.Command, mbm module.BasicManager,
	genBalIterator banktypes.GenesisBalancesIterator) *Setting {
	outputDir, _ := cmd.Flags().GetString(flagOutputDir)
	minGasPrices, _ := cmd.Flags().GetString(server.FlagMinGasPrices)
	nodeDirPrefix, _ := cmd.Flags().GetString(flagNodeDirPrefix)
	nodeDaemonHome, _ := cmd.Flags().GetString(flagNodeDaemonHome)
	algo, _ := cmd.Flags().GetString(flags.Algo)
	numValidators, _ := cmd.Flags().GetInt(flagNumValidators)
	genesisFolder, _ := cmd.Flags().GetString(flags.GenesisFolder)

	supportedChainsArr := getSupportedChains(cmd, genesisFolder)

	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		panic(err)
	}

	setting := &Setting{
		clientCtx:      clientCtx,
		cmd:            cmd,
		mbm:            mbm,
		genBalIterator: genBalIterator,
		outputDir:      outputDir,
		minGasPrices:   minGasPrices,
		nodeDirPrefix:  nodeDirPrefix,
		nodeDaemonHome: nodeDaemonHome,
		algoStr:        algo,
		numValidators:  numValidators,
		params: &types.Params{
			MajorityThreshold: int32(math.Ceil(float64(numValidators) * 2 / 3)),
			SupportedChains:   supportedChainsArr,
			CommissionRate:    10,  // 0.1%
			ExpirationBlock:   100, // For testing, make it only 100 blocks
		},
		tokens: getTokens(filepath.Join(genesisFolder, "tokens.json")),
		chains: helper.GetChains(filepath.Join(genesisFolder, "chains.json")),
		vaults: getVaults(filepath.Join(genesisFolder, "vault.json")),
	}

	// Check if solana is enabled
	if helper.IsSolanaEnabled(genesisFolder) {
		// Read Solana config file.
		solanaConfig := &config.SolanaConfig{}

		file, _ := ioutil.ReadFile(filepath.Join(genesisFolder, "solana.json"))
		err := json.Unmarshal([]byte(file), solanaConfig)
		if err != nil {
			panic(err)
		}

		setting.solanaConfig = solanaConfig
	}

	return setting
}

func getDeyesChains(cmd *cobra.Command, genesisFolder string) econfig.Deyes {
	deyesCfg := econfig.Load(filepath.Join(genesisFolder, "deyes.toml"))
	chains := helper.GetChains(filepath.Join(genesisFolder, "chains.json"))

	// Verify that all the ETH hains in chains.json are present in the deyesCfg
	for _, chain := range chains {
		if _, ok := deyesCfg.Chains[chain.Id]; !ok {
			panic(fmt.Errorf("Chain %s is not present in the deyes toml config", chain))
		}
	}

	// Cardano config
	cardanoCfg := helper.ReadCardanoConfig(genesisFolder)
	if cardanoCfg.Enable {
		newCfg, ok := deyesCfg.Chains[cardanoCfg.Chain]
		if !ok {
			panic(fmt.Errorf("Cardano chain is not present in the deyes config"))
		}

		var syncDbConfig econfig.SyncDbConfig
		var clientType econfig.ClientType

		if len(cardanoCfg.UseSyncDb) > 0 {
			clientType = econfig.ClientTypeSelfHost
		} else {
			clientType = econfig.ClientTypeBlockFrost
		}

		newCfg.ClientType = clientType
		newCfg.RpcSecret = cardanoCfg.Secret
		newCfg.SyncDB = syncDbConfig

		deyesCfg.Chains[cardanoCfg.Chain] = newCfg
	} else {
		delete(deyesCfg.Chains, cardanoCfg.Chain)
	}

	// Solana config
	solanaCfg := helper.ReadSolanaConfig(genesisFolder)
	if solanaCfg.Enable {
		_, ok := deyesCfg.Chains[solanaCfg.Chain]
		if !ok {
			panic(fmt.Errorf("Solana chain is not present in the deyes config"))
		}
	} else {
		delete(deyesCfg.Chains, solanaCfg.Chain)
	}

	// Lisk config
	liskCfg := helper.ReadLiskConfig(genesisFolder)
	if liskCfg.Enable {
		if _, ok := deyesCfg.Chains[liskCfg.Chain]; !ok {
			panic(fmt.Errorf("Lisk chain is not present in the deyes config"))
		}
	} else {
		delete(deyesCfg.Chains, liskCfg.Chain)
	}

	return deyesCfg
}

func getSupportedChains(cmd *cobra.Command, genesisFolder string) []string {
	chains := helper.GetChains(filepath.Join(genesisFolder, "chains.json"))
	supportedChainsArr := make([]string, 0)
	for _, chain := range chains {
		supportedChainsArr = append(supportedChainsArr, chain.Id)
	}
	sort.Strings(supportedChainsArr)

	// Add Cardano config
	if helper.IsCardanoEnabled(genesisFolder) {
		cardanoConfig := helper.ReadCardanoConfig(genesisFolder)
		supportedChainsArr = append(supportedChainsArr, cardanoConfig.Chain)
		chains = append(chains, &types.Chain{
			Id:          cardanoConfig.Chain,
			NativeToken: "ADA",
		})
	}

	if helper.IsSolanaEnabled(genesisFolder) {
		solanaCfg := helper.ReadSolanaConfig(genesisFolder)
		supportedChainsArr = append(supportedChainsArr, solanaCfg.Chain)
		chains = append(chains, &types.Chain{
			Id:          solanaCfg.Chain,
			NativeToken: "SOL",
		})
	}

	if helper.IsLiskEnabled(genesisFolder) {
		liskCfg := helper.ReadLiskConfig(genesisFolder)
		supportedChainsArr = append(supportedChainsArr, liskCfg.Chain)
		chains = append(chains, &types.Chain{
			Id:          liskCfg.Chain,
			NativeToken: "LSK",
		})
	}

	return supportedChainsArr
}
