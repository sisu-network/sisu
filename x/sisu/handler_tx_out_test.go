package sisu

import (
	"fmt"
	"testing"

	"github.com/sisu-network/sisu/x/sisu/background"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
	db "github.com/tendermint/tm-db"
)

func mockForHandlerTxOut() (sdk.Context, background.ManagerContainer) {
	ctx := testmock.TestContext()
	k := testmock.KeeperTestGenesis(ctx)
	pmm := components.NewPostedMessageManager(k)
	valsManager := &components.MockValidatorManager{
		GetAssignedValidatorFunc: func(ctx sdk.Context, hash string) (*types.Node, error) {
			return &types.Node{
				AccAddress: "signer",
			}, nil
		},
	}
	mockAppKeys := components.NewMockAppKeys()

	mc := background.MockManagerContainer(k, pmm, &MockTxOutQueue{}, valsManager, mockAppKeys,
		keeper.NewPrivateDb(".", db.MemDBBackend), &background.MockBackground{})
	return ctx, mc
}

func TestTxOut_MultipleSigners(t *testing.T) {
	ctx, mc := mockForHandlerTxOut()
	k := mc.Keeper()

	voteCount := 0
	bg := mc.Background().(*background.MockBackground)
	bg.AddVoteTxOutFunc = func(height int64, msg *types.TxOutMsg) {
		voteCount++
	}

	params := k.GetParams(ctx)
	params.MajorityThreshold = 4
	k.SaveParams(ctx, params)

	destChain := "ganache2"

	txOutMsg1 := &types.TxOutMsg{
		Signer: "signer1",
		Data: &types.TxOut{
			TxType: types.TxOutType_TRANSFER,
			Content: &types.TxOutContent{
				OutChain: destChain,
				OutBytes: []byte{},
			},
			Input: &types.TxOutInput{
				TransferUniqIds: []string{fmt.Sprintf("%s__%s__1", "ganache1", "hash1")},
			},
		},
	}

	transfers := []*types.TransferDetails{
		{
			Id: fmt.Sprintf("%s__%s", "ganache1", "hash1"),
		},
		{
			Id: fmt.Sprintf("%s__%s", "ganache1", "hash2"),
		},
		{
			Id: fmt.Sprintf("%s__%s", "ganache1", "hash3"),
		},
	}

	k.AddTransfers(ctx, transfers)
	k.SetTransferQueue(ctx, destChain, transfers)

	valManager := mc.ValidatorManager().(*components.MockValidatorManager)
	valManager.GetAssignedValidatorFunc = func(ctx sdk.Context, hash string) (*types.Node, error) {
		return &types.Node{
			AccAddress: "signer1",
		}, nil
	}

	handler := NewHandlerTxOutProposal(mc)

	for i := 1; i <= 4; i++ {
		msg := *txOutMsg1
		msg.Signer = fmt.Sprintf("signer%d", i)
		_, err := handler.DeliverMsg(ctx, &msg)
		require.Nil(t, err)
	}

	// Make sure the background only add 1 TxOut for vote.
	require.Equal(t, 1, voteCount)
}
