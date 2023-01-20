package utils

import (
	"github.com/stretchr/testify/require"

	"math/big"
	"testing"
)

func TestLovelaceToETHTokens(t *testing.T) {
	t.Run("two ada to tokens", func(t *testing.T) {
		// 2 ADA = 2*10^6 lovelace
		twoAda := big.NewInt(2_000_000)

		tokens := LovelaceToWei(twoAda)
		require.Equal(t, big.NewInt(2_000_000_000_000_000_000), tokens)
	})

	t.Run("0.5 ada to tokens", func(t *testing.T) {
		halfAda := big.NewInt(500_000)
		tokens := LovelaceToWei(halfAda)
		require.Equal(t, big.NewInt(500_000_000_000_000_000), tokens)
	})

	t.Run("0 ada to tokens", func(t *testing.T) {
		zero := big.NewInt(0)
		require.Equal(t, zero, LovelaceToWei(zero))
	})
}

func TestETHTokensToLovelace(t *testing.T) {
	t.Run("10^18 tokens to 1 ada", func(t *testing.T) {
		n := big.NewInt(1_000_000_000_000_000_000)

		require.Equal(t, big.NewInt(OneAdaInLoveLace), WeiToLovelace(n))
	})
}

func TestConvertInt(t *testing.T) {
	x := uint32(8)
	bz := Uint32ToBytes(x)
	y := BytesToUint32(bz)
	require.Equal(t, x, y)

	x = uint32(100)
	bz = Uint32ToBytes(x)
	y = BytesToUint32(bz)
	require.Equal(t, x, y)
}
