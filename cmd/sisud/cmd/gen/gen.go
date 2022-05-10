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
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/spf13/cobra"

	tmconfig "github.com/tendermint/tendermint/config"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	ttypes "github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

const nodeDirPerm = 0755

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
	numCandidates     int
	isLocalMultiNode  bool

	nodeConfigs []config.Config
	tokens      []*types.Token // tokens in the genesis data
	chains      []*types.Chain // chains in the genesis data
	liquidities []*types.Liquidity
	params      *types.Params
}

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
	numCandidates := settings.numCandidates
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

	totalNodes := numValidators + numCandidates

	memos := make([]string, totalNodes)
	p2ps := make([]string, totalNodes)
	rpcs := make([]string, totalNodes)
	proxies := make([]string, totalNodes)
	pprofs := make([]string, totalNodes)
	nodes := make([]*types.Node, totalNodes)
	valAndCanPubKeys := make([]cryptotypes.PubKey, totalNodes)
	//canPubKeys := make([]cryptotypes.PubKey, numCandidates)

	// Temporary set os.Stdin to nil to read the keyring-passphrase from created buffer
	if settings.keyringBackend == keyring.BackendFile {
		oldStdin := os.Stdin
		defer func() {
			os.Stdin = oldStdin
		}()

		os.Stdin = nil
	}

	// generate private keys, node IDs, and initial transactions
	for i := 0; i < totalNodes; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outputDir, nodeDirName)
		mainAppDir := filepath.Join(nodeDir, nodeDaemonHome)

		tmConfig.SetRoot(mainAppDir)
		rpcs[i] = fmt.Sprintf("tcp://0.0.0.0:%d", 36656+i)
		p2ps[i] = fmt.Sprintf("tcp://0.0.0.0:%d", 26656+i)
		proxies[i] = fmt.Sprintf("tcp://0.0.0.0:%d", 16656+i)
		pprofs[i] = fmt.Sprintf("localhost:%d", 6060+i)

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

		isValidator := true
		if i >= numValidators {
			isValidator = false
		}

		node, secret, tendermintKey := getNode(kb, algoStr, nodeDirName, outputDir, tmConfig, isValidator)
		nodes[i] = node
		valAndCanPubKeys[i] = tendermintKey

		var memo string
		if settings.isLocalMultiNode {
			memo = fmt.Sprintf("%s@%s:%d", node.Id, ip, 26656+i)
		} else {
			memo = fmt.Sprintf("%s@%s:26656", node.Id, ip)
		}
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
		simappConfig.API.Address = fmt.Sprintf("tcp://0.0.0.0:%d", 1317+i)
		simappConfig.GRPC.Address = fmt.Sprintf("0.0.0.0:%d", 9090+i)
		srvconfig.WriteConfigFile(filepath.Join(mainAppDir, "config/app.toml"), simappConfig)

		// Genreate sisu.toml
		generateSisuToml(settings, i, nodeDir)
	}

	if err := initGenFiles(
		clientCtx, mbm, chainID, genAccounts, genBalances, genFiles, nodes,
		settings.tokens, settings.chains, settings.liquidities, settings.params,
	); err != nil {
		return nil, err
	}

	err := collectGenFiles(
		clientCtx, tmConfig, chainID, numValidators, numCandidates,
		outputDir, nodeDirPrefix, nodeDaemonHome, memos, rpcs, p2ps, proxies, pprofs,
	)
	if err != nil {
		return nil, err
	}

	cmd.PrintErrf("Successfully initialized %d node directories\n", totalNodes)
	return valAndCanPubKeys, nil
}

func getNode(kb keyring.Keyring, algoStr string, nodeDirName string, outputDir string, tmConfig *tmconfig.Config, isValidator bool) (*types.Node, string, cryptotypes.PubKey) {
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

	nodeId, cosmosPubKey, err := InitializeNodeValidatorFilesFromMnemonic(tmConfig, secret)
	if err != nil {
		_ = os.RemoveAll(outputDir)
		panic(err)
	}

	status := types.NodeStatus_Candidate
	if isValidator {
		status = types.NodeStatus_Validator
	}

	return &types.Node{
		Id: nodeId,
		ConsensusKey: &types.Pubkey{
			Type:  cosmosPubKey.Type(),
			Bytes: cosmosPubKey.Bytes(),
		},
		AccAddress:  addr.String(),
		IsValidator: isValidator,
		Status:      status,
	}, secret, cosmosPubKey
}

func initGenFiles(
	clientCtx client.Context, mbm module.BasicManager, chainID string,
	genAccounts []authtypes.GenesisAccount, genBalances []banktypes.Balance,
	genFiles []string, nodes []*types.Node, tokens []*types.Token, chains []*types.Chain, liquids []*types.Liquidity, params *types.Params,
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

	sisuGenState := types.DefaultGenesis()
	sisuGenState.Nodes = nodes
	sisuGenState.Tokens = tokens
	sisuGenState.Chains = chains
	sisuGenState.Liquids = liquids
	sisuGenState.Params = params

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
	numValidators int, numCandidates int, outputDir, nodeDirPrefix, nodeDaemonHome string, memos []string, rpcs []string, p2ps []string, proxies []string, pprofs []string,
) error {

	var appState json.RawMessage
	genTime := tmtime.Now()

	persistentPeers := make([][]string, numValidators)
	for i := 0; i < numValidators; i++ {
		peers := make([]string, 0, numValidators-1)
		for j := 0; j < numValidators; j++ {
			if i == j {
				continue
			}

			peers = append(peers, memos[j])
		}

		persistentPeers[i] = peers
	}

	// set all validator nodes as persistent peers of candidates
	canPersistentPeers := make([][]string, numCandidates)
	for i := 0; i < numCandidates; i++ {
		canPersistentPeers[i] = memos[:numValidators]
	}

	persistentPeers = append(persistentPeers, canPersistentPeers...)

	log.Debug("persistent peers: ", persistentPeers)

	for i := 0; i < numValidators+numCandidates; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outputDir, nodeDirName, nodeDaemonHome)
		nodeConfig.Moniker = nodeDirName
		nodeConfig.P2P.ListenAddress = p2ps[i]
		nodeConfig.RPC.ListenAddress = rpcs[i]
		nodeConfig.ProxyApp = proxies[i]
		nodeConfig.RPC.PprofListenAddress = pprofs[i]

		nodeConfig.SetRoot(nodeDir)

		genDoc, err := ttypes.GenesisDocFromFile(nodeConfig.GenesisFile())
		if err != nil {
			return err
		}

		nodeAppState, err := GenAppStateFromConfig(
			clientCtx.JSONMarshaler,
			nodeConfig,
			*genDoc,
			strings.Join(persistentPeers[i], ","),
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
