package utils

import (
	"github.com/stretchr/testify/require"

	"math/big"
	"testing"
)

func TestLovelaceToETHTokens(t *testing.T) {
	t.Parallel()

	t.Run("two ada to tokens", func(t *testing.T) {
		t.Parallel()
		// 2 ADA = 2*10^6 lovelace
		twoAda := big.NewInt(2_000_000)

		tokens := LovelaceToETHTokens(twoAda)
		require.Equal(t, big.NewInt(2_000_000_000_000_000_000), tokens)
	})

	t.Run("0.5 ada to tokens", func(t *testing.T) {
		t.Parallel()

		halfAda := big.NewInt(500_000)
		tokens := LovelaceToETHTokens(halfAda)
		require.Equal(t, big.NewInt(500_000_000_000_000_000), tokens)
	})

	t.Run("0 ada to tokens", func(t *testing.T) {
		t.Parallel()

		zero := big.NewInt(0)
		require.Equal(t, zero, LovelaceToETHTokens(zero))
	})
}

func TestETHTokensToLovelace(t *testing.T) {
	t.Parallel()

	t.Run("10^18 tokens to 1 ada", func(t *testing.T) {
		n := big.NewInt(1_000_000_000_000_000_000)

		require.Equal(t, ONE_ADA_IN_LOVELACE, ETHTokensToLovelace(n))
	})
}
