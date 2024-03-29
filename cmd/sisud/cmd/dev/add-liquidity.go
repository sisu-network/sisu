package dev

import (
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/helper"
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
			chainString, _ := cmd.Flags().GetString(flags.Chains)
			mnemonic, _ := cmd.Flags().GetString(flags.Mnemonic)
			tokenAddrString, _ := cmd.Flags().GetString(flags.Erc20Addrs)
			genesisFolder, _ := cmd.Flags().GetString(flags.GenesisFolder)

			c := &AddLiquidityCmd{}
			chains := strings.Split(chainString, ",")
			vaults := helper.ReadVaults(genesisFolder, chains)
			c.approveAndAddLiquidity(mnemonic, genesisFolder, chains, tokenAddrString, vaults)

			return nil
		},
	}

	cmd.Flags().String(flags.Chains, "ganache1,ganache2", "Names of all chains we want to fund.")
	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum", "Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.GenesisFolder, "./misc/dev", "The genesis folder that contains config files to generate data.")
	cmd.Flags().String(flags.Erc20Addrs, fmt.Sprintf("%s,%s", ExpectedSisuAddress, ExpectedSisuAddress), "Token address.")

	return cmd
}

func (c *AddLiquidityCmd) approveAndAddLiquidity(mnemonic, genesisFolder string, chains []string,
	tokenAddrString string, vaultAddrs []string) {
	tokenAddrs := strings.Split(tokenAddrString, ",")
	clients := getEthClients(chains, genesisFolder)
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
				log.Infof("Adding liquidity of token %s to the pool at %s for chain %s", tokenAddrs[i],
					vaultAddrs[i], chains[i])
				transferErc20(client, mnemonic, tokenAddrs[i], vaultAddrs[i])

				balance, _ = queryErc20Balance(client, tokenAddrs[i], vaultAddrs[i])
			} else {
				log.Infof("Vault received %s tokens (%s)", balance.String(), tokenAddrs[i])
			}
		}(i, client)
	}
	wg.Wait()
}
