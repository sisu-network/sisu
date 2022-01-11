package tss

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
	mocktssclients "github.com/sisu-network/sisu/tests/mock/tss/tssclients"
	"github.com/sisu-network/sisu/x/tss/types"
	"github.com/stretchr/testify/require"
)

func TestDeliverTxOut_Normal(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

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

	mockPrivateDb := mocktss.NewMockPrivateDb(ctrl)
	mockPrivateDb.EXPECT().SaveTxOut(gomock.Any()).Times(1)

	mockGlobalData := mockcommon.NewMockGlobalData(ctrl)
	mockGlobalData.EXPECT().IsCatchingUp().Return(false).Times(1)

	p := &Processor{
		keeper:       mockKeeper,
		privateDb:    mockPrivateDb,
		partyManager: mockPartyManager,
		dheartClient: mockDheartClient,
		globalData:   mockGlobalData,
	}

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

	mockPrivateDb := mocktss.NewMockPrivateDb(ctrl)
	mockPrivateDb.EXPECT().SaveTxOut(gomock.Any()).Times(1)

	mockGlobalData := mockcommon.NewMockGlobalData(ctrl)
	mockGlobalData.EXPECT().IsCatchingUp().Return(true).Times(1) // block is catching up.

	mockDheartClient := mocktssclients.NewMockDheartClient(ctrl)
	mockDheartClient.EXPECT().KeySign(gomock.Any(), gomock.Any()).Return(nil).Times(0)

	p := &Processor{
		privateDb:  mockPrivateDb,
		keeper:     mockKeeper,
		globalData: mockGlobalData,
	}

	bytes, err := p.deliverTxOut(ctx, txOutWithSigner)
	require.NoError(t, err)
	require.Empty(t, bytes)
}
