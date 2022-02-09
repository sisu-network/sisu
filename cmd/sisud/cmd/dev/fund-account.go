package dev

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"time"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/contracts/eth/erc20"
	"github.com/sisu-network/sisu/x/sisu"
	tssTypes "github.com/sisu-network/sisu/x/sisu/types"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	ExpectedErc20Address = "0x3A84fBbeFD21D6a5ce79D54d348344EE11EBd45C"
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

			// Deploy ERC20
			c := &fundAccountCmd{}
			for i := 0; i < len(args); i += 2 {
				if i == len(args)-1 {
					break
				}

				port, err := strconv.Atoi(args[i+1])
				if err != nil {
					return err
				}
				url := "http://0.0.0.0:" + strconv.Itoa(port)
				c.deployErc20(url, "sisu", "SISU")
			}

			// Fund the accounts
			allPubKeys := queryPubKeys(cmd)
			for i := 0; i < len(args); i += 2 {
				if i == len(args)-1 {
					break
				}

				// Get chain and local chain URL
				chain := args[i]
				pubKeyBytes := allPubKeys[libchain.KEY_TYPE_ECDSA]

				if pubKeyBytes == nil {
					return fmt.Errorf("cannot find pubkey for chain %s", chain)
				}

				pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
				addr := crypto.PubkeyToAddress(*pubKey).Hex()

				port, err := strconv.Atoi(args[i+1])
				if err != nil {
					return err
				}
				url := "http://0.0.0.0:" + strconv.Itoa(port)

				log.Info("Sending ETH To address ", addr, " of chain", chain)

				c.transferEth(url, addr, amount)
			}

			// Now we wait until all gateway contracts have been deployed.
			log.Info("Now we wait until gateway contracts are deployed on all chains.")
			c.waitForGatewayContract(cmd.Context())

			return nil
		},
	}

	return cmd
}

// deployErc20 deploys ERC20 contracts to dev chains
func (c *fundAccountCmd) deployErc20(url string, tokenName string, tokenSymbol string) string {
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Error("please check the ganache is up and running, url = ", url)
		panic(err)
	}

	auth, err := c.getAuthTransactor(client, common.HexToAddress("0xbeF23B2AC7857748fEA1f499BE8227c5fD07E70c"))
	if err != nil {
		panic(err)
	}

	if auth.Nonce.Cmp(big.NewInt(0)) != 0 {
		panic(fmt.Errorf("valid nonce, the account0 nonce should be zero. Please restart your ganache and try again."))
	}

	address, tx, instance, err := erc20.DeployErc20(auth, client, "sisu", "SISU")
	_ = instance
	if err != nil {
		panic(err)
	}

	if address.String() != ExpectedErc20Address {
		panic(fmt.Errorf("Invalid ERC20 address. You need to update the expected address (both in this file and the tokens_dev.json."))
	}

	log.Info("Deploying erc20....")
	bind.WaitDeployed(context.Background(), client, tx)
	time.Sleep(time.Second * 2)
	log.Info("Deployment done! Contract address: ", address.String())

	return address.String()
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
func (c *fundAccountCmd) transferEth(url, recipient string, amount int) {
	client, err := ethclient.Dial(url)
	if err != nil {
		panic(err)
	}

	account := c.getAccountAddress()

	log.Info("from Account.Address = ", account.String(), " recipient = ", recipient)

	nonce, err := client.PendingNonceAt(context.Background(), account)
	if err != nil {
		panic(err)
	}

	value := new(big.Int).Mul(big.NewInt(1000000000000000000), big.NewInt(int64(amount))) // in wei (10 eth)
	gasLimit := uint64(21000)                                                             // in units
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

func (c *fundAccountCmd) waitForGatewayContract(context context.Context) []string {
	var gateway1, gateway2 string

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

		if len(gateway1) == 0 {
			res, err := queryClient.QueryContract(context, &tssTypes.QueryContractRequest{
				Chain: "ganache1",
				Hash:  sisu.SupportedContracts[sisu.ContractErc20Gateway].AbiHash,
			})
			if err != nil || len(res.Contract.Address) == 0 {
				log.Verbose("we have not found contract address for gateway 1 yet. Keep sleeping...")
				time.Sleep(time.Second * 3)
				continue
			}
			gateway1 = res.Contract.Address
		}

		if len(gateway2) == 0 {
			res, err := queryClient.QueryContract(context, &tssTypes.QueryContractRequest{
				Chain: "ganache2",
				Hash:  sisu.SupportedContracts[sisu.ContractErc20Gateway].AbiHash,
			})
			if err != nil || len(res.Contract.Address) == 0 {
				log.Verbose("we have not found contract address for gateway 2 yet. Keep sleeping...")
				time.Sleep(time.Second * 3)
				continue
			}
			gateway2 = res.Contract.Address
		}

		break
	}

	log.Info("Gateway contract addresses = ", gateway1, " ", gateway2)

	return []string{gateway1, gateway2}
}
