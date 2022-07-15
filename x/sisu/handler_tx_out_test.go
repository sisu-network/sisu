package sisu

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForHandlerTxOut() (sdk.Context, ManagerContainer) {
	ctx := testContext()
	k := keeperTestGenesis(ctx)
	txTracker := &MockTxTracker{}
	globalData := &common.MockGlobalData{}
	pmm := NewPostedMessageManager(k)

	mc := MockManagerContainer(k, pmm, globalData, txTracker, &MockTxOutQueue{})
	return ctx, mc
}

func TestHandlerTxOut_TransferOut(t *testing.T) {
	t.Parallel()

	destChain := "ganache2"
	txOutMsg1 := &types.TxOutMsg{
		Signer: "signer",
		Data: &types.TxOut{
			InHashes: []string{fmt.Sprintf("%s__%s", "ganache1", "hash1")},
			TxType:   types.TxOutType_TRANSFER_OUT,
			OutChain: destChain,
			OutBytes: []byte{},
		},
	}

	t.Run("transfer_out_successful", func(t *testing.T) {
		ctx, mc := mockForHandlerTxOut()
		mc.Keeper().SetTransferQueue(ctx, destChain, []*types.Transfer{
			{
				Id: fmt.Sprintf("%s__%s", "ganache1", "hash1"),
			},
			{
				Id: fmt.Sprintf("%s__%s", "ganache1", "hash2"),
			},
			{
				Id: fmt.Sprintf("%s__%s", "ganache1", "hash3"),
			},
		})

		addTxCount := 0
		// txOutQueue := mc.TxOutQueue()
		// txOutQueue.(*MockTxOutQueue).AddTxOutFunc = func(txOut *types.TxOut) {
		// 	addTxCount++
		// }

		handler := NewHandlerTxOut(mc)
		_, err := handler.DeliverMsg(ctx, txOutMsg1)
		require.NoError(t, err)
		require.Equal(t, 1, addTxCount)

		// We are not processing the second request since we have some tx in the pending transfer queue.
		txOutMsg2 := &(*txOutMsg1)
		txOutMsg2.Data.InHashes = []string{fmt.Sprintf("%s__%s", "ganache1", "hash2")}
		handler = NewHandlerTxOut(mc)
		_, err = handler.DeliverMsg(ctx, txOutMsg2)
		require.NoError(t, err)
		require.Equal(t, 1, addTxCount)

		// // Clear the pending queue and we should be able to transfer again
		// mc.Keeper().SetPendingTransfers(ctx, destChain, make([]*types.Transfer, 0))
		// txOutMsg3 := &(*txOutMsg1)
		// txOutMsg2.Data.InHashes = []string{fmt.Sprintf("%s__%s", "ganache1", "hash3")}
		// handler = NewHandlerTxOut(mc)
		// _, err = handler.DeliverMsg(ctx, txOutMsg3)
		// require.NoError(t, err)
		// require.Equal(t, 2, addTxCount)
	})

	t.Run("node_is_catching_up", func(t *testing.T) {
		ctx, mc := mockForHandlerTxOut()
		addTxCount := 0
		// txOutQueue := mc.TxOutQueue()
		// txOutQueue.(*MockTxOutQueue).AddTxOutFunc = func(txOut *types.TxOut) {
		// 	addTxCount = 1
		// }

		globalData := mc.GlobalData().(*common.MockGlobalData)
		globalData.IsCatchingUpFunc = func() bool {
			return true
		}

		handler := NewHandlerTxOut(mc)
		_, err := handler.DeliverMsg(ctx, txOutMsg1)
		require.NoError(t, err)
		require.Equal(t, 0, addTxCount)
	})
}
