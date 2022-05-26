package sisu

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForHandlerTxIn() (sdk.Context, ManagerContainer) {
	txSubmit := &common.MockTxSubmit{}
	txTracker := &MockTxTracker{}

	txOutputProducer := &MockTxOutputProducer{}
	txOutputProducer.GetTxOutsFunc = func(ctx sdk.Context, height int64, tx *types.TxIn) []*types.TxOutWithSigner {
		txout := types.NewMsgTxOutWithSigner("signer", types.TxOutType_TRANSFER_OUT, 0, "ganache1",
			"inHash", "ganache2", "outHash", []byte{}, "")

		return []*types.TxOutWithSigner{txout}
	}

	k, ctx := keeper.GetTestKeeperAndContext()
	pmm := NewPostedMessageManager(k)
	k.SaveParams(ctx, &types.Params{
		MajorityThreshold: 1,
	})

	globalData := &common.MockGlobalData{}

	mc := MockManagerContainer(txSubmit, txTracker, txOutputProducer, ctx, k, pmm, globalData)

	return ctx, mc
}

func TestHandlerTxIn_HappyCase(t *testing.T) {
	t.Parallel()
	t.Run("output_is_broadcasted", func(t *testing.T) {
		t.Parallel()

		submitCount := 0
		ctx, mc := mockForHandlerTxIn()
		txSubmit := mc.TxSubmit().(*common.MockTxSubmit)
		txSubmit.SubmitMessageAsyncFunc = func(msg sdk.Msg) error {
			submitCount = 1
			return nil
		}

		handler := NewHandlerTxIn(mc)

		msg := types.NewTxInWithSigner("signer", "ganache1", "", 0, []byte{})

		_, err := handler.DeliverMsg(ctx, msg)
		require.Nil(t, err)

		require.Equal(t, 1, submitCount)
	})

	t.Run("output_is_not_broadcasted", func(t *testing.T) {
		t.Parallel()

		submitCount := 0
		ctx, mc := mockForHandlerTxIn()
		globalData := mc.GlobalData().(*common.MockGlobalData)
		globalData.IsCatchingUpFunc = func() bool {
			return true
		}

		txSubmit := mc.TxSubmit().(*common.MockTxSubmit)
		txSubmit.SubmitMessageAsyncFunc = func(msg sdk.Msg) error {
			submitCount = 1
			return nil
		}

		handler := NewHandlerTxIn(mc)

		msg := types.NewTxInWithSigner("signer", "ganache1", "", 0, []byte{})

		_, err := handler.DeliverMsg(ctx, msg)
		require.Nil(t, err)

		require.Equal(t, 0, submitCount)
	})
}
