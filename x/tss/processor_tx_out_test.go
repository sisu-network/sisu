package tss

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	ctypes "github.com/sisu-network/cosmos-sdk/crypto/types"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/sisu/tests/mock"
	mockcommon "github.com/sisu-network/sisu/tests/mock/common"
	mocktss "github.com/sisu-network/sisu/tests/mock/tss"
	mocktssclients "github.com/sisu-network/sisu/tests/mock/tss/tssclients"
	"github.com/sisu-network/sisu/x/tss/types"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"
	"github.com/stretchr/testify/require"
)

func TestDeliverTxOut(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	mockDb := mock.NewMockDatabase(ctrl)
	mockDb.EXPECT().UpdateTxOutStatus("eth", gomock.Any(), tssTypes.TxOutStatusPreSigning, gomock.Any()).Return(nil).Times(1)
	mockDb.EXPECT().UpdateTxOutStatus("eth", gomock.Any(), tssTypes.TxOutStatusSigning, gomock.Any()).Return(nil).Times(1)
	mockDb.EXPECT().UpdateTxOutStatus("eth", gomock.Any(), tssTypes.TxOutStatusSigned, gomock.Any()).Return(nil).Times(1)

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

	ctx := sdk.Context{}
	txOutWithSigner := &types.TxOutWithSigner{
		Data: &types.TxOut{
			OutChain: "eth",
			OutBytes: binary,
		},
	}

	mockKeeper := mocktss.NewMockKeeper(ctrl)
	mockKeeper.EXPECT().IsTxOutExisted(gomock.Any(), gomock.Any()).Return(false).Times(1)
	mockKeeper.EXPECT().SaveTxOut(gomock.Any(), gomock.Any()).Times(1)

	mockGlobalData := mockcommon.NewMockGlobalData(ctrl)
	mockGlobalData.EXPECT().IsCatchingUp().Return(false).Times(1)

	p := &Processor{
		partyManager: mockPartyManager,
		dheartClient: mockDheartClient,
		db:           mockDb,
		keeper:       mockKeeper,
		globalData:   mockGlobalData,
	}
	p.currentHeight.Store(int64(0))

	bytes, err := p.deliverTxOut(ctx, txOutWithSigner)
	require.NoError(t, err)
	require.Empty(t, bytes)
}

func TestDeliverTxOut_BlockCatchingUp(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	ctx := sdk.Context{}
	txOutWithSigner := &types.TxOutWithSigner{
		Data: &types.TxOut{
			OutChain: "eth",
		},
	}

	mockKeeper := mocktss.NewMockKeeper(ctrl)
	mockKeeper.EXPECT().IsTxOutExisted(gomock.Any(), gomock.Any()).Return(false).Times(1)
	mockKeeper.EXPECT().SaveTxOut(gomock.Any(), gomock.Any()).Times(1)

	mockGlobalData := mockcommon.NewMockGlobalData(ctrl)
	mockGlobalData.EXPECT().IsCatchingUp().Return(true).Times(1) // block is catching up.

	// This is the case when a node is catching up with the network, no TSS call is made.
	mockDb := mock.NewMockDatabase(ctrl)
	mockDb.EXPECT().UpdateTxOutStatus("eth", gomock.Any(), tssTypes.TxOutStatusPreSigning, gomock.Any()).Return(nil).Times(0)
	mockDb.EXPECT().InsertTxOuts(gomock.Any()).Return(nil).Times(1)

	p := &Processor{
		db:         mockDb,
		keeper:     mockKeeper,
		globalData: mockGlobalData,
	}

	bytes, err := p.deliverTxOut(ctx, txOutWithSigner)
	require.NoError(t, err)
	require.Empty(t, bytes)
}
