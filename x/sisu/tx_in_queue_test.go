package sisu

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
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
	q.queue = []*types.TxIn{{}}
	q.processTxIns()
	require.Equal(t, 1, submitCount)

	contracts = mc.Keeper().GetPendingContracts(ctx, chain)
	require.Equal(t, 0, len(contracts))
}
