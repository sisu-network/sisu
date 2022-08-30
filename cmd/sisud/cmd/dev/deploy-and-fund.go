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
			vaultString, _ := cmd.Flags().GetString(flags.ExpectedVaultAddrs)
			cardanoSecret, _ := cmd.Flags().GetString(flags.CardanoSecret)
			cardanoNetwork, _ := cmd.Flags().GetString(flags.CardanoNetwork)

			log.Info("chainUrls = ", chainUrls)

			log.Info("========= Deploy ERC20 =========")

			// Deploy ERC20 And liquidity pool
			deployContractCmd := &DeployContractCmd{}
			// Deploy Sisu and ADA tokens
			sisuAddrs := deployContractCmd.doDeployment(chainUrls, "erc20", mnemonic, fmt.Sprintf("%s,%s", ExpectedSisuAddress, ExpectedSisuAddress), "Sisu Token", "SISU")
			adaAddrs := deployContractCmd.doDeployment(chainUrls, "erc20", mnemonic, fmt.Sprintf("%s,%s", ExpectedAdaAddress, ExpectedAdaAddress), "Ada Token", "ADA")
			time.Sleep(time.Second * 3)

			log.Info("========= Adding Liquidity to the Vault =========")
			allTokenAddrs := [][]string{sisuAddrs, adaAddrs}
			tokenSymbols := []string{"SISU", "ADA"}
			// Add support token to the pool
			for _, tokenAddrs := range allTokenAddrs {
				// tokenAddrs is an array of token address on different chains
				tokenAddrString := strings.Join(tokenAddrs, ",")

				// Add liquidity to the pool
				addLiquidityCmd := &AddLiquidityCmd{}
				addLiquidityCmd.approveAndAddLiquidity(chainUrls, mnemonic, tokenAddrString, vaultString)

				// Wait for block to mine
				time.Sleep(time.Second * 3)
			}

			// Fund Sisu's account
			log.Info("========= Fund token to sisu's account =========")
			fundSisuCmd := &fundAccountCmd{}
			fundSisuCmd.fundSisuAccounts(cmd.Context(), chainString, chainUrls, mnemonic,
				tokenSymbols, vaultString, sisuRpc, cardanoNetwork, cardanoSecret)

			return nil
		},
	}

	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum",
		"Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.ChainUrls, "http://0.0.0.0:7545,http://0.0.0.0:8545", "RPCs of all the chains we want to fund.")
	cmd.Flags().String(flags.Chains, "ganache1,ganache2", "Names of all chains we want to fund.")
	cmd.Flags().String(flags.SisuRpc, "0.0.0.0:9090", "URL to connect to Sisu. Please do NOT include http:// prefix")
	cmd.Flags().String(flags.ExpectedVaultAddrs, fmt.Sprintf("%s,%s", ExpectedVaultAddress, ExpectedVaultAddress), "Expected addressed of the contract after deployment. Empty string means do not check for address match.")
	cmd.Flags().String(flags.CardanoSecret, "", "The blockfrost secret to interact with cardano network.")
	cmd.Flags().String(flags.CardanoNetwork, "cardano-testnet", "The Cardano network that we are interacting with.")

	return cmd
}
