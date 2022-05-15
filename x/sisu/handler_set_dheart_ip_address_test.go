package sisu

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func TestHandlerSetDheartIPAddress(t *testing.T) {
	t.Parallel()

	t.Run("Set DHeart IP address successfully", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		ctx := testContext()
		keeper := keeperTestGenesis(ctx)
		appKeys := common.NewMockAppKeys()

		mc := MockManagerContainer(keeper)
		handler := NewHandlerSetDheartIPAddress(mc)

		msg := types.NewSetDheartIPAddressMsg(appKeys.GetSignerAddress().String(), "192.168.1.1")
		_, err := handler.DeliverMsg(ctx, msg)
		require.NoError(t, err)

		ips := keeper.GetAllDheartIPAddresses(ctx)
		require.Len(t, ips, 1)
		require.Equal(t, appKeys.GetSignerAddress(), ips[0].Addr)
		require.Equal(t, "192.168.1.1", ips[0].IP)
	})
}
