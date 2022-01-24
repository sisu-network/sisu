package sisu

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/golang/mock/gomock"
	eyesTypes "github.com/sisu-network/deyes/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/sisu/tests/mock"
	mocktss "github.com/sisu-network/sisu/tests/mock/tss"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"
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

		hashWithSig := utils.RandomHeximalString(32)
		hashNoSig := utils.RandomHeximalString(32)

		mockPrivateDb := mocktss.NewMockStorage(ctrl)
		mockPrivateDb.EXPECT().GetTxOutSig(observedChain, hashWithSig).Return(&types.TxOutSig{
			Chain:       "eth",
			HashWithSig: hashWithSig,
			HashNoSig:   hashNoSig,
		}).Times(1)

		mockPublicDb := mocktss.NewMockStorage(ctrl)
		mockPublicDb.EXPECT().GetTxOut(observedChain, hashNoSig).Return(
			&types.TxOut{
				TxType:   types.TxOutType_NORMAL, // non-deployment tx
				OutChain: observedChain,
				OutHash:  hashNoSig,
			},
		).Times(1)

		mockPublicDb.EXPECT().IsKeygenAddress(libchain.KEY_TYPE_ECDSA, keygenAddress).Return(true).Times(1)

		priv := ed25519.GenPrivKey()
		addr := sdk.AccAddress(priv.PubKey().Address())
		appKeysMock := mock.NewMockAppKeys(ctrl)
		appKeysMock.EXPECT().GetSignerAddress().Return(addr).MinTimes(1)

		done := make(chan bool)
		mockTxSubmit := mock.NewMockTxSubmit(ctrl)
		mockTxSubmit.EXPECT().SubmitMessageAsync(gomock.Any()).Return(nil).Do(func(id interface{}) {
			done <- true
		}).Times(1)

		txs := &eyesTypes.Txs{
			Chain: observedChain,
			Block: int64(utils.RandomNaturalNumber(1000)),
			Arr: []*eyesTypes.Tx{{
				Hash:       hashWithSig,
				Serialized: []byte{},
				From:       keygenAddress,
			}},
		}
		processor := &Processor{
			publicDb:  mockPublicDb,
			privateDb: mockPrivateDb,
			appKeys:   appKeysMock,
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
		mockTxSubmit.EXPECT().SubmitMessageAsync(gomock.Any()).Return(nil).Do(func(id interface{}) {
			done <- true
		}).Times(1)

		observedChain := "eth"
		toAddress := utils.RandomHeximalString(64)
		fromAddres := utils.RandomHeximalString(64)

		mockPublicDb := mocktss.NewMockStorage(ctrl)
		mockPublicDb.EXPECT().IsKeygenAddress(libchain.KEY_TYPE_ECDSA, fromAddres).Return(false).Times(1)

		priv := ed25519.GenPrivKey()
		addr := sdk.AccAddress(priv.PubKey().Address())
		appKeysMock := mock.NewMockAppKeys(ctrl)
		appKeysMock.EXPECT().GetSignerAddress().Return(addr).MinTimes(1)

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
			publicDb: mockPublicDb,
			appKeys:  appKeysMock,
			txSubmit: mockTxSubmit,
		}

		err := processor.OnTxIns(txs)
		<-done

		require.NoError(t, err)
	})
}
