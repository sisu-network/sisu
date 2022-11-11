package solana

import (
	"context"
	"fmt"
	"math/big"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/sisu-network/lib/log"
)

// querySolanaAccountBalance queries token balance on a solana chain
func QuerySolanaAccountBalance(client *rpc.Client, tokenAddr string) (*big.Int, error) {
	result, err := client.GetTokenAccountBalance(
		context.Background(),
		solanago.MustPublicKeyFromBase58(tokenAddr),
		rpc.CommitmentConfirmed,
	)

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("querySolanaAccountBalance: result is nil")
	}

	amount, ok := new(big.Int).SetString(result.Value.Amount, 10)
	if !ok {
		return nil, fmt.Errorf("Invalid returned amount %s", result.Value.Amount)
	}

	return amount, err
}

func SignAndSubmit(client *rpc.Client, wsClient *ws.Client,
	ixs []solanago.Instruction, feePayer solanago.PrivateKey) error {
	// Get blockhash
	result, err := client.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}

	tx, err := solanago.NewTransaction(
		ixs,
		result.Value.Blockhash,
		solanago.TransactionPayer(feePayer.PublicKey()),
	)
	if err != nil {
		panic(err)
	}

	tx.Sign(
		func(key solanago.PublicKey) *solanago.PrivateKey {
			if feePayer.PublicKey().Equals(key) {
				return &feePayer
			}

			return nil
		},
	)
	log.Verbose("tx sig = ", tx.Signatures[0])

	// Send transaction, and wait for confirmation
	sig, err := confirm.SendAndConfirmTransaction(
		context.Background(),
		client,
		wsClient,
		tx,
	)
	log.Verbose("sig = ", sig)

	return err
}
