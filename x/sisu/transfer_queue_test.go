package sisu

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForTransferQueue() (sdk.Context, ManagerContainer) {
	ctx := testContext()
	k := keeperTestGenesis(ctx)
	params := k.GetParams(ctx)
	params.TransferOutParams = []*types.TransferOutParams{
		{
			Chain:       "ganache2",
			MaxBatching: 1,
		},
	}
	k.SaveParams(ctx, params)

	txOutputProducer := &MockTxOutputProducer{}
	globalData := &common.MockGlobalData{
		GetReadOnlyContextFunc: func() sdk.Context {
			return ctx
		},
	}
	txSubmit := &common.MockTxSubmit{}

	mc := MockManagerContainer(ctx, k, txOutputProducer, globalData, txSubmit)

	return ctx, mc
}

func TestTransferQueue(t *testing.T) {
	t.Run("transfer_is_saved", func(t *testing.T) {
		t.Parallel()

		ctx, mc := mockForTransferQueue()
		txOutProducer := mc.TxOutProducer().(*MockTxOutputProducer)
		txSubmit := mc.TxSubmit().(*common.MockTxSubmit)
		txSubmitCount := 0
		keeper := mc.Keeper()

		queue := NewTransferQueue(mc.Keeper(), mc.TxOutProducer(), mc.TxSubmit(),
			mc.Config(), nil).(*defaultTransferQueue)
		transfer := &types.Transfer{
			Id:        "ganache1__hash1",
			Recipient: "0x98Fa8Ab1dd59389138B286d0BeB26bfa4808EC80",
			Token:     "SISU",
			Amount:    utils.EthToWei.String(),
		}

		keeper.SetTransferQueue(ctx, "ganache2", []*types.Transfer{transfer})
		txOutProducer.GetTxOutsFunc = func(ctx sdk.Context, chain string,
			transfers []*types.Transfer) ([]*types.TxOutMsg, error) {
			ret := make([]*types.TxOutMsg, len(transfers))
			for i := range transfers {
				ret[i] = &types.TxOutMsg{
					Signer: "signer",
				}
			}
			return ret, nil
		}

		txSubmit.SubmitMessageAsyncFunc = func(msg sdk.Msg) error {
			txSubmitCount++
			return nil
		}

		queue.processBatch(ctx)
		require.Equal(t, 1, txSubmitCount)

		transfer2 := &types.Transfer{
			Id:        "ganache1__hash2",
			Recipient: "0x98Fa8Ab1dd59389138B286d0BeB26bfa4808EC80",
			Token:     "SISU",
			Amount:    utils.EthToWei.String(),
		}
		keeper.SetTransferQueue(ctx, "ganache2", []*types.Transfer{transfer, transfer2})

		queue.processBatch(ctx)
		require.Equal(t, 2, txSubmitCount) // Only 1 more txout is created
	})
}
