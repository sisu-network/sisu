package dev

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/sisu-network/sisu/contracts/eth/erc20"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/spf13/cobra"
)

type queryCommand struct{}

func Query() *cobra.Command {
	cmd := &cobra.Command{
		Use: "query",
		Long: `Query ERC20 token balance.
Usage:
./sisu dev query --account 0x2d532C099CA476780c7703610D807948ae47856A

./sisu dev query --erc20-symbol SISU --chain ganache1 --chain-url http://127.0.0.1:7545 --account 0x2d532C099CA476780c7703610D807948ae47856A
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			chain, _ := cmd.Flags().GetString(flags.Chain)
			chainUrl, _ := cmd.Flags().GetString(flags.ChainUrl)
			tokenSymbol, _ := cmd.Flags().GetString(flags.Erc20Symbol)
			account, _ := cmd.Flags().GetString(flags.Account)
			sisuRpc, _ := cmd.Flags().GetString(flags.SisuRpc)

			log.Infof("Querying token at address %s on chain %s", tokenSymbol, chainUrl)

			client, err := ethclient.Dial(chainUrl)
			if err != nil {
				log.Error("cannot connect to chain, url = ", chainUrl)
				panic(err)
			}
			defer client.Close()

			if len(account) == 0 {
				panic(flags.Account + " cannot be empty")
			}

			token := queryToken(cmd.Context(), sisuRpc, tokenSymbol)
			if token == nil {
				panic("cannot find token " + tokenSymbol)
			}
			addr := ""
			for i := range token.Addresses {
				if token.Chains[i] == chain {
					addr = token.Addresses[i]
					break
				}
			}
			if addr == "" {
				panic("cannot find address on chain " + chain)
			}

			store, err := erc20.NewErc20(common.HexToAddress(addr), client)
			if err != nil {
				panic(err)
			}

			balance, err := store.BalanceOf(nil, common.HexToAddress(account))
			if err != nil {
				panic(err)
			}

			log.Info("Balance = ", balance)

			return nil
		},
	}

	cmd.Flags().String(flags.Chain, "ganache2", "Source chain where the token is transferred from")
	cmd.Flags().String(flags.ChainUrl, "http://127.0.0.1:8545", "Source chain url")
	cmd.Flags().String(flags.Erc20Symbol, ExpectedSisuAddress, "Id of the token to be queried")
	cmd.Flags().String(flags.Account, "", "account address that we want to query")
	cmd.Flags().String(flags.SisuRpc, "0.0.0.0:9090", "URL to connect to Sisu. Please do NOT include http:// prefix")

	return cmd
}

func (c *queryCommand) getTokens(genesisFolder string) []*types.Token {
	tokens := []*types.Token{}

	dat, err := os.ReadFile(filepath.Join(genesisFolder, "tokens.json"))
	if err != nil {
		panic(err)
	}

	json.Unmarshal(dat, &tokens)

	return tokens
}
