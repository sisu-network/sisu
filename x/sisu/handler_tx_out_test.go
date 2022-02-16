package sisu_test

import (
	"math/big"
	"testing"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	mockcommon "github.com/sisu-network/sisu/tests/mock/common"
	mocktss "github.com/sisu-network/sisu/tests/mock/tss"
	mock "github.com/sisu-network/sisu/tests/mock/x/sisu"
	mocktssclients "github.com/sisu-network/sisu/tests/mock/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func TestHandlerTxOut_Normal(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	mockPmm := mock.NewMockPostedMessageManager(ctrl)
	mockPmm.EXPECT().ShouldProcessMsg(gomock.Any(), gomock.Any()).Return(true, []byte("")).Times(1)

	mockPartyManager := mocktss.NewMockPartyManager(ctrl)
	mockPartyManager.EXPECT().GetActivePartyPubkeys().Return([]ctypes.PubKey{}).Times(1)

	mockDheartClient := mocktssclients.NewMockDheartClient(ctrl)
	mockDheartClient.EXPECT().KeySign(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	amount := big.NewInt(100)
	gasLimit := uint64(100)
	gasPrice := big.NewInt(100)
	ethTransaction := ethTypes.NewTx(&ethTypes.LegacyTx{
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &common.Address{},
		Value:    amount,
	})
	binary, err := ethTransaction.MarshalBinary()
	require.NoError(t, err)

	txOutWithSigner := &types.TxOutWithSigner{
		Signer: "signer",
		Data: &types.TxOut{
			OutChain: "eth",
			OutBytes: binary,
		},
	}

	mockPublicDb := mocktss.NewMockStorage(ctrl)
	mockPublicDb.EXPECT().SaveTxOut(gomock.Any()).Times(1)
	mockPublicDb.EXPECT().ProcessTxRecord(gomock.Any()).Times(1)

	mockGlobalData := mockcommon.NewMockGlobalData(ctrl)
	mockGlobalData.EXPECT().IsCatchingUp().Return(false).Times(1)

	mc := sisu.MockManagerContainer(mockPmm, mockGlobalData, mockDheartClient, mockPartyManager, mockPublicDb)

	handler := sisu.NewHandlerTxOut(mc)
	_, err = handler.DeliverMsg(sdk.Context{}, txOutWithSigner)
	require.NoError(t, err)
}
