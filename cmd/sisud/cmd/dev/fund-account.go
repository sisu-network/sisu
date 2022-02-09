package dev

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	hdwallet "github.com/sisu-network/sisu/utils/hdwallet"
	tssTypes "github.com/sisu-network/sisu/x/sisu/types"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func FundAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use: "fund-account",
		Short: `Fund localhost accounts. Example:
fund-account ganache1 7545 ganache2 8545 10
`,

		RunE: func(cmd *cobra.Command, args []string) error {
			// Get all the pubkey
			allPubKeys := queryPubKeys(cmd)

			localWallet, err := hdwallet.NewFromMnemonic(defaultMnemonic)

			amount, err := strconv.Atoi(args[len(args)-1])
			if err != nil {
				panic(err)
			}

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

				log.Info("Sending ETH To address", addr, "of chain", chain)

				transferEth(url, addr, localWallet, amount)
			}

			return nil
		},
	}

	return cmd
}

// transferEth transfers a specific ETH amount to an address.
func transferEth(url, recipient string, wallet *hdwallet.Wallet, amount int) {
	client, err := ethclient.Dial(url)
	if err != nil {
		panic(err)
	}
	path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", 0))
	fromAccount, err := wallet.Derive(path, true)
	if err != nil {
		panic(err)
	}

	log.Info("from Account.Address = ", fromAccount.Address, " recipient = ", recipient)

	nonce, err := client.PendingNonceAt(context.Background(), fromAccount.Address)
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
	tx := ethTypes.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		panic(err)
	}

	signedTx, err := wallet.SignTx(fromAccount, tx, chainID)
	if err != nil {
		panic(err)
	}

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
