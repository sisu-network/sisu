package sisu

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForTestHandlerFundGateway() (sdk.Context, ManagerContainer) {
	ctx := testContext()
	k := keeperTestGenesis(ctx)
	pmm := NewPostedMessageManager(k)
	txSubmit := &common.MockTxSubmit{}
	appKeys := common.NewMockAppKeys()

	mc := MockManagerContainer(k, pmm, txSubmit, appKeys)

	return ctx, mc
}

func TestHandlerFundGateway(t *testing.T) {
	ctx, mc := mockForTestHandlerFundGateway()
	txSubmit := mc.TxSubmit().(*common.MockTxSubmit)

	count := 0
	txSubmit.SubmitMessageAsyncFunc = func(msg sdk.Msg) error {
		txOutMsg, ok := msg.(*types.TxOutMsg)

		require.True(t, ok, "Submitted mesasge should be a txout message")
		require.Equal(t, types.TxOutType_CONTRACT_DEPLOYMENT, txOutMsg.Data.TxType,
			"This must be a deployment tx")

		count++
		return nil
	}

	k := mc.Keeper()
	erc20 := SupportedContracts[ContractErc20Gateway]
	k.SaveContract(ctx, &types.Contract{
		Chain: "ganache1",
		Hash:  erc20.AbiHash,
	}, false)

	h := NewHandlerFundGateway(mc)
	h.DeliverMsg(ctx, &types.FundGatewayMsg{
		Data: &types.FundGateway{
			Chain: "ganache1",
		},
	})

	require.Equal(t, 1, count, "A new tx must be submitted")

	// Try to fund the gateway second time. There should be no new tx created.
	h.DeliverMsg(ctx, &types.FundGatewayMsg{
		Data: &types.FundGateway{
			Chain: "ganache1",
		},
	})
	require.Equal(t, 1, count, "No new tx is submitted")
}
