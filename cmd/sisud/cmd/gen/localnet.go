package gen

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/types"
)

var (
	flagNodeDirPrefix     = "node-dir-prefix"
	flagNumValidators     = "v"
	flagOutputDir         = "output-dir"
	flagNodeDaemonHome    = "node-daemon-home"
	flagStartingIPAddress = "starting-ip-address"
	flagTmpDir            = "tmp-dir"
	flagChainId           = "chain-id"
	flagConfigString      = "config-string"
)

type localnetGenerator struct{}

// get cmd to initialize all files for tendermint localnet and application
func LocalnetCmd(mbm module.BasicManager, genBalIterator banktypes.GenesisBalancesIterator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "localnet",
		Short: "Initialize files for a simapp localnet",
		Long: `localnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).
Note, strict routability for addresses is turned off in the config file.
Example:
	For running single instance:
		./sisu localnet
	`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			serverCtx := server.GetServerContextFromCmd(cmd)
			tmConfig := serverCtx.Config
			tmConfig.LogLevel = ""
			tmConfig.Consensus.TimeoutCommit = time.Second * 3

			generator := &localnetGenerator{}

			outputDir, _ := cmd.Flags().GetString(flagOutputDir)
			minGasPrices, _ := cmd.Flags().GetString(server.FlagMinGasPrices)
			nodeDirPrefix, _ := cmd.Flags().GetString(flagNodeDirPrefix)
			nodeDaemonHome, _ := cmd.Flags().GetString(flagNodeDaemonHome)
			startingIPAddress, _ := cmd.Flags().GetString(flagStartingIPAddress)
			numValidators, _ := cmd.Flags().GetInt(flagNumValidators)
			algo, _ := cmd.Flags().GetString(flags.FlagKeyAlgorithm)

			// Get Chain id and keyring backend from .env file.
			chainID := "eth-sisu-local"
			keyringBackend := keyring.BackendFile

			nodeConfig := config.Config{
				Mode: "dev",
				Sisu: config.SisuConfig{
					ChainId:        chainID,
					KeyringBackend: keyringBackend,
					ApiHost:        "0.0.0.0",
					ApiPort:        25456,
				},
				Tss: config.TssConfig{
					DheartHost:        "0.0.0.0",
					DheartPort:        5678,
					DeyesUrl:          "http://0.0.0.0:31001",
					SupportedChains: map[string]config.TssChainConfig{
						"ganache1": {
							Id:    "ganache1",
							Token: "NATIVE_GANACHE1",
						},
						"ganache2": {
							Id:    "ganache2",
							Token: "NATIVE_GANACHE2",
						},
					},
				},
			}

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
				ips:            generator.getLocalIps(startingIPAddress, numValidators),
				keyringBackend: keyringBackend,
				algoStr:        algo,
				numValidators:  numValidators,
				nodeConfigs:    []config.Config{nodeConfig},
				tokens:         getTokens("./misc/tokens_dev.json"),
				chains:         getChains("./misc/chains.json"),
				liquidities:    getLiquidity("./misc/liquid_dev.json"),
				params:         &types.Params{MajorityThreshold: int32(math.Ceil(float64(numValidators) * 2 / 3))},
			}

			_, err = InitNetwork(settings)
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

func (g *localnetGenerator) getLocalIps(startingIPAddress string, count int) []string {
	ips := make([]string, count)
	for i := 0; i < count; i++ {
		ip, err := g.getIP(i, startingIPAddress)
		if err != nil {
			panic(err)
		}
		ips[i] = ip
	}

	return ips
}

func (g *localnetGenerator) getIP(i int, startingIPAddr string) (ip string, err error) {
	if len(startingIPAddr) == 0 {
		ip, err = server.ExternalIP()
		if err != nil {
			return "", err
		}
		return ip, nil
	}
	return g.calculateIP(startingIPAddr, i)
}

func (g *localnetGenerator) calculateIP(ip string, i int) (string, error) {
	ipv4 := net.ParseIP(ip).To4()
	if ipv4 == nil {
		return "", fmt.Errorf("%v: non ipv4 address", ip)
	}

	for j := 0; j < i; j++ {
		ipv4[3]++
	}

	return ipv4.String(), nil
}

func (g *localnetGenerator) getAuthTransactor(client *ethclient.Client, address common.Address) (*bind.TransactOpts, error) {
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	// This is the private key of the accounts0
	bz, err := hex.DecodeString("9f575b88940d452da46a6ceec06a108fcd5863885524aec7fb0bc4906eb63ab1")
	if err != nil {
		panic(err)
	}

	privateKey, err := ethcrypto.ToECDSA(bz)
	if err != nil {
		panic(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasPrice = gasPrice

	auth.GasLimit = uint64(10_000_000)

	return auth, nil
}
