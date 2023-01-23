package sisu

import (
	"github.com/sisu-network/sisu/x/sisu/components"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/testmock"
)

func mockForHandlerTransferFailure() (sdk.Context, keeper.Keeper, components.PostedMessageManager) {
	ctx := testmock.TestContext()
	k := testmock.KeeperTestGenesis(ctx)
	pmm := components.NewPostedMessageManager(k)

	return ctx, k, pmm
}

func TestHandlerTransferFailure(t *testing.T) {
	// TODO: Add back this test.
}
