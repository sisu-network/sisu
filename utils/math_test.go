package utils

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubtractCommissionRate(t *testing.T) {
	amount := big.NewInt(2000)
	output := SubtractCommissionRate(amount, 10) // 0.1%
	require.Equal(t, big.NewInt(1998), output)
}
