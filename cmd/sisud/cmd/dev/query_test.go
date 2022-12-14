package dev

import (
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/chains/solana"
	"github.com/stretchr/testify/require"
)

// Integration test for sanity checking. Disabled when running on test system.
func TestQuerySolanaBalance(t *testing.T) {
	t.Skip()

	endpoint := rpc.LocalNet_RPC
	client := rpc.New(endpoint)

	amount, err := solana.QuerySolanaAccountBalance([]*rpc.Client{client},
		"BJ9ArHvbeUhVLChS2yksw8xqvoRpWYLtGkg7CVHNa31a")
	require.Nil(t, err)

	log.Verbose("amount = ", amount)
}
