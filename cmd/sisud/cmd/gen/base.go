package gen

import (
	"encoding/json"
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
	cardanoChain      string
	cardanoDbConfig   *econfig.SyncDbConfig
	vaults            []*types.Vault

	nodeConfigs []config.Config
	tokens      []*types.Token // tokens in the genesis data
	chains      []*types.Chain // chains in the genesis data
	params      *types.Params

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
	cardanoSecret, _ := cmd.Flags().GetString(flags.CardanoSecret)
	cardanoChain, _ := cmd.Flags().GetString(flags.CardanoChain)

	supportedChainsArr := getSupportedChains(cmd, genesisFolder)

	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		panic(err)
	}
	deyesChains := getDeyesChains(cmd, genesisFolder)
	pendingTxOutHeights := getPendingTxTimeoutHeight(deyesChains)

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
			MajorityThreshold:       int32(math.Ceil(float64(numValidators) * 2 / 3)),
			SupportedChains:         supportedChainsArr,
			PendingTxTimeoutHeights: pendingTxOutHeights,
			CommissionRate:          10, // 0.1%
		},
		cardanoSecret: cardanoSecret,
		cardanoChain:  cardanoChain,
		tokens:        getTokens(filepath.Join(genesisFolder, "tokens.json")),
		chains:        helper.GetChains(filepath.Join(genesisFolder, "chains.json")),
		vaults:        getVaults(filepath.Join(genesisFolder, "vault.json")),
	}

	return setting
}

func getDeyesChains(cmd *cobra.Command, genesisFolder string) []econfig.Chain {
	cardanoSecret, _ := cmd.Flags().GetString(flags.CardanoSecret)
	cardanoDbConfig, _ := cmd.Flags().GetString(flags.CardanoDbConfig)
	deyesChains := readDeyesChainConfigs(filepath.Join(genesisFolder, "deyes_chains.json"))

	chains := helper.GetChains(filepath.Join(genesisFolder, "chains.json"))
	// Add Cardano config
	if len(cardanoSecret) > 0 || len(cardanoDbConfig) > 0 {
		chains = append(chains, &types.Chain{
			Id: "cardano-testnet",
		})

		var syncDbConfig econfig.SyncDbConfig
		var clientType econfig.ClientType

		if len(cardanoDbConfig) > 0 {
			err := json.Unmarshal([]byte(cardanoDbConfig), &syncDbConfig)
			if err != nil {
				panic(err)
			}
			clientType = econfig.ClientTypeSelfHost
		} else {
			clientType = econfig.ClientTypeBlockFrost
		}

		// Add cardano configuration
		deyesChains = append(deyesChains, econfig.Chain{
			Chain:      "cardano-testnet",
			BlockTime:  10000,
			AdjustTime: 1000,
			ClientType: clientType,
			RpcSecret:  cardanoSecret,
			SyncDB:     syncDbConfig,
		})
	}

	// Check if solana is enable

	return deyesChains
}

func getSupportedChains(cmd *cobra.Command, genesisFolder string) []string {
	cardanoSecret, _ := cmd.Flags().GetString(flags.CardanoSecret)
	cardanoDbConfig, _ := cmd.Flags().GetString(flags.CardanoDbConfig)

	chains := helper.GetChains(filepath.Join(genesisFolder, "chains.json"))
	supportedChainsArr := make([]string, 0)
	for _, chain := range chains {
		supportedChainsArr = append(supportedChainsArr, chain.Id)
	}
	sort.Strings(supportedChainsArr)

	// Add Cardano config
	if len(cardanoSecret) > 0 || len(cardanoDbConfig) > 0 {
		supportedChainsArr = append(supportedChainsArr, "cardano-testnet")
		chains = append(chains, &types.Chain{
			Id:          "cardano-testnet",
			NativeToken: "ADA",
		})
	}

	return supportedChainsArr
}

func readDeyesChainConfigs(path string) []econfig.Chain {
	deyesChains := make([]econfig.Chain, 0)
	file, _ := ioutil.ReadFile(path)
	err := json.Unmarshal([]byte(file), &deyesChains)
	if err != nil {
		panic(err)
	}

	return deyesChains
}

func getPendingTxTimeoutHeight(deyesChains []econfig.Chain) []int64 {
	heights := make([]int64, len(deyesChains))
	for i, c := range deyesChains {
		heights[i] = int64(c.BlockTime) * 12000 / 5000 // 5s is the estimated Sisu's blocktime
	}

	return heights
}
