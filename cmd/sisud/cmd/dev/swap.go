package dev

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/sisu-network/sisu/contracts/eth/vault"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"
	tssTypes "github.com/sisu-network/sisu/x/sisu/types"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
)

type swapCommand struct{}

func Swap() *cobra.Command {
	cmd := &cobra.Command{
		Use: "swap",
		Long: `Swap ERC20 token.
Usage:
./sisu dev swap --amount 10 --account 0x2d532C099CA476780c7703610D807948ae47856A

for swapping token from chain ganache1 to ganache2.

Full command swap tokens between 2 chains:

./sisu dev swap --src ganache1 --src-url http://0.0.0.0:7545 --dst ganache2 --erc20-symbol SISU --amount 10 --account 0x2d532C099CA476780c7703610D807948ae47856A

Please note that the amount is the number of whole unit. amount 1 is equivalent to 10^18 in the
transfer params.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			mnemonic, _ := cmd.Flags().GetString(flags.Mnemonic)
			src, _ := cmd.Flags().GetString(flags.Src)
			srcUrl, _ := cmd.Flags().GetString(flags.SrcUrl)
			dst, _ := cmd.Flags().GetString(flags.Dst)
			tokenSymbol, _ := cmd.Flags().GetString(flags.Erc20Symbol)
			recipient, _ := cmd.Flags().GetString(flags.Account)
			amount, _ := cmd.Flags().GetInt(flags.Amount)
			sisuRpc, _ := cmd.Flags().GetString(flags.SisuRpc)
			cardanoChain, _ := cmd.Flags().GetString(flags.CardanoChain)
			cardanoMnemonic, _ := cmd.Flags().GetString(flags.CardanoMnemonic)
			cardanoSecret, _ := cmd.Flags().GetString(flags.CardanoSecret)
			deyesUrl, _ := cmd.Flags().GetString(flags.DeyesApiUrl)
			if len(cardanoMnemonic) == 0 {
				cardanoMnemonic = mnemonic
			}

			c := &swapCommand{}

			if len(recipient) == 0 {
				panic(flags.Account + " cannot be empty")
			}

			log.Info("srcUrl = ", srcUrl)

			client, err := ethclient.Dial(srcUrl)
			if err != nil {
				log.Error("cannot dial source chain, url = ", srcUrl)
				panic(err)
			}
			defer client.Close()

			token, srcToken, dstToken := c.getTokenAddrs(tokenSymbol, src, dst, sisuRpc)

			log.Verbosef("srcToken = %s, dstToken = %s", srcToken, dstToken)

			// Swapping from ETH chain
			if libchain.IsETHBasedChain(src) {
				vault := c.getVaultAddress(cmd.Context(), src, sisuRpc)
				log.Info("Vault address = ", vault)
				amountBigInt := big.NewInt(int64(amount))
				amountBigInt = new(big.Int).Mul(amountBigInt, utils.EthToWei)

				c.swapFromEth(client, mnemonic, vault, dst, srcToken, dstToken, recipient, amountBigInt)
			} else if libchain.IsCardanoChain(src) {
				vault := c.getCardanoVault(cmd.Context(), sisuRpc)
				log.Info("Cardano gateway = ", vault)

				amountBigInt := big.NewInt(int64(amount))
				amountBigInt = new(big.Int).Mul(amountBigInt, utils.ONE_ADA_IN_LOVELACE)

				c.swapFromCardano(src, dst, token, recipient, vault, amountBigInt, cardanoChain,
					cardanoSecret, cardanoMnemonic, deyesUrl)
			}

			return nil
		},
	}

	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum", "Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.Src, "ganache1", "Source chain where the token is transferred from")
	cmd.Flags().String(flags.SrcUrl, "http://127.0.0.1:7545", "Source chain url")
	cmd.Flags().String(flags.SisuRpc, "0.0.0.0:9090", "URL to connect to Sisu. Please do NOT include http:// prefix")
	cmd.Flags().String(flags.Dst, "ganache2", "Destination chain where the token is transferred to")
	cmd.Flags().String(flags.Erc20Symbol, "SISU", "ID of the ERC20 to transferred")
	cmd.Flags().String(flags.Account, "", "Recipient address in the destination chain")
	cmd.Flags().Int(flags.Amount, 1, "The amount of token to be transferred")
	cmd.Flags().String(flags.DeyesApiUrl, "http://127.0.0.1:31001", "Url to deyes api server.")
	cmd.Flags().String(flags.CardanoChain, "", "Cardano chain.")
	cmd.Flags().String(flags.CardanoMnemonic, "", "Cardano mnemonic.")
	cmd.Flags().String(flags.CardanoSecret, "", "The blockfrost secret to interact with cardano network.")

	return cmd
}

func (c *swapCommand) getTokenAddrs(tokenId string, srcChain string, dstChain string, sisuRpc string) (*types.Token, string, string) {
	grpcConn, err := grpc.Dial(
		sisuRpc,
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
	src, dest := token.GetAddressForChain(srcChain), token.GetAddressForChain(dstChain)

	if len(src) == 0 && !libchain.IsCardanoChain(srcChain) {
		log.Info("source chain = ", srcChain)
		panic(fmt.Errorf("cannot find token address, available token addresses = %v", token.Addresses))
	}

	if len(dest) == 0 && !libchain.IsCardanoChain(dstChain) {
		log.Info("dest chain = ", dstChain)
		panic(fmt.Errorf("cannot find token address, available token addresses = %v", token.Addresses))
	}

	return token, src, dest
}

func (c *swapCommand) getVaultAddress(context context.Context, chain string, sisuRpc string) string {
	grpcConn, err := grpc.Dial(
		sisuRpc,
		grpc.WithInsecure(),
	)
	defer grpcConn.Close()
	if err != nil {
		panic(err)
	}

	queryClient := tssTypes.NewTssQueryClient(grpcConn)
	res, err := queryClient.QueryVault(context, &tssTypes.QueryVaultRequest{
		Chain: chain,
	})

	if err != nil {
		panic(err)
	}

	if len(res.Vault.Address) == 0 {
		panic("gateway contract address is empty")
	}

	return res.Vault.Address
}

// swapFromEth creates an ETH transaction and sends to gateway contract.
func (c *swapCommand) swapFromEth(client *ethclient.Client, mnemonic string, vaultAddr string, dstChain string,
	srcToken string, dstToken string, recipient string, amount *big.Int) {
	v := common.HexToAddress(vaultAddr)
	contract, err := vault.NewVault(v, client)
	if err != nil {
		panic(err)
	}

	opts, err := getAuthTransactor(client, mnemonic)
	if err != nil {
		panic(err)
	}

	srcTokenAddr := common.HexToAddress(srcToken)

	log.Verbosef("destination = %s, recipientAddr %s, srcTokenAddr = %s, amount = %s",
		dstChain, recipient, srcTokenAddr.String(), amount)

	var tx *ethtypes.Transaction
	if libchain.IsETHBasedChain(dstChain) {
		recipientAddr := common.HexToAddress(recipient)
		tx, err = contract.TransferOut(opts, srcTokenAddr, libchain.GetChainIntFromId(dstChain),
			recipientAddr, amount)
		if err != nil {
			panic(err)
		}
	} else {
		tx, err = contract.TransferOutNonEvm(opts, srcTokenAddr, libchain.GetChainIntFromId(dstChain),
			recipient, amount)
		if err != nil {
			panic(err)
		}
	}

	waitTx, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		panic(err)
	}

	log.Info("txHash = ", waitTx.TxHash.Hex())

	time.Sleep(time.Second * 3)
}
