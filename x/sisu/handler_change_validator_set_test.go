package sisu

import (
	"encoding/base64"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sisu-network/sisu/common"
	m "github.com/sisu-network/sisu/tests/mock/common"
	mock "github.com/sisu-network/sisu/tests/mock/x/sisu"
	mocktssclients "github.com/sisu-network/sisu/tests/mock/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func TestHandlerChangeValidatorSet(t *testing.T) {
	t.Parallel()

	t.Run("Send reshare request to dheart", func(t *testing.T) {
		t.Parallel()

		ctx := testContext()
		keeper := keeperTestGenesis(ctx)
		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		appKeys := common.NewMockAppKeys()
		global := m.NewMockGlobalData(ctrl)
		global.EXPECT().IsCatchingUp().Return(false)

		mockDheartClient := mocktssclients.NewMockDheartClient(ctrl)
		mockDheartClient.EXPECT().Reshare(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		mockValidatorManager := mock.NewMockValidatorManager(ctrl)
		mockValidatorManager.EXPECT().HasConsensus(gomock.Any(), gomock.Any()).Return(true).Times(1)

		pmm := NewPostedMessageManager(keeper, mockValidatorManager)
		mc := MockManagerContainer(keeper, mockDheartClient, global, mockValidatorManager, pmm)

		oldVal, err := base64.StdEncoding.DecodeString("1jPHjoWahm5WDES2ud3zJbzmRzCPLFacQsrl/pbO/Wo=")
		require.NoError(t, err)
		newVal, err := base64.StdEncoding.DecodeString("qsXeJ51BGalR2V2Zz9ugh3ofsIS58Kjya9pDgfKH018=")
		require.NoError(t, err)
		msg := types.NewChangeValidatorSetMsg(appKeys.GetSignerAddress().String(), [][]byte{oldVal}, [][]byte{newVal})
		handler := NewHandlerChangeValidatorSet(mc)
		_, err = handler.DeliverMsg(ctx, msg)
		require.NoError(t, err)
	})
}
