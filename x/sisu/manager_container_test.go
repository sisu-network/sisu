package sisu

import (
	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/golang/mock/gomock"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	mockcommon "github.com/sisu-network/sisu/tests/mock/common"
	mocktss "github.com/sisu-network/sisu/tests/mock/tss"
	mock "github.com/sisu-network/sisu/tests/mock/x/sisu"
	mocktssclients "github.com/sisu-network/sisu/tests/mock/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
)

func MockManagerContainer(args ...interface{}) ManagerContainer {
	mc := &DefaultManagerContainer{}

	for _, arg := range args {
		switch t := arg.(type) {
		case PostedMessageManager:
			mc.pmm = t
		case keeper.Storage:
			mc.publicDb = t
		case common.GlobalData:
			mc.globalData = t
		case tssclients.DeyesClient:
			mc.deyesClient = t
		case tssclients.DheartClient:
			mc.dheartClient = t
		case config.TssConfig:
			mc.config = t
		case common.TxSubmit:
			mc.txSubmit = t
		case common.AppKeys:
			mc.appKeys = t
		case PartyManager:
			mc.partyManager = t
		case TxTracker:
			mc.txTracker = t
		}
	}

	return mc
}

func CreateManagerContainer(ctrl *gomock.Controller, args ...interface{}) ManagerContainer {
	var (
		mockPmm          *mock.MockPostedMessageManager
		mockPublicDb     *mocktss.MockStorage
		mockGlobalData   *mockcommon.MockGlobalData
		mockDheartClient *mocktssclients.MockDheartClient
		mockPartyManager *mocktss.MockPartyManager
	)

	for _, arg := range args {
		switch arg.(type) {
		case *mock.MockPostedMessageManager:
			mockPmm = arg.(*mock.MockPostedMessageManager)
		case *mocktss.MockStorage:
			mockPublicDb = arg.(*mocktss.MockStorage)
		case *mockcommon.MockGlobalData:
			mockGlobalData = arg.(*mockcommon.MockGlobalData)
		case *mocktssclients.MockDheartClient:
			mockDheartClient = arg.(*mocktssclients.MockDheartClient)
		case *mocktss.MockPartyManager:
			mockPartyManager = arg.(*mocktss.MockPartyManager)
		}
	}

	if mockPmm == nil {
		mockPmm = mock.NewMockPostedMessageManager(ctrl)
		mockPmm.EXPECT().ShouldProcessMsg(gomock.Any(), gomock.Any()).Return(true, []byte("")).Times(1)
	}

	if mockPublicDb == nil {
		mockPublicDb = mocktss.NewMockStorage(ctrl)
		mockPublicDb.EXPECT().SaveKeygen(gomock.Any()).Times(1)
		mockPublicDb.EXPECT().ProcessTxRecord(gomock.Any()).Times(1)
	}

	if mockGlobalData == nil {
		mockGlobalData = mockcommon.NewMockGlobalData(ctrl)
		mockGlobalData.EXPECT().IsCatchingUp().Return(false).Times(1)

	}

	if mockDheartClient == nil {

		mockDheartClient = mocktssclients.NewMockDheartClient(ctrl)
		mockDheartClient.EXPECT().KeyGen(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
	}

	if mockPartyManager == nil {
		mockPartyManager = mocktss.NewMockPartyManager(ctrl)
		mockPartyManager.EXPECT().GetActivePartyPubkeys().Return([]ctypes.PubKey{}).Times(1)

	}

	return MockManagerContainer(mockPmm, mockGlobalData, mockDheartClient, mockPartyManager, mockPublicDb)
}
