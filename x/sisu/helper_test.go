package sisu

import (
	"github.com/golang/mock/gomock"
	mocktss "github.com/sisu-network/sisu/tests/mock/tss"
)

func mockCheckTxRecord(mockPublicDb *mocktss.MockStorage) {
	mockPublicDb.EXPECT().SaveTxRecord(gomock.Any(), "signer").Return(1).Times(1)
	mockPublicDb.EXPECT().IsTxRecordProcessed(gomock.Any()).Return(false).Times(1)
	mockPublicDb.EXPECT().ProcessTxRecord(gomock.Any()).Times(1)
}
