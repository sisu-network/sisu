package dev

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
			tokenSymbol, _ := cmd.Flags().GetString(flags.Erc20Symbol)
			tokenAddrString, _ := cmd.Flags().GetString(flags.Erc20Addrs)
			liquidityAddrString, _ := cmd.Flags().GetString(flags.LiquidityAddrs)

			c := &AddPoolTokenCommand{}
			c.addToken(urlString, mnemonic, tokenSymbol, tokenAddrString, liquidityAddrString)

			return nil
		},
	}

	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum", "Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.ChainUrls, "http://0.0.0.0:7545,http://0.0.0.0:8545", "RPCs of all the chains we want to fund.")
	cmd.Flags().String(flags.Erc20Symbol, "SISU", "Token symbol.")
	cmd.Flags().String(flags.LiquidityAddrs, fmt.Sprintf("%s,%s", ExpectedLiquidPoolAddress, ExpectedLiquidPoolAddress), "Token symbol.")
	cmd.Flags().String(flags.Erc20Addrs, fmt.Sprintf("%s,%s", ExpectedSisuAddress, ExpectedSisuAddress), "Token address.")

	return cmd
}

func (c *AddPoolTokenCommand) addToken(urlString, mnemonic, tokenSymbol, tokenAddrString, liquidityAddrString string) {
	liquidityAddrs := strings.Split(liquidityAddrString, ",")
	tokenAddrs := strings.Split(tokenAddrString, ",")

	clients := getEthClients(urlString)
	defer func() {
		for _, client := range clients {
			client.Close()
		}
	}()

	log.Infof("Adding token %s (address %s) for the pool %s", tokenSymbol, tokenAddrs[0], liquidityAddrString)

	wg := &sync.WaitGroup{}
	wg.Add(len(clients))
	for i, liquidityAddr := range liquidityAddrs {
		go func(i int, liquidityAddr string) {
			liquidInstance, err := liquidity.NewLiquiditypool(common.HexToAddress(liquidityAddr), clients[i])
			if err != nil {
				panic(err)
			}

			auth, err := getAuthTransactor(clients[i], mnemonic)
			if err != nil {
				panic(err)
			}

			tx, err := liquidInstance.AddToken(auth, []common.Address{common.HexToAddress(tokenAddrs[i])}, []string{tokenSymbol})
			if err != nil {
				panic(err)
			}

			log.Info("Tx hash = ", tx.Hash())

			bind.WaitMined(context.Background(), clients[i], tx)
			wg.Done()
		}(i, liquidityAddr)
	}
	wg.Wait()
}
