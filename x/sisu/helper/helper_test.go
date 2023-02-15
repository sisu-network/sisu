package helper

import (
	"math/big"
	"testing"

	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func TestGetChainGasCostInToken(t *testing.T) {
	ctx := testmock.TestContext()
	k := keeper.NewKeeper(testmock.TestKeyStore)

	nativeToken := "NATIVE_GANACHE1"
	token := "SISU"

	chain := "ganache1"
	k.SaveChain(ctx, &types.Chain{
		Id:          chain,
		NativeToken: nativeToken,
		EthConfig: &types.ChainEthConfig{
			GasPrice: 10 * 1_000_000_000,
		},
	})

	k.SetTokens(ctx, map[string]*types.Token{
		"NATIVE_GANACHE1": {
			Id:       nativeToken,
			Chains:   []string{"ganache1"},
			Decimals: []uint32{18},
		},
		"SISU": {
			Id:        token,
			Chains:    []string{"ganache1", "ganache2"},
			Decimals:  []uint32{18, 18},
			Addresses: []string{"", ""},
		},
	})

	mockDeyes := &external.MockDeyesClient{
		GetTokenPriceFunc: func(id string) (*big.Int, error) {
			if id == nativeToken {
				// native token price is 2 ETH
				return big.NewInt(utils.OneEtherInWei * 2), nil
			}

			// SISU token price is 2 ETH
			return big.NewInt(utils.OneEtherInWei * 4), nil
		},
	}
	gas := big.NewInt(8_000_000)
	amount, err := GetChainGasCostInToken(ctx, k, mockDeyes, token, chain,
		gas.Mul(gas, big.NewInt(10*1_000_000_000)))

	require.Nil(t, err)

	// amount = 0.008 * 10 * 2 / 4 ~ 0.04. Since 1 ETH = 10^18 wei, 0.04 ETH is 40_000_000_000_000_000 wei.
	require.Equal(t, big.NewInt(40_000_000_000_000_000), amount)
}

func TestGasCostInToken(t *testing.T) {
	gas := big.NewInt(8_000_000)
	gasCost := gas.Mul(gas, big.NewInt(10*1_000_000_000))
	tokenPrice := big.NewInt(utils.OneEtherInWei * 4)
	nativePriceToken := big.NewInt(utils.OneEtherInWei * 2)
	amount, err := GasCostInToken(gasCost, tokenPrice, nativePriceToken)

	require.Nil(t, err)

	// amount = 0.008 * 10 * 2 / 4 ~ 0.04. Since 1 ETH = 10^18 wei, 0.04 ETH is 40_000_000_000_000_000 wei.
	require.Equal(t, big.NewInt(40_000_000_000_000_000), amount)
}

func TestCheckRatioThreshol(t *testing.T) {
	t.Run("b_is_zero", func(t *testing.T) {
		a := new(big.Int).SetInt64(1)
		b := new(big.Int).SetInt64(0)
		_, ok := CheckRatioThreshold(a, b, 1)
		require.False(t, ok)
	})
	t.Run("a_is_three_times_of_b", func(t *testing.T) {
		a := new(big.Int).SetInt64(30)
		b := new(big.Int).SetInt64(10)
		_, ok := CheckRatioThreshold(a, b, 3)
		require.True(t, ok)
	})
	t.Run("a_is_larger_than_three_times_of_b", func(t *testing.T) {
		a := new(big.Int).SetInt64(31)
		b := new(big.Int).SetInt64(10)
		_, ok := CheckRatioThreshold(a, b, 3)
		require.False(t, ok)
	})
	t.Run("a_is_one_third_of_b", func(t *testing.T) {
		a := new(big.Int).SetInt64(10)
		b := new(big.Int).SetInt64(30)
		_, ok := CheckRatioThreshold(a, b, 3)
		require.True(t, ok)
	})
	t.Run("a_is_smaller_than_one_third_of_b", func(t *testing.T) {
		a := new(big.Int).SetInt64(9)
		b := new(big.Int).SetInt64(30)
		_, ok := CheckRatioThreshold(a, b, 3)
		require.False(t, ok)
	})
	t.Run("a_is_1.1_of_b", func(t *testing.T) {
		a := new(big.Int).SetInt64(110)
		b := new(big.Int).SetInt64(100)
		_, ok := CheckRatioThreshold(a, b, 1.1)
		require.True(t, ok)
	})
	t.Run("a_is_larger_than_1.1_of_b", func(t *testing.T) {
		a := new(big.Int).SetInt64(111)
		b := new(big.Int).SetInt64(100)
		_, ok := CheckRatioThreshold(a, b, 1.1)
		require.False(t, ok)
	})
	t.Run("a_is_0.9_of_b", func(t *testing.T) {
		a := new(big.Int).SetInt64(91)
		b := new(big.Int).SetInt64(100)
		_, ok := CheckRatioThreshold(a, b, 1.1)
		require.True(t, ok)
	})
	t.Run("a_is_smaller_than_0.9_of_b", func(t *testing.T) {
		a := new(big.Int).SetInt64(90)
		b := new(big.Int).SetInt64(100)
		_, ok := CheckRatioThreshold(a, b, 1.1)
		require.False(t, ok)
	})
}
