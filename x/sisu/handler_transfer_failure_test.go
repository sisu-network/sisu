package sisu

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForHandlerTransferFailure() (sdk.Context, keeper.Keeper, PostedMessageManager) {
	ctx := TestContext()
	k := KeeperTestGenesis(ctx)
	pmm := NewPostedMessageManager(k)

	return ctx, k, pmm
}

func TestHandlerTransferFailure(t *testing.T) {
	ctx, k, pmm := mockForHandlerTransferFailure()
	chain := "ganache1"
	queue := []*types.Transfer{
		{
			Id:      "1",
			ToChain: chain,
		},
		{
			Id:      "2",
			ToChain: chain,
		},
	}

	k.AddTransfers(ctx, queue)
	k.SetTransferQueue(ctx, chain, queue)

	h := NewHanlderTransferFailure(k, pmm)
	h.DeliverMsg(ctx, types.NewTransferFailureMsg("signer", &types.TransferFailure{
		Chain: chain,
		Ids:   []string{"2"},
	}))

	queue = h.keeper.GetTransferQueue(ctx, chain)
	require.Equal(t, []*types.Transfer{
		{
			Id:      "1",
			ToChain: chain,
		},
	}, queue)
}
