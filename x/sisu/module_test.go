package sisu

import (
	"testing"

	"github.com/sisu-network/sisu/x/sisu/background"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	db "github.com/tendermint/tm-db"
)

func mockForTestModule() (sdk.Context, background.ManagerContainer) {
	ctx := testmock.TestContext()
	k := testmock.KeeperTestGenesis(ctx)
	globalData := &components.MockGlobalData{}
	txOutQueue := &MockTxOutQueue{}
	privateDb := keeper.NewPrivateDb(".", db.MemDBBackend)

	mc := background.MockManagerContainer(k, txOutQueue, globalData, privateDb)
	return ctx, mc
}

func TestModule_signTxOut(t *testing.T) {
	// ctx, mc := mockForTestModule()
	// kpr := mc.Keeper()
	// privateDb := mc.PrivateDb()
	// module := NewAppModule(nil, nil, mc.Keeper(), nil, nil, nil, mc)

	// rawTx := ethTypes.NewContractCreation(
	// 	0,
	// 	big.NewInt(0),
	// 	100,
	// 	big.NewInt(100),
	// 	nil,
	// )
	// bz, err := rawTx.MarshalBinary()
	// require.Nil(t, err)

	// txOut1_1 := &types.TxOut{
	// 	Content: &types.TxOutContent{
	// 		OutChain: "ganache1",
	// 		OutHash:  "hash1_1",
	// 		OutBytes: bz,
	// 	},
	// }
	// txOut1_2 := &types.TxOut{
	// 	Content: &types.TxOutContent{
	// 		OutChain: "ganache1",
	// 		OutHash:  "hash1_2",
	// 		OutBytes: bz,
	// 	},
	// }
	// txOut2_1 := &types.TxOut{
	// 	Content: &types.TxOutContent{
	// 		OutChain: "ganache2",
	// 		OutHash:  "hash2_1",
	// 		OutBytes: bz,
	// 	},
	// }

	// kpr.SetTxOutQueue(ctx, "ganache1", []*types.TxOut{txOut1_1, txOut1_2})
	// kpr.SetTxOutQueue(ctx, "ganache2", []*types.TxOut{txOut2_1})
	// privateDb.SetPendingTxOut("ganache1", &types.PendingTxOutInfo{
	// 	TxOut:        txOut1_1,
	// 	ExpiredBlock: 50,
	// 	State:        types.PendingTxOutInfo_IN_QUEUE,
	// })
	// privateDb.SetPendingTxOut("ganache2", &types.PendingTxOutInfo{
	// 	TxOut:        txOut2_1,
	// 	ExpiredBlock: 50,
	// 	State:        types.PendingTxOutInfo_IN_QUEUE,
	// })

	// txOutQueue1 := kpr.GetTxOutQueue(ctx, "ganache1")
	// require.Equal(t, []*types.TxOut{txOut1_2}, txOutQueue1)
	// txOutQueue2 := kpr.GetTxOutQueue(ctx, "ganache2")
	// require.Equal(t, []*types.TxOut{}, txOutQueue2)

	// pending1 := privateDb.GetPendingTxOut("ganache1")
	// require.Equal(t, &types.PendingTxOutInfo{
	// 	TxOut:        txOut1_1,
	// 	ExpiredBlock: 50,
	// 	State:        types.PendingTxOutInfo_SIGNING,
	// }, pending1)
	// pending2 := privateDb.GetPendingTxOut("ganache2")
	// require.Equal(t, &types.PendingTxOutInfo{
	// 	TxOut:        txOut2_1,
	// 	ExpiredBlock: 50,
	// 	State:        types.PendingTxOutInfo_SIGNING,
	// }, pending2)

	// // Clone ctx with height = 60. The pending transaction expires. We should add it back to the
	// // queue.
	// cloneCtx := sdk.Context{}
	// cacheMS := ctx.MultiStore().CacheMultiStore()
	// header := ctx.BlockHeader()
	// header.Height = 60
	// cloneCtx = sdk.NewContext(
	// 	cacheMS, header, ctx.IsCheckTx(), nil,
	// )

	// // The pending tx should be empty
	// pending1 = privateDb.GetPendingTxOut("ganache1")
	// require.Nil(t, pending1)
	// pending2 = privateDb.GetPendingTxOut("ganache2")
	// require.Nil(t, pending2)

	// // The tx is added back to the queue.
	// txOutQueue1 = kpr.GetTxOutQueue(ctx, "ganache1")
	// require.Equal(t, []*types.TxOut{txOut1_2}, txOutQueue1)
	// txOutQueue2 = kpr.GetTxOutQueue(ctx, "ganache2")
	// require.Equal(t, []*types.TxOut{}, txOutQueue2)
}
