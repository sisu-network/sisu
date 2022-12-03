package dev

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"path/filepath"
	"strings"
	"sync"
	"time"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	econfig "github.com/sisu-network/deyes/config"

	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/helper"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/chains/solana"
	"github.com/sisu-network/sisu/x/sisu/types"
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

		"avaxc-testnet__fantom-testnet",
		"avaxc-testnet__binance-testnet",
		"avaxc-testnet__goerli-testnet",
		"avaxc-testnet__polygon-testnet",
		"avaxc-testnet__solana-devnet",

		// "binance-testnet__fantom-testnet",
		// "binance-testnet__avaxc-testnet",
		// "binance-testnet__goerli-testnet",
		// "binance-testnet__polygon-testnet",
		// "binance-testnet__solana-devnet",
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
			amount, _ := cmd.Flags().GetInt(flags.Amount)
			sisuRpc, _ := cmd.Flags().GetString(flags.SisuRpc)
			genesisFolder, _ := cmd.Flags().GetString(flags.GenesisFolder)
			deyesUrl, _ := cmd.Flags().GetString(flags.DeyesUrl)

			tokens := strings.Split(tokenString, ",")

			// Generate some random addresses
			ethAddr, solanaAddr := getRandomAddresses()
			log.Verbosef("Eth address = %s, solana address = %s", ethAddr.Hex(), solanaAddr.String())

			// Read RPC from the genesis folder.
			deyesChains := helper.ReadDeyesChainConfigs(filepath.Join(genesisFolder, "deyes_chains.json"))

			hasSolana := false
			swapList := make(map[string][]string)
			for _, networkPair := range NetworkPairs {
				src, dst := parseNetworkPair(networkPair)
				if swapList[src] == nil {
					swapList[src] = make([]string, 0)
				}
				swapList[src] = append(swapList[src], dst)

				if libchain.IsSolanaChain(dst) {
					hasSolana = true
				}
			}

			// Generate random addresses.
			c := stressSwapCmd{}
			// If there is some solana chain, we have to create a new account for the address and all the
			// token ATA accounts.
			if hasSolana {
				for _, eyesCfg := range deyesChains {
					if libchain.IsSolanaChain(eyesCfg.Chain) {
						client := rpc.New(eyesCfg.Rpcs[0])
						wsClient, err := ws.Connect(context.Background(), eyesCfg.Wss[0])
						if err != nil {
							panic(err)
						}

						transferSOL(client, wsClient, mnemonic, solanaAddr.String(), uint64(1_000_000)) // 0.01 SOL
						for _, tokenId := range tokens {
							token := queryToken(context.Background(), sisuRpc, tokenId)
							tokenAddr := token.GetAddressForChain(eyesCfg.Chain)

							// Create the ata address
							ata, _, err := createSolanaAta(client, wsClient, mnemonic, solanaAddr,
								solanago.MustPublicKeyFromBase58(tokenAddr))
							if err != nil {
								panic(err)
							}

							log.Verbosef("ATA address created for token %s: %s", tokenId, ata.String())
						}
						client.Close()
						wsClient.Close()
					}
				}
			}

			for _, tokenId := range tokens {
				token := queryToken(context.Background(), sisuRpc, tokenId)

				log.Verbosef("Swapping token %s", tokenId)
				wg := &sync.WaitGroup{}
				wg.Add(len(swapList))

				for src, list := range swapList {
					// Run each swap on each list in a separate go routine.
					go func(src string, list []string) {
						log.Verbosef("Swap src = %s, destinations = %v", src, list)
						for _, dst := range list {
							var recipient string
							if libchain.IsETHBasedChain(dst) {
								recipient = ethAddr.Hex()
							} else if libchain.IsSolanaChain(dst) {
								recipient = solanaAddr.String()
							} else {
								panic(fmt.Errorf("Unsupported chain dst = %s", dst))
							}

							c.doSwap(mnemonic, sisuRpc, genesisFolder, deyesUrl, deyesChains, src, dst, token,
								recipient, amount)

							// Sleep a few second for remote rpc to update sender's nonce.
							time.Sleep(time.Second * 3)
						}
						wg.Done()
					}(src, list)
				}

				wg.Wait()
			}

			sleepSec := 40
			log.Verbose(fmt.Sprintf(
				"Sleeping %d second for all trannsaction to finalize. You can increase sleep time if needed ...",
				sleepSec,
			))
			time.Sleep(time.Second * time.Duration(sleepSec))

			// Verify balances on all chains
			for _, tokenId := range tokens {
				token := queryToken(context.Background(), sisuRpc, tokenId)
				assertBalance(swapList, genesisFolder, token, ethAddr, solanaAddr, amount)
			}

			return nil
		},
	}

	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum", "Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.SisuRpc, "0.0.0.0:9090", "URL to connect to Sisu. Please do NOT include http:// prefix")
	cmd.Flags().String(flags.Tokens, "TIGER", "The list of tokens to be transferred.")
	cmd.Flags().Int(flags.Amount, 1, "The amount of token to be transferred")
	cmd.Flags().String(flags.GenesisFolder, "./misc/dev", "Genesis folder that contains configuration files.")

	return cmd
}

func getRandomAddresses() (ethcommon.Address, solanago.PublicKey) {
	// ETH
	ethKey, err := ethcrypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	publicKey := ethKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("Invalid public key")
	}
	address := ethcrypto.PubkeyToAddress(*publicKeyECDSA)

	// Solana
	solanaKey, err := solanago.NewRandomPrivateKey()
	if err != nil {
		panic(err)
	}

	return address, solanaKey.PublicKey()
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
	deyesChains []econfig.Chain, src, dst string, token *types.Token, recipient string, amount int) {
	srcToken := token.GetAddressForChain(src)
	dstToken := token.GetAddressForChain(src)
	solanaCfg, err := helper.ReadCmdSolanaConfig(filepath.Join(genesisFolder, "solana.json"))
	if err != nil {
		log.Errorf("Failed to read solana config, err = %v", err)
		return
	}

	// If the destination chain is a solana chain, create an ata account for it.
	if libchain.IsSolanaChain(dst) {
		tokenAddr := token.GetAddressForChain(solanaCfg.Chain)
		ata, _, err := solanago.FindAssociatedTokenAddress(
			solanago.MustPublicKeyFromBase58(recipient),
			solanago.MustPublicKeyFromBase58(tokenAddr),
		)
		if err != nil {
			panic(err)
		}

		recipient = ata.String()
	}

	log.Verbosef("Trying to swap from %s to %s, recipient = %s", src, dst, recipient)

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
		tokenAddr := token.GetAddressForChain(solanaCfg.Chain)

		decimal := token.GetDecimalsForChain(solanaCfg.Chain)
		amountBigInt := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimal)), nil)
		amountBigInt = amountBigInt.Mul(amountBigInt, big.NewInt(int64(amount)))

		dstChainId := libchain.GetChainIntFromId(dst)

		swapFromSolana(genesisFolder, solanaCfg.Chain, mnemonic, tokenAddr, recipient,
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

func (c *stressSwapCmd) getSolanaClientAndWss(genesisFolder string) (*rpc.Client, *ws.Client) {
	deyesChains := helper.ReadDeyesChainConfigs(filepath.Join(genesisFolder, "deyes_chains.json"))
	for _, cfg := range deyesChains {
		if libchain.IsSolanaChain(cfg.Chain) {
			client := rpc.New(cfg.Rpcs[0])
			wsClient, err := ws.Connect(context.Background(), cfg.Wss[0])
			if err != nil {
				panic(err)
			}

			return client, wsClient
		}
	}

	panic(fmt.Errorf("Cannot find config for solaan chain, genesis folder = %s", genesisFolder))
}

func getEthClient(genesisFolder string, chain string) *ethclient.Client {
	deyesChains := helper.ReadDeyesChainConfigs(filepath.Join(genesisFolder, "deyes_chains.json"))
	for _, cfg := range deyesChains {
		if cfg.Chain == chain {
			fmt.Println("cfg.Rpcs[0] = ", cfg.Rpcs[0])
			client, err := ethclient.Dial(cfg.Rpcs[0])
			if err != nil {
				panic(err)
			}
			return client
		}
	}

	panic(fmt.Errorf("Cannot find chain %s in the config", chain))
}

func getSolanaClient(genesisFolder string) *rpc.Client {
	solanaConfig, err := helper.ReadCmdSolanaConfig(filepath.Join(genesisFolder, "solana.json"))
	if err != nil {
		panic(err)
	}

	return rpc.New(solanaConfig.Rpc)
}

func assertBalance(swapList map[string][]string, genesisFolder string, token *types.Token,
	ethAddr ethcommon.Address, solanaAddr solanago.PublicKey, amount int) {
	expectedAmounts := make(map[string]int64)
	for _, list := range swapList {
		for _, dst := range list {
			if libchain.IsETHBasedChain(dst) {
				expectedAmounts[dst] = expectedAmounts[dst] + 1
			}
		}
	}

	for chain, value := range expectedAmounts {
		tokenAddr := token.GetAddressForChain(chain)
		maxValue, err := token.GetUnits(chain, int(value))
		if err != nil {
			panic(err)
		}

		if libchain.IsETHBasedChain(chain) {
			client := getEthClient(genesisFolder, chain)
			// Query ERC20 balance
			balance, err := queryErc20Balance(client, tokenAddr, ethAddr.String())
			if err != nil {
				panic(err)
			}

			log.Verbosef("Balance on chain %s = %s, max = %s", chain, balance.String(), maxValue.String())
		}

		if libchain.IsSolanaChain(chain) {
			// Query ATA balance
			client := getSolanaClient(genesisFolder)
			ata, _, err := solanago.FindAssociatedTokenAddress(
				solanaAddr,
				solanago.MustPublicKeyFromBase58(tokenAddr),
			)
			if err != nil {
				panic(err)
			}

			balance, err := solana.QuerySolanaAccountBalance(client, ata.String())
			if err != nil {
				panic(err)
			}

			log.Verbosef("Balance on chain %s = %s, max = %s", chain, balance.String(), maxValue.String())
		}
	}
}
