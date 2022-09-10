package dev

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	libchain "github.com/sisu-network/lib/chain"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/echovl/cardano-go"
	"github.com/echovl/cardano-go/blockfrost"
	"github.com/echovl/cardano-go/wallet"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	hutils "github.com/sisu-network/dheart/utils"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/sisu-network/sisu/contracts/eth/erc20"
	"github.com/sisu-network/sisu/contracts/eth/vault"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"
	tssTypes "github.com/sisu-network/sisu/x/sisu/types"
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
			urlString, _ := cmd.Flags().GetString(flags.ChainUrls)
			mnemonic, _ := cmd.Flags().GetString(flags.Mnemonic)
			tokenString, _ := cmd.Flags().GetString(flags.Erc20Symbols)
			vaultString, _ := cmd.Flags().GetString(flags.VaultAddrs)
			cardanoSecret, _ := cmd.Flags().GetString(flags.CardanoSecret)
			cardanoNetwork, _ := cmd.Flags().GetString(flags.CardanoNetwork)

			sisuRpc, _ := cmd.Flags().GetString(flags.SisuRpc)
			tokens := strings.Split(tokenString, ",")

			c := &fundAccountCmd{}
			c.fundSisuAccounts(cmd.Context(), chainString, urlString, mnemonic, tokens, vaultString,
				sisuRpc, cardanoNetwork, cardanoSecret)

			return nil
		},
	}

	cmd.Flags().String(flags.Chains, "ganache1,ganache2", "Names of all chains we want to fund.")
	cmd.Flags().String(flags.ChainUrls, "http://0.0.0.0:7545,http://0.0.0.0:8545", "RPCs of all the chains we want to fund.")
	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum", "Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.SisuRpc, "0.0.0.0:9090", "URL to connect to Sisu. Please do NOT include http:// prefix")
	cmd.Flags().String(flags.VaultAddrs, fmt.Sprintf("%s,%s", ExpectedVaultAddress, ExpectedVaultAddress), "List of vault addresses")
	cmd.Flags().String(flags.Erc20Symbols, "SISU,ADA", "List of ERC20 to approve")
	cmd.Flags().String(flags.CardanoSecret, "", "The blockfrost secret to interact with cardano network.")
	cmd.Flags().String(flags.CardanoNetwork, "cardano-testnet", "The Cardano network that we are interacting with.")
	cmd.Flags().String(flags.GenesisFolder, "./misc/dev", "Relative path to the folder that contains genesis configuration.")

	return cmd
}

func (c *fundAccountCmd) fundSisuAccounts(ctx context.Context, chainString, urlString, mnemonic string,
	tokens []string, vaultString string, sisuRpc, cardanoNetwork, cardanoSecret string) {
	chains := strings.Split(chainString, ",")
	vaults := strings.Split(vaultString, ",")

	wg := &sync.WaitGroup{}

	clients := getEthClients(urlString)
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
	if len(cardanoSecret) > 0 {
		cardanoKey, ok := allPubKeys[libchain.KEY_TYPE_EDDSA]
		if !ok {
			panic("can not find cardano pub key")
		}

		cardanoAddr := hutils.GetAddressFromCardanoPubkey(cardanoKey)
		log.Info("Sisu Cardano Gateway = ", cardanoAddr)
		c.fundCardano(cardanoAddr, mnemonic, cardanoNetwork, cardanoSecret, sisuRpc, tokens)
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

func (c *fundAccountCmd) getMultiAsset(sisuRpc, cardanoNetwork string, tokens []string, amt uint64) *cardano.MultiAsset {
	tokenAddrs := c.getTokenAddrs(context.Background(), sisuRpc, tokens, cardanoNetwork)
	m := make(map[string]*cardano.Assets)
	for _, tokenAddr := range tokenAddrs {
		index := strings.Index(tokenAddr, ":")
		policyString := tokenAddr[:index]
		assetName := tokenAddr[index+1:]

		if m[policyString] == nil {
			m[policyString] = cardano.NewAssets()
		}

		asset := cardano.NewAssetName(assetName)
		m[policyString].Set(asset, cardano.BigNum(amt*CardanoDecimals))
	}

	multiAsset := cardano.NewMultiAsset()
	for policy, assets := range m {
		policyHash, err := cardano.NewHash28(policy)
		if err != nil {
			err := fmt.Errorf("error when parsing policyID hash: %v", err)
			panic(err)
		}
		policyID := cardano.NewPolicyIDFromHash(policyHash)
		multiAsset.Set(policyID, assets)
	}

	return multiAsset
}

func (c *fundAccountCmd) fundCardano(receiver cardano.Address, funderMnemonic string,
	cardanoNetwork, blockfrostSecret string, sisuRpc string, tokens []string) {
	node := blockfrost.NewNode(cardano.Testnet, blockfrostSecret)
	opts := &wallet.Options{
		Node: node,
	}
	client := wallet.NewClient(opts)
	funderWallet, err := c.getWalletFromMnemonic(client, DefaultCardanoWalletName, DefaultCardanoPassword, funderMnemonic)
	if err != nil {
		panic(err)
	}

	addrs, err := funderWallet.Addresses()
	if err != nil {
		panic(err)
	}
	funderAddr := addrs[0]
	log.Info("Cardano funder address = ", funderAddr.String())

	// fund 30 ADA and 1000 WRAP_ADA
	txHash, err := funderWallet.Transfer(receiver, cardano.NewValueWithAssets(30*CardanoDecimals,
		c.getMultiAsset(sisuRpc, cardanoNetwork, tokens, 1e3)), nil) // 30ADA
	if err != nil {
		panic(err)
	}

	log.Infof("Address funded = %s, txHash = %s, explorer: https://testnet.cardanoscan.io/transaction/%s\n",
		receiver, txHash.String(), txHash.String())
}

func (c *fundAccountCmd) getWalletFromMnemonic(client *wallet.Client, name, password, mnemonic string) (*wallet.Wallet, error) {
	w, err := client.RestoreWallet(name, password, mnemonic)
	if err != nil {
		return nil, err
	}

	return w, nil
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

// isContractDeployed checks if a contract has been deployed at a specific address so that
// we do not have to deploy again.
func (c *fundAccountCmd) isContractDeployed(client *ethclient.Client, tokenAddress common.Address) bool {
	bz, err := client.CodeAt(context.Background(), tokenAddress, nil)
	if err != nil {
		log.Error("Cannot get code at ", tokenAddress.String(), " err = ", err)
		return false
	}

	return len(bz) > 10
}

// transferEth transfers a specific ETH amount to an address.
func (c *fundAccountCmd) transferEth(client *ethclient.Client, sisuRpc, chain, mnemonic, recipient string) {
	ch, err := queryChain(context.Background(), sisuRpc, chain)
	if err != nil {
		panic(fmt.Errorf("failed to get chain, err = %v", err))
	}
	genesisGas := big.NewInt(ch.GasPrice)

	_, account := getPrivateKey(mnemonic)

	log.Info("from address = ", account.String(), " to Address = ", recipient)

	nonce, err := client.PendingNonceAt(context.Background(), account)
	if err != nil {
		panic(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}

	if gasPrice.Cmp(genesisGas) < 0 {
		gasPrice = genesisGas
	}

	// Add some 10% premimum to the gas price
	gasPrice = gasPrice.Mul(gasPrice, big.NewInt(110))
	gasPrice = gasPrice.Quo(gasPrice, big.NewInt(100))

	log.Info("Gas price = ", gasPrice, " on chain ", chain)

	// 0.05 ETH
	amount := new(big.Int).Div(utils.EthToWei, big.NewInt(20))

	gasLimit := uint64(22000) // in units

	amountFloat := new(big.Float).Quo(new(big.Float).SetInt(amount), new(big.Float).SetInt(utils.ONE_ETHER_IN_WEI))
	log.Info("Amount in ETH: ", amountFloat, " on chain ", chain)

	toAddress := common.HexToAddress(recipient)
	var data []byte
	tx := ethtypes.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, data)

	privateKey, _ := getPrivateKey(mnemonic)
	signedTx, err := ethtypes.SignTx(tx, getSigner(client), privateKey)

	log.Info("Tx hash = ", signedTx.Hash())

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(fmt.Errorf("Failed to transfer ETH on chain %s, err = %s", chain, err))
	}

	bind.WaitDeployed(context.Background(), client, signedTx)
	time.Sleep(time.Second * 3)
}

func queryChain(ctx context.Context, sisuRpc string, chain string) (*types.Chain, error) {
	grpcConn, err := grpc.Dial(
		sisuRpc,
		grpc.WithInsecure(),
	)
	defer grpcConn.Close()
	if err != nil {
		panic(err)
	}

	queryClient := tssTypes.NewTssQueryClient(grpcConn)
	res, err := queryClient.QueryChain(ctx, &tssTypes.QueryChainRequest{
		Chain: chain,
	})
	if err != nil {
		return nil, err
	}

	return res.Chain, nil
}

func queryPubKeys(ctx context.Context, sisuRpc string) map[string][]byte {
	grpcConn, err := grpc.Dial(
		sisuRpc,
		grpc.WithInsecure(),
	)
	defer grpcConn.Close()
	if err != nil {
		panic(err)
	}

	queryClient := tssTypes.NewTssQueryClient(grpcConn)

	res, err := queryClient.AllPubKeys(ctx, &tssTypes.QueryAllPubKeysRequest{})
	if err != nil {
		panic(err)
	}

	return res.Pubkeys
}

func (c *fundAccountCmd) addVaultSpender(client *ethclient.Client, mnemonic string, vaultAddr, spender common.Address) {
	fmt.Println("Add vault spender, vault = ", vaultAddr.String(), " spender = ", spender.String())
	vaultInstance, err := vault.NewVault(vaultAddr, client)
	if err != nil {
		panic(err)
	}

	auth, err := getAuthTransactor(client, mnemonic)
	if err != nil {
		panic(err)
	}

	tx, err := vaultInstance.AddSpender(auth, spender)
	if err != nil {
		panic(err)
	}

	bind.WaitDeployed(context.Background(), client, tx)
	time.Sleep(time.Second * 5)
}

func (c *fundAccountCmd) queryAllownace(client *ethclient.Client,
	tokenAddr, owner, target string) *big.Int {
	store, err := erc20.NewErc20(common.HexToAddress(tokenAddr), client)
	if err != nil {
		panic(err)
	}

	balance, err := store.Allowance(nil, common.HexToAddress(owner), common.HexToAddress(target))
	if err != nil {
		panic(err)
	}

	return balance
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
