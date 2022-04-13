package dev

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	libchain "github.com/sisu-network/lib/chain"

	"github.com/ethereum/go-ethereum/crypto"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/sisu-network/sisu/contracts/eth/erc20"
	liquidity "github.com/sisu-network/sisu/contracts/eth/liquiditypool"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu"
	tssTypes "github.com/sisu-network/sisu/x/sisu/types"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	ExpectedErc20Address      = "0x3A84fBbeFD21D6a5ce79D54d348344EE11EBd45C"
	ExpectedLiquidPoolAddress = "0xf0D676183dD5ae6b370adDdbE770235F23546f9d"
)

var (
	ContributionAmount = new(big.Int).Mul(big.NewInt(500), utils.EthToWei)
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
			erc20AddrString, _ := cmd.Flags().GetString(flags.Erc20Addrs)
			liquidityAddrString, _ := cmd.Flags().GetString(flags.LiquidityAddrs)

			amount, _ := cmd.Flags().GetInt(flags.Amount)
			sisuRpc, _ := cmd.Flags().GetString(flags.SisuRpc)
			log.Info("Amount = ", amount)

			c := &fundAccountCmd{}
			c.fundSisuAccounts(cmd.Context(), chainString, urlString, erc20AddrString, liquidityAddrString, sisuRpc, amount)

			return nil
		},
	}

	cmd.Flags().String(flags.Chains, "ganache1,ganache2", "Names of all chains we want to fund.")
	cmd.Flags().String(flags.ChainUrls, "http://0.0.0.0:7545,http://0.0.0.0:8545", "RPCs of all the chains we want to fund.")
	cmd.Flags().String(flags.SisuRpc, "0.0.0.0:9090", "URL to connect to Sisu. Please do NOT include http:// prefix")
	cmd.Flags().String(flags.Erc20Addrs, fmt.Sprintf("%s,%s", ExpectedErc20Address, ExpectedErc20Address), "List of erc20 addresses")
	cmd.Flags().String(flags.LiquidityAddrs, fmt.Sprintf("%s,%s", ExpectedLiquidPoolAddress, ExpectedLiquidPoolAddress), "List of liquidity pool addresses")

	cmd.Flags().Int(flags.Amount, 100, "The amount that gateway addresses will receive")

	return cmd
}

func (c *fundAccountCmd) fundSisuAccounts(ctx context.Context, chainString, urlString, erc20AddrString, liquidityAddrString, sisuRpc string, amount int) {
	chains := strings.Split(chainString, ",")
	urls := strings.Split(urlString, ",")
	tokenAddrs := strings.Split(erc20AddrString, ",")
	liquidityAddrs := strings.Split(liquidityAddrString, ",")

	log.Info("tokenAddrs = ", tokenAddrs)
	log.Info("liquidityAddrs = ", liquidityAddrs)

	wg := &sync.WaitGroup{}
	clients := make([]*ethclient.Client, 0)

	// Get all urls from command arguments.
	for i := 0; i < len(urls); i++ {
		client, err := ethclient.Dial(urls[i])
		if err != nil {
			log.Error("please check chain is up and running, url = ", urls[i])
			panic(err)
		}
		clients = append(clients, client)
	}
	defer func() {
		for _, client := range clients {
			client.Close()
		}
	}()

	// Approve the contract with some preallocated token from account0
	wg.Add(len(clients))
	for i, client := range clients {
		go func(i int, client *ethclient.Client) {
			c.approveAddress(client, tokenAddrs[i], liquidityAddrs[i])
			wg.Done()
		}(i, client)
	}
	wg.Wait()
	log.Info("Liquidity approval done!")

	time.Sleep(time.Second * 3)

	// Add liquidity to the pool
	wg.Add(len(clients))
	for i, client := range clients {
		go func(i int, client *ethclient.Client) {
			defer wg.Done()

			balance, err := c.queryErc20Balance(client, tokenAddrs[i], liquidityAddrs[i])
			if err != nil {
				panic(err)
			}

			if balance.Cmp(big.NewInt(0)) == 0 {
				log.Infof("Adding liquidity of token %s to the pool at %s", tokenAddrs[i], liquidityAddrs[i])
				c.addLiquidity(client, liquidityAddrs[i], tokenAddrs[i])
			} else {
				log.Infof("Liquidity pool has received %s tokens (%s) \n", balance.String(), tokenAddrs[i])
			}
		}(i, client)
	}
	wg.Wait()

	// Waits for Sisu to create contract instance in its database. At this stage, the contract is
	// not deployed yet.
	c.waitForGatewayCreationInSisuDb(ctx, chains, sisuRpc)

	time.Sleep(time.Second * 3)

	// Fund the accounts with some native ETH
	allPubKeys := queryPubKeys(ctx, sisuRpc)
	var tssPubAddr common.Address
	wg.Add(len(clients))
	for i, client := range clients {
		go func(i int, client *ethclient.Client) {
			defer wg.Done()

			// Get chain and local chain URL
			pubKeyBytes := allPubKeys[libchain.KEY_TYPE_ECDSA]
			pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
			if err != nil {
				panic(err)
			}
			tssPubAddr = crypto.PubkeyToAddress(*pubKey)

			c.transferEth(client, tssPubAddr.Hex(), amount)
		}(i, client)
	}
	wg.Wait()

	// Waits until all gateway contracts are deployed.
	log.Info("Now we wait until gateway contracts are deployed on all chains.")
	gateways := c.waitForGatewayDeployed(ctx, chains, sisuRpc)

	// Approve the contract with some preallocated token from account0
	wg.Add(len(clients))
	for i, client := range clients {
		go func(i int, client *ethclient.Client) {
			c.approveAddress(client, tokenAddrs[i], gateways[i])
			wg.Done()
		}(i, client)
	}
	wg.Wait()
	log.Info("Gateway approval done!")

	// Grant permission for gateway to use liquidity pool' funds
	log.Info("Set gateway address for liqiuitidy pool")
	wg.Add(len(gateways))
	for i, client := range clients {
		go func(i int, client *ethclient.Client) {
			c.setGatewayForLiquidity(client, liquidityAddrs[i], gateways[i])
			wg.Done()
		}(i, client)
	}
	wg.Wait()

	// Transfer ownership of liquidity pool to TSS public address
	log.Info("Transferring ownership of liquidity pool for ", tssPubAddr.Hex())
	wg.Add(len(clients))
	for i, client := range clients {
		go func(i int, client *ethclient.Client) {
			if err := c.transferLiquidityOwnership(client, liquidityAddrs[i], tssPubAddr.String()); err != nil {
				panic(err)
			}

			wg.Done()
		}(i, client)
	}
	wg.Wait()
	log.Info("Transferred ownership of liquidity pool to tss public address")

	c.doSanityCheck(clients, tokenAddrs, liquidityAddrs, gateways)
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

func (c *fundAccountCmd) approveAddress(client *ethclient.Client, erc20Addr string, target string) {
	contract, err := erc20.NewErc20(common.HexToAddress(erc20Addr), client)
	if err != nil {
		panic(err)
	}

	opts, err := c.getAuthTransactor(client, account0.Address)
	if err != nil {
		panic(err)
	}

	tx, err := contract.Approve(opts, common.HexToAddress(target), ContributionAmount)
	bind.WaitDeployed(context.Background(), client, tx)
	time.Sleep(time.Second * 3)
}

func (c *fundAccountCmd) setGatewayForLiquidity(client *ethclient.Client, liquidityAddr, gatewayAddr string) {
	log.Infof("Setting gateway for liquidity pool, gateway address: %s, liquidity pool address %s\n", gatewayAddr, liquidityAddr)

	contract, err := liquidity.NewLiquiditypool(common.HexToAddress(liquidityAddr), client)
	if err != nil {
		panic(err)
	}

	opts, err := c.getAuthTransactor(client, account0.Address)
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

func (c *fundAccountCmd) deployLiquid(client *ethclient.Client, tokens []common.Address, names []string) common.Address {
	auth, err := c.getAuthTransactor(client, account0.Address)
	if err != nil {
		panic(err)
	}

	if auth.Nonce.Cmp(big.NewInt(1)) != 0 {
		panic("invalid nonce, the account0 nonce should be zero. Please restart your ganache and try again.")
	}

	_, tx, _, err := liquidity.DeployLiquiditypool(auth, client, tokens, names)
	if err != nil {
		panic(err)
	}

	log.Info("Deploying liquidity ... ")
	addr, err := bind.WaitDeployed(context.Background(), client, tx)
	if err != nil {
		panic(err)
	}

	if addr.String() != ExpectedLiquidPoolAddress {
		panic(fmt.Errorf(`Unmatched Liquid pool address. We expect address %s but get %s.
You need to update the expected address (both in this file and the tokens_dev.json).`,
			ExpectedLiquidPoolAddress, addr.String()))
	}

	log.Info("Deployed liquidity successfully, addr: ", addr.String())

	return addr
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

// deployErc20 deploys ERC20 contracts to dev chains
func (c *fundAccountCmd) deployErc20(client *ethclient.Client) (common.Address, string) {
	auth, err := c.getAuthTransactor(client, account0.Address)
	if err != nil {
		panic(err)
	}

	// seed 1000 * 10^18 for msg.sender
	address, tx, instance, err := erc20.DeployErc20(auth, client, "Sisu Token", "SISU")
	_ = instance
	if err != nil {
		panic(err)
	}

	if address.String() != ExpectedErc20Address {
		panic(fmt.Errorf(`Unmatched ERC20 address. We expect address %s but get %s.
You need to update the expected address (both in this file and the tokens_dev.json.`,
			ExpectedErc20Address, address.String()))
	}

	log.Info("Deploying erc20....")
	bind.WaitDeployed(context.Background(), client, tx)
	time.Sleep(time.Second * 3)
	log.Info("Deployment done! ERC20 Contract address: ", address.String())

	return address, "Sisu"
}

// getAuthTransactor returns transaction opts for creating transaction object.
func (c *fundAccountCmd) getAuthTransactor(client *ethclient.Client, address common.Address) (*bind.TransactOpts, error) {
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	// This is the private key of the accounts0
	privateKey := c.getPrivateKey()

	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasPrice = gasPrice

	auth.GasLimit = uint64(10_000_000)

	return auth, nil
}

// transferEth transfers a specific ETH amount to an address.
func (c *fundAccountCmd) transferEth(client *ethclient.Client, recipient string, amount int) {
	account := c.getAccountAddress()

	log.Info("from address = ", account.String(), " to Address = ", recipient)

	nonce, err := client.PendingNonceAt(context.Background(), account)
	if err != nil {
		panic(err)
	}

	value := new(big.Int).Mul(utils.EthToWei, big.NewInt(int64(amount)))
	gasLimit := uint64(21000) // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}

	toAddress := common.HexToAddress(recipient)
	var data []byte
	tx := ethtypes.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	privateKey := c.getPrivateKey()
	signedTx, err := ethtypes.SignTx(tx, ethtypes.HomesteadSigner{}, privateKey)

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

func (c *fundAccountCmd) getPrivateKey() *ecdsa.PrivateKey {
	// This is the private key for account 0xbeF23B2AC7857748fEA1f499BE8227c5fD07E70c
	bz, err := hex.DecodeString("9f575b88940d452da46a6ceec06a108fcd5863885524aec7fb0bc4906eb63ab1")
	if err != nil {
		panic(err)
	}

	privateKey, err := ethcrypto.ToECDSA(bz)
	if err != nil {
		panic(err)
	}

	return privateKey
}

func (c *fundAccountCmd) getAccountAddress() common.Address {
	privateKey := c.getPrivateKey()
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("error casting public key to ECDSA")
	}

	return ethcrypto.PubkeyToAddress(*publicKeyECDSA)
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

func (c *fundAccountCmd) getNonceAndGas(client *ethclient.Client, addr common.Address) (uint64, uint64, *big.Int, error) {
	nonce, err := client.PendingNonceAt(context.Background(), account0.Address)
	if err != nil {
		log.Error("failed to get nonce, err = ", err)
		return 0, 0, nil, err
	}

	gasLimit := uint64(100000) // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Error("failed to get gas price, err = ", err)
		return 0, 0, nil, err
	}

	return nonce, gasLimit, gasPrice, nil
}

func (c *fundAccountCmd) addLiquidity(client *ethclient.Client, liquidAddr, tokenAddress string) {
	liquidInstance, err := liquidity.NewLiquiditypool(common.HexToAddress(liquidAddr), client)
	if err != nil {
		panic(err)
	}

	auth, err := c.getAuthTransactor(client, account0.Address)
	if err != nil {
		panic(err)
	}

	_, err = liquidInstance.AddLiquidity(auth, common.HexToAddress(tokenAddress), ContributionAmount)
	if err != nil {
		panic(err)
	}
}

func (c *fundAccountCmd) setGateway(client *ethclient.Client, liquidAddr, gateway common.Address) {
	liquidInstance, err := liquidity.NewLiquiditypool(liquidAddr, client)
	if err != nil {
		panic(err)
	}

	auth, err := c.getAuthTransactor(client, account0.Address)
	if err != nil {
		panic(err)
	}

	_, err = liquidInstance.SetGateway(auth, gateway)
	if err != nil {
		panic(err)
	}
}

func (c *fundAccountCmd) queryErc20Balance(
	client *ethclient.Client,
	tokenAddr string,
	target string,
) (*big.Int, error) {
	store, err := erc20.NewErc20(common.HexToAddress(tokenAddr), client)
	if err != nil {
		return nil, err
	}

	balance, err := store.BalanceOf(nil, common.HexToAddress(target))

	return balance, err
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

func (c *fundAccountCmd) transferLiquidityOwnership(
	client *ethclient.Client, liquidAddr, newOwner string) error {
	liquidInstance, err := liquidity.NewLiquiditypool(common.HexToAddress(liquidAddr), client)
	if err != nil {
		return err
	}

	auth, err := c.getAuthTransactor(client, account0.Address)
	if err != nil {
		panic(err)
	}

	tx, err := liquidInstance.TransferOwnership(auth, common.HexToAddress(newOwner))
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		return err
	}

	return nil
}

func (c *fundAccountCmd) doSanityCheck(clients []*ethclient.Client, tokenAddrs, liquidityAddrs, gateways []string) {
	// Query balance
	for i, client := range clients {
		balance, err := c.queryErc20Balance(client, tokenAddrs[i], liquidityAddrs[i])
		if err != nil {
			panic(err)
		}

		if balance.Cmp(ContributionAmount) != 0 {
			panic(fmt.Sprintf("balance does not match: expected %s, actual %s", ContributionAmount.String(), balance.String()))
		}

		// Check allowance
		allowance := c.queryAllownace(client, tokenAddrs[i], account0.Address.String(), gateways[i])
		if allowance.Cmp(ContributionAmount) != 0 {
			panic(fmt.Sprintf("Allowance to gateway should not be 0"))
		}
	}
}
