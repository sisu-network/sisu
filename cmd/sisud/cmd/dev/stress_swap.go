package dev

import (
	"context"
	"fmt"
	"math/big"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	solanago "github.com/gagliardetto/solana-go"
	econfig "github.com/sisu-network/deyes/config"

	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/helper"
	"github.com/sisu-network/sisu/utils"
	"github.com/spf13/cobra"

	libchain "github.com/sisu-network/lib/chain"
)

var (
	NetworkPairs = []string{
		"fantom-testnet__avaxc-testnet",
		"fantom-testnet__binance-testnet",
		"fantom-testnet__goerli-testnet",
		"fantom-testnet__polygon-testnet",
		"fantom-testnet__solana-devnet",

		// "avaxc-testnet__fantom-testnet",
		// "avaxc-testnet__binance-testnet",
		// "avaxc-testnet__goerli-testnet",
		// "avaxc-testnet__polygon-testnet",
		// "avaxc-testnet__solana-devnet",

		// "binance-testnet__fantom-testnet",
		// "binance-testnet__avaxc-testnet",
		// "binance-testnet__goerli-testnet",
		// "binance-testnet__polygon-testnet",
	}

	Recipients = []string{
		"avaxc-testnet__0x1C388F170af377C79b577e676712F1363CceaeeD",
		"fantom-testnet__0x1C388F170af377C79b577e676712F1363CceaeeD",
		"binance-testnet__0x1C388F170af377C79b577e676712F1363CceaeeD",
		"polygon-testnet__0x1C388F170af377C79b577e676712F1363CceaeeD",
		"goerli-testnet__0x1C388F170af377C79b577e676712F1363CceaeeD",
		"solana-testnet__2nDgb2py4cjHy8dxWGPUUR5641SSHCBEjnGGcuyqCKCh",
	}
)

type stressSwapCmd struct{}

func StressSwap() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "stress-swap",
		Long: ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			mnemonic, _ := cmd.Flags().GetString(flags.Mnemonic)
			tokenString, _ := cmd.Flags().GetString(flags.Tokens)
			recipientString, _ := cmd.Flags().GetString(flags.Recipients)
			amount, _ := cmd.Flags().GetInt(flags.Amount)
			sisuRpc, _ := cmd.Flags().GetString(flags.SisuRpc)
			genesisFolder, _ := cmd.Flags().GetString(flags.GenesisFolder)
			deyesUrl, _ := cmd.Flags().GetString(flags.DeyesUrl)

			tokens := strings.Split(tokenString, ",")

			var recipients []string
			if len(recipientString) != 0 {
				recipients = strings.Split(recipientString, ",")
			} else {
				recipients = Recipients
			}

			// Read RPC from the genesis folder.
			deyesChains := helper.ReadDeyesChainConfigs(filepath.Join(genesisFolder, "deyes_chains.json"))

			swapList := make(map[string][]string)
			for _, networkPair := range NetworkPairs {
				src, dst := parseNetworkPair(networkPair)
				if swapList[src] == nil {
					swapList[src] = make([]string, 0)
				}
				swapList[src] = append(swapList[src], dst)
			}

			c := stressSwapCmd{}
			for _, token := range tokens {
				log.Verbosef("Swapping token %s", token)
				wg := &sync.WaitGroup{}
				wg.Add(len(swapList))

				for src, list := range swapList {
					// Run each swap on each list in a separate go routine.
					go func(src string, list []string) {
						log.Verbosef("Swap src = %s, destinations = %v", src, list)
						for _, dst := range list {
							recipient := c.getRecipient(dst, recipients)

							c.doSwap(mnemonic, sisuRpc, genesisFolder, deyesUrl, deyesChains, src, dst, token,
								recipient, amount)

							// Sleep a few second for remote rpc to update sender's nonce.
							time.Sleep(time.Second * 5)
						}
						wg.Done()
					}(src, list)
				}

				wg.Wait()
			}

			return nil
		},
	}

	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum", "Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.SisuRpc, "0.0.0.0:9090", "URL to connect to Sisu. Please do NOT include http:// prefix")
	cmd.Flags().String(flags.Tokens, "TIGER", "The list of tokens to be transferred.")
	cmd.Flags().String(flags.Recipients, "", "List of recipient on each chain. For example 'fantom-testnet__0x1C388F170af377C79b577e676712F1363CceaeeD,avaxc-testnet__0x1C388F170af377C79b577e676712F1363CceaeeD'")
	cmd.Flags().Int(flags.Amount, 1, "The amount of token to be transferred")
	cmd.Flags().String(flags.GenesisFolder, "./misc/dev", "Genesis folder that contains configuration files.")

	return cmd
}

func parseNetworkPair(networkPair string) (string, string) {
	index := strings.Index(networkPair, "__")
	if index < 0 {
		log.Verbose("Invalid network pair: ", networkPair)
		return "", ""
	}

	src := networkPair[:index]
	dst := networkPair[index+2:]

	return src, dst
}

func (c *stressSwapCmd) getRecipient(chain string, recipients []string) string {
	for _, recipient := range recipients {
		index := strings.Index(recipient, "__")
		if index < 0 {
			panic(fmt.Errorf("Invalid recipient: %s", recipient))
		}

		c := recipient[:index]
		if c == chain {
			return recipient[index+2:]
		}
	}

	panic(fmt.Errorf("Cannot fidn recipient for chain : %s", chain))
}

func (c *stressSwapCmd) doSwap(mnemonic, sisuRpc, genesisFolder, deyesUrl string,
	deyesChains []econfig.Chain, src, dst string, tokenSymbol string, recipient string, amount int) {
	token, srcToken, dstToken := getTokenAddrsFromSisu(tokenSymbol, src, dst, sisuRpc)

	log.Verbosef("Trying to swap from %s to %s", src, dst)

	if libchain.IsETHBasedChain(src) {
		// Swap from ETH
		eyesChainCfg := c.getChainConfig(src, deyesChains)
		if eyesChainCfg == nil {
			log.Error("Cannot find config for chain ")
			return
		}

		client, err := ethclient.Dial(eyesChainCfg.Rpcs[0])
		if err != nil {
			log.Verbosef("cannot dial chain %s with url %s, err = %v", src, eyesChainCfg.Rpcs[0], err)
			return
		}

		amountBigInt := big.NewInt(int64(amount))
		amountBigInt = new(big.Int).Mul(amountBigInt, utils.EthToWei)
		vault := getEthVaultAddress(context.Background(), src, sisuRpc)
		swapFromEth(client, mnemonic, vault, dst, srcToken, dstToken, recipient, amountBigInt)
	} else if libchain.IsCardanoChain(src) {
		// Swap from Cardano
		vault := getCardanoVault(context.Background(), sisuRpc)

		amountBigInt := big.NewInt(int64(amount))
		amountBigInt = new(big.Int).Mul(amountBigInt, utils.ONE_ADA_IN_LOVELACE)

		cardanoCfg := helper.ReadCardanoConfig(genesisFolder)

		swapFromCardano(src, dst, token, recipient, vault, amountBigInt, cardanoCfg.Chain,
			cardanoCfg.Secret, mnemonic, deyesUrl)
	} else if libchain.IsSolanaChain(src) {
		solanaCfg, err := helper.ReadCmdSolanaConfig(genesisFolder)
		if err != nil {
			log.Errorf("Failed to read solana config, err = %v", err)
			return
		}

		tokenAddr := token.GetAddressForChain(solanaCfg.Chain)
		// Build ata address from recipient's wallet
		ataAddr, _, err := solanago.FindAssociatedTokenAddress(
			solanago.MustPublicKeyFromBase58(recipient),
			solanago.MustPublicKeyFromBase58(tokenAddr),
		)
		if err != nil {
			log.Errorf("Failed to find ata address for address %s and token %s, err = %v", recipient,
				solanaCfg.Chain)
			return
		}

		decimal := token.GetDecimalsForChain(solanaCfg.Chain)
		amountBigInt := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimal)), nil)
		amountBigInt = amountBigInt.Mul(amountBigInt, big.NewInt(int64(amount)))

		dstChainId := libchain.GetChainIntFromId(dst)

		swapFromSolana(genesisFolder, solanaCfg.Chain, mnemonic, tokenAddr, ataAddr.String(),
			dstChainId.Uint64(), amountBigInt.Uint64())
	}
}

func (c *stressSwapCmd) getChainConfig(chain string, deyesChains []econfig.Chain) *econfig.Chain {
	for _, cfg := range deyesChains {
		if cfg.Chain == chain {
			return &cfg
		}
	}

	return nil
}
