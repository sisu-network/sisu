package gen

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sisu-network/cosmos-sdk/client"
	"github.com/sisu-network/cosmos-sdk/client/tx"
	"github.com/sisu-network/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/sisu-network/cosmos-sdk/crypto/types"
	"github.com/sisu-network/cosmos-sdk/server"
	srvconfig "github.com/sisu-network/cosmos-sdk/server/config"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/cosmos-sdk/types/module"
	authtypes "github.com/sisu-network/cosmos-sdk/x/auth/types"
	banktypes "github.com/sisu-network/cosmos-sdk/x/bank/types"
	"github.com/sisu-network/cosmos-sdk/x/genutil"
	genutiltypes "github.com/sisu-network/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/sisu-network/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"

	tmconfig "github.com/sisu-network/tendermint/config"
	tmos "github.com/sisu-network/tendermint/libs/os"
	tmrand "github.com/sisu-network/tendermint/libs/rand"
	"github.com/sisu-network/tendermint/types"
	tmtime "github.com/sisu-network/tendermint/types/time"
)

const nodeDirPerm = 0755

type Setting struct {
	clientCtx      client.Context
	cmd            *cobra.Command
	nodeConfig     *tmconfig.Config
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
	enableTss      bool
}

// Initialize the localnet
func InitNetwork(settings *Setting) error {
	clientCtx := settings.clientCtx
	cmd := settings.cmd
	nodeConfig := settings.nodeConfig
	mbm := settings.mbm
	genBalIterator := settings.genBalIterator
	outputDir := settings.outputDir
	chainID := settings.chainID
	minGasPrices := settings.minGasPrices
	nodeDirPrefix := settings.nodeDirPrefix
	nodeDaemonHome := settings.nodeDaemonHome
	ips := settings.ips
	keyringBackend := settings.keyringBackend
	algoStr := settings.algoStr
	numValidators := settings.numValidators
	monikers := settings.monikers

	if chainID == "" {
		chainID = "chain-" + tmrand.NewRand().Str(6)
	}

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

		nodeConfig.SetRoot(mainAppDir)
		nodeConfig.RPC.ListenAddress = "tcp://0.0.0.0:26657"

		if err := os.MkdirAll(filepath.Join(mainAppDir, "config"), nodeDirPerm); err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		if monikers == nil || len(monikers) == 0 {
			nodeConfig.Moniker = nodeDirName
		} else {
			nodeConfig.Moniker = monikers[i]
		}

		ip := ips[i]

		var err error
		nodeIDs[i], valPubKeys[i], err = genutil.InitializeNodeValidatorFiles(nodeConfig)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		memo := fmt.Sprintf("%s@%s:26656", nodeIDs[i], ip)
		genFiles = append(genFiles, nodeConfig.GenesisFile())

		kb, err := keyring.New(sdk.KeyringServiceName(), keyringBackend, mainAppDir, inBuf)
		if err != nil {
			return err
		}

		keyringAlgos, _ := kb.SupportedAlgorithms()
		algo, err := keyring.NewSigningAlgoFromString(algoStr, keyringAlgos)
		if err != nil {
			return err
		}

		addr, secret, err := server.GenerateSaveCoinKey(kb, nodeDirName, true, algo)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		info := map[string]string{"secret": secret}

		cliPrint, err := json.Marshal(info)
		if err != nil {
			return err
		}

		// save private key seed words
		if err := writeFile(fmt.Sprintf("%v.json", "key_seed"), mainAppDir, cliPrint); err != nil {
			return err
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
			return err
		}

		txBuilder := clientCtx.TxConfig.NewTxBuilder()
		if err := txBuilder.SetMsgs(createValMsg); err != nil {
			return err
		}

		txBuilder.SetMemo(memo)

		txFactory := tx.Factory{}
		txFactory = txFactory.
			WithChainID(chainID).
			WithMemo(memo).
			WithKeybase(kb).
			WithTxConfig(clientCtx.TxConfig)

		if err := tx.Sign(txFactory, nodeDirName, txBuilder, true); err != nil {
			return err
		}

		txBz, err := clientCtx.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
		if err != nil {
			return err
		}

		if err := writeFile(fmt.Sprintf("%v.json", nodeDirName), gentxsDir, txBz); err != nil {
			return err
		}

		srvconfig.WriteConfigFile(filepath.Join(mainAppDir, "config/app.toml"), simappConfig)

		if settings.enableTss {
			generategTssToml(settings, nodeDir)
		}
	}

	if err := initGenFiles(clientCtx, mbm, chainID, genAccounts, genBalances, genFiles, numValidators); err != nil {
		return err
	}

	err := collectGenFiles(
		clientCtx, nodeConfig, chainID, nodeIDs, valPubKeys, numValidators,
		outputDir, nodeDirPrefix, nodeDaemonHome, genBalIterator,
	)
	if err != nil {
		return err
	}

	cmd.PrintErrf("Successfully initialized %d node directories\n", numValidators)
	return nil
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

// Generate TSS Toml file. This should only work in local dev mode. In production mode, we should
// create the tss.toml file and put in the tss folder of the app.
func generategTssToml(settings *Setting, nodeDir string) {
	content := `enable = true
dheart_host = "localhost"
dheart_port = 5678
[supported_chains]
[supported_chains.eth]
	symbol = "eth"
	id = 1
	deyes_url = "http://0.0.0.0:31001"
[supported_chains.sisu-eth]
	symbol = "sisu-eth"
	id = 36767
	deyes_url = "http://0.0.0.0:31001"
`

	if err := os.MkdirAll(filepath.Join(nodeDir, "tss"), nodeDirPerm); err != nil {
		_ = os.RemoveAll(settings.outputDir)
		panic(err)
	}

	dir := filepath.Join(nodeDir, "tss")
	writeFile("tss.toml", dir, []byte(content))
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
