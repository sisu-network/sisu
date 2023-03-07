package dev

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/gogo/protobuf/proto"
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

	lisktypes "github.com/sisu-network/deyes/chains/lisk/types"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/helper"
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
			dst, _ := cmd.Flags().GetString(flags.Dst)
			tokenSymbol, _ := cmd.Flags().GetString(flags.Erc20Symbol)
			recipient, _ := cmd.Flags().GetString(flags.Recipient)
			amount, _ := cmd.Flags().GetInt(flags.Amount)
			sisuRpc, _ := cmd.Flags().GetString(flags.SisuRpc)
			deyesUrl, _ := cmd.Flags().GetString(flags.DeyesUrl)
			genesisFolder, _ := cmd.Flags().GetString(flags.GenesisFolder)

			if len(recipient) == 0 {
				panic(flags.Recipient + " cannot be empty")
			}

			token, srcToken, dstToken := getTokenAddrsFromSisu(tokenSymbol, src, dst, sisuRpc)

			log.Verbosef("srcToken = %s, dstToken = %s", srcToken, dstToken)

			// Swapping from ETH chain
			if libchain.IsETHBasedChain(src) {
				clients := getEthClients([]string{src}, genesisFolder)
				if len(clients) == 0 {
					panic("There is no healthy client")
				}

				client := clients[0]

				vault := getEthVaultAddress(cmd.Context(), src, sisuRpc)
				log.Info("Vault address = ", vault)
				amountBigInt := big.NewInt(int64(amount))
				amountBigInt = new(big.Int).Mul(amountBigInt, utils.EthToWei)

				swapFromEth(client, mnemonic, vault, dst, srcToken, dstToken, recipient, amountBigInt)
			} else if libchain.IsCardanoChain(src) {
				vault := getCardanoVault(cmd.Context(), sisuRpc)
				log.Info("Cardano gateway = ", vault)

				amountBigInt := big.NewInt(int64(amount))
				amountBigInt = new(big.Int).Mul(amountBigInt, big.NewInt(utils.OneAdaInLoveLace))

				cardanoconfig := helper.ReadCardanoConfig(genesisFolder)

				swapFromCardano(src, dst, token, recipient, vault, amountBigInt, src,
					cardanoconfig.Secret, mnemonic, deyesUrl)

			} else if libchain.IsSolanaChain(src) {
				swapFromSolana(genesisFolder, src, mnemonic, srcToken, recipient,
					libchain.GetChainIntFromId(dst).Uint64(), uint64(amount*100_000_000))

			} else if libchain.IsLiskChain(src) {
				allPubKeys := queryPubKeys(cmd.Context(), sisuRpc)

				swapFromLisk(genesisFolder, mnemonic, dst, allPubKeys[libchain.KEY_TYPE_EDDSA], recipient,
					uint64(amount*100_000_000))
			}

			return nil
		},
	}

	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum", "Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.Src, "ganache1", "Source chain where the token is transferred from")
	cmd.Flags().String(flags.SisuRpc, "0.0.0.0:9090", "URL to connect to Sisu. Please do NOT include http:// prefix")
	cmd.Flags().String(flags.Dst, "ganache2", "Destination chain where the token is transferred to")
	cmd.Flags().String(flags.Erc20Symbol, "SISU", "ID of the ERC20 to transferred")
	cmd.Flags().String(flags.Recipient, "", "Recipient address in the destination chain")
	cmd.Flags().Int(flags.Amount, 1, "The amount of token to be transferred")
	cmd.Flags().String(flags.DeyesUrl, "http://127.0.0.1:31001", "Url to deyes api server.")
	cmd.Flags().String(flags.GenesisFolder, "./misc/dev", "Genesis folder that contains configuration files.")

	return cmd
}

func getTokenAddrsFromSisu(tokenId string, srcChain string, dstChain string, sisuRpc string) (*types.Token, string, string) {
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

	if len(src) == 0 && !libchain.IsCardanoChain(srcChain) && !libchain.IsLiskChain(srcChain) {
		log.Info("source chain = ", srcChain)
		panic(fmt.Errorf("cannot find token address, available token addresses = %v", token.Addresses))
	}

	if len(dest) == 0 && !libchain.IsCardanoChain(dstChain) && !libchain.IsLiskChain(dstChain) {
		log.Info("dest chain = ", dstChain)
		panic(fmt.Errorf("cannot find token address, available token addresses = %v", token.Addresses))
	}

	return token, src, dest
}

// swapFromEth creates an ETH transaction and sends to gateway contract.
func swapFromEth(client *ethclient.Client, mnemonic string, vaultAddr string, dstChain string,
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

func swapFromLisk(genesisFolder, mnemonic string, toChain string, mpcPubKey []byte,
	recipient string, amount uint64) {
	toChainId := libchain.GetChainIntFromId(toChain).Uint64()
	if toChainId == 0 {
		panic(fmt.Errorf("Invalid chain %s", toChain))
	}

	var recipientBytes []byte
	if libchain.IsETHBasedChain(toChain) {
		var err error
		recipientBytes, err = hex.DecodeString(recipient[2:])
		if err != nil {
			panic(err)
		}
	} else {
		panic(fmt.Errorf("Unsupported chain %s", toChain))
	}

	payload := lisktypes.TransferData{
		ChainId:   &toChainId,
		Recipient: recipientBytes,
		Amount:    &amount,
	}

	bz, err := proto.Marshal(&payload)
	if err != nil {
		panic(err)
	}

	msg := base64.StdEncoding.EncodeToString(bz)
	transferLisk(genesisFolder, mnemonic, mpcPubKey, 100_000, msg)
}
