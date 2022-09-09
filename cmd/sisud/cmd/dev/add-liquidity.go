package dev

import (
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/spf13/cobra"
)

type AddLiquidityCmd struct {
}

func AddLiquidity() *cobra.Command {
	cmd := &cobra.Command{
		Use: "add-liquidity",
		Long: `Add liquidity(erc20) to a (list of) liquidity pool(s).
Usage:
./sisu dev add-liquidity --chain-urls [List of CHAINS] --erc20-addrs [List of ERC20_ADDRS] --liquidity-addrs [List of Liquidity Addresses] --amount [AMOUNT]

Example:
./sisu dev add-liquidity --erc20-addrs 0x3A84fBbeFD21D6a5ce79D54d348344EE11EBd45C,0x3A84fBbeFD21D6a5ce79D54d348344EE11EBd45C --liquidity-addrs 0xf0D676183dD5ae6b370adDdbE770235F23546f9d,0xf0D676183dD5ae6b370adDdbE770235F23546f9d --chain-urls http://localhost:7545,http://localhost:8545

Short:
./sisu dev add-liquidity
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			urlString, _ := cmd.Flags().GetString(flags.ChainUrls)
			mnemonic, _ := cmd.Flags().GetString(flags.Mnemonic)
			tokenAddrString, _ := cmd.Flags().GetString(flags.Erc20Addrs)
			vaultAddrsString, _ := cmd.Flags().GetString(flags.VaultAddrs)

			c := &AddLiquidityCmd{}
			c.approveAndAddLiquidity(urlString, mnemonic, tokenAddrString, vaultAddrsString)

			return nil
		},
	}

	cmd.Flags().String(flags.ChainUrls, "http://0.0.0.0:7545,http://0.0.0.0:8545", "RPCs of all the chains we want to fund.")
	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum", "Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.Erc20Addrs, fmt.Sprintf("%s,%s", ExpectedSisuAddress, ExpectedSisuAddress), "Token address.")
	cmd.Flags().String(flags.VaultAddrs, fmt.Sprintf("%s,%s", ExpectedVaultAddress, ExpectedVaultAddress), "Liquidity addresses.")

	return cmd
}

func (c *AddLiquidityCmd) approveAndAddLiquidity(urlString, mnemonic, tokenAddrString, vaultAddrString string) {
	tokenAddrs := strings.Split(tokenAddrString, ",")
	vaultAddrs := strings.Split(vaultAddrString, ",")
	urls := strings.Split(urlString, ",")
	clients := getEthClients(urlString)
	defer func() {
		for _, client := range clients {
			client.Close()
		}
	}()

	wg := &sync.WaitGroup{}
	// Approve the contract with some preallocated token from account0
	wg.Add(len(clients))
	for i, client := range clients {
		go func(i int, client *ethclient.Client) {
			approveAddress(client, mnemonic, tokenAddrs[i], vaultAddrs[i])
			wg.Done()
		}(i, client)
	}
	wg.Wait()
	log.Info("Vault approval done!")

	// Add liquidity to the vault
	wg.Add(len(clients))
	for i, client := range clients {
		go func(i int, client *ethclient.Client) {
			defer wg.Done()

			balance, err := queryErc20Balance(client, tokenAddrs[i], vaultAddrs[i])
			if err != nil {
				panic(err)
			}

			if balance.Cmp(big.NewInt(0)) == 0 {
				log.Infof("Adding liquidity of token %s to the pool at %s for chain url %s", tokenAddrs[i], vaultAddrs[i], urls[i])
				transferErc20(client, mnemonic, tokenAddrs[i], vaultAddrs[i])

				balance, _ = queryErc20Balance(client, tokenAddrs[i], vaultAddrs[i])
			} else {
				log.Infof("Vault received %s tokens (%s)", balance.String(), tokenAddrs[i])
			}
		}(i, client)
	}
	wg.Wait()
}
