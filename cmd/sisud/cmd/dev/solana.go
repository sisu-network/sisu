package dev

import (
	"context"
	"fmt"
	"math/big"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// querySolanaAccountBalance queries token balance on a solana chain
func querySolanaAccountBalance(client *rpc.Client, tokenAddr string) (*big.Int, error) {
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
