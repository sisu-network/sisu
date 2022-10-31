package chains

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForBridgeManagerTest() (sdk.Context, keeper.Keeper) {
	ctx := testmock.TestContext()
	k := testmock.KeeperTestAfterContractDeployed(ctx)

	return ctx, k
}

func TestGetBridge(t *testing.T) {
	ctx, k := mockForBridgeManagerTest()

	params := &types.Params{
		SupportedChains: []string{"ganache1"},
	}

	k.SaveParams(ctx, params)

	bm := NewBridgeManager("", k, nil, config.Config{})

	bridge := bm.GetBridge(ctx, "ganache1")
	require.NotNil(t, bridge)

	bridge = bm.GetBridge(ctx, "ganache2")
	require.Nil(t, bridge)
}
