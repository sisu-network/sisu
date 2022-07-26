package dev

import (
	"context"
	"fmt"
	"math/big"
	"path/filepath"
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
	"github.com/sisu-network/sisu/cmd/sisud/cmd/helper"
	"github.com/sisu-network/sisu/contracts/eth/erc20"
	liquidity "github.com/sisu-network/sisu/contracts/eth/liquiditypool"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu"
	tssTypes "github.com/sisu-network/sisu/x/sisu/types"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	ExpectedSisuAddress       = "0x3A84fBbeFD21D6a5ce79D54d348344EE11EBd45C"
	ExpectedAdaAddress        = "0xf0D676183dD5ae6b370adDdbE770235F23546f9d"
	ExpectedLiquidPoolAddress = "0x3DeaCe7E9C8b6ee632bb71663315d6330914f915"
	CardanoDecimals           = 1000 * 1000
	DefaultCardanoWalletName  = "sisu"
	DefaultCardanoPassword    = "12345678910"
)

var (
	WrapADA = cardano.NewAssetName("WRAP_ADA")
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
			liquidityAddrString, _ := cmd.Flags().GetString(flags.LiquidityAddrs)
			cardanoSecret, _ := cmd.Flags().GetString(flags.CardanoSecret)
			cardanoFunderMnemonic, _ := cmd.Flags().GetString(flags.CardanoFunderMnemonic)
			genesisFolder, _ := cmd.Flags().GetString(flags.GenesisFolder)

			sisuRpc, _ := cmd.Flags().GetString(flags.SisuRpc)
			tokens := strings.Split(tokenString, ",")

			c := &fundAccountCmd{}
			c.fundSisuAccounts(cmd.Context(), chainString, urlString, mnemonic, genesisFolder, tokens,
				liquidityAddrString, sisuRpc, cardanoSecret, cardanoFunderMnemonic)

			return nil
		},
	}

	cmd.Flags().String(flags.Chains, "ganache1,ganache2", "Names of all chains we want to fund.")
	cmd.Flags().String(flags.ChainUrls, "http://0.0.0.0:7545,http://0.0.0.0:8545", "RPCs of all the chains we want to fund.")
	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum", "Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.SisuRpc, "0.0.0.0:9090", "URL to connect to Sisu. Please do NOT include http:// prefix")
	cmd.Flags().String(flags.LiquidityAddrs, fmt.Sprintf("%s,%s", ExpectedLiquidPoolAddress, ExpectedLiquidPoolAddress), "List of liquidity pool addresses")
	cmd.Flags().String(flags.Erc20Symbols, "SISU,ADA", "List of ERC20 to approve")
	cmd.Flags().String(flags.CardanoSecret, "", "The blockfrost secret to interact with cardano network.")
	cmd.Flags().String(flags.CardanoFunderMnemonic, "", "Mnemonic of funder wallet which already has a lot of test tokens")
	cmd.Flags().String(flags.GenesisFolder, "./misc/dev", "Relative path to the folder that contains genesis configuration.")

	return cmd
}

func (c *fundAccountCmd) fundSisuAccounts(ctx context.Context, chainString, urlString, mnemonic string,
	genesisFolder string, tokens []string, liquidityAddrString, sisuRpc, cardanoSecret, cardanoFunderMnemonic string) {
	chains := strings.Split(chainString, ",")
	liquidityAddrs := strings.Split(liquidityAddrString, ",")

	log.Info("liquidityAddrs = ", liquidityAddrs)

	wg := &sync.WaitGroup{}

	clients := getEthClients(urlString)
	defer func() {
		for _, client := range clients {
			client.Close()
		}
	}()

	// Waits for Sisu to create contract instance in its database. At this stage, the contract is
	// not deployed yet.
	c.waitForGatewayCreationInSisuDb(ctx, chains, sisuRpc)
	time.Sleep(time.Second * 3)
	allPubKeys := queryPubKeys(ctx, sisuRpc)

	// Fund native cardano.
	if len(cardanoFunderMnemonic) > 0 && len(cardanoSecret) > 0 {
		cardanoKey, ok := allPubKeys[libchain.KEY_TYPE_EDDSA]
		if !ok {
			panic("can not find cardano pub key")
		}

		cardanoAddr := hutils.GetAddressFromCardanoPubkey(cardanoKey)
		log.Info("Sisu Cardano Gateway = ", cardanoAddr)
		c.fundCardano(cardanoAddr, cardanoFunderMnemonic, cardanoSecret)
	}

	// Fund the accounts with some native ETH
	var tssPubAddr common.Address
	wg.Add(len(clients))
	for i, client := range clients {
		go func(i int, client *ethclient.Client, chain string) {
			defer wg.Done()

			// Get chain and local chain URL
			pubKeyBytes := allPubKeys[libchain.KEY_TYPE_ECDSA]
			pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
			if err != nil {
				panic(err)
			}
			tssPubAddr = crypto.PubkeyToAddress(*pubKey)

			c.transferEth(client, genesisFolder, chain, mnemonic, tssPubAddr.Hex())
		}(i, client, chains[i])
	}
	wg.Wait()

	// Waits until all gateway contracts are deployed.
	log.Info("Now we wait until gateway contracts are deployed on all chains.")
	gateways := c.waitForGatewayDeployed(ctx, chains, sisuRpc)

	// Approve the contract with some preallocated token from account0
	wg.Add(len(clients))
	for i, client := range clients {
		go func(i int, client *ethclient.Client) {
			addrs := c.getTokenAddrs(ctx, sisuRpc, tokens, chains[i])
			for _, addr := range addrs {
				log.Infof("Approve token with address %s for gateway %s", addr, gateways[i])
				approveAddress(client, mnemonic, addr, gateways[i])
			}

			wg.Done()
		}(i, client)
	}
	wg.Wait()
	log.Info("Gateway approval done!")

	// Set gateway for the liquidity
	log.Info("Set gateway address for liquidity pool")
	wg.Add(len(gateways))
	for i, client := range clients {
		go func(i int, client *ethclient.Client) {
			c.setGatewayForLiquidity(client, mnemonic, liquidityAddrs[i], gateways[i])
			wg.Done()
		}(i, client)
	}
	wg.Wait()

}

// Get WRAP_ADA (hex: 575241505f414441) token.
// PolicyID (dc89700b3adf88f6b520aba2f3cfa4c26fa7a19bd8eadf430d73b9d4) got from there:
// https://testnet.cardanoscan.io/transaction/31c019b737edc7b54ae60a87f372f860715e8bb02b979ed853395ccbf4bf0209
func getMultiAsset(amt uint64) *cardano.MultiAsset {
	policyHash, err := cardano.NewHash28("dc89700b3adf88f6b520aba2f3cfa4c26fa7a19bd8eadf430d73b9d4")
	if err != nil {
		err := fmt.Errorf("error when parsing policyID hash: %v", err)
		panic(err)
	}

	policyID := cardano.NewPolicyIDFromHash(policyHash)
	asset := cardano.NewAssets().Set(WrapADA, cardano.BigNum(amt*CardanoDecimals))

	return cardano.NewMultiAsset().Set(policyID, asset)
}

func (c *fundAccountCmd) fundCardano(receiver cardano.Address, funderMnemonic string, blockfrostSecret string) {
	node := blockfrost.NewNode(cardano.Testnet, blockfrostSecret)
	opts := &wallet.Options{
		Node: node,
	}
	client := wallet.NewClient(opts)

	// funder address: addr_test1vqyqp03az6w8xuknzpfup3h7ghjwu26z7xa6gk7l9j7j2gs8zfwcy
	funderWallet, err := c.getWalletFromMnemonic(client, DefaultCardanoWalletName, DefaultCardanoPassword, funderMnemonic)
	if err != nil {
		panic(err)
	}

	funderAddr, err := funderWallet.AddAddress()
	if err != nil {
		panic(err)
	}
	log.Info("Cardano funder address = ", funderAddr.String())

	// fund 10 ADA and 1000 WRAP_ADA
	txHash, err := funderWallet.Transfer(receiver, cardano.NewValueWithAssets(10*CardanoDecimals, getMultiAsset(1e3)), nil) // 10 ADA
	if err != nil {
		panic(err)
	}

	log.Infof("Funded 10 ADA and 1000 WRAP_ADA for address %s, txHash = %s, "+
		"explorer: https://testnet.cardanoscan.io/transaction/%s\n", receiver, txHash.String(), txHash.String())
}

func (c *fundAccountCmd) getWalletFromMnemonic(client *wallet.Client, name, password, mnemonic string) (*wallet.Wallet, error) {
	w, err := client.RestoreWallet(name, password, mnemonic)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (c *fundAccountCmd) waitForGatewayCreationInSisuDb(goCtx context.Context, chains []string, sisuRpc string) []string {
	log.Info("Waiting for all contract created in Sisu's db")

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

		queryClient := tssTypes.NewTssQueryClient(grpcConn)

		done := true
		for i, chain := range chains {
			res, err := queryClient.QueryContract(goCtx, &tssTypes.QueryContractRequest{
				Chain: chain,
				Hash:  sisu.SupportedContracts[sisu.ContractErc20Gateway].AbiHash,
			})

			if err != nil {
				log.Error("err = ", err)
				log.Verbose("Contract has not been created yet ", i, " yet. Keep sleeping...")
				time.Sleep(time.Second * 3)
				done = false
				break
			}

			contractAddrs[i] = res.Contract.Address
		}

		if done {
			break
		}
	}

	log.Info("All contracts have been created in Sisu db.")
	return contractAddrs
}

func (c *fundAccountCmd) setGatewayForLiquidity(client *ethclient.Client, mnemonic string, liquidityAddr, gatewayAddr string) {
	log.Infof("Setting gateway for liquidity pool, gateway address: %s, liquidity pool address %s\n", gatewayAddr, liquidityAddr)

	contract, err := liquidity.NewLiquiditypool(common.HexToAddress(liquidityAddr), client)
	if err != nil {
		panic(err)
	}

	opts, err := getAuthTransactor(client, mnemonic)
	if err != nil {
		panic(err)
	}

	tx, err := contract.SetGateway(opts, common.HexToAddress(gatewayAddr))
	if err != nil {
		panic(err)
	}
	txReceipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		panic(err)
	}

	if txReceipt.Status != ethtypes.ReceiptStatusSuccessful {
		log.Info("Tx Hash = ", txReceipt.TxHash)
		panic("tx grant liquidity pool access failed")
	}

	log.Info("Grant access for gateway successfully")
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
func (c *fundAccountCmd) transferEth(client *ethclient.Client, genesisFolder, chain, mnemonic, recipient string) {
	chains := helper.GetChains(filepath.Join(genesisFolder, "chains.json"))
	var genesisGas *big.Int
	for _, c := range chains {
		if c.Id == chain {
			genesisGas = big.NewInt(c.GasPrice)
		}
	}
	if genesisGas == nil {
		panic(fmt.Errorf("Genesis gas is nil on chain %s", chain))
	}

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

	amount := new(big.Int).Mul(big.NewInt(8_000_000), gasPrice)
	// amount = amount * 1.4
	amount = amount.Mul(amount, big.NewInt(14))
	amount = amount.Quo(amount, big.NewInt(10))

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
		panic(err)
	}

	bind.WaitDeployed(context.Background(), client, signedTx)
	time.Sleep(time.Second * 3)
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

func (c *fundAccountCmd) waitForGatewayDeployed(goCtx context.Context, chains []string,
	sisuRpc string) []string {
	addrs := make([]string, len(chains))

	for {
		grpcConn, err := grpc.Dial(
			sisuRpc,
			grpc.WithInsecure(),
		)
		defer grpcConn.Close()
		if err != nil {
			panic(err)
		}

		queryClient := tssTypes.NewTssQueryClient(grpcConn)

		done := true
		for i, chain := range chains {
			if len(addrs[i]) > 0 {
				continue
			}

			res, err := queryClient.QueryContract(goCtx, &tssTypes.QueryContractRequest{
				Chain: chain,
				Hash:  sisu.SupportedContracts[sisu.ContractErc20Gateway].AbiHash,
			})

			if err != nil || len(res.Contract.Address) == 0 {
				log.Verbose("we have not found gateway contract address for chain ", chain, " yet. Keep sleeping...")
				time.Sleep(time.Second * 3)
				done = false
				break
			}

			log.Info("GatewayContract Address = ", res.Contract.Address)
			addrs[i] = res.Contract.Address
		}

		if done {
			break
		}
	}

	gateways := make([]string, len(addrs))
	for i, addr := range addrs {
		gateways[i] = addr
	}

	return gateways
}

func (c *fundAccountCmd) setGateway(client *ethclient.Client, mnemonic string, liquidAddr, gateway common.Address) {
	liquidInstance, err := liquidity.NewLiquiditypool(liquidAddr, client)
	if err != nil {
		panic(err)
	}

	auth, err := getAuthTransactor(client, mnemonic)
	if err != nil {
		panic(err)
	}

	_, err = liquidInstance.SetGateway(auth, gateway)
	if err != nil {
		panic(err)
	}
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
