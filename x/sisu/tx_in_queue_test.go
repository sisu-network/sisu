package sisu

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForTxInQueue() (sdk.Context, ManagerContainer) {
	ctx := testContext()
	k := keeperTestGenesis(ctx)
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

func TestTxInQueue_AddTransaction(t *testing.T) {
	t.Parallel()

	ctx, mc := mockForTxInQueue()

	submitCount := 0
	txSubmit := mc.TxSubmit().(*common.MockTxSubmit)
	txSubmit.SubmitMessageAsyncFunc = func(msg sdk.Msg) error {
		submitCount = 1
		return nil
	}

	chain := "ganache2"
	txOutputProducer := mc.TxOutProducer().(*MockTxOutputProducer)
	txOutputProducer.GetTxOutsFunc = func(ctx sdk.Context, height int64, tx []*types.TxIn) []*types.TxOutWithSigner {
		return []*types.TxOutWithSigner{{
			Data: &types.TxOut{
				OutChain:     chain,
				ContractHash: "hash",
				TxType:       types.TxOutType_CONTRACT_DEPLOYMENT,
			},
		}}
	}

	mc.Keeper().SaveContract(ctx, &types.Contract{
		Chain: chain,
		Hash:  "hash",
	}, false)
	contracts := mc.Keeper().GetPendingContracts(ctx, chain)
	require.Equal(t, 1, len(contracts))

	q := NewTxInQueue(mc.Keeper(), mc.TxOutProducer(), mc.GlobalData(), mc.TxSubmit()).(*defaultTxInQueue)
	q.AddTxIn(ctx.BlockHeight(), &types.TxIn{})
	q.processTxIns(ctx)
	require.Equal(t, 1, submitCount)

	contracts = mc.Keeper().GetPendingContracts(ctx, chain)
	require.Equal(t, 0, len(contracts))
}

func TestTxInQueue_ProcessTxInDifferenBlock(t *testing.T) {
	t.Parallel()

	ctx, mc := mockForTxInQueue()
	ctx2 := utils.CloneSdkContext(ctx)
	header := ctx2.BlockHeader()
	header.Height = ctx.BlockHeight() + 1
	ctx2 = ctx2.WithBlockHeader(header)

	submitCount := 0
	txSubmit := mc.TxSubmit().(*common.MockTxSubmit)
	txSubmit.SubmitMessageAsyncFunc = func(msg sdk.Msg) error {
		submitCount = submitCount + 1
		return nil
	}

	chain := "ganache2"
	txOutputProducer := mc.TxOutProducer().(*MockTxOutputProducer)

	q := NewTxInQueue(mc.Keeper(), mc.TxOutProducer(), mc.GlobalData(), mc.TxSubmit()).(*defaultTxInQueue)
	q.AddTxIn(ctx.BlockHeight(), &types.TxIn{
		TxHash: "Hash1",
	})
	q.AddTxIn(ctx2.BlockHeight(), &types.TxIn{
		TxHash: "Hash2",
	})

	txOutputProducer.GetTxOutsFunc = func(ctx sdk.Context, height int64, tx []*types.TxIn) []*types.TxOutWithSigner {
		require.Equal(t, 1, len(tx))
		require.Equal(t, "Hash1", tx[0].TxHash)

		return []*types.TxOutWithSigner{{
			Data: &types.TxOut{
				OutChain: chain,
				TxType:   types.TxOutType_TRANSFER_OUT,
			},
		}}
	}
	q.processTxIns(ctx)
	require.Equal(t, 1, submitCount)

	txOutputProducer.GetTxOutsFunc = func(ctx sdk.Context, height int64, tx []*types.TxIn) []*types.TxOutWithSigner {
		require.Equal(t, 1, len(tx))
		require.Equal(t, "Hash2", tx[0].TxHash)

		return []*types.TxOutWithSigner{{
			Data: &types.TxOut{
				OutChain: chain,
				TxType:   types.TxOutType_TRANSFER_OUT,
			},
		}}
	}
	q.processTxIns(ctx2)
	require.Equal(t, 2, submitCount)
}
