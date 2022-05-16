package sisu

import (
	"encoding/base64"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func TestHandlerReshareResult(t *testing.T) {
	t.Parallel()

	t.Run("Reshare result success", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		ctx := testContext()
		keeper := keeperTestGenesis(ctx)
		keeper.SaveParams(ctx, &types.Params{MajorityThreshold: 2})

		consensusKey1, err := base64.StdEncoding.DecodeString("1jPHjoWahm5WDES2ud3zJbzmRzCPLFacQsrl/pbO/Wo=")
		require.NoError(t, err)
		consensusKey2, err := base64.StdEncoding.DecodeString("qsXeJ51BGalR2V2Zz9ugh3ofsIS58Kjya9pDgfKH018=")
		require.NoError(t, err)

		signer1 := "cosmos14z5sxua3m2hxda6e0c7d9j2jzuxfh2fr94xjkw"
		signer2 := "cosmos1nec9fjd7dp0aph6xqp9rqnv426yyrutrj40rt3"
		validatorManager := NewValidatorManager(keeper)
		validatorManager.AddNode(ctx, &types.Node{
			ConsensusKey: &types.Pubkey{
				Type:  "ed25519",
				Bytes: consensusKey1,
			},
			AccAddress:  signer1,
			IsValidator: true,
			Status:      types.NodeStatus_Validator,
		})
		validatorManager.AddNode(ctx, &types.Node{
			ConsensusKey: &types.Pubkey{
				Type:  "ed25519",
				Bytes: consensusKey2,
			},
			AccAddress:  signer2,
			IsValidator: false,
			Status:      types.NodeStatus_Candidate,
		})

		newValSet := [][]byte{consensusKey1, consensusKey2}
		reshareMsg1 := types.NewReshareResultWithSigner(signer1, newValSet, types.ReshareData_SUCCESS)

		pmm := NewPostedMessageManager(keeper, validatorManager)
		mc := MockManagerContainer(validatorManager, keeper, pmm)
		handler := NewHandlerReshareResult(mc)

		// Signed by signer 1
		_, err = handler.DeliverMsg(ctx, reshareMsg1)
		require.NoError(t, err)

		// Signed by signer 2
		reshareMsg2 := types.NewReshareResultWithSigner(signer2, newValSet, types.ReshareData_SUCCESS)
		_, err = handler.DeliverMsg(ctx, reshareMsg2)
		require.NoError(t, err)

		incomingValidateUpdates := keeper.GetIncomingValidatorUpdates(ctx)
		require.NotEmpty(t, incomingValidateUpdates)
	})
}
