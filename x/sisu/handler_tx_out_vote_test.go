package sisu

import (
	"fmt"
	"testing"

	"github.com/sisu-network/sisu/x/sisu/background"
	"github.com/sisu-network/sisu/x/sisu/components"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
	db "github.com/tendermint/tm-db"
)

func MockForHandlerTxOutVote() (sdk.Context, background.ManagerContainer) {
	ctx := testmock.TestContext()
	k := testmock.KeeperTestGenesis(ctx)
	pmm := components.NewPostedMessageManager(k)
	privateDb := keeper.NewPrivateDb(".", db.MemDBBackend)

	mc := background.MockManagerContainer(k, pmm, privateDb)

	return ctx, mc
}

func TestHandlerTxOutVote(t *testing.T) {
	t.Run("approve", func(t *testing.T) {
		ctx, mc := MockForHandlerTxOutVote()
		k := mc.Keeper()

		toChain := "ganache2"
		transfers := []*types.TransferDetails{
			{
				Id:      fmt.Sprintf("%s__%s", "ganache1", "hash1"),
				ToChain: toChain,
			},
			{
				Id:      fmt.Sprintf("%s__%s", "ganache1", "hash2"),
				ToChain: toChain,
			},
			{
				Id:      fmt.Sprintf("%s__%s", "ganache1", "hash3"),
				ToChain: toChain,
			},
		}
		outHash := "TxOutHash"
		proposedTxOut := &types.TxOutMsg{
			Signer: "signer1",
			Data: &types.TxOut{
				TxType: types.TxOutType_TRANSFER_OUT,
				Content: &types.TxOutContent{
					OutChain: toChain,
					OutBytes: []byte{0x00},
					OutHash:  outHash,
				},
				Input: &types.TxOutInput{
					TransferIds: []string{fmt.Sprintf("%s__%s", "ganache1", "hash1")},
				},
			},
		}
		k.AddProposedTxOut(ctx, "signer1", proposedTxOut.Data)

		k.AddTransfers(ctx, transfers)
		k.SetTransferQueue(ctx, "ganache2", transfers)

		params := k.GetParams(ctx)
		params.MajorityThreshold = 2
		k.SaveParams(ctx, params)

		h := NewHandlerTxOutConsensed(mc.PostedMessageManager(), mc.Keeper(), mc.PrivateDb())
		msg1 := types.NewTxOutVoteMsg("signer1", &types.TxOutVote{
			AssignedValidator: "signer1",
			TxOutId:           proposedTxOut.Data.GetId(),
			Vote:              types.VoteResult_APPROVE,
		})

		// The TxOut is not processed yet. It needs second vote
		h.DeliverMsg(ctx, msg1)
		txOut := k.GetFinalizedTxOut(ctx, types.GetTxOutIdFromChainAndHash(toChain, outHash))
		require.Nil(t, txOut)

		msg2 := *msg1
		msg2.Signer = "signer2"
		h.DeliverMsg(ctx, &msg2)
		txOut = k.GetFinalizedTxOut(ctx, types.GetTxOutIdFromChainAndHash(toChain, outHash))
		require.NotNil(t, txOut)
		require.Equal(t, proposedTxOut.Data.Input, txOut.Input)
		require.Equal(t, proposedTxOut.Data.Content, txOut.Content)
	})
}
