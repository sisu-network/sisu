package dev

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/contracts/eth/erc20"
	erc20Gateway "github.com/sisu-network/sisu/contracts/eth/erc20gateway"
	"github.com/sisu-network/sisu/x/tss"
	"github.com/sisu-network/sisu/x/tss/types"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// WIP. TODO: Complete and clean up this.
func TransferOut() *cobra.Command {
	cmd := &cobra.Command{
		Use: "transfer-out",
		Long: `Transfer an ERC20 or ERC721 asset.
Usage:
transfer-out [ContractType] [FromChain] [Port] [TokenAddress] [ToChain] [RecipientAddress]

Example:
transfer-out erc20 ganache1 7545 0xf0D676183dD5ae6b370adDdbE770235F23546f9d ganache2 0xE8382821BD8a0F9380D88e2c5c33bc89Df17E466
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get the contract address of token
			log.Info("args = ", args)
			contractType := args[0]
			fromChain := args[1]
			port, err := strconv.Atoi(args[2])
			if err != nil {
				panic(err)
			}

			tokenAddressString := args[3]
			toChain := args[4]
			recipient := args[5]

			switch contractType {
			case "erc20":
				client, err := getEthClient(port)
				if err != nil {
					panic(err)
				}

				hash := tss.SupportedContracts[contractType].AbiHash
				contract := queryContract(cmd, fromChain, hash)
				if contract == nil {
					panic(fmt.Errorf("cannot find contract"))
				}

				// fmt.Println("Contract is not nil, contract hash = ", contract.Hash, contract.Address)
				contract.Address = "0x9d64dc6c7c9e6df3c08b345be8859ead38154b9f"

				// return nil

				gatewayAddress := common.HexToAddress(contract.Address)
				gateway, err := erc20Gateway.NewErc20gateway(gatewayAddress, client)
				if err != nil {
					panic(err)
				}

				log.Info("gatewayAddress = ", gatewayAddress.String())

				tokenAddress := common.HexToAddress(tokenAddressString)
				erc20Contract, err := erc20.NewErc20(tokenAddress, client)
				if err != nil {
					panic(err)
				}

				log.Info("Approving gateway address...")
				amount := big.NewInt(1)
				approveAddress(erc20Contract, gatewayAddress, amount, client)

				// Check the allowance
				allowance, err := erc20Contract.Allowance(&bind.CallOpts{Pending: true}, account0.Address, gatewayAddress)
				if err != nil {
					panic(err)
				}
				if allowance.Cmp(amount) != 0 {
					panic(fmt.Errorf("Invalid balance: expected %s, actual %s", amount, allowance))
				}

				log.Info("Transferring token out....")
				auth, err := getAuthTransactor(client, account0.Address)
				if err != nil {
					panic(err)
				}

				tx, err := gateway.TransferOutFromContract(auth, tokenAddress, toChain, recipient, amount)
				if err != nil {
					panic(err)
				}
				bind.WaitDeployed(context.Background(), client, tx)

				time.Sleep(Blocktime)

				gatewayBalance := getBalance(erc20Contract, gatewayAddress)
				log.Info("gatewayBalance = ", gatewayBalance)
			}

			return nil
		},
	}
	return cmd
}

func queryContract(cmd *cobra.Command, chain string, hash string) *types.Contract {
	grpcConn, err := grpc.Dial(
		"0.0.0.0:9090",
		grpc.WithInsecure(),
	)
	defer grpcConn.Close()
	if err != nil {
		panic(err)
	}

	queryClient := types.NewQueryClient(grpcConn)
	req := &types.QueryContractRequest{
		Chain: chain,
		Hash:  hash,
	}
	res, err := queryClient.QueryContract(cmd.Context(), req)
	if err != nil {
		panic(err)
	}

	return res.Contract
}

func approveAddress(erc20Contract *erc20.Erc20, recipient common.Address, amount *big.Int, client *ethclient.Client) {
	auth, err := getAuthTransactor(client, account0.Address)
	if err != nil {
		panic(err)
	}

	_, err = erc20Contract.Approve(auth, recipient, amount)
	if err != nil {
		panic(err)
	}

	time.Sleep(Blocktime)
}

func getBalance(erc20Contract *erc20.Erc20, address common.Address) *big.Int {
	tokenBalance, err := erc20Contract.BalanceOf(&bind.CallOpts{Pending: true}, address)
	if err != nil {
		panic(err)
	}

	return tokenBalance
}
