package dev

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/contracts/eth/erc20gateway"
	"github.com/sisu-network/sisu/x/sisu"
	tssTypes "github.com/sisu-network/sisu/x/sisu/types"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/sisu-network/sisu/cmd/sisud/cmd/helper"
)

type swapCommand struct{}

func Swap() *cobra.Command {
	cmd := &cobra.Command{
		Use: "swap",
		Long: `Swap ERC20 token.
Usage:
./sisu dev swap --token SISU --amount 10 --recipient 0x2d532C099CA476780c7703610D807948ae47856A

for swapping token from chain ganache1 to ganache2. To swap tokens between 2 chains:

./sisu dev swap --src ganache2 --src-url http://0.0.0.0:8545 --dst ganache1 --token SISU --amount 10 --recipient 0x2d532C099CA476780c7703610D807948ae47856A

Please note that the amount is the number of whole unit. amount 1 is equivalent to 10^18 in the
transfer params.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			src, _ := cmd.Flags().GetString(flagSrc)
			srcUrl, _ := cmd.Flags().GetString(flagSrcUrl)
			dst, _ := cmd.Flags().GetString(flagDst)
			token, _ := cmd.Flags().GetString(flagToken)
			recipient, _ := cmd.Flags().GetString(flagRecipient)
			unit, _ := cmd.Flags().GetInt(flagAmount)

			c := &swapCommand{}

			if len(src) == 0 {
				panic("src chain cannot be empty")
			}

			if len(srcUrl) == 0 {
				srcUrl = getDefaultChainUrl(src)
			}

			log.Info("srcUrl = ", srcUrl)

			client, err := ethclient.Dial(srcUrl)
			if err != nil {
				log.Error("cannot dial source chain, url = ", srcUrl)
				panic(err)
			}
			defer client.Close()

			srcToken, dstToken := c.getTokenAddrs(token, src, dst)

			amount := big.NewInt(int64(unit))
			amount = amount.Exp(amount, big.NewInt(18), nil)

			gateway := c.getGatewayAddresses(cmd.Context(), src)
			c.swap(client, gateway, dst, srcToken, dstToken, recipient, amount)

			return nil
		},
	}

	cmd.Flags().String(flagSrc, "ganache1", "Source chain where the token is transferred from")
	cmd.Flags().String(flagSrcUrl, "", "Source chain url")
	cmd.Flags().String(flagDst, "ganache2", "Destination chain where the token is transferred to")
	cmd.Flags().String(flagToken, "SISU", "ID of the ERC20 to transferred")
	cmd.Flags().String(flagRecipient, "", "Recipient address in the destination chain")
	cmd.Flags().Int(flagAmount, 1, "The amount of token to be transferred")

	return cmd
}

func (c *swapCommand) getTokenAddrs(tokenId string, srcChain string, dstChain string) (string, string) {
	grpcConn, err := grpc.Dial(
		"0.0.0.0:9090",
		grpc.WithInsecure(),
	)
	defer grpcConn.Close()
	if err != nil {
		panic(err)
	}

	queryClient := tssTypes.NewTssQueryClient(grpcConn)
	res, err := queryClient.QueryToken(context.Background(), &tssTypes.QueryTokenRequest{
		Id: tokenId,
	})
	if err != nil {
		panic(err)
	}

	token := res.Token
	if len(token.Addresses[srcChain]) == 0 || len(token.Addresses[dstChain]) == 0 {
		panic(fmt.Errorf("cannot find token address, available token addresses = %v", token.Addresses))
	}

	return token.Addresses[srcChain], token.Addresses[dstChain]
}

func (c *swapCommand) getAuthTransactor(client *ethclient.Client, address common.Address) (*bind.TransactOpts, error) {
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	// This is the private key of the accounts0
	privateKey := helper.GetDevPrivateKey()

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasPrice = gasPrice

	auth.GasLimit = uint64(10_000_000)

	return auth, nil
}

func (c *swapCommand) getGatewayAddresses(context context.Context, chain string) string {
	grpcConn, err := grpc.Dial(
		"0.0.0.0:9090",
		grpc.WithInsecure(),
	)
	defer grpcConn.Close()
	if err != nil {
		panic(err)
	}

	queryClient := tssTypes.NewTssQueryClient(grpcConn)

	res, err := queryClient.QueryContract(context, &tssTypes.QueryContractRequest{
		Chain: chain,
		Hash:  sisu.SupportedContracts[sisu.ContractErc20Gateway].AbiHash,
	})

	if err != nil {
		panic(err)
	}

	if len(res.Contract.Address) == 0 {
		panic("gateway contract address is empty")
	}

	return res.Contract.Address
}

func (c *swapCommand) swap(client *ethclient.Client, gateay string, dstChain string,
	srcToken string, dstToken string, recipient string, amount *big.Int) {
	gatewayAddr := common.HexToAddress(gateay)
	contract, err := erc20gateway.NewErc20gateway(gatewayAddr, client)
	if err != nil {
		panic(err)
	}

	opts, err := c.getAuthTransactor(client, account0.Address)
	if err != nil {
		panic(err)
	}

	recipientAddr := common.HexToAddress(recipient)
	srcTokenAddr := common.HexToAddress(srcToken)
	dstTokenAddr := common.HexToAddress(dstToken)

	log.Verbose("destination, recipientAddr, srcTokenAddr, dstTokenAddr, amount = %s %s %s %s %s",
		dstChain, recipientAddr.String(), srcTokenAddr.String(), dstTokenAddr.String(), amount)

	tx, err := contract.TransferOut(opts, dstChain, recipientAddr, srcTokenAddr, dstTokenAddr, amount)
	bind.WaitDeployed(context.Background(), client, tx)

	time.Sleep(time.Second * 3)
}
