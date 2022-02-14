package dev

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/lib/log"
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
./sisu dev query --token SISU --src ganache1 --account 0x2d532C099CA476780c7703610D807948ae47856A
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			tokenId, _ := cmd.Flags().GetString(flagToken)
			src, _ := cmd.Flags().GetString(flagSrc)
			srcUrl, _ := cmd.Flags().GetString(flagSrcUrl)
			account, _ := cmd.Flags().GetString(flagAccount)

			log.Infof("Querying token %s on chain %s", tokenId, src)
			if len(srcUrl) == 0 {
				srcUrl = getDefaultChainUrl(src)
			}

			client, err := ethclient.Dial(srcUrl)
			if err != nil {
				log.Error("cannot connect to chain, url = ", srcUrl)
				panic(err)
			}
			defer client.Close()

			c := &queryCommand{}

			tokens := c.getTokens()
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

			tokenAddr := token.Addresses[src]
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

	cmd.Flags().String(flagSrc, "ganache1", "Source chain where the token is transferred from")
	cmd.Flags().String(flagSrcUrl, "", "Source chain url")
	cmd.Flags().String(flagToken, "SISU", "Id of the token to be queried")
	cmd.Flags().String(flagAccount, "account", "account address that we want to query")

	return cmd
}

func (c *queryCommand) getTokens() []*types.Token {
	tokens := []*types.Token{}

	dat, err := os.ReadFile("./tokens_dev.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal(dat, &tokens)

	return tokens
}
