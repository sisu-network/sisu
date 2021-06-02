package gen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sisu-network/sisu/utils"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

const (
	flagNodeIps = "node-ips"
	flagChainId = "chain-id"
)

// get cmd to initialize all files for tendermint localnet and application
func TestnetCmd(mbm module.BasicManager, genBalIterator banktypes.GenesisBalancesIterator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testnet",
		Short: "Initialize files for a simapp localnet",
		Long: `privatenet creates configuration for a network with N validators.
Example:
	For multiple nodes (running with docker):
	  sisu testnet --v 2 --output-dir ./output --chain-id testnet --node-ips 192.168.10.2,192.168.10.3
	`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			outputDir, _ := cmd.Flags().GetString(flagOutputDir)
			minGasPrices, _ := cmd.Flags().GetString(server.FlagMinGasPrices)
			nodeDirPrefix, _ := cmd.Flags().GetString(flagNodeDirPrefix)
			nodeDaemonHome, _ := cmd.Flags().GetString(flagNodeDaemonHome)
			nodeIps, _ := cmd.Flags().GetString(flagNodeIps)
			chainId, _ := cmd.Flags().GetString(flagChainId)
			numValidators, _ := cmd.Flags().GetInt(flagNumValidators)
			algo, _ := cmd.Flags().GetString(flags.FlagKeyAlgorithm)

			// Get Chain id and keyring backend from .env file.
			keyringBackend := keyring.BackendFile

			monikers := make([]string, numValidators)
			for i := 0; i < numValidators; i++ {
				monikers[i] = "node-talon-" + strconv.Itoa(i)
			}

			ips := strings.Split(strings.TrimSpace(nodeIps), ",")
			utils.LogDebug("ips = ", ips)

			settings := &Setting{
				clientCtx:      clientCtx,
				cmd:            cmd,
				nodeConfig:     config,
				mbm:            mbm,
				genBalIterator: genBalIterator,
				outputDir:      outputDir,
				chainID:        chainId,
				minGasPrices:   minGasPrices,
				nodeDirPrefix:  nodeDirPrefix,
				nodeDaemonHome: nodeDaemonHome,
				ips:            ips,
				keyringBackend: keyringBackend,
				algoStr:        algo,
				numValidators:  numValidators,
			}

			return InitNetwork(settings)
		},
	}

	cmd.Flags().Int(flagNumValidators, 1, "Number of validators to initialize the localnet with")
	cmd.Flags().StringP(flagOutputDir, "o", "./output", "Directory to store initialization data for the localnet")
	cmd.Flags().String(flagNodeDirPrefix, "node", "Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().String(flagNodeDaemonHome, "main", "Home directory of the node's daemon configuration")
	cmd.Flags().String(flagNodeIps, "192.168.10.2,192.168.10.3", "List of ip addresses of validators")
	cmd.Flags().String(flagChainId, "fang-01", "Name of the chain")
	cmd.Flags().String(flagStartingIPAddress, "127.0.0.1", "Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:46656, ID1@192.168.0.2:46656, ...)")
	cmd.Flags().String(server.FlagMinGasPrices, fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom), "Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01photino,0.001stake)")
	cmd.Flags().String(flags.FlagKeyAlgorithm, string(hd.Secp256k1Type), "Key signing algorithm to generate keys for")

	return cmd
}
