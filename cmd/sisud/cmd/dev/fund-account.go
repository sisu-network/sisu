package dev

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
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
	"github.com/sisu-network/sisu/contracts/eth/erc20"
	"github.com/sisu-network/sisu/contracts/eth/liquidity"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu"
	tssTypes "github.com/sisu-network/sisu/x/sisu/types"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	ExpectedErc20Address = "0xf0D676183dD5ae6b370adDdbE770235F23546f9d"
	ExpectedLiquidPoolAddress = "0x3A84fBbeFD21D6a5ce79D54d348344EE11EBd45C"
	Erc20MethodTransfer  = "transfer"
)

type fundAccountCmd struct{}

func FundAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use: "fund-account",
		Short: `Fund localhost accounts. Example:
fund-account ganache1 7545 ganache2 8545 10
`,

		RunE: func(cmd *cobra.Command, args []string) error {
			amount, err := strconv.Atoi(args[len(args)-1])
			if err != nil {
				panic(err)
			}
			log.Verbose("Amount = ", amount)

			c := &fundAccountCmd{}
			wg := &sync.WaitGroup{}
			chains := make([]string, 0)
			clients := make([]*ethclient.Client, 0)

			// Get all urls from command arguments.
			urls := make([]string, 0)
			for i := 0; i < len(args)-1; i += 2 {
				chains = append(chains, args[i])
				port, err := strconv.Atoi(args[i+1])
				if err != nil {
					return err
				}
				url := "http://0.0.0.0:" + strconv.Itoa(port)
				urls = append(urls, url)
				client, err := ethclient.Dial(url)
				if err != nil {
					log.Error("please check the ganache is up and running, url = ", url)
					panic(err)
				}
				clients = append(clients, client)
			}
			defer func() {
				for _, client := range clients {
					client.Close()
				}
			}()

			// Deploy liquidity contract
			liquidityAddrs := make([]common.Address, len(urls))
			wg.Add(len(clients))
			for i, client := range clients {
				go func(i int, client *ethclient.Client) {
					liquidityAddrs[i] = c.deployLiquid(client)
					wg.Done()
				}(i, client)
			}
			wg.Wait()

			// Deploy ERC20 contract
			tokenAddrs := make([]common.Address, len(urls))
			wg.Add(len(clients))
			for i, client := range clients {
				go func(i int, client *ethclient.Client) {
					tokenAddrs[i] = c.deployErc20(client)
					wg.Done()
				}(i, client)
			}
			wg.Wait()

			// Waits for Sisu to create contract instance in its database. At this stage, the contract is
			// not deployed yet.
			c.waitForContractCreation(cmd.Context(), chains)

			// Fund the accounts
			allPubKeys := queryPubKeys(cmd)
			for _, client := range clients {
				// Get chain and local chain URL
				pubKeyBytes := allPubKeys[libchain.KEY_TYPE_ECDSA]
				pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
				if err != nil {
					panic(err)
				}
				addr := crypto.PubkeyToAddress(*pubKey).Hex()

				c.transferEth(client, addr, amount)
			}

			// Waits until all gateway contracts are deployed.
			log.Info("Now we wait until gateway contracts are deployed on all chains.")
			gateways := c.waitForGatewayDeployed(cmd.Context(), chains)

			// Grant permission for gateway to use liquidity pool' funds
			wg.Add(len(gateways))
			for i, client := range clients {
				go func(i int, client *ethclient.Client) {
					c.grantLiquidityPoolAccess(client, liquidityAddrs[i], gateways[i])
					wg.Done()
				}(i, client)
			}
			wg.Wait()

			// Approve the contract with some preallocated token
			wg.Add(len(gateways))
			for i, client := range clients {
				log.Verbose("gateway = ", gateways[i])
				go func(i int, client *ethclient.Client) {
					c.approveGateway(client, tokenAddrs[i], gateways[i])
					// c.approveGateway(client, common.HexToAddress("0x3A84fBbeFD21D6a5ce79D54d348344EE11EBd45C"), gateways[i])
					wg.Done()
				}(i, client)
			}
			wg.Wait()
			log.Info("Gateway approval done!")

			// Transfer ERC20 tokens
			erc20Amount := new(big.Int).Mul(big.NewInt(500), utils.EthToWei)
			for i, client := range clients {
				c.transferErc20Tokens(client, tokenAddrs[i], liquidityAddrs[i], erc20Amount)
			}

			time.Sleep(3 * time.Second)

			log.Info("Transfer ERC20 token done")

			// Query balance
			for i, client := range clients {
				balance, err := c.queryErc20Balance(client, tokenAddrs[i], liquidityAddrs[i])
				if err != nil {
					panic(err)
				}

				if balance.Cmp(erc20Amount) != 0 {
					panic(fmt.Sprintf("balance does not match: expected %s, actual %s", erc20Amount.String(), balance.String()))
				}
			}

			return nil
		},
	}

	return cmd
}

func (c *fundAccountCmd) waitForContractCreation(context context.Context, chains []string) []string {
	log.Info("Waiting for all contract created in Sisu's db")

	contractAddrs := make([]string, len(chains))
	for {
		grpcConn, err := grpc.Dial(
			"0.0.0.0:9090",
			grpc.WithInsecure(),
		)
		defer grpcConn.Close()
		if err != nil {
			panic(err)
		}

		queryClient := tssTypes.NewTssQueryClient(grpcConn)

		done := true
		for i, chain := range chains {
			res, err := queryClient.QueryContract(context, &tssTypes.QueryContractRequest{
				Chain: chain,
				Hash:  sisu.SupportedContracts[sisu.ContractErc20Gateway].AbiHash,
			})

			if err != nil {
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

func (c *fundAccountCmd) approveGateway(client *ethclient.Client, erc20Addr common.Address, addr common.Address) {
	log.Info("Approving gateway: ", addr)

	contract, err := erc20.NewErc20(erc20Addr, client)
	if err != nil {
		panic(err)
	}

	opts, err := c.getAuthTransactor(client, account0.Address)
	if err != nil {
		panic(err)
	}

	amount := new(big.Int).Mul(big.NewInt(500), utils.EthToWei)
	tx, err := contract.Approve(opts, addr, amount)
	bind.WaitDeployed(context.Background(), client, tx)
	time.Sleep(time.Second * 3)
}

func (c *fundAccountCmd) grantLiquidityPoolAccess(client *ethclient.Client, liquidityAddr, gatewayAddr common.Address) {
	log.Infof("Granting access for gatewayAddr to call liquidity pool, gateway address: %s\n", gatewayAddr.String())

	contract, err := liquidity.NewLiquidity(liquidityAddr, client)
	if err != nil {
		panic(err)
	}

	opts, err := c.getAuthTransactor(client, account0.Address)
	if err != nil {
		panic(err)
	}

	tx, err := contract.SetGateway(opts, gatewayAddr)
	if err != nil {
		panic(err)
	}
	txReceipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		panic(err)
	}

	if txReceipt.Status != ethtypes.ReceiptStatusSuccessful {
		panic("tx grant liquidity pool access failed")
	}

	log.Info("Grant access for gateway successfully")
}

func (c *fundAccountCmd) deployLiquid(client *ethclient.Client) common.Address {
	auth, err := c.getAuthTransactor(client, account0.Address)
	if err != nil {
		panic(err)
	}

	if auth.Nonce.Cmp(big.NewInt(0)) != 0 {
		panic("valid nonce, the account0 nonce should be zero. Please restart your ganache and try again.")
	}

	_, tx, _, err := liquidity.DeployLiquidity(auth, client, []common.Address{}, []common.Address{})
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

// deployErc20 deploys ERC20 contracts to dev chains
func (c *fundAccountCmd) deployErc20(client *ethclient.Client) common.Address {
	auth, err := c.getAuthTransactor(client, account0.Address)
	if err != nil {
		panic(err)
	}

	if auth.Nonce.Cmp(big.NewInt(1)) != 0 {
		panic("valid nonce, the account0 nonce should be 1. Please restart your ganache and try again.")
	}

	// seed 1000 * 10^18 for msg.sender
	address, tx, instance, err := erc20.DeployErc20(auth, client)
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

	return address
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

	auth := bind.NewKeyedTransactor(privateKey)
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

func queryPubKeys(cmd *cobra.Command) map[string][]byte {
	grpcConn, err := grpc.Dial(
		"0.0.0.0:9090",
		grpc.WithInsecure(),
	)
	defer grpcConn.Close()
	if err != nil {
		panic(err)
	}

	queryClient := tssTypes.NewTssQueryClient(grpcConn)

	res, err := queryClient.AllPubKeys(cmd.Context(), &tssTypes.QueryAllPubKeysRequest{})
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

func (c *fundAccountCmd) waitForGatewayDeployed(context context.Context, chains []string) []common.Address {
	addrs := make([]string, len(chains))

	for {
		grpcConn, err := grpc.Dial(
			"0.0.0.0:9090",
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

			res, err := queryClient.QueryContract(context, &tssTypes.QueryContractRequest{
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

	gateways := make([]common.Address, len(addrs))
	for i, addr := range addrs {
		gateways[i] = common.HexToAddress(addr)
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

// transferErc20Tokens transfer ERC20 from account0 to gateway address.
func (c *fundAccountCmd) transferErc20Tokens(
	client *ethclient.Client,
	tokenAddress common.Address,
	gatewateAddr common.Address,
	amount *big.Int,
) {
	opts, err := c.getAuthTransactor(client, account0.Address)
	if err != nil {
		panic(err)
	}

	store, err := erc20.NewErc20(tokenAddress, client)
	if err != nil {
		panic(err)
	}

	tx, err := store.Transfer(opts, gatewateAddr, amount)
	if err != nil {
		panic(err)
	}

	bind.WaitDeployed(context.Background(), client, tx)
	time.Sleep(time.Second * 3)
}

func (c *fundAccountCmd) queryErc20Balance(
	client *ethclient.Client,
	tokenAddr common.Address,
	target common.Address,
) (*big.Int, error) {
	store, err := erc20.NewErc20(tokenAddr, client)
	if err != nil {
		return nil, err
	}

	balance, err := store.BalanceOf(nil, target)

	return balance, err
}
