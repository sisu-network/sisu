package sisu

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForHandlerTxOut() (sdk.Context, ManagerContainer) {
	ctx := testmock.TestContext()
	k := testmock.KeeperTestGenesis(ctx)
	txTracker := &MockTxTracker{}
	pmm := NewPostedMessageManager(k)

	mc := MockManagerContainer(k, pmm, txTracker, &MockTxOutQueue{})
	return ctx, mc
}

func TestHandlerTxOut_TransferOut(t *testing.T) {
	destChain := "ganache2"
	txOutMsg1 := &types.TxOutMsg{
		Signer: "signer",
		Data: &types.TxOutOld{
			TxType: types.TxOutType_TRANSFER_OUT,
			Content: &types.TxOutContent{
				OutChain: destChain,
				OutBytes: []byte{},
			},
			Input: &types.TxOutInput{
				TransferIds: []string{fmt.Sprintf("%s__%s", "ganache1", "hash1")},
			},
		},
	}

	t.Run("transfer_out_successful", func(t *testing.T) {
		ctx, mc := mockForHandlerTxOut()
		kpr := mc.Keeper()

		transfers := []*types.Transfer{
			{
				Id: fmt.Sprintf("%s__%s", "ganache1", "hash1"),
			},
			{
				Id: fmt.Sprintf("%s__%s", "ganache1", "hash2"),
			},
			{
				Id: fmt.Sprintf("%s__%s", "ganache1", "hash3"),
			},
		}

		kpr.AddTransfers(ctx, transfers)
		kpr.SetTransferQueue(ctx, destChain, transfers)

		handler := NewHandlerTxOut(mc)
		_, err := handler.DeliverMsg(ctx, txOutMsg1)
		require.NoError(t, err)
		transferQueue := kpr.GetTransferQueue(ctx, txOutMsg1.Data.Content.OutChain)
		require.Equal(t, []*types.Transfer{
			{
				Id: fmt.Sprintf("%s__%s", "ganache1", "hash2"),
			},
			{
				Id: fmt.Sprintf("%s__%s", "ganache1", "hash3"),
			},
		}, transferQueue)

		// We are not processing the second request since we have some tx in the pending transfer queue.
		txOutMsg2 := &(*txOutMsg1)
		txOutMsg2.Data.Input.TransferIds = []string{fmt.Sprintf("%s__%s", "ganache1", "hash2")}
		handler = NewHandlerTxOut(mc)
		_, err = handler.DeliverMsg(ctx, txOutMsg2)
		require.NoError(t, err)
		transferQueue = kpr.GetTransferQueue(ctx, txOutMsg1.Data.Content.OutChain)
		require.Equal(t, []*types.Transfer{
			{
				Id: fmt.Sprintf("%s__%s", "ganache1", "hash3"),
			},
		}, transferQueue)
	})
}
