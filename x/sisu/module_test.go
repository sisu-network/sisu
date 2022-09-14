package sisu

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForTestModule() (sdk.Context, ManagerContainer) {
	ctx := TestContext()
	k := KeeperTestGenesis(ctx)
	globalData := &common.MockGlobalData{}
	txOutQueue := &MockTxOutQueue{}

	mc := MockManagerContainer(k, txOutQueue, globalData)
	return ctx, mc
}

func TestModule_signTxOut(t *testing.T) {
	ctx, mc := mockForTestModule()
	kpr := mc.Keeper()
	module := NewAppModule(nil, nil, mc.Keeper(), nil, nil, mc)

	rawTx := ethTypes.NewContractCreation(
		0,
		big.NewInt(0),
		100,
		big.NewInt(100),
		nil,
	)
	bz, err := rawTx.MarshalBinary()
	require.Nil(t, err)

	txOut1_1 := &types.TxOut{
		Content: &types.TxOutContent{
			OutChain: "ganache1",
			OutHash:  "hash1_1",
			OutBytes: bz,
		},
	}
	txOut1_2 := &types.TxOut{
		Content: &types.TxOutContent{
			OutChain: "ganache1",
			OutHash:  "hash1_2",
			OutBytes: bz,
		},
	}
	txOut2_1 := &types.TxOut{
		Content: &types.TxOutContent{
			OutChain: "ganache2",
			OutHash:  "hash2_1",
			OutBytes: bz,
		},
	}

	kpr.SetTxOutQueue(ctx, "ganache1", []*types.TxOut{txOut1_1, txOut1_2})
	kpr.SetTxOutQueue(ctx, "ganache2", []*types.TxOut{txOut2_1})

	module.signTxOut(ctx)

	txOutQueue1 := kpr.GetTxOutQueue(ctx, "ganache1")
	require.Equal(t, []*types.TxOut{txOut1_2}, txOutQueue1)
	txOutQueue2 := kpr.GetTxOutQueue(ctx, "ganache2")
	require.Equal(t, []*types.TxOut{}, txOutQueue2)

	pending1 := kpr.GetPendingTxOutInfo(ctx, "ganache1")
	require.Equal(t, &types.PendingTxOutInfo{
		TxOut:        txOut1_1,
		ExpiredBlock: 10,
	}, pending1)
	pending2 := kpr.GetPendingTxOutInfo(ctx, "ganache2")
	require.Equal(t, &types.PendingTxOutInfo{
		TxOut:        txOut2_1,
		ExpiredBlock: 10,
	}, pending2)

	// Clone ctx with height = 20. The pending transaction expires. We should add it back to the
	// queue.
	cloneCtx := sdk.Context{}
	cacheMS := ctx.MultiStore().CacheMultiStore()
	header := ctx.BlockHeader()
	header.Height = 20
	cloneCtx = sdk.NewContext(
		cacheMS, header, ctx.IsCheckTx(), nil,
	)

	// The pending tx should be empty
	module.signTxOut(cloneCtx)
	pending1 = kpr.GetPendingTxOutInfo(cloneCtx, "ganache1")
	require.Nil(t, pending1)
	pending2 = kpr.GetPendingTxOutInfo(cloneCtx, "ganache2")
	require.Nil(t, pending2)

	// The tx is added back to the queue.
	txOutQueue1 = kpr.GetTxOutQueue(ctx, "ganache1")
	require.Equal(t, []*types.TxOut{txOut1_2}, txOutQueue1)
	txOutQueue2 = kpr.GetTxOutQueue(ctx, "ganache2")
	require.Equal(t, []*types.TxOut{}, txOutQueue2)
}
