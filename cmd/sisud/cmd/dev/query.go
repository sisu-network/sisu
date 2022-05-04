package dev

import (
	"encoding/json"
	"fmt"
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
			tokenId, _ := cmd.Flags().GetString(flags.Erc20Symbol)
			src, _ := cmd.Flags().GetString(flags.Chain)
			srcUrl, _ := cmd.Flags().GetString(flags.ChainUrl)
			account, _ := cmd.Flags().GetString(flags.Account)
			genesisFolder, _ := cmd.Flags().GetString(flags.GenesisFolder)

			log.Infof("Querying token %s on chain %s", tokenId, src)

			client, err := ethclient.Dial(srcUrl)
			if err != nil {
				log.Error("cannot connect to chain, url = ", srcUrl)
				panic(err)
			}
			defer client.Close()

			if len(account) == 0 {
				panic(flags.Account + " cannot be empty")
			}

			c := &queryCommand{}

			tokens := c.getTokens(genesisFolder)
			var token *types.Token

			for _, t := range tokens {
				if t.Id == tokenId {
					token = t
					break
				}
			}
			if token == nil {
				panic(fmt.Errorf("cannot find token %s", tokenId))
			}

			if token.Addresses == nil {
				panic(fmt.Errorf("this is not an ERC20 token"))
			}

			tokenAddr := token.GetAddressForChain(src)
			if len(tokenAddr) == 0 {
				panic(fmt.Errorf("cannot find address for tokne %s on chain %s", tokenId, src))
			}

			store, err := erc20.NewErc20(common.HexToAddress(tokenAddr), client)
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
	cmd.Flags().String(flags.Erc20Symbol, "SISU", "Id of the token to be queried")
	cmd.Flags().String(flags.Account, "", "account address that we want to query")
	cmd.Flags().String(flags.GenesisFolder, "./misc/dev", "Location of genesis folder. This is used to load the list of tokens.")

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
