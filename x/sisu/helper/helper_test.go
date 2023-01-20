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

func TestGasCostInToken(t *testing.T) {
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

	require.Equal(t, nil, err)

	// amount = 0.008 * 10 * 2 / 4 ~ 0.04. Since 1 ETH = 10^18 wei, 0.04 ETH is 40_000_000_000_000_000 wei.
	require.Equal(t, big.NewInt(40_000_000_000_000_000), amount)
}
