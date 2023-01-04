package sisu

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
	db "github.com/tendermint/tm-db"
)

func mockForTransferQueue() (sdk.Context, ManagerContainer) {
	ctx := testmock.TestContext()
	k := testmock.KeeperTestGenesis(ctx)
	privateDb := keeper.NewPrivateDb(".", db.MemDBBackend)
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

	mc := MockManagerContainer(ctx, k, txOutputProducer, globalData, txSubmit, privateDb)

	return ctx, mc
}

func TestTransferQueue(t *testing.T) {
	t.Run("transfer_is_saved", func(t *testing.T) {
		ctx, mc := mockForTransferQueue()
		txOutProducer := mc.TxOutProducer().(*MockTxOutputProducer)
		txSubmit := mc.TxSubmit().(*common.MockTxSubmit)
		txSubmitCount := 0
		k := mc.Keeper()

		queue := NewTransferQueue(k, mc.TxOutProducer(), txSubmit,
			mc.Config().Tss, nil, mc.PrivateDb()).(*defaultTransferQueue)
		transfer := &types.TransferDetails{
			Id:          "ganache1__hash1",
			ToRecipient: "0x98Fa8Ab1dd59389138B286d0BeB26bfa4808EC80",
			Token:       "SISU",
			Amount:      utils.EthToWei.String(),
		}

		k.AddTransfers(ctx, []*types.TransferDetails{transfer})
		k.SetTransferQueue(ctx, "ganache2", []*types.TransferDetails{transfer})

		txOutProducer.GetTxOutsFunc = func(ctx sdk.Context, chain string,
			transfers []*types.TransferDetails) ([]*types.TxOutMsg, error) {
			ret := make([]*types.TxOutMsg, len(transfers))
			for i := range transfers {
				ret[i] = &types.TxOutMsg{
					Signer: "signer",
					Data: &types.TxOutOld{
						Content: &types.TxOutContent{
							OutChain: "ganache2",
						},
					},
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
	})
}
