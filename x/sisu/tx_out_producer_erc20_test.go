package sisu

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/helper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func TestTxOutProducerErc20_getGasCostInToken(t *testing.T) {
	ctx := testContext()
	k := keeperTestGenesis(ctx)
	deyesClient := &tssclients.MockDeyesClient{}

	worldState := defaultWorldStateTest(ctx, k, deyesClient)

	chain := "ganache1"
	token := &types.Token{
		Id:    "SISU",
		Price: int64(4 * utils.DecinmalUnit),
	}
	worldState.SetTokens(map[string]*types.Token{
		"SISU": token,
	})

	gas := big.NewInt(8_000_000)
	gasPrice := big.NewInt(10 * 1_000_000_000) // 10 gwei
	nativeTokenPrice, err := worldState.GetNativeTokenPriceForChain(chain)
	require.NoError(t, err)
	amount, err := helper.GetGasCostInToken(gas, gasPrice, big.NewInt(token.Price), big.NewInt(nativeTokenPrice))

	require.Equal(t, nil, err)

	// amount = 0.008 * 10 * 2 / 4 ~ 0.04. Since 1 ETH = 10^18 wei, 0.04 ETH is 40_000_000_000_000_000 wei.
	require.Equal(t, big.NewInt(40_000_000_000_000_000), amount)
}
