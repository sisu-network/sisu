package sisu

import (
	"testing"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForHandlerGasPrice() (sdk.Context, ManagerContainer) {
	ctx := testmock.TestContext()
	k := testmock.KeeperTestGenesis(ctx)
	k.SaveParams(ctx, &types.Params{MajorityThreshold: 1})

	globalData := &common.MockGlobalData{}
	pmm := NewPostedMessageManager(k)

	partyManager := &MockPartyManager{}
	partyManager.GetActivePartyPubkeysFunc = func() []ctypes.PubKey {
		return []ctypes.PubKey{}
	}

	dheartClient := &external.MockDheartClient{}
	appKeys := common.NewMockAppKeys()
	mc := MockManagerContainer(k, pmm, globalData, partyManager, dheartClient, appKeys)

	return ctx, mc
}

func TestHandlerGasPrice(t *testing.T) {
	t.Run("set_gas_price_successfully", func(t *testing.T) {
		ctx, mc := mockForHandlerGasPrice()

		chains := []string{"ganache1", "ganache2"}
		prices := []int64{1, 2}
		signer := mc.AppKeys().GetSignerAddress().String()
		msg := types.NewGasPriceMsg(signer, chains, prices, []int64{0, 0}, []int64{0, 0})

		handler := NewHandlerGasPrice(mc)
		_, err := handler.DeliverMsg(ctx, msg)
		require.NoError(t, err)

		pricesMap := mc.Keeper().GetGasPrices(ctx, "ganache1")
		require.Equal(t, 1, len(pricesMap))
		require.Equal(t, int64(1), pricesMap[signer].GasPrice)

		pricesMap = mc.Keeper().GetGasPrices(ctx, "ganache")
		require.Equal(t, 1, len(pricesMap))
		require.Equal(t, int64(2), pricesMap[signer].GasPrice)
	})

	t.Run("multiple_signers_set_gas_price_successfully", func(t *testing.T) {
		ctx, mc := mockForHandlerGasPrice()

		chains := []string{"ganache1", "ganache2"}
		prices1 := []int64{1, 10}
		prices2 := []int64{2, 11}
		prices3 := []int64{1, 9}

		signer1 := mc.AppKeys().GetSignerAddress().String()
		signer2 := "cosmos1zf2ssujzp6y577gzwn457tnxy7yj44yq37t05z"
		signer3 := "cosmos1g64vzyutdjfdvw5kyae73fc39sksg3r7gzmrzy"
		// require.NoError(t, err)

		msg1 := types.NewGasPriceMsg(signer1, chains, prices1, []int64{0, 0}, []int64{0, 0})
		msg2 := types.NewGasPriceMsg(signer2, chains, prices2, []int64{0, 0}, []int64{0, 0})
		msg3 := types.NewGasPriceMsg(signer3, chains, prices3, []int64{0, 0}, []int64{0, 0})

		handler := NewHandlerGasPrice(mc)
		_, err := handler.DeliverMsg(ctx, msg1)
		require.NoError(t, err)
		_, err = handler.DeliverMsg(ctx, msg2)
		require.NoError(t, err)
		_, err = handler.DeliverMsg(ctx, msg3)
		require.NoError(t, err)

		pricesMap := mc.Keeper().GetGasPrices(ctx, "ganache1")
		require.Equal(t, 3, len(pricesMap))
		require.Equal(t, int64(1), pricesMap[signer1].GasPrice)
		require.Equal(t, int64(2), pricesMap[signer2].GasPrice)
		require.Equal(t, int64(1), pricesMap[signer3].GasPrice)

		pricesMap = mc.Keeper().GetGasPrices(ctx, "ganache2")
		require.Equal(t, 3, len(pricesMap))
		require.Equal(t, int64(10), pricesMap[signer1].GasPrice)
		require.Equal(t, int64(11), pricesMap[signer2].GasPrice)
		require.Equal(t, int64(9), pricesMap[signer3].GasPrice)
	})
}
