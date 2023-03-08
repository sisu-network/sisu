package dev

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/sisu-network/sisu/contracts/eth/example"
	"github.com/sisu-network/sisu/contracts/eth/vault"
	"github.com/sisu-network/sisu/utils"
	"github.com/spf13/cobra"
)

type deployExampleContractCommand struct {
}

func DeployExampleContract() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deploy-example-contract",
		Long: `Deploy an example contract on a chain and adds it into vaults of other chains.
Usage:
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			src, _ := cmd.Flags().GetString(flags.Src)
			sisuRpc, _ := cmd.Flags().GetString(flags.SisuRpc)
			mnemonic, _ := cmd.Flags().GetString(flags.Mnemonic)
			genesisFolder, _ := cmd.Flags().GetString(flags.GenesisFolder)
			chainString, _ := cmd.Flags().GetString(flags.Chains)
			chains := []string{src}
			if chainString != "" {
				chains = append(chains, strings.Split(chainString, ",")...)
			}

			clients := getEthClients(chains, genesisFolder)
			if len(clients) == 0 {
				panic("There is no healthy client")
			}

			c := &deployExampleContractCommand{}
			c.doDeployment(clients, mnemonic, chains, sisuRpc)

			return nil
		},
	}

	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum", "Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.Src, "ganache1", "Name of chain we want to deploy.")
	cmd.Flags().String(flags.Chains, "ganache2", "Names of all chains we want to add into vault, separated by comma.")
	cmd.Flags().String(flags.SisuRpc, "0.0.0.0:9090", "URL to connect to Sisu. Please do NOT include http:// prefix")
	cmd.Flags().String(flags.GenesisFolder, "./misc/dev", "The genesis folder that contains config files to generate data.")

	return cmd
}

func (c *deployExampleContractCommand) doDeployment(
	clients []*ethclient.Client, mnemonic string, chains []string, sisuRpc string,
) {
	engine, err := NewEngine(clients[0], mnemonic)
	if err != nil {
		panic(err)
	}

	contractAddr := engine.Deploy(func(opts *bind.TransactOpts) *types.Transaction {
		_, tx, _, err := example.DeployExample(opts, engine.Client)
		if err != nil {
			panic(err)
		}
		return tx
	})

	v := common.HexToAddress(getEthVaultAddress(context.Background(), chains[0], sisuRpc))
	contract, err := vault.NewVault(v, clients[0])
	if err != nil {
		panic(err)
	}

	log.Infof("Deposit 0.01ETH to vault for mnemonic address")
	engine.SetValue(new(big.Int).Div(utils.EthToWei, big.NewInt(100))) // 0.01ETH
	engine.Run(func(opts *bind.TransactOpts) *types.Transaction {
		tx, err := contract.DepositNative(opts)
		if err != nil {
			panic(err)
		}
		return tx
	})

	for i, client := range clients[1:] {
		v := common.HexToAddress(getEthVaultAddress(context.Background(), chains[i+1], sisuRpc))
		contract, err := vault.NewVault(v, client)
		if err != nil {
			panic(err)
		}

		engine, err = NewEngine(client, mnemonic)
		if err != nil {
			panic(err)
		}

		log.Infof("Creating App in vault, chain = %s ... ", chains[i+1])
		engine.Run(func(opts *bind.TransactOpts) *types.Transaction {
			tx, err := contract.CreateApp(opts, contractAddr, opts.From, []common.Address{})
			if err != nil {
				panic(err)
			}

			return tx
		})

		log.Infof("Seting any caller to true, chain = %s ... ", chains[i+1])
		engine.Run(func(opts *bind.TransactOpts) *types.Transaction {
			tx, err := contract.SetAppAnyCaller(opts, contractAddr, true)
			if err != nil {
				panic(err)
			}

			return tx
		})
	}
}
