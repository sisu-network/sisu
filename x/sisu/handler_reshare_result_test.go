package sisu

import (
	"encoding/base64"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sisu-network/sisu/common"
	mock "github.com/sisu-network/sisu/tests/mock/x/sisu"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func TestHandlerReshareResult(t *testing.T) {
	t.Parallel()

	t.Run("Reshare result success", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		ctx := testContext()
		keeper := keeperTestGenesis(ctx)
		validatorManager := NewValidatorManager(keeper)
		appKeys := common.NewMockAppKeys()
		mockPmm := mock.NewMockPostedMessageManager(ctrl)
		mockPmm.EXPECT().ShouldProcessMsg(gomock.Any(), gomock.Any()).Return(true, []byte("mock_hash")).Times(1)
		mockPmm.EXPECT().IsReachedThreshold(gomock.Any(), gomock.Any(), 2).Return(true, []byte("")).Times(1)

		mc := MockManagerContainer(validatorManager, keeper, mockPmm)
		handler := NewHandlerReshareResult(mc)

		newVal1, err := base64.StdEncoding.DecodeString("1jPHjoWahm5WDES2ud3zJbzmRzCPLFacQsrl/pbO/Wo=")
		require.NoError(t, err)
		newVal2, err := base64.StdEncoding.DecodeString("qsXeJ51BGalR2V2Zz9ugh3ofsIS58Kjya9pDgfKH018=")
		require.NoError(t, err)
		newValSet := [][]byte{newVal1, newVal2}
		reshareMsg := types.NewReshareResultWithSigner(appKeys.GetSignerAddress().String(), newValSet, types.ReshareData_SUCCESS)

		_, err = handler.DeliverMsg(ctx, reshareMsg)
		require.NoError(t, err)

		incomingValidateUpdates := keeper.GetIncomingValidatorUpdates(ctx)
		require.NotEmpty(t, incomingValidateUpdates)
	})
}
