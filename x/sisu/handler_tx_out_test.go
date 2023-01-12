package sisu

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForHandlerTxOut() (sdk.Context, ManagerContainer) {
	ctx := testmock.TestContext()
	k := testmock.KeeperTestGenesis(ctx)
	pmm := NewPostedMessageManager(k)

	mc := MockManagerContainer(k, pmm, &MockTxOutQueue{}, &common.MockTxSubmit{},
		&MockValidatorManager{})
	return ctx, mc
}

func TestHandlerTxOut_TransferOut(t *testing.T) {
	destChain := "ganache2"
	txOutMsg1 := &types.TxOutMsg{
		Signer: "signer",
		Data: &types.TxOut{
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

		transfers := []*types.TransferDetails{
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
		require.Equal(t, []*types.TransferDetails{
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
		require.Equal(t, []*types.TransferDetails{
			{
				Id: fmt.Sprintf("%s__%s", "ganache1", "hash3"),
			},
		}, transferQueue)
	})
}

func TestTxOut_MultipleSigners(t *testing.T) {
	ctx, mc := mockForHandlerTxOut()
	k := mc.Keeper()
	txSubmit := mc.TxSubmit().(*common.MockTxSubmit)
	submitCount := 0

	txSubmit.SubmitMessageAsyncFunc = func(msg sdk.Msg) error {
		submitCount++
		return nil
	}

	params := k.GetParams(ctx)
	params.MajorityThreshold = 4
	k.SaveParams(ctx, params)

	destChain := "ganache2"

	txOutMsg1 := &types.TxOutMsg{
		Signer: "signer1",
		Data: &types.TxOut{
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

	transfers := []*types.TransferDetails{
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

	k.AddTransfers(ctx, transfers)
	k.SetTransferQueue(ctx, destChain, transfers)

	valManager := mc.ValidatorManager().(*MockValidatorManager)
	valManager.GetAssignedValidatorFunc = func(ctx sdk.Context, hash string) *types.Node {
		return &types.Node{
			AccAddress: "signer1",
		}
	}

	handler := NewHandlerTxOut(mc)

	for i := 1; i <= 4; i++ {
		msg := *txOutMsg1
		msg.Signer = fmt.Sprintf("signer%d", i)
		_, err := handler.DeliverMsg(ctx, &msg)
		require.Nil(t, err)

		if i < 4 {
			require.Equal(t, i, submitCount)
		} else {
			// There is no fourth message since with 3 messages, the TxOut is already processed.
			require.Equal(t, 3, submitCount)
		}
	}
}
