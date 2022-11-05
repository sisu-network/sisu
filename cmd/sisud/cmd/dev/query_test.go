package dev

import (
	"testing"

	"github.com/sisu-network/lib/log"
	"github.com/stretchr/testify/require"
)

// Integration test for sanity checking. Disabled when running on test system.
func TestQuerySolanaBalance(t *testing.T) {
	t.Skip()

	c := new(queryCommand)

	amount, err := c.querySolanaAccountBalance("http://127.0.0.1:8899", "BJ9ArHvbeUhVLChS2yksw8xqvoRpWYLtGkg7CVHNa31a")
	require.Nil(t, err)

	log.Verbose("amount = ", amount)
}
