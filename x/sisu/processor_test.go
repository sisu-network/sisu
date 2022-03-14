package sisu

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/golang/mock/gomock"
	eyesTypes "github.com/sisu-network/deyes/types"
	libchain "github.com/sisu-network/lib/chain"
	mock "github.com/sisu-network/sisu/tests/mock/common"
	mockkeeper "github.com/sisu-network/sisu/tests/mock/x/sisu/keeper"
	"github.com/sisu-network/sisu/utils"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestProcessor_OnTxIns(t *testing.T) {
	t.Parallel()

	t.Run("empty_tx", func(t *testing.T) {
		t.Parallel()

		mc := MockManagerContainer(sdk.Context{})

		processor := &Processor{
			mc: mc,
		}
		require.NoError(t, processor.OnTxIns(&eyesTypes.Txs{}))
	})

	t.Run("success_to_our_key", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		mockTxSubmit := mock.NewMockTxSubmit(ctrl)
		mockTxSubmit.EXPECT().SubmitMessageAsync(gomock.Any()).Return(nil).Times(1)

		observedChain := "eth"
		toAddress := utils.RandomHeximalString(64)
		fromAddres := utils.RandomHeximalString(64)

		mockKeeper := mockkeeper.NewMockKeeper(ctrl)
		mockKeeper.EXPECT().IsKeygenAddress(gomock.Any(), libchain.KEY_TYPE_ECDSA, fromAddres).Return(false).Times(1)

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

		mc := MockManagerContainer(mockKeeper, sdk.Context{})

		// Init processor with mocks
		processor := &Processor{
			mc:       mc,
			keeper:   mockKeeper,
			appKeys:  appKeysMock,
			txSubmit: mockTxSubmit,
		}

		err := processor.OnTxIns(txs)
		// <-done

		require.NoError(t, err)
	})
}
