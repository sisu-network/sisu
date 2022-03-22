package cmd

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/helper"
	"github.com/sisu-network/sisu/contracts/eth/liquidity"
	"github.com/spf13/cobra"
)

type EmergencyWithdrawFundCmd struct {
	privateKey    *ecdsa.PrivateKey
	senderAddress common.Address
}

func ContractEmergencyWithdrawFund() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "withdraw",
		Long: "Emergency withdraw fund. Use it when found trouble in security risk",
		RunE: func(cmd *cobra.Command, args []string) error {
			rpcEndpoint, _ := cmd.Flags().GetString(flags.ChainUrl)
			contractAddr, _ := cmd.Flags().GetString(flags.ContractAddress)
			tokenString, _ := cmd.Flags().GetString(flags.Tokens)
			newOwner, _ := cmd.Flags().GetString(flags.NewOwner)
			mnemonic, _ := cmd.Flags().GetString(flags.Mnemonic)

			if len(contractAddr) == 0 {
				panic("invalid contractAddr hash. Example: --contractAddr-address 0xdac17f958d2ee523a2206206994597c13d831ec7")
			}

			if len(tokenString) == 0 {
				panic("invalid tokens. Example: --token 0xdac17f958d2ee523a2206206994597c13d831ec7,0xdac17f958d2ee523a2206206994597c13d831ec7")
			}

			if len(newOwner) == 0 {
				panic("invalid newOwner. Example: --new-owner 0xdac17f958d2ee523a2206206994597c13d831ec7")
			}

			if len(mnemonic) == 0 {
				panic("invalid mnemonic")
			}

			tokens := strings.Split(tokenString, ",")
			tokenHashes := make([]common.Address, 0, len(tokens))
			for _, token := range tokens {
				tokenHashes = append(tokenHashes, common.HexToAddress(token))
			}

			client, err := ethclient.Dial(rpcEndpoint)
			if err != nil {
				log.Error("please check chain is up and running, url = ", rpcEndpoint)
				panic(err)
			}
			defer func() {
				client.Close()
			}()

			c := &EmergencyWithdrawFundCmd{}
			c.privateKey, c.senderAddress = helper.GetPrivateKeyFromMnemonic(mnemonic)

			for _, token := range tokens {
				if !c.isContractDeployed(client, common.HexToAddress(token)) {
					panic(fmt.Sprintf("token address %s is not found", token))
				}
			}

			if !c.isContractDeployed(client, common.HexToAddress(contractAddr)) {
				panic(fmt.Sprintf("liquidity contractAddr address %s at endpoint %s is not deployed", contractAddr, rpcEndpoint))
			}

			return c.withdrawFund(client, common.HexToAddress(contractAddr), tokenHashes, common.HexToAddress(newOwner))
		},
	}

	cmd.Flags().String(flags.ContractAddress, "liquidity", "Liquidity contract address which we want to interact.")
	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum", "Mnemonic used to interact with contract.")
	cmd.Flags().String(flags.ChainUrl, "http://0.0.0.0:7545", "RPC endpoint of chain")
	cmd.Flags().String(flags.Tokens, "0xf0D676183dD5ae6b370adDdbE770235F23546f9d,0xf0D676183dD5ae6b370adDdbE770235F23546f9d", "Token addresses to withdraw")
	cmd.Flags().String(flags.NewOwner, "0x215375950B138B9f5aDfaEb4dc172E8AD1dDe7f5", "New fund's owner")
	return cmd
}

func (c *EmergencyWithdrawFundCmd) withdrawFund(client *ethclient.Client, liquidityAddr common.Address,
	tokens []common.Address, newOwner common.Address) error {
	liquidInstance, err := liquidity.NewLiquidity(liquidityAddr, client)
	if err != nil {
		return err
	}

	auth, err := c.getAuthTransactor(client)
	if err != nil {
		return err
	}

	tx, err := liquidInstance.EmergencyWithdrawFunds(auth, tokens, newOwner)
	if err != nil {
		log.Error("error when calling EmergencyWithdrawFunds fund ", err)
		return err
	}
	_, err = bind.WaitMined(context.Background(), client, tx)
	return err
}

func (c *EmergencyWithdrawFundCmd) getAuthTransactor(client *ethclient.Client) (*bind.TransactOpts, error) {
	nonce, err := client.PendingNonceAt(context.Background(), c.senderAddress)
	if err != nil {
		log.Error("error when get pending nonce ", err)
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Error("error when get suggest gas price ", err)
		return nil, err
	}

	auth := bind.NewKeyedTransactor(c.privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasPrice = gasPrice

	auth.GasLimit = uint64(10_000_000)

	return auth, nil
}

func (c *EmergencyWithdrawFundCmd) isContractDeployed(client *ethclient.Client, liquidityAddr common.Address) bool {
	bz, err := client.CodeAt(context.Background(), liquidityAddr, nil)
	if err != nil {
		log.Error("Cannot get code at ", liquidityAddr.String(), " err = ", err)
		return false
	}

	return len(bz) > 10
}
