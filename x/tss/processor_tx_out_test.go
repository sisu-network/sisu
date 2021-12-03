package tss

import (
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/sisu/tests/mock"
	"github.com/sisu-network/sisu/x/tss/types"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"
	"github.com/stretchr/testify/require"
	"math/big"
	"math/rand"
	"testing"
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

	amount := big.NewInt(rand.Int63())
	gasLimit := uint64(rand.Int63())
	gasPrice := big.NewInt(rand.Int63())
	ethTransaction := ethTypes.NewTx(&ethTypes.LegacyTx{
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &common.Address{},
		Value:    amount,
	})
	binary, err := ethTransaction.MarshalBinary()
	require.NoError(t, err)

	ctx := sdk.Context{}
	txOut := types.TxOut{
		OutChain: "eth",
		OutBytes: binary,
	}

	p := &Processor{
		partyManager: &mock.PartyManager{},
		dheartClient: &mock.DheartClient{},
		db:           mockDb,
	}

	bytes, err := p.DeliverTxOut(ctx, &txOut)
	require.NoError(t, err)
	require.Empty(t, bytes)
}
