package gen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/types"

	tmconfig "github.com/tendermint/tendermint/config"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	ttypes "github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

const nodeDirPerm = 0755

// Initialize the localnet
func InitNetwork(settings *Setting) ([]cryptotypes.PubKey, error) {
	clientCtx := settings.clientCtx
	cmd := settings.cmd
	tmConfig := settings.tmConfig
	mbm := settings.mbm
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

	memos := make([]string, numValidators)
	nodes := make([]*types.Node, numValidators)
	valPubKeys := make([]cryptotypes.PubKey, numValidators)

	// Temporary set os.Stdin to nil to read the keyring-passphrase from created buffer
	if settings.keyringBackend == keyring.BackendFile {
		oldStdin := os.Stdin
		defer func() {
			os.Stdin = oldStdin
		}()

		os.Stdin = nil
	}

	// generate private keys, node IDs, and initial transactions
	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outputDir, nodeDirName)
		mainAppDir := filepath.Join(nodeDir, nodeDaemonHome)

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

		inBuf := io.Reader(os.Stdin)
		if keyringBackend == keyring.BackendFile {
			buf := bytes.NewBufferString(settings.keyringPassphrase)
			buf.WriteByte('\n')
			buf.WriteString(settings.keyringPassphrase)
			buf.WriteByte('\n')

			inBuf = buf
		}

		kb, err := keyring.New(sdk.KeyringServiceName(), keyringBackend, mainAppDir, inBuf)
		if err != nil {
			return nil, err
		}

		node, secret, tendermintKey := getNode(kb, algoStr, nodeDirName, outputDir, tmConfig)
		nodes[i] = node
		valPubKeys[i] = tendermintKey

		memo := fmt.Sprintf("%s@%s:26656", node.Id, ip)
		genFiles = append(genFiles, tmConfig.GenesisFile())
		memos[i] = memo

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

		coins := sdk.Coins{
			sdk.NewCoin(common.SisuCoinName, accTokens),
		}

		genBalances = append(genBalances, banktypes.Balance{Address: node.AccAddress, Coins: coins.Sort()})
		acc, err := sdk.AccAddressFromBech32(node.AccAddress)
		if err != nil {
			panic(err)
		}
		genAccounts = append(genAccounts, authtypes.NewBaseAccount(acc, nil, 0, 0))

		// Write config/app.toml
		srvconfig.WriteConfigFile(filepath.Join(mainAppDir, "config/app.toml"), simappConfig)

		// Generate sisu.toml
		generateSisuToml(settings, i, nodeDir)
	}

	if err := initGenFiles(
		clientCtx, mbm, chainID, genAccounts, genBalances, genFiles, nodes,
		settings.tokens, settings.chains, settings.liquidities, settings.params,
	); err != nil {
		return nil, err
	}

	err := collectGenFiles(
		clientCtx, tmConfig, chainID, numValidators,
		outputDir, nodeDirPrefix, nodeDaemonHome, memos,
	)
	if err != nil {
		return nil, err
	}

	cmd.PrintErrf("Successfully initialized %d node directories\n", numValidators)
	return valPubKeys, nil
}

func getNode(kb keyring.Keyring, algoStr string, nodeDirName string, outputDir string, tmConfig *tmconfig.Config) (*types.Node, string, cryptotypes.PubKey) {
	keyringAlgos, _ := kb.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(algoStr, keyringAlgos)
	if err != nil {
		panic(err)
	}

	addr, secret, err := server.GenerateSaveCoinKey(kb, nodeDirName, true, algo)
	if err != nil {
		_ = os.RemoveAll(outputDir)
		panic(err)
	}

	nodeId, tendermintPubKey, err := InitializeNodeValidatorFilesFromMnemonic(tmConfig, secret)
	if err != nil {
		_ = os.RemoveAll(outputDir)
		panic(err)
	}

	return &types.Node{
		Id: nodeId,
		ConsensusKey: &types.Pubkey{
			Type:  tendermintPubKey.Type(),
			Bytes: tendermintPubKey.Bytes(),
		},
		AccAddress:  addr.String(),
		IsValidator: true,
	}, secret, tendermintPubKey
}

func initGenFiles(
	clientCtx client.Context, mbm module.BasicManager, chainID string,
	genAccounts []authtypes.GenesisAccount, genBalances []banktypes.Balance,
	genFiles []string, nodes []*types.Node, tokens []*types.Token, chains []*types.Chain,
	liquids []*types.Liquidity, params *types.Params,
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

	checkPoints := make([]*types.GatewayCheckPoint, 0)
	for _, chain := range chains {
		checkPoints = append(checkPoints, &types.GatewayCheckPoint{
			Chain:       chain.Id,
			BlockHeight: 1,
			Nonce:       0,
		})
	}

	sisuGenState := types.DefaultGenesis()
	sisuGenState.Nodes = nodes
	sisuGenState.Tokens = tokens
	sisuGenState.Chains = chains
	sisuGenState.Liquids = liquids
	sisuGenState.Params = params
	sisuGenState.Checkpoints = checkPoints

	appGenState[types.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(sisuGenState)

	/////////////

	appGenStateJSON, err := json.MarshalIndent(appGenState, "", "  ")
	if err != nil {
		return err
	}

	genDoc := ttypes.GenesisDoc{
		ChainID:    chainID,
		AppState:   appGenStateJSON,
		Validators: nil,
	}

	// generate empty genesis files for each validator and save
	for i := 0; i < len(nodes); i++ {
		if err := genDoc.SaveAs(genFiles[i]); err != nil {
			return err
		}
	}
	return nil
}

func collectGenFiles(
	clientCtx client.Context, nodeConfig *tmconfig.Config, chainID string,
	numValidators int, outputDir, nodeDirPrefix, nodeDaemonHome string, memos []string,
) error {

	var appState json.RawMessage
	genTime := tmtime.Now()

	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outputDir, nodeDirName, nodeDaemonHome)
		nodeConfig.Moniker = nodeDirName

		nodeConfig.SetRoot(nodeDir)

		genDoc, err := ttypes.GenesisDocFromFile(nodeConfig.GenesisFile())
		if err != nil {
			return err
		}

		persistenPeers := make([]string, 0, numValidators-1)
		for j, memo := range memos {
			if i == j {
				continue
			}

			persistenPeers = append(persistenPeers, memo)
		}

		nodeAppState, err := GenAppStateFromConfig(
			clientCtx.JSONMarshaler,
			nodeConfig,
			*genDoc,
			strings.Join(persistenPeers, ","),
		)
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
	if index == 0 {
		cfg.Sisu.EmailAlert = settings.emailAlert
	}

	if len(settings.cardanoSecret) > 0 {
		cfg.Cardano.BlockfrostSecret = settings.cardanoSecret
	}

	if settings.cardanoDbConfig != nil {
		cfg.Cardano.SyncDB = *settings.cardanoDbConfig
	}

	if len(settings.cardanoClientType) > 0 {
		cfg.Cardano.ClientType = settings.cardanoClientType
	}

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
