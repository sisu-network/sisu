package sisu

import (
	"testing"

	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func TestDefaultValidatorManager_GetExceedSlashThresholdValidators(t *testing.T) {
	t.Parallel()

	ctx := testContext()
	keeper := keeperTestGenesis(ctx)
	validatorManager := NewValidatorManager(keeper)

	t.Run("emtpy", func(t *testing.T) {
		t.Parallel()

		slashValidators, err := validatorManager.GetExceedSlashThresholdValidators(ctx)
		require.NoError(t, err)
		require.Empty(t, slashValidators)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		pk := []byte("pubkey1")
		validatorManager.AddNode(ctx, &types.Node{
			ConsensusKey: &types.Pubkey{
				Type:  "ed25519",
				Bytes: pk,
			},
			AccAddress:  "0x1",
			IsValidator: true,
			Status:      types.NodeStatus_Validator,
		})

		// Slash points exceed threshold then node should be slash
		require.NoError(t, keeper.IncSlashToken(ctx, pk, SlashPointThreshold+1))
		slashValidators, err := validatorManager.GetExceedSlashThresholdValidators(ctx)
		require.NoError(t, err)
		require.Len(t, slashValidators, 1)

		require.NoError(t, keeper.DecSlashToken(ctx, pk, SlashPointThreshold+1))
		slashValidators, err = validatorManager.GetExceedSlashThresholdValidators(ctx)
		require.NoError(t, err)
		require.Empty(t, slashValidators)
	})
}

func TestDefaultValidatorManager_UpdateNodeStatus(t *testing.T) {
	t.Parallel()

	ctx := testContext()
	keeper := keeperTestGenesis(ctx)
	validatorManager := NewValidatorManager(keeper)

	t.Run("from_validator_to_candidate", func(t *testing.T) {
		t.Parallel()

		pk := []byte("pubkey1")
		accAddr := "0x1"
		validatorManager.AddNode(ctx, &types.Node{
			ConsensusKey: &types.Pubkey{
				Type:  "ed25519",
				Bytes: pk,
			},
			AccAddress:  accAddr,
			IsValidator: true,
			Status:      types.NodeStatus_Validator,
		})

		validatorManager.UpdateNodeStatus(ctx, accAddr, pk, types.NodeStatus_Candidate)
		vals := validatorManager.GetNodesByStatus(ctx, types.NodeStatus_Candidate)
		node := vals[accAddr]
		require.Equal(t, types.NodeStatus_Candidate, node.Status)
		require.False(t, node.IsValidator)
	})

	t.Run("from_candidate_to_validator", func(t *testing.T) {
		t.Parallel()

		pk := []byte("pubkey2")
		accAddr := "0x2"
		validatorManager.AddNode(ctx, &types.Node{
			ConsensusKey: &types.Pubkey{
				Type:  "ed25519",
				Bytes: pk,
			},
			AccAddress:  accAddr,
			IsValidator: false,
			Status:      types.NodeStatus_Candidate,
		})

		validatorManager.UpdateNodeStatus(ctx, accAddr, pk, types.NodeStatus_Validator)
		vals := validatorManager.GetNodesByStatus(ctx, types.NodeStatus_Candidate)
		node := vals[accAddr]
		require.NotEmpty(t, node)
		require.Equal(t, types.NodeStatus_Validator, node.Status)
		require.True(t, node.IsValidator)
	})
}
