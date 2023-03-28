package dev

import (
	"context"
	"math/big"
	"strings"
	"time"

	liskcrypto "github.com/sisu-network/deyes/chains/lisk/crypto"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/spf13/cobra"
)

type LiskBotCmd struct {
}

func LiskBot() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lisk-bot",
		Long: `Swap between lisk and other chains`,
		RunE: func(cmd *cobra.Command, args []string) error {
			mnemonic, _ := cmd.Flags().GetString(flags.Mnemonic)
			chainString, _ := cmd.Flags().GetString(flags.Chains)
			sisuRpc, _ := cmd.Flags().GetString(flags.SisuRpc)
			genesisFolder, _ := cmd.Flags().GetString(flags.GenesisFolder)
			chains := strings.Split(chainString, ",")

			allPubKeys := queryPubKeys(cmd.Context(), sisuRpc)
			c := &LiskBotCmd{}
			c.Run(mnemonic, genesisFolder, chains, sisuRpc, allPubKeys)

			return nil
		},
	}

	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum", "Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.Chains, "ganache1,ganache2", "Names of all chains we want to fund.")
	cmd.Flags().String(flags.GenesisFolder, "./misc/dev", "The genesis folder that contains config files to generate data.")
	cmd.Flags().String(flags.SisuRpc, "0.0.0.0:9090", "URL to connect to Sisu. Please do NOT include http:// prefix")

	return cmd
}

func (cmd *LiskBotCmd) Run(mnemonic, genesisFolder string, chains []string, sisuRpc string,
	allPubKeys map[string][]byte) {
	liskChain := "lisk-testnet"
	_, ethAddress := getPrivateKey(mnemonic)

	// Get original balances
	originalBalances := make(map[string]*big.Int)
	for _, chain := range chains {
		clients := getEthClients([]string{chain}, genesisFolder)
		if len(clients) == 0 {
			log.Error("None of the clients in the genesis folder is healthy")
			continue
		}

		token := queryToken(context.Background(), sisuRpc, "LSK")
		srcToken := token.GetAddressForChain(chain)
		if token == nil {
			continue
		}

		client := clients[0]
		liskBalance, err := queryErc20Balance(client, srcToken, ethAddress.String())
		if err != nil {
			log.Errorf("Failed to get lisk ERC20 balance, err = %s", err)
			continue
		}

		originalBalances[chain] = liskBalance
	}

	for {
		for _, chain := range chains {
			// Swap from Lisk
			amount := 1
			swapFromLisk(genesisFolder, mnemonic, chain, allPubKeys[libchain.KEY_TYPE_EDDSA], ethAddress.String(),
				uint64(amount*100_000_000))

			// Swap to Lisk
			clients := getEthClients([]string{chain}, genesisFolder)
			if len(clients) == 0 {
				log.Error("None of the clients in the genesis folder is healthy")
				continue
			}

			client := clients[0]
			token := queryToken(context.Background(), sisuRpc, "LSK")
			srcToken := token.GetAddressForChain(chain)
			if token == nil {
				continue
			}
			vault := getEthVaultAddress(context.Background(), chains[0], sisuRpc)
			faucetPubKey := liskcrypto.GetPublicKeyFromSecret(mnemonic)
			lisk32 := liskcrypto.GetLisk32AddressFromPublickey(faucetPubKey)

			time.Sleep(time.Second * 60)

			// Get the Lisk balance
			liskBalance, err := queryErc20Balance(client, srcToken, ethAddress.String())
			if err != nil {
				log.Errorf("Failed to get lisk ERC20 balance, err = %s", err)
				continue
			}

			origBalance := originalBalances[chain]
			if liskBalance.Cmp(origBalance) < 0 {
				continue
			}

			diff := liskBalance.Sub(liskBalance, origBalance)
			log.Verbosef("Lisk balance diff on chain %s = %s\n", chain, diff)
			swapFromEth(client, mnemonic, vault, liskChain, srcToken, "", lisk32, diff)
		}

		time.Sleep(time.Minute * 2)
	}
}
