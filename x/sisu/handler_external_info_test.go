package sisu

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForHandlerExternalInfo() (sdk.Context, ManagerContainer) {
	ctx := testContext()
	k := keeperTestGenesis(ctx)
	mc := MockManagerContainer(k)

	return ctx, mc
}

func TestHandlerExternalInfo_BlockHeight(t *testing.T) {
	ctx, mc := mockForHandlerExternalInfo()

	h := NewHandlerExternalInfo(mc)

	// Msg1 - chain ganache1
	sdkMsg1 := types.NewExternalInfoBlockHeight("signer1", &types.BlockHeight{
		Chain:  "ganache1",
		Height: 5,
	})
	_, err := h.DeliverMsg(ctx, sdkMsg1)
	require.Nil(t, err)

	// Msg2 - chain ganache1
	sdkMsg2 := types.NewExternalInfoBlockHeight("signer1", &types.BlockHeight{
		Chain:  "ganache1",
		Height: 10,
	})
	_, err = h.DeliverMsg(ctx, sdkMsg2)
	require.Nil(t, err)

	record := h.keeper.GetBlockHeightRecord(ctx, "signer1")
	require.Equal(t, []*types.BlockHeight{
		{
			Chain:  "ganache1",
			Height: 10,
		},
	}, record.BlockHeights)

	// Msg3 - chain ganache2
	sdkMsg3 := types.NewExternalInfoBlockHeight("signer1", &types.BlockHeight{
		Chain:  "ganache2",
		Height: 15,
	})
	_, err = h.DeliverMsg(ctx, sdkMsg3)
	require.Nil(t, err)

	record = h.keeper.GetBlockHeightRecord(ctx, "signer1")
	require.Equal(t, []*types.BlockHeight{
		{
			Chain:  "ganache1",
			Height: 10,
		},
		{
			Chain:  "ganache2",
			Height: 15,
		},
	}, record.BlockHeights)
}
