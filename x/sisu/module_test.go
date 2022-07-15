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
	ctx := testContext()
	k := keeperTestGenesis(ctx)
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
		OutChain: "ganache1",
		OutHash:  "hash1_1",
		OutBytes: bz,
	}
	txOut1_2 := &types.TxOut{
		OutChain: "ganache1",
		OutHash:  "hash1_2",
		OutBytes: bz,
	}
	txOut2_1 := &types.TxOut{
		OutChain: "ganache2",
		OutHash:  "hash2_1",
		OutBytes: bz,
	}

	kpr.SetTxOutQueue(ctx, "ganache1", []*types.TxOut{txOut1_1, txOut1_2})
	kpr.SetTxOutQueue(ctx, "ganache2", []*types.TxOut{txOut2_1})

	module.signTxOut(ctx)

	txOutQueue1 := kpr.GetTxOutQueue(ctx, "ganache1")
	require.Equal(t, []*types.TxOut{txOut1_2}, txOutQueue1)
	txOutQueue2 := kpr.GetTxOutQueue(ctx, "ganache2")
	require.Equal(t, []*types.TxOut{}, txOutQueue2)

	pending1 := kpr.GetPendingTxOut(ctx, "ganache1")
	require.Equal(t, txOut1_1, pending1)
	pending2 := kpr.GetPendingTxOut(ctx, "ganache2")
	require.Equal(t, txOut2_1, pending2)
}
