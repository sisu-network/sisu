package gen

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/sisu-network/sisu/config"
	"github.com/spf13/cobra"

	tmconfig "github.com/tendermint/tendermint/config"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	"github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

const nodeDirPerm = 0755

type Setting struct {
	clientCtx      client.Context
	cmd            *cobra.Command
	tmConfig       *tmconfig.Config
	mbm            module.BasicManager
	genBalIterator banktypes.GenesisBalancesIterator
	outputDir      string
	chainID        string
	minGasPrices   string
	nodeDirPrefix  string
	nodeDaemonHome string
	ips            []string
	monikers       []string
	keyringBackend string
	algoStr        string
	numValidators  int

	nodeConfigs []config.Config
}

// Initialize the localnet
func InitNetwork(settings *Setting) ([]cryptotypes.PubKey, error) {
	clientCtx := settings.clientCtx
	cmd := settings.cmd
	tmConfig := settings.tmConfig
	mbm := settings.mbm
	genBalIterator := settings.genBalIterator
	outputDir := settings.outputDir
	minGasPrices := settings.minGasPrices
	nodeDirPrefix := settings.nodeDirPrefix
	nodeDaemonHome := settings.nodeDaemonHome
	ips := settings.ips
	keyringBackend := settings.keyringBackend
	algoStr := settings.algoStr
	numValidators := settings.numValidators
	monikers := settings.monikers

	if settings.chainID == "" {
		settings.chainID = "chain-" + tmrand.NewRand().Str(6)
	}
	chainID := settings.chainID

	nodeIDs := make([]string, numValidators)
	valPubKeys := make([]cryptotypes.PubKey, numValidators)

	simappConfig := srvconfig.DefaultConfig()
	simappConfig.MinGasPrices = minGasPrices
	simappConfig.API.Enable = true
	simappConfig.Telemetry.Enabled = true
	simappConfig.Telemetry.PrometheusRetentionTime = 60
	simappConfig.Telemetry.EnableHostnameLabel = false
	simappConfig.Telemetry.GlobalLabels = [][]string{{"chain_id", chainID}}

	var (
		genAccounts []authtypes.GenesisAccount
		genBalances []banktypes.Balance
		genFiles    []string
	)

	inBuf := bufio.NewReader(cmd.InOrStdin())
	// generate private keys, node IDs, and initial transactions
	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outputDir, nodeDirName)
		mainAppDir := filepath.Join(nodeDir, nodeDaemonHome)
		gentxsDir := filepath.Join(outputDir, "gentxs")

		tmConfig.SetRoot(mainAppDir)
		tmConfig.RPC.ListenAddress = "tcp://0.0.0.0:26657"

		if err := os.MkdirAll(filepath.Join(mainAppDir, "config"), nodeDirPerm); err != nil {
			_ = os.RemoveAll(outputDir)
			return nil, err
		}

		if monikers == nil || len(monikers) == 0 {
			tmConfig.Moniker = nodeDirName
		} else {
			tmConfig.Moniker = monikers[i]
		}

		ip := ips[i]

		var err error
		nodeIDs[i], valPubKeys[i], err = genutil.InitializeNodeValidatorFiles(tmConfig)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return nil, err
		}

		memo := fmt.Sprintf("%s@%s:26656", nodeIDs[i], ip)
		genFiles = append(genFiles, tmConfig.GenesisFile())

		kb, err := keyring.New(sdk.KeyringServiceName(), keyringBackend, mainAppDir, inBuf)
		if err != nil {
			return nil, err
		}

		keyringAlgos, _ := kb.SupportedAlgorithms()
		algo, err := keyring.NewSigningAlgoFromString(algoStr, keyringAlgos)
		if err != nil {
			return nil, err
		}

		addr, secret, err := server.GenerateSaveCoinKey(kb, nodeDirName, true, algo)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return nil, err
		}

		info := map[string]string{"secret": secret}

		cliPrint, err := json.Marshal(info)
		if err != nil {
			return nil, err
		}

		// save private key seed words
		if err := writeFile(fmt.Sprintf("%v.json", "key_seed"), mainAppDir, cliPrint); err != nil {
			return nil, err
		}

		accTokens := sdk.TokensFromConsensusPower(1000)
		accStakingTokens := sdk.TokensFromConsensusPower(500)
		coins := sdk.Coins{
			sdk.NewCoin(fmt.Sprintf("%stoken", nodeDirName), accTokens),
			sdk.NewCoin(sdk.DefaultBondDenom, accStakingTokens),
		}

		genBalances = append(genBalances, banktypes.Balance{Address: addr.String(), Coins: coins.Sort()})
		genAccounts = append(genAccounts, authtypes.NewBaseAccount(addr, nil, 0, 0))

		valTokens := sdk.TokensFromConsensusPower(100)
		createValMsg, err := stakingtypes.NewMsgCreateValidator(
			sdk.ValAddress(addr),
			valPubKeys[i],
			sdk.NewCoin(sdk.DefaultBondDenom, valTokens),
			stakingtypes.NewDescription(nodeDirName, "", "", "", ""),
			stakingtypes.NewCommissionRates(sdk.OneDec(), sdk.OneDec(), sdk.OneDec()),
			sdk.OneInt(),
		)
		if err != nil {
			return nil, err
		}

		txBuilder := clientCtx.TxConfig.NewTxBuilder()
		if err := txBuilder.SetMsgs(createValMsg); err != nil {
			return nil, err
		}

		txBuilder.SetMemo(memo)

		txFactory := tx.Factory{}
		txFactory = txFactory.
			WithChainID(chainID).
			WithMemo(memo).
			WithKeybase(kb).
			WithTxConfig(clientCtx.TxConfig)

		if err := tx.Sign(txFactory, nodeDirName, txBuilder, true); err != nil {
			return nil, err
		}

		txBz, err := clientCtx.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
		if err != nil {
			return nil, err
		}

		if err := writeFile(fmt.Sprintf("%v.json", nodeDirName), gentxsDir, txBz); err != nil {
			return nil, err
		}

		srvconfig.WriteConfigFile(filepath.Join(mainAppDir, "config/app.toml"), simappConfig)

		generateSisuToml(settings, i, nodeDir)
	}

	if err := initGenFiles(clientCtx, mbm, chainID, genAccounts, genBalances, genFiles, numValidators); err != nil {
		return nil, err
	}

	err := collectGenFiles(
		clientCtx, tmConfig, chainID, nodeIDs, valPubKeys, numValidators,
		outputDir, nodeDirPrefix, nodeDaemonHome, genBalIterator,
	)
	if err != nil {
		return nil, err
	}

	cmd.PrintErrf("Successfully initialized %d node directories\n", numValidators)
	return valPubKeys, nil
}

func initGenFiles(
	clientCtx client.Context, mbm module.BasicManager, chainID string,
	genAccounts []authtypes.GenesisAccount, genBalances []banktypes.Balance,
	genFiles []string, numValidators int,
) error {

	appGenState := mbm.DefaultGenesis(clientCtx.JSONMarshaler)

	// set the accounts in the genesis state
	var authGenState authtypes.GenesisState
	clientCtx.JSONMarshaler.MustUnmarshalJSON(appGenState[authtypes.ModuleName], &authGenState)

	accounts, err := authtypes.PackAccounts(genAccounts)
	if err != nil {
		return err
	}

	authGenState.Accounts = accounts
	appGenState[authtypes.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(&authGenState)

	// set the balances in the genesis state
	var bankGenState banktypes.GenesisState
	clientCtx.JSONMarshaler.MustUnmarshalJSON(appGenState[banktypes.ModuleName], &bankGenState)

	bankGenState.Balances = genBalances
	appGenState[banktypes.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(&bankGenState)

	appGenStateJSON, err := json.MarshalIndent(appGenState, "", "  ")
	if err != nil {
		return err
	}

	genDoc := types.GenesisDoc{
		ChainID:    chainID,
		AppState:   appGenStateJSON,
		Validators: nil,
	}

	// generate empty genesis files for each validator and save
	for i := 0; i < numValidators; i++ {
		if err := genDoc.SaveAs(genFiles[i]); err != nil {
			return err
		}
	}
	return nil
}

func collectGenFiles(
	clientCtx client.Context, nodeConfig *tmconfig.Config, chainID string,
	nodeIDs []string, valPubKeys []cryptotypes.PubKey, numValidators int,
	outputDir, nodeDirPrefix, nodeDaemonHome string, genBalIterator banktypes.GenesisBalancesIterator,
) error {

	var appState json.RawMessage
	genTime := tmtime.Now()

	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outputDir, nodeDirName, nodeDaemonHome)
		gentxsDir := filepath.Join(outputDir, "gentxs")
		nodeConfig.Moniker = nodeDirName

		nodeConfig.SetRoot(nodeDir)

		nodeID, valPubKey := nodeIDs[i], valPubKeys[i]
		initCfg := genutiltypes.NewInitConfig(chainID, gentxsDir, nodeID, valPubKey)

		genDoc, err := types.GenesisDocFromFile(nodeConfig.GenesisFile())
		if err != nil {
			return err
		}

		nodeAppState, err := genutil.GenAppStateFromConfig(clientCtx.JSONMarshaler,
			clientCtx.TxConfig, nodeConfig, initCfg, *genDoc, genBalIterator)
		if err != nil {
			return err
		}

		if appState == nil {
			// set the canonical application state (they should not differ)
			appState = nodeAppState
		}

		genFile := nodeConfig.GenesisFile()

		// overwrite each validator's genesis file to have a canonical genesis time
		if err := genutil.ExportGenesisFileWithTime(genFile, chainID, nil, appState, genTime); err != nil {
			return err
		}
	}

	return nil
}

func generateSisuToml(settings *Setting, index int, nodeDir string) {
	appDir := filepath.Join(nodeDir, "main")
	configDir := filepath.Join(appDir, "config")

	if err := os.MkdirAll(configDir, nodeDirPerm); err != nil {
		_ = os.RemoveAll(settings.outputDir)
		panic(err)
	}

	// Create tss folder
	if err := os.MkdirAll(filepath.Join(nodeDir, "tss"), nodeDirPerm); err != nil {
		panic(err)
	}

	cfg := settings.nodeConfigs[index]

	config.WriteConfigFile(filepath.Join(configDir, "sisu.toml"), &cfg)
}

func writeFile(name string, dir string, contents []byte) error {
	writePath := filepath.Join(dir)
	file := filepath.Join(writePath, name)

	err := tmos.EnsureDir(writePath, 0755)
	if err != nil {
		return err
	}

	err = tmos.WriteFile(file, contents, 0644)
	if err != nil {
		return err
	}

	return nil
}
