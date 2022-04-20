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
			amount, _ := cmd.Flags().GetInt(flags.Amount)

			expectedErc20String := fmt.Sprintf("%s,%s", ExpectedErc20Address, ExpectedErc20Address)
			expectedLiquidityString := fmt.Sprintf("%s,%s", ExpectedLiquidPoolAddress, ExpectedLiquidPoolAddress)

			log.Info("========= Deploy ERC20 and Liquidity Pool =========")

			// Deploy ERC20 And liquidity pool
			deployContractCmd := &DeployContractCmd{}
			erc20Addrs := deployContractCmd.doDeployment(chainUrls, "erc20", mnemonic, expectedErc20String, "Sisu Token", "SISU")
			liquidityAddrs := deployContractCmd.doDeployment(chainUrls, "liquidity", mnemonic, expectedLiquidityString, "", "")

			time.Sleep(time.Second * 3)

			log.Info("========= Adding support token to the pool =========")

			// Add support token to the pool
			tokenAddrString := strings.Join(erc20Addrs, ",")
			liquidityAddrString := strings.Join(liquidityAddrs, ",")
			addPoolTokenCmd := &AddPoolTokenCommand{}
			addPoolTokenCmd.addToken(chainUrls, mnemonic, "SISU,SISU", tokenAddrString, liquidityAddrString)

			log.Info("========= Adding liquidity to the pool =========")

			// Add liquidity to the pool
			addLiquidityCmd := &AddLiquidityCmd{}
			addLiquidityCmd.approveAndAddLiquidity(chainUrls, mnemonic, tokenAddrString, liquidityAddrString, amount)

			log.Info("========= Fund sisu's account and gateway =========")

			// Fund Sisu's account
			fundSisuCmd := &fundAccountCmd{}
			fundSisuCmd.fundSisuAccounts(cmd.Context(), chainString, chainUrls, mnemonic, tokenAddrString, liquidityAddrString, sisuRpc, amount)

			return nil
		},
	}

	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum", "Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.ChainUrls, "http://0.0.0.0:7545,http://0.0.0.0:8545", "RPCs of all the chains we want to fund.")
	cmd.Flags().String(flags.Chains, "ganache1,ganache2", "Names of all chains we want to fund.")
	cmd.Flags().String(flags.SisuRpc, "0.0.0.0:9090", "URL to connect to Sisu. Please do NOT include http:// prefix")

	cmd.Flags().Int(flags.Amount, 100, "The amount that gateway addresses will receive")

	return cmd
}
