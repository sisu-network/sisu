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

			liquidityAddrString, _ := cmd.Flags().GetString(flags.LiquidityAddrs)

			amount, _ := cmd.Flags().GetInt(flags.Amount)
			sisuRpc, _ := cmd.Flags().GetString(flags.SisuRpc)
			log.Info("Amount = ", amount)

			c := &fundAccountCmd{}
			c.fundSisuAccounts(cmd.Context(), chainString, urlString, mnemonic, liquidityAddrString, sisuRpc, amount)

			return nil
		},
	}

	cmd.Flags().String(flags.Chains, "ganache1,ganache2", "Names of all chains we want to fund.")
	cmd.Flags().String(flags.ChainUrls, "http://0.0.0.0:7545,http://0.0.0.0:8545", "RPCs of all the chains we want to fund.")
	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum", "Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.SisuRpc, "0.0.0.0:9090", "URL to connect to Sisu. Please do NOT include http:// prefix")
	cmd.Flags().String(flags.LiquidityAddrs, fmt.Sprintf("%s,%s", ExpectedLiquidPoolAddress, ExpectedLiquidPoolAddress), "List of liquidity pool addresses")

	cmd.Flags().Int(flags.Amount, 100, "The amount that gateway addresses will receive")

	return cmd
}

func (c *fundAccountCmd) fundSisuAccounts(ctx context.Context, chainString, urlString, mnemonic, liquidityAddrString, sisuRpc string, amount int) {
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

			c.transferEth(client, mnemonic, tssPubAddr.Hex(), amount)
		}(i, client)
	}
	wg.Wait()

	// Waits until all gateway contracts are deployed.
	log.Info("Now we wait until gateway contracts are deployed on all chains.")
	gateways := c.waitForGatewayDeployed(ctx, chains, sisuRpc)

	// Set gateway for the liquidity
	log.Info("Set gateway address for liqiuitidy pool")
	wg.Add(len(gateways))
	for i, client := range clients {
		go func(i int, client *ethclient.Client) {
			c.setGatewayForLiquidity(client, mnemonic, liquidityAddrs[i], gateways[i])
			wg.Done()
		}(i, client)
	}
	wg.Wait()
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
func (c *fundAccountCmd) transferEth(client *ethclient.Client, mnemonic, recipient string, amount int) {
	_, account := getPrivateKey(mnemonic)

	log.Info("from address = ", account.String(), " to Address = ", recipient)

	nonce, err := client.PendingNonceAt(context.Background(), account)
	if err != nil {
		panic(err)
	}

	// 0.1 ETH
	value := new(big.Int).Mul(utils.EthToWei, big.NewInt(int64(amount)))
	value = new(big.Int).Div(value, big.NewInt(10))
	gasLimit := uint64(21000) // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}

	toAddress := common.HexToAddress(recipient)
	var data []byte
	tx := ethtypes.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	privateKey, _ := getPrivateKey(mnemonic)
	signedTx, err := ethtypes.SignTx(tx, getSigner(client), privateKey)

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

func (c *fundAccountCmd) transferLiquidityOwnership(
	client *ethclient.Client, mnemonic string, liquidAddr, newOwner string) error {
	liquidInstance, err := liquidity.NewLiquiditypool(common.HexToAddress(liquidAddr), client)
	if err != nil {
		return err
	}

	auth, err := getAuthTransactor(client, mnemonic)
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
