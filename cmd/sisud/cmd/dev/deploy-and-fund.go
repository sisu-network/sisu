package dev

import (
	"fmt"
	"strings"
	"time"

	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/spf13/cobra"
)

type DeployAndFundCmd struct {
}

func DeployAndFund() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deploy-and-fund",
		Short: `Deploy ERC20 tokens, liquidity contracts and fund Sisu's account. Example:
./sisu dev deploy-and-fund
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			chainString, _ := cmd.Flags().GetString(flags.Chains)
			chainUrls, _ := cmd.Flags().GetString(flags.ChainUrls)
			mnemonic, _ := cmd.Flags().GetString(flags.Mnemonic)
			sisuRpc, _ := cmd.Flags().GetString(flags.SisuRpc)
			expectedLiquidityString, _ := cmd.Flags().GetString(flags.ExpectedLiquidityAddrs)
			cardanoSecret, _ := cmd.Flags().GetString(flags.CardanoSecret)
			cardanoFunderMnemonic, _ := cmd.Flags().GetString(flags.CardanoFunderMnemonic)
			genesisFolder, _ := cmd.Flags().GetString(flags.GenesisFolder)

			log.Info("chainUrls = ", chainUrls)

			log.Info("========= Deploy ERC20 and Liquidity Pool =========")

			// Deploy ERC20 And liquidity pool
			deployContractCmd := &DeployContractCmd{}
			// Deploy Sisu and ADA tokens
			sisuAddrs := deployContractCmd.doDeployment(chainUrls, "erc20", mnemonic, fmt.Sprintf("%s,%s", ExpectedSisuAddress, ExpectedSisuAddress), "Sisu Token", "SISU")
			adaAddrs := deployContractCmd.doDeployment(chainUrls, "erc20", mnemonic, fmt.Sprintf("%s,%s", ExpectedAdaAddress, ExpectedAdaAddress), "Ada Token", "ADA")
			liquidityAddrs := deployContractCmd.doDeployment(chainUrls, "liquidity", mnemonic, expectedLiquidityString, "", "")
			liquidityAddrString := strings.Join(liquidityAddrs, ",")
			time.Sleep(time.Second * 3)

			log.Info("========= Adding support token to the pool =========")
			allTokenAddrs := [][]string{sisuAddrs, adaAddrs}
			tokenSymbols := []string{"SISU", "ADA"}

			// Add support token to the pool
			for i, tokenAddrs := range allTokenAddrs {
				// tokenAddrs is an array of token address on different chains
				tokenAddrString := strings.Join(tokenAddrs, ",")
				addPoolTokenCmd := &AddPoolTokenCommand{}

				log.Infof("========= Adding %s liquidity to the pool =========", tokenSymbols[i])
				addPoolTokenCmd.addToken(chainUrls, mnemonic, tokenSymbols[i], tokenAddrString, liquidityAddrString)

				// Add liquidity to the pool
				addLiquidityCmd := &AddLiquidityCmd{}
				addLiquidityCmd.approveAndAddLiquidity(chainUrls, mnemonic, tokenAddrString, liquidityAddrString)

				// Wait for block to mine
				time.Sleep(time.Second * 3)
			}

			// Fund Sisu's account
			log.Info("========= Fund token to sisu's account and gateway =========")
			fundSisuCmd := &fundAccountCmd{}
			fundSisuCmd.fundSisuAccounts(cmd.Context(), chainString, chainUrls, mnemonic, genesisFolder,
				tokenSymbols, liquidityAddrString, sisuRpc, cardanoSecret, cardanoFunderMnemonic)

			return nil
		},
	}

	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum",
		"Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.ChainUrls, "http://0.0.0.0:7545,http://0.0.0.0:8545", "RPCs of all the chains we want to fund.")
	cmd.Flags().String(flags.Chains, "ganache1,ganache2", "Names of all chains we want to fund.")
	cmd.Flags().String(flags.SisuRpc, "0.0.0.0:9090", "URL to connect to Sisu. Please do NOT include http:// prefix")
	cmd.Flags().String(flags.ExpectedLiquidityAddrs, fmt.Sprintf("%s,%s", ExpectedLiquidPoolAddress, ExpectedLiquidPoolAddress),
		"Expected addressed of the liquidity contract after deployment. Empty string means do not check for address match.")
	cmd.Flags().String(flags.CardanoSecret, "", "The blockfrost secret to interact with cardano network.")
	cmd.Flags().String(flags.CardanoFunderMnemonic, "", "Mnemonic of funder wallet which already has a lot of test tokens")
	cmd.Flags().String(flags.GenesisFolder, "./misc/dev", "Relative path to the folder that contains genesis configuration.")

	return cmd
}
