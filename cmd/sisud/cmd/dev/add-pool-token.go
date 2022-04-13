package dev

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	liquidity "github.com/sisu-network/sisu/contracts/eth/liquiditypool"
	"github.com/spf13/cobra"
)

type AddPoolTokenCommand struct {
}

func AddPoolToken() *cobra.Command {
	cmd := &cobra.Command{
		Use: "add-pool-token",
		Long: `Add a SINGLE token to a list of chains provided by chain urls and liquidity addresses.
Usage:
./sisu dev add-pool-token  --erc20-symbols [List of TOKEN_SYMBOL] --chain-urls [List of CHAINS]

Example:
./sisu dev add-pool-token  --erc20-symbols SISU,SISU --chain-urls http://localhost:7545,http://localhost:8545

Short:
./sisu dev add-pool-token
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			urlString, _ := cmd.Flags().GetString(flags.ChainUrls)
			mnemonic, _ := cmd.Flags().GetString(flags.Mnemonic)
			tokenSymbolString, _ := cmd.Flags().GetString(flags.Erc20Symbols)
			tokenAddrString, _ := cmd.Flags().GetString(flags.Erc20Addrs)
			liquidityAddrString, _ := cmd.Flags().GetString(flags.LiquidityAddrs)

			c := &AddPoolTokenCommand{}
			c.addToken(urlString, mnemonic, tokenSymbolString, tokenAddrString, liquidityAddrString)

			return nil
		},
	}

	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum", "Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.ChainUrls, "http://0.0.0.0:7545,http://0.0.0.0:8545", "RPCs of all the chains we want to fund.")
	cmd.Flags().String(flags.Erc20Symbols, "SISU,SISU", "Token symbol.")
	cmd.Flags().String(flags.LiquidityAddrs, fmt.Sprintf("%s,%s", ExpectedLiquidPoolAddress, ExpectedLiquidPoolAddress), "Token symbol.")
	cmd.Flags().String(flags.Erc20Addrs, fmt.Sprintf("%s,%s", ExpectedErc20Address, ExpectedErc20Address), "Token address.")

	return cmd
}

func (c *AddPoolTokenCommand) addToken(urlString, mnemonic, tokenSymbol, tokenAddrString, liquidityAddrString string) {
	urls := strings.Split(urlString, ",")
	liquidityAddrs := strings.Split(liquidityAddrString, ",")
	tokenAddrs := strings.Split(tokenAddrString, ",")

	clients := make([]*ethclient.Client, 0)

	// Get all urls from command arguments.
	for i := 0; i < len(urls); i++ {
		client, err := ethclient.Dial(urls[i])
		if err != nil {
			log.Error("please check chain is up and running, url = ", urls[i])
			panic(err)
		}
		clients = append(clients, client)
	}
	defer func() {
		for _, client := range clients {
			client.Close()
		}
	}()

	log.Infof("Adding pool token %s (address %s) for the pool %s", tokenSymbol, tokenAddrs[0], liquidityAddrString)

	wg := &sync.WaitGroup{}
	wg.Add(len(clients))
	for i, liquidityAddr := range liquidityAddrs {
		go func(i int, liquidityAddr string) {
			liquidInstance, err := liquidity.NewLiquiditypool(common.HexToAddress(liquidityAddr), clients[i])
			if err != nil {
				panic(err)
			}

			auth, err := getAuthTransactor(clients[i], mnemonic, account0.Address)
			if err != nil {
				panic(err)
			}

			tx, err := liquidInstance.AddToken(auth, []common.Address{common.HexToAddress(tokenAddrs[i])}, []string{tokenSymbol})
			if err != nil {
				panic(err)
			}

			bind.WaitMined(context.Background(), clients[i], tx)
			wg.Done()
		}(i, liquidityAddr)
	}
	wg.Wait()
}
