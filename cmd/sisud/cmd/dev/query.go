package dev

import (
	"context"
	"fmt"
	"math/big"

	libchain "github.com/sisu-network/lib/chain"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/sisu-network/sisu/contracts/eth/erc20"
	"github.com/spf13/cobra"
	"github.com/ybbus/jsonrpc/v3"
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

			if len(account) == 0 {
				panic(flags.Account + " cannot be empty")
			}

			c := new(queryCommand)

			if libchain.IsETHBasedChain(chain) {
				c.queryEth(sisuRpc, chain, chainUrl, tokenSymbol, account)
			} else if libchain.IsSolanaChain(chain) {
			}

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

// queryEth returns an account balance on an ETH chain.
func (c *queryCommand) queryEth(sisuRpc, chain, chainUrl, tokenSymbol, account string) {
	addr := c.getTokenAddress(sisuRpc, chain, tokenSymbol)

	client, err := ethclient.Dial(chainUrl)
	if err != nil {
		log.Error("cannot connect to chain, url = ", chainUrl)
		panic(err)
	}
	defer client.Close()

	store, err := erc20.NewErc20(common.HexToAddress(addr), client)
	if err != nil {
		panic(err)
	}

	balance, err := store.BalanceOf(nil, common.HexToAddress(account))
	if err != nil {
		panic(err)
	}

	log.Info("Balance = ", balance)
}

func (c *queryCommand) getTokenAddress(sisuRpc, chain, tokenSymbol string) string {
	token := queryToken(context.Background(), sisuRpc, tokenSymbol)
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

	return addr
}

// querySolana query token balance on a solana chain
func (c *queryCommand) querySolanaAccountBalance(url, tokenAddr string) (*big.Int, error) {
	rpcClient := jsonrpc.NewClient(url)

	type ResponseValue struct {
		Amount         string `json:"amount,omitempty"`
		Decimals       int64  `json:"decimals,omitempty"`
		UiAmountString string `json:"uiAmountString,omitempty"`
	}

	type QueryResponse struct {
		Value ResponseValue `json:"value,omitempty"`
	}

	response := new(QueryResponse)

	res, err := rpcClient.Call(context.Background(), "getTokenAccountBalance", tokenAddr)
	if err != nil {
		return nil, err
	}

	if res.Error != nil {
		return nil, res.Error
	}

	err = res.GetObject(response)
	if err != nil {
		return nil, err
	}

	ret, ok := new(big.Int).SetString(response.Value.Amount, 10)
	if !ok {
		return nil, fmt.Errorf("Invalid respnose amount: %s", response.Value.Amount)
	}

	return ret, nil
}
