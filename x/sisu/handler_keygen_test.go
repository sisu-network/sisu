package sisu_test

import (
	"testing"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	libchain "github.com/sisu-network/lib/chain"
	mock "github.com/sisu-network/sisu/tests/mock/x/sisu"
	"github.com/sisu-network/sisu/x/sisu"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"

	mockcommon "github.com/sisu-network/sisu/tests/mock/common"
	mocktss "github.com/sisu-network/sisu/tests/mock/tss"
	mockkeeper "github.com/sisu-network/sisu/tests/mock/x/sisu/keeper"
	mocktssclients "github.com/sisu-network/sisu/tests/mock/x/sisu/tssclients"
)

func createManagerContainer_TestHandlerKeygen(ctrl *gomock.Controller) sisu.ManagerContainer {
	mockPublicDb := mocktss.NewMockStorage(ctrl)
	mockPublicDb.EXPECT().SaveKeygen(gomock.Any()).Times(1)
	mockPublicDb.EXPECT().ProcessTxRecord(gomock.Any()).Times(1)

	mockPmm := mock.NewMockPostedMessageManager(ctrl)
	mockPmm.EXPECT().ShouldProcessMsg(gomock.Any(), gomock.Any()).Return(true, []byte("")).Times(1)

	mockGlobalData := mockcommon.NewMockGlobalData(ctrl)
	mockGlobalData.EXPECT().IsCatchingUp().Return(false).Times(1)

	mockDheartClient := mocktssclients.NewMockDheartClient(ctrl)
	mockDheartClient.EXPECT().KeyGen(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockPartyManager := mocktss.NewMockPartyManager(ctrl)
	mockPartyManager.EXPECT().GetActivePartyPubkeys().Return([]ctypes.PubKey{}).Times(1)

	return sisu.MockManagerContainer(mockPmm, mockGlobalData, mockDheartClient, mockPartyManager, mockPublicDb)
}

func TestHandlerKeygen_normal(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	mockKeeper := mockkeeper.NewMockKeeper(ctrl)
	mockKeeper.EXPECT().SaveKeygen(gomock.Any(), gomock.Any()).Times(1)
	mockKeeper.EXPECT().ProcessTxRecord(gomock.Any(), gomock.Any()).Times(1)

	mockPmm := mock.NewMockPostedMessageManager(ctrl)
	mockPmm.EXPECT().ShouldProcessMsg(gomock.Any(), gomock.Any()).Return(true, []byte("")).Times(1)

	mockGlobalData := mockcommon.NewMockGlobalData(ctrl)
	mockGlobalData.EXPECT().IsCatchingUp().Return(false).Times(1)

	mockDheartClient := mocktssclients.NewMockDheartClient(ctrl)
	mockDheartClient.EXPECT().KeyGen(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockPartyManager := mocktss.NewMockPartyManager(ctrl)
	mockPartyManager.EXPECT().GetActivePartyPubkeys().Return([]ctypes.PubKey{}).Times(1)

	mc := sisu.MockManagerContainer(mockPmm, mockGlobalData, mockDheartClient, mockPartyManager, mockKeeper)

	msg := &types.KeygenWithSigner{
		Signer: "signer",
		Data: &types.Keygen{
			KeyType: libchain.KEY_TYPE_ECDSA,
			Index:   0,
		},
	}

	handler := sisu.NewHandlerKeygen(mc)
	_, err := handler.DeliverMsg(sdk.Context{}, msg)

	require.NoError(t, err)
}

func TestHandlerKeygen_CatchingUp(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	mockKeeper := mockkeeper.NewMockKeeper(ctrl)
	mockKeeper.EXPECT().SaveKeygen(gomock.Any(), gomock.Any()).Times(1)
	mockKeeper.EXPECT().ProcessTxRecord(gomock.Any(), gomock.Any()).Times(1)

	mockPmm := mock.NewMockPostedMessageManager(ctrl)
	mockPmm.EXPECT().ShouldProcessMsg(gomock.Any(), gomock.Any()).Return(true, []byte("")).Times(1)

	mockGlobalData := mockcommon.NewMockGlobalData(ctrl)
	mockGlobalData.EXPECT().IsCatchingUp().Return(true).Times(1) // We are catching up

	mockDheartClient := mocktssclients.NewMockDheartClient(ctrl)
	mockDheartClient.EXPECT().KeyGen(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(0) // We don't do keygen

	mc := sisu.MockManagerContainer(mockPmm, mockGlobalData, mockDheartClient, mockKeeper)

	msg := &types.KeygenWithSigner{
		Signer: "signer",
		Data: &types.Keygen{
			KeyType: libchain.KEY_TYPE_ECDSA,
			Index:   0,
		},
	}

	handler := sisu.NewHandlerKeygen(mc)
	_, err := handler.DeliverMsg(sdk.Context{}, msg)

	require.NoError(t, err)
}
