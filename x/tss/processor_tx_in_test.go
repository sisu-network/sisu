package tss

import (
	"testing"

	"github.com/golang/mock/gomock"
	eyesTypes "github.com/sisu-network/deyes/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/sisu/tests/mock"
	mocktss "github.com/sisu-network/sisu/tests/mock/tss"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
	"github.com/stretchr/testify/require"
)

func TestProcessor_OnTxIns(t *testing.T) {
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

		mockPrivateDb := mocktss.NewMockPrivateDb(ctrl)
		mockPrivateDb.EXPECT().IsKeygenAddress(libchain.KEY_TYPE_ECDSA, keygenAddress).Return(true).Times(1)
		mockPrivateDb.EXPECT().PrintStoreKeys(gomock.Any())
		mockPrivateDb.EXPECT().GetTxOutFromSigHash(observedChain, gomock.Any()).Return(&types.TxOut{
			TxType:   types.TxOutType_NORMAL, // non-deployment tx
			OutChain: "eth2",
			OutHash:  utils.RandomHeximalString(32),
		}).Times(1)
		mockPrivateDb.EXPECT().SaveTxOutConfirm(gomock.Any()).Times(1)

		done := make(chan bool)
		mockTxSubmit := mock.NewMockTxSubmit(ctrl)
		mockTxSubmit.EXPECT().SubmitMessage(gomock.Any()).Return(nil).Do(func(id interface{}) {
			done <- true
		}).Times(1)

		txs := &eyesTypes.Txs{
			Chain: observedChain,
			Block: int64(utils.RandomNaturalNumber(1000)),
			Arr: []*eyesTypes.Tx{{
				Hash:       utils.RandomHeximalString(64),
				Serialized: []byte{},
				From:       keygenAddress,
			}},
		}
		processor := &Processor{
			privateDb: mockPrivateDb,
			appKeys:   getMockAppKey(ctrl),
			txSubmit:  mockTxSubmit,
		}

		err := processor.OnTxIns(txs)

		<-done

		require.NoError(t, err)
	})

	t.Run("success_to_our_key", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		done := make(chan bool)
		mockTxSubmit := mock.NewMockTxSubmit(ctrl)
		mockTxSubmit.EXPECT().SubmitMessage(gomock.Any()).Return(nil).Do(func(id interface{}) {
			done <- true
		}).Times(1)

		observedChain := "eth"
		toAddress := utils.RandomHeximalString(64)
		fromAddres := utils.RandomHeximalString(64)

		mockPrivateDb := mocktss.NewMockPrivateDb(ctrl)
		mockPrivateDb.EXPECT().IsKeygenAddress(libchain.KEY_TYPE_ECDSA, fromAddres).Return(false).Times(1)
		mockPrivateDb.EXPECT().SaveTxIn(gomock.Any()).Times(1)

		txs := &eyesTypes.Txs{
			Chain: observedChain,
			Block: int64(utils.RandomNaturalNumber(1000)),
			Arr: []*eyesTypes.Tx{{
				Hash:       utils.RandomHeximalString(64),
				Serialized: []byte{},
				To:         toAddress,
				From:       fromAddres,
			}},
		}

		// Init processor with mocks
		processor := &Processor{
			privateDb: mockPrivateDb,
			appKeys:   getMockAppKey(ctrl),
			txSubmit:  mockTxSubmit,
		}

		err := processor.OnTxIns(txs)
		<-done

		require.NoError(t, err)
	})
}
