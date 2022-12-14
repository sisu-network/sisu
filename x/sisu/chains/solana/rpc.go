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
func QuerySolanaAccountBalance(clients []*rpc.Client, ataAccount string) (*big.Int, error) {
	result, err := clients[0].GetTokenAccountBalance(
		context.Background(),
		solanago.MustPublicKeyFromBase58(ataAccount),
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

func SignAndSubmit(clients []*rpc.Client, wsClients []*ws.Client,
	ixs []solanago.Instruction, feePayer solanago.PrivateKey) error {
	opts := rpc.TransactionOpts{
		SkipPreflight:       false,
		PreflightCommitment: rpc.CommitmentFinalized,
	}

	return SignAndSubmitWithOptions(clients, wsClients, ixs, feePayer, opts)
}

func SignAndSubmitWithOptions(clients []*rpc.Client, wsClients []*ws.Client,
	ixs []solanago.Instruction, feePayer solanago.PrivateKey, opts rpc.TransactionOpts) error {
	for i, client := range clients {
		sig, err := trySubmit(client, wsClients[i], ixs, feePayer, opts)
		if err == nil {
			log.Verbose("Final sig = ", sig)
			return nil
		}
	}

	return fmt.Errorf("Failed to submit Solana transaction.")
}

func trySubmit(client *rpc.Client, wsClient *ws.Client,
	ixs []solanago.Instruction, feePayer solanago.PrivateKey, opts rpc.TransactionOpts) (solanago.Signature, error) {
	// Get blockhash
	result, err := client.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return solanago.Signature{}, err
	}

	tx, err := solanago.NewTransaction(
		ixs,
		result.Value.Blockhash,
		solanago.TransactionPayer(feePayer.PublicKey()),
	)
	if err != nil {
		return solanago.Signature{}, err
	}

	tx.Sign(
		func(key solanago.PublicKey) *solanago.PrivateKey {
			if feePayer.PublicKey().Equals(key) {
				return &feePayer
			}

			return nil
		},
	)

	// Send transaction, and wait for confirmation
	sig, err := confirm.SendAndConfirmTransactionWithOpts(context.Background(), client, wsClient, tx,
		opts, nil)
	if err != nil {
		log.Verbose("Cannot send transaction, err = ", err)
		return solanago.Signature{}, err
	}

	return sig, nil
}
