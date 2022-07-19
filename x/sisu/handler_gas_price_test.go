package sisu

import (
	"testing"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForHandlerGasPrice() (sdk.Context, ManagerContainer) {
	ctx := testContext()
	k := keeperTestGenesis(ctx)
	k.SaveParams(ctx, &types.Params{MajorityThreshold: 1})

	globalData := &common.MockGlobalData{}
	pmm := NewPostedMessageManager(k)

	partyManager := &MockPartyManager{}
	partyManager.GetActivePartyPubkeysFunc = func() []ctypes.PubKey {
		return []ctypes.PubKey{}
	}

	dheartClient := &tssclients.MockDheartClient{}
	appKeys := common.NewMockAppKeys()
	mc := MockManagerContainer(k, pmm, globalData, partyManager, dheartClient, appKeys)

	return ctx, mc
}

func TestHandlerGasPrice(t *testing.T) {
	t.Parallel()

	t.Run("set_gas_price_successfully", func(t *testing.T) {
		ctx, mc := mockForHandlerGasPrice()

		chains := []string{"ETH", "BSC", "POLYGON"}
		prices := []int64{1, 2, 3}
		signer := mc.AppKeys().GetSignerAddress().String()
		msg := types.NewGasPriceMsg(signer, chains, 100, prices)

		handler := NewHandlerGasPrice(mc)
		_, err := handler.DeliverMsg(ctx, msg)
		require.NoError(t, err)

		eth := mc.Keeper().GetChain(ctx, "ETH")
		require.Equal(t, int64(1), eth.GasPrice)

		bsc := mc.Keeper().GetChain(ctx, "BSC")
		require.Equal(t, int64(2), bsc.GasPrice)

		polygon := mc.Keeper().GetChain(ctx, "POLYGON")
		require.Equal(t, int64(3), polygon.GasPrice)
	})

	t.Run("multiple_signers_set_gas_price_successfully", func(t *testing.T) {
		t.Parallel()

		ctx, mc := mockForHandlerGasPrice()

		chains := []string{"ETH", "BSC", "POLYGON"}
		prices1 := []int64{1, 10, 20}
		prices2 := []int64{2, 11, 19}
		prices3 := []int64{1, 9, 21}

		signer1 := mc.AppKeys().GetSignerAddress().String()
		signer2, err := sdk.AccAddressFromBech32("cosmos1zf2ssujzp6y577gzwn457tnxy7yj44yq37t05z")
		require.NoError(t, err)
		signer3, err := sdk.AccAddressFromBech32("cosmos1g64vzyutdjfdvw5kyae73fc39sksg3r7gzmrzy")
		require.NoError(t, err)

		msg1 := types.NewGasPriceMsg(signer1, chains, 100, prices1)
		msg2 := types.NewGasPriceMsg(signer2.String(), chains, 100, prices2)
		msg3 := types.NewGasPriceMsg(signer3.String(), chains, 100, prices3)

		handler := NewHandlerGasPrice(mc)
		_, err = handler.DeliverMsg(ctx, msg1)
		require.NoError(t, err)
		_, err = handler.DeliverMsg(ctx, msg2)
		require.NoError(t, err)
		_, err = handler.DeliverMsg(ctx, msg3)
		require.NoError(t, err)

		eth := mc.Keeper().GetChain(ctx, "ETH")
		require.Equal(t, int64(1), eth.GasPrice) // median of [1, 2, 1]

		bsc := mc.Keeper().GetChain(ctx, "BSC")
		require.Equal(t, int64(10), bsc.GasPrice) // median of [10, 11, 9]

		polygon := mc.Keeper().GetChain(ctx, "POLYGON")
		require.Equal(t, int64(20), polygon.GasPrice) // median of [20, 19, 21]
	})
}
