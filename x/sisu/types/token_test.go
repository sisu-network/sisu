package types

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertAmountToSisuAmount(t *testing.T) {
	token := &Token{
		Chains: []string{
			"ganache1",
			"cardano",
			"solana-devnet",
		},
		Decimals: []byte{
			18, 6, 8,
		},
	}

	expected := big.NewInt(1_500_000_000_000_000_000)

	amount, err := token.ConvertAmountToSisuAmount("solana-devnet", big.NewInt(150_000_000))
	require.Nil(t, err)
	require.Equal(t, expected, amount)

	amount, err = token.ConvertAmountToSisuAmount("ganache1", big.NewInt(1_500_000_000_000_000_000))
	require.Nil(t, err)
	require.Equal(t, expected, amount)

	amount, err = token.ConvertAmountToSisuAmount("cardano", big.NewInt(1_500_000))
	require.Nil(t, err)
	require.Equal(t, expected, amount)

	amount, err = token.ConvertAmountToSisuAmount("unknow", big.NewInt(1_500_000))
	require.NotNil(t, err)
}
