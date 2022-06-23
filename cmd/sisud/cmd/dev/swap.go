package dev

import (
	"context"
	"fmt"
	"math/big"
	"time"

	cgblockfrost "github.com/echovl/cardano-go/blockfrost"

	"github.com/echovl/cardano-go"
	"github.com/echovl/cardano-go/wallet"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	hutils "github.com/sisu-network/dheart/utils"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/contracts/eth/erc20gateway"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu"
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
			token, _ := cmd.Flags().GetString(flags.Erc20Symbol)
			recipient, _ := cmd.Flags().GetString(flags.Account)
			unit, _ := cmd.Flags().GetInt(flags.Amount)
			sisuRpc, _ := cmd.Flags().GetString(flags.SisuRpc)
			cardanoSecret, _ := cmd.Flags().GetString(flags.CardanoSecret)
			cardanoMnemonic, _ := cmd.Flags().GetString(flags.CardanoFunderMnemonic)

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

			srcToken, dstToken := c.getTokenAddrs(token, src, dst, sisuRpc)

			log.Verbosef("srcToken = %s, dstToken = %s", srcToken, dstToken)

			// Swapping from ETH chain
			if libchain.IsETHBasedChain(src) {
				gateway := c.getEthGatewayAddresses(cmd.Context(), src, sisuRpc)
				amount := big.NewInt(int64(unit))
				amount = new(big.Int).Mul(amount, utils.EthToWei)

				c.swapFromEth(client, mnemonic, gateway, dst, srcToken, dstToken, recipient, amount)
			} else if libchain.IsCardanoChain(src) {
				gateway := c.getCardanoGateway(cmd.Context(), sisuRpc)
				log.Info("Cardano gateway = ", gateway)

				amount := big.NewInt(int64(unit))
				amount = new(big.Int).Mul(amount, utils.ONE_ADA_IN_LOVELACE)

				c.swapFromCardano(dst, dstToken, recipient, gateway, amount, cardanoSecret, cardanoMnemonic)
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
	cmd.Flags().String(flags.CardanoSecret, "", "The blockfrost secret to interact with cardano network.")
	cmd.Flags().String(flags.CardanoFunderMnemonic, "", "Mnemonic of funder wallet which already has a lot of test tokens")

	return cmd
}

func (c *swapCommand) getTokenAddrs(tokenId string, srcChain string, dstChain string, sisuRpc string) (string, string) {
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

	// If destination chain is cardano, set destToken address to same as srcToken address
	if libchain.IsCardanoChain(dstChain) {
		dest = src
	}

	if len(src) == 0 && !libchain.IsCardanoChain(srcChain) {
		log.Info("source chain = ", srcChain)
		panic(fmt.Errorf("cannot find token address, available token addresses = %v", token.Addresses))
	}

	if len(dest) == 0 && !libchain.IsCardanoChain(dstChain) {
		log.Info("dest chain = ", dstChain)
		panic(fmt.Errorf("cannot find token address, available token addresses = %v", token.Addresses))
	}

	return src, dest
}

func (c *swapCommand) getEthGatewayAddresses(context context.Context, chain string, sisuRpc string) string {
	grpcConn, err := grpc.Dial(
		sisuRpc,
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

// swapFromEth creates an ETH transaction and sends to gateway contract.
func (c *swapCommand) swapFromEth(client *ethclient.Client, mnemonic string, gateway string, dstChain string,
	srcToken string, dstToken string, recipient string, amount *big.Int) {
	gatewayAddr := common.HexToAddress(gateway)
	contract, err := erc20gateway.NewErc20gateway(gatewayAddr, client)
	if err != nil {
		panic(err)
	}

	opts, err := getAuthTransactor(client, mnemonic)
	if err != nil {
		panic(err)
	}

	srcTokenAddr := common.HexToAddress(srcToken)
	dstTokenAddr := common.HexToAddress(dstToken)

	log.Verbosef("destination = %s, recipientAddr %s, srcTokenAddr = %s, dstTokenAddr = %s, amount = %s",
		dstChain, recipient, srcTokenAddr.String(), dstTokenAddr.String(), amount)

	tx, err := contract.TransferOut(opts, dstChain, recipient, srcTokenAddr, dstTokenAddr, amount)
	if err != nil {
		panic(err)
	}
	waitTx, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		panic(err)
	}

	log.Info("txHash = ", waitTx.TxHash.Hex())

	time.Sleep(time.Second * 3)
}

func (c *swapCommand) getCardanoGateway(ctx context.Context, sisuRpc string) string {
	allPubKeys := queryPubKeys(ctx, sisuRpc)
	cardanoKey, ok := allPubKeys[libchain.KEY_TYPE_EDDSA]
	if !ok {
		panic("can not find cardano pub key")
	}

	return hutils.GetAddressFromCardanoPubkey(cardanoKey).String()
}

func (c *swapCommand) swapFromCardano(destChain, destToken, destRecipient, cardanoGwAddr string,
	value *big.Int, blockfrostSecret, cardanoMnemonic string) {
	node := cgblockfrost.NewNode(cardano.Testnet, blockfrostSecret)
	opts := &wallet.Options{Node: node}
	client := wallet.NewClient(opts)

	wallet, err := client.RestoreWallet(DefaultCardanoWalletName, DefaultCardanoPassword, cardanoMnemonic)
	if err != nil {
		panic(err)
	}

	sender, err := wallet.AddAddress()
	if err != nil {
		panic(err)
	}
	log.Info("sender = ", sender)

	receiver, err := cardano.NewAddress(cardanoGwAddr)
	if err != nil {
		panic(err)
	}

	metadata := cardano.Metadata{
		0: map[string]interface{}{
			"chain":         destChain,
			"recipient":     destRecipient,
			"token_address": destToken,
		},
	}

	// TODO: use NewValueWithAssets when multiassets is ready
	hash, err := wallet.Transfer(receiver, cardano.NewValue(cardano.Coin(value.Uint64())), metadata)
	if err != nil {
		panic(err)
	}

	log.Info("Cardano tx hash = ", hash.String())
}
