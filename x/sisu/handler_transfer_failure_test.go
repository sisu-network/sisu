package sisu

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForHandlerTransferFailure() (sdk.Context, keeper.Keeper, PostedMessageManager) {
	ctx := testContext()
	k := keeperTestGenesis(ctx)
	pmm := NewPostedMessageManager(k)

	return ctx, k, pmm
}

func TestHandlerTransferFailure(t *testing.T) {
	ctx, k, pmm := mockForHandlerTransferFailure()
	queue := []*types.Transfer{
		{
			Id: "1",
		},
		{
			Id: "2",
		},
	}
	chain := "ganache1"
	k.SetTransferQueue(ctx, chain, queue)

	h := NewHanlderTransferFailure(k, pmm)
	h.DeliverMsg(ctx, types.NewTransferFailureMsg("signer", &types.TransferFailure{
		Chain: chain,
		Ids:   []string{"2"},
	}))

	queue = h.keeper.GetTransferQueue(ctx, chain)
	require.Equal(t, []*types.Transfer{
		{
			Id: "1",
		},
	}, queue)
}
