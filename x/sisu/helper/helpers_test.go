package helper

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCardanoTxFeeInToken(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			name                              string
			adaPrice, tokenPrice, adaForTxFee *big.Int
			expect                            *big.Int
		}{
			{
				name:        "success",
				adaPrice:    big.NewInt(1),
				tokenPrice:  big.NewInt(2),
				adaForTxFee: big.NewInt(1_300_000),
				expect:      big.NewInt(650_000),
			},
			{
				name:        "zero all",
				adaPrice:    big.NewInt(0),
				tokenPrice:  big.NewInt(0),
				adaForTxFee: big.NewInt(0),
				expect:      big.NewInt(0),
			},
			{
				name:        "zero token price",
				adaPrice:    big.NewInt(1),
				tokenPrice:  big.NewInt(0),
				adaForTxFee: big.NewInt(2),
				expect:      big.NewInt(0),
			},
		}

		for _, tc := range testCases {
			require.Equal(t, tc.expect.Uint64(), GetCardanoTxFeeInToken(tc.adaPrice, tc.tokenPrice, tc.adaForTxFee).Uint64())
		}
	})
}
