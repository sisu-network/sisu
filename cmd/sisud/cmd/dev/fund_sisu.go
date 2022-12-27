package dev

import (
	"context"
	"path/filepath"
	"strings"
	"sync"
	"time"

	libchain "github.com/sisu-network/lib/chain"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	hutils "github.com/sisu-network/dheart/utils"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/helper"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	ExpectedVaultAddress     = "0x3a84fbbefd21d6a5ce79d54d348344ee11ebd45c"
	ExpectedSisuAddress      = "0xf0d676183dd5ae6b370adddbe770235f23546f9d"
	ExpectedAdaAddress       = "0x3deace7e9c8b6ee632bb71663315d6330914f915"
	CardanoDecimals          = 1000 * 1000
	DefaultCardanoWalletName = "sisu"
	DefaultCardanoPassword   = "12345678910"
)

type fundAccountCmd struct{}

func FundSisu() *cobra.Command {
	cmd := &cobra.Command{
		Use: "fund-sisu",
		Short: `Fund accounts with on a list of chains. Example:
./sisu dev fund-sisu --amount 10
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			chainString, _ := cmd.Flags().GetString(flags.Chains)
			mnemonic, _ := cmd.Flags().GetString(flags.Mnemonic)
			genesisFolder, _ := cmd.Flags().GetString(flags.GenesisFolder)

			sisuRpc, _ := cmd.Flags().GetString(flags.SisuRpc)

			c := &fundAccountCmd{}
			c.fundSisuAccounts(cmd.Context(), chainString, mnemonic,
				sisuRpc, genesisFolder)

			return nil
		},
	}

	cmd.Flags().String(flags.Chains, "ganache1,ganache2", "Names of all chains we want to fund.")
	cmd.Flags().String(flags.ChainUrls, "http://0.0.0.0:7545,http://0.0.0.0:8545", "RPCs of all the chains we want to fund.")
	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum", "Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.SisuRpc, "0.0.0.0:9090", "URL to connect to Sisu. Please do NOT include http:// prefix")
	cmd.Flags().String(flags.GenesisFolder, "./misc/dev", "The genesis folder that contains config files to generate data.")

	return cmd
}

func (c *fundAccountCmd) fundSisuAccounts(ctx context.Context, chainString, mnemonic string,
	sisuRpc, genesisFolder string) {
	chains := strings.Split(chainString, ",")
	vaults := helper.ReadVaults(genesisFolder, chains)

	wg := &sync.WaitGroup{}

	clients := getEthClients(chains, genesisFolder)
	defer func() {
		for _, client := range clients {
			client.Close()
		}
	}()

	// Waits for Sisu to create contract instance in its database. At this stage, the contract is
	// not deployed yet.
	c.waitForPubkeys(ctx, chains, sisuRpc)
	time.Sleep(time.Second * 3)
	allPubKeys := queryPubKeys(ctx, sisuRpc)

	// Fund native cardano.
	if helper.IsCardanoEnabled(genesisFolder) {
		cardanoConfig := helper.ReadCardanoConfig(genesisFolder)
		cardanoKey, ok := allPubKeys[libchain.KEY_TYPE_EDDSA]
		if !ok {
			panic("can not find cardano pub key")
		}

		cardanoAddr := hutils.GetAddressFromCardanoPubkey(cardanoKey)
		log.Info("Sisu Cardano Gateway = ", cardanoAddr)

		tokens := helper.GetTokens(filepath.Join(genesisFolder, "tokens.json"))
		c.fundCardano(genesisFolder, cardanoAddr, mnemonic, cardanoConfig.Secret, sisuRpc, tokens)
	}

	// Fund solana
	if helper.IsSolanaEnabled(genesisFolder) {
		wg.Add(1)

		go func() {
			log.Verbose("Funding on solana chain...")
			c.fundSolana(genesisFolder, mnemonic, allPubKeys[libchain.KEY_TYPE_EDDSA])
			wg.Done()
		}()
	}

	// Fund the accounts with some native ETH and other tokens
	var sisuAccount common.Address
	// Get chain and local chain URL
	pubKeyBytes := allPubKeys[libchain.KEY_TYPE_ECDSA]
	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		panic(err)
	}
	sisuAccount = crypto.PubkeyToAddress(*pubKey)
	log.Info("Sisu account = ", sisuAccount)

	log.Verbose("Funding Sisu's account with some native ETH....")
	wg.Add(len(clients))
	for i, client := range clients {
		go func(i int, client *ethclient.Client, chain string) {
			defer wg.Done()

			c.transferEth(client, sisuRpc, chain, mnemonic, sisuAccount.Hex())
		}(i, client, chains[i])
	}
	wg.Wait()

	log.Verbose("Setting spender for the vault...")

	// Add vault spender
	wg.Add(len(clients))
	for i, client := range clients {
		go func(i int, client *ethclient.Client) {
			c.addVaultSpender(client, mnemonic, common.HexToAddress(vaults[i]), sisuAccount)
			wg.Done()
		}(i, client)
	}
	wg.Wait()
}

func (c *fundAccountCmd) waitForPubkeys(goCtx context.Context, chains []string, sisuRpc string) []string {
	log.Info("Waiting for public keys to be generated in Sisu's db")

	contractAddrs := make([]string, len(chains))
	for {
		grpcConn, err := grpc.Dial(
			sisuRpc,
			grpc.WithInsecure(),
		)
		defer grpcConn.Close()
		if err != nil {
			panic(err)
		}

		allPubKeys := queryPubKeys(goCtx, sisuRpc)
		if allPubKeys == nil || len(allPubKeys) == 0 {
			time.Sleep(time.Second * 3)
			continue
		}

		break
	}

	log.Info("All public keys have been created in Sisu db.")
	return contractAddrs
}

func (c *fundAccountCmd) getTokenAddrs(ctx context.Context, sisuRpc string, tokenSymbols []string, chain string) []string {
	addrs := make([]string, len(tokenSymbols))

	for i, tokenSymbol := range tokenSymbols {
		token := queryToken(ctx, sisuRpc, tokenSymbol)
		for j, addr := range token.Addresses {
			if token.Chains[j] == chain {
				addrs[i] = addr
				break
			}
		}
	}

	return addrs
}
