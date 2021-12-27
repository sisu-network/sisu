package tss

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sisu-network/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	eyesTypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/sisu/tests/mock"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
	"github.com/stretchr/testify/require"
)

func TestProcessor_OnObservedTxs(t *testing.T) {
	t.Parallel()

	t.Run("empty_tx", func(t *testing.T) {
		t.Parallel()

		processor := &Processor{}
		require.NoError(t, processor.OnTxIns(&eyesTypes.Txs{}))
	})

	t.Run("success_from_our_key", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		observedChain := "eth"
		keygenAddress := utils.RandomHeximalString(64)

		// Define mock db
		mockDb := mock.NewMockDatabase(ctrl)
		mockDb.EXPECT().IsChainKeyAddress(gomock.Any(), gomock.Any()).Return(true).MinTimes(1)
		mockDb.EXPECT().UpdateTxOutStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).MinTimes(1)
		mockDb.EXPECT().GetTxOutWithHash(gomock.Any(), gomock.Any(), gomock.Any()).Return(
			&types.TxOut{}).MinTimes(1)

		txs := &eyesTypes.Txs{
			Chain: observedChain,
			Block: int64(utils.RandomNaturalNumber(1000)),
			Arr: []*eyesTypes.Tx{{
				Hash:       utils.RandomHeximalString(64),
				Serialized: []byte{},
				From:       keygenAddress,
			}},
		}
		processor := &Processor{}
		processor.db = mockDb

		require.NoError(t, processor.OnTxIns(txs))
	})

	t.Run("success_to_our_key", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		txSubmitterMock := mock.NewMockTxSubmit(ctrl)
		txSubmitterMock.EXPECT().SubmitMessage(gomock.Any()).Return(nil).AnyTimes()

		priv := ed25519.GenPrivKey()
		addr := sdk.AccAddress(priv.PubKey().Address())
		appKeysMock := mock.NewMockAppKeys(ctrl)
		appKeysMock.EXPECT().GetSignerAddress().Return(addr).MinTimes(1)

		observedChain := "eth"
		mockDb := mock.NewMockDatabase(ctrl)
		mockDb.EXPECT().IsChainKeyAddress(gomock.Any(), gomock.Any()).Return(false).MinTimes(1)
		mockDb.EXPECT().UpdateTxOutStatus(observedChain, gomock.Any(), types.TxOutStatusPreBroadcast, gomock.Any()).Return(nil).AnyTimes()
		mockDb.EXPECT().InsertTxIn(gomock.Any()).Return(nil).MinTimes(1)
		keygenAddress := utils.RandomHeximalString(64)

		txs := &eyesTypes.Txs{
			Chain: observedChain,
			Block: int64(utils.RandomNaturalNumber(1000)),
			Arr: []*eyesTypes.Tx{{
				Hash:       utils.RandomHeximalString(64),
				Serialized: []byte{},
				To:         keygenAddress,
				From:       utils.RandomHeximalString(64),
			}},
		}

		// Init processor with mocks
		processor := &Processor{
			db:       mockDb,
			appKeys:  appKeysMock,
			txSubmit: txSubmitterMock,
		}

		require.NoError(t, processor.OnTxIns(txs))
	})
}
