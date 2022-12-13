package dev

import (
	"fmt"
	"strings"
	"time"

	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/helper"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/spf13/cobra"
)

type DeployAndFundCmd struct {
}

func DeployAndFund() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deploy-and-fund",
		Short: `Deploy ERC20 tokens, vault contracts and fund Sisu's account. Example:
./sisu dev deploy-and-fund
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			chainString, _ := cmd.Flags().GetString(flags.Chains)
			chainUrls, _ := cmd.Flags().GetString(flags.ChainUrls)
			mnemonic, _ := cmd.Flags().GetString(flags.Mnemonic)
			sisuRpc, _ := cmd.Flags().GetString(flags.SisuRpc)
			genesisFolder, _ := cmd.Flags().GetString(flags.GenesisFolder)
			chains := strings.Split(chainString, ",")

			log.Info("chainUrls = ", chainUrls)

			log.Info("========= Deploy ERC20 =========")

			// Deploy Vault & ERC20
			deployContractCmd := &DeployContractCmd{}

			// Deploy vault
			log.Info("========= Deploying Vault =========")
			expectedVaults := helper.ReadVaults(genesisFolder, chains)
			vaultAddrs := deployContractCmd.doDeployment(chainUrls, "vault", mnemonic, expectedVaults, "", "")

			// Deploy Sisu and ADA tokens
			tokens := helper.ReadToken(genesisFolder)
			filteredTokens := filterTokenAddressForChains(tokens, chains)
			allTokenAddrs := make([][]string, 0)
			for _, token := range filteredTokens {
				log.Verbosef("Deploying token %s", token.Id)
				expectedAddrs := getExpectedAddressForTokens(token, chains)
				addrs := deployContractCmd.doDeployment(chainUrls, "erc20", mnemonic, expectedAddrs, "", token.Id)
				allTokenAddrs = append(allTokenAddrs, addrs)
				fmt.Println("allTokenAddrs = ", allTokenAddrs)
			}
			time.Sleep(time.Second * 3)

			log.Info("========= Adding Liquidity to the Vault =========")
			// Add support token to the pool
			for _, tokenAddrs := range allTokenAddrs {
				// tokenAddrs is an array of token address on different chains
				tokenAddrString := strings.Join(tokenAddrs, ",")

				// Add liquidity to the pool
				addLiquidityCmd := &AddLiquidityCmd{}
				addLiquidityCmd.approveAndAddLiquidity(chainUrls, mnemonic, tokenAddrString, vaultAddrs)

				// Wait for block to mine
				time.Sleep(time.Second * 3)
			}

			// Fund Sisu's account
			log.Info("========= Fund token to sisu's account =========")
			fundSisuCmd := &fundAccountCmd{}
			fundSisuCmd.fundSisuAccounts(cmd.Context(), chainString, chainUrls, mnemonic,
				sisuRpc, genesisFolder)

			return nil
		},
	}

	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum",
		"Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.ChainUrls, "http://0.0.0.0:7545,http://0.0.0.0:8545", "RPCs of all the chains we want to fund.")
	cmd.Flags().String(flags.Chains, "ganache1,ganache2", "Names of all chains we want to fund.")
	cmd.Flags().String(flags.SisuRpc, "0.0.0.0:9090", "URL to connect to Sisu. Please do NOT include http:// prefix")
	cmd.Flags().String(flags.GenesisFolder, "./misc/dev", "The genesis folder that contains config files to generate data.")

	return cmd
}

func filterTokenAddressForChains(tokens []*types.Token, chains []string) []*types.Token {
	chainMap := make(map[string]bool)
	for _, chain := range chains {
		chainMap[chain] = true
	}

	filteredTokens := make([]*types.Token, 0)
	for _, token := range tokens {
		if token.Addresses == nil {
			continue
		}

		for i, chain := range token.Chains {
			if chainMap[chain] && token.Addresses[i] != "" {
				filteredTokens = append(filteredTokens, token)
				break
			}
		}
	}

	return filteredTokens
}

func getExpectedAddressForTokens(token *types.Token, chains []string) []string {
	tokenAddrs := make([]string, 0)
	for _, targetChain := range chains {
		found := false
		for i, chain := range token.Chains {
			if chain == targetChain {
				tokenAddrs = append(tokenAddrs, token.Addresses[i])
				found = true
				break
			}
		}

		if !found {
			panic(fmt.Errorf("Cannot find address for token %s on chain %s in tokens file", token.Id, targetChain))
		}
	}

	return tokenAddrs
}
