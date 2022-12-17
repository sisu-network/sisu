package sisu

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForHandlerTxOutResult() (sdk.Context, ManagerContainer) {
	ctx := testmock.TestContext()
	k := testmock.KeeperTestGenesis(ctx)
	pmm := NewPostedMessageManager(k)
	transferQ := MockTransferQueue{}
	storage := keeper.GetTestStorage()

	mc := MockManagerContainer(k, pmm, transferQ, &MockTxOutQueue{}, storage)
	return ctx, mc
}

func getTransfers() []*types.Transfer {
	srcChain := "ganache1"
	destChain := "ganache2"
	recipient := "0x8095f5b69F2970f38DC6eBD2682ed71E4939f988"
	token := "SISU"
	hash1 := "123"
	amount := "10000"
	hash2 := "456"
	recipient2 := "0x98Fa8Ab1dd59389138B286d0BeB26bfa4808EC80"
	token2 := "ADA"

	return []*types.Transfer{
		{
			Id:              fmt.Sprintf("%s__%s", srcChain, hash1),
			FromChain:       srcChain,
			FromBlockHeight: 10,
			ToChain:         destChain,
			Token:           token,
			FromHash:        hash1,
			ToRecipient:     recipient,
			Amount:          amount,
		},
		{
			Id:              fmt.Sprintf("%s__%s", srcChain, hash2),
			FromChain:       srcChain,
			FromBlockHeight: 11,
			ToChain:         destChain,
			Token:           token2,
			FromHash:        hash2,
			ToRecipient:     recipient2,
			Amount:          amount,
		},
	}
}

func TestHandlerTxOutResult(t *testing.T) {
	t.Run("tx_included_in_block_successfully", func(t *testing.T) {
		ctx, mc := mockForHandlerTxOutResult()
		k := mc.Keeper()
		privateDb := mc.PrivateDb()

		txOut := &types.TxOut{
			TxType: types.TxOutType_TRANSFER_OUT,
			Content: &types.TxOutContent{
				OutChain: "ganache2",
				OutHash:  "TxOutHash",
			},
			Input: &types.TxOutInput{},
		}
		k.SaveTxOut(ctx, txOut)
		privateDb.SetPendingTxOut("ganache2", &types.PendingTxOutInfo{
			TxOut:        txOut,
			ExpiredBlock: 0,
		})
		txOutId, err := txOut.GetId()
		require.Nil(t, err)

		handler := NewHandlerTxOutResult(mc)
		handler.DeliverMsg(ctx, &types.TxOutResultMsg{
			Data: &types.TxOutResult{
				TxOutId:  txOutId,
				Result:   types.TxOutResultType_IN_BLOCK_SUCCESS,
				OutChain: "ganache2",
				OutHash:  "TxOutHash",
			},
		})

		pending := privateDb.GetPendingTxOut("ganache2")
		require.Nil(t, pending)
	})

	t.Run("tx_result_failure", func(t *testing.T) {
		ctx, mc := mockForHandlerTxOutResult()
		k := mc.Keeper()
		transfers := getTransfers()
		privateDb := mc.PrivateDb()
		k.AddTransfers(ctx, transfers)

		txOut := &types.TxOut{
			TxType: types.TxOutType_TRANSFER_OUT,
			Content: &types.TxOutContent{
				OutChain: "ganache2",
				OutHash:  "TxOutHash",
			},
			Input: &types.TxOutInput{
				TransferIds: []string{transfers[0].Id, transfers[1].Id},
			},
		}
		k.SaveTxOut(ctx, txOut)
		privateDb.SetPendingTxOut("ganache2", &types.PendingTxOutInfo{
			TxOut:        txOut,
			ExpiredBlock: 0,
		})

		txOutId, err := txOut.GetId()
		require.Nil(t, err)
		handler := NewHandlerTxOutResult(mc)
		handler.DeliverMsg(ctx, &types.TxOutResultMsg{
			Data: &types.TxOutResult{
				TxOutId:  txOutId,
				Result:   types.TxOutResultType_IN_BLOCK_FAILURE,
				OutChain: "ganache2",
				OutHash:  "TxOutHash",
			},
		})

		transfers2 := k.GetTransfers(ctx, []string{transfers[0].Id, transfers[1].Id})
		require.Equal(t, 2, len(transfers2))
		for _, transfer := range transfers2 {
			require.Equal(t, 1, int(transfer.RetryNum))
		}
	})
}
