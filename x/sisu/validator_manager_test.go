package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		addr, err := sdk.AccAddressFromBech32("cosmos1g64vzyutdjfdvw5kyae73fc39sksg3r7gzmrzy")
		require.NoError(t, err)
		validatorManager.AddNode(ctx, &types.Node{
			ConsensusKey: &types.Pubkey{
				Type:  "ed25519",
				Bytes: pk,
			},
			AccAddress:  addr.String(),
			IsValidator: true,
			Status:      types.NodeStatus_Validator,
		})

		// Slash points exceed threshold then node should be slashed
		require.NoError(t, keeper.IncSlashToken(ctx, addr, SlashPointThreshold+1))
		slashValidators, err := validatorManager.GetExceedSlashThresholdValidators(ctx)
		require.NoError(t, err)
		require.Len(t, slashValidators, 1)

		require.NoError(t, keeper.DecSlashToken(ctx, addr, SlashPointThreshold+1))
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

		consensusKey := []byte("pubkey1")
		addr, err := sdk.AccAddressFromBech32("cosmos1g64vzyutdjfdvw5kyae73fc39sksg3r7gzmrzy")
		require.NoError(t, err)
		validatorManager.AddNode(ctx, &types.Node{
			ConsensusKey: &types.Pubkey{
				Type:  "ed25519",
				Bytes: consensusKey,
			},
			AccAddress:  addr.String(),
			IsValidator: true,
			Status:      types.NodeStatus_Validator,
		})

		validatorManager.UpdateNodeStatus(ctx, consensusKey, types.NodeStatus_Candidate)
		candidates := validatorManager.GetNodesByStatus(types.NodeStatus_Candidate)
		node := candidates[string(consensusKey)]
		require.NotEmpty(t, node)
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

		validatorManager.UpdateNodeStatus(ctx, pk, types.NodeStatus_Validator)
		candidates := validatorManager.GetNodesByStatus(types.NodeStatus_Validator)
		node := candidates[string(pk)]
		require.NotEmpty(t, node)
		require.Equal(t, types.NodeStatus_Validator, node.Status)
		require.True(t, node.IsValidator)
	})
}

func TestTestDefaultValidatorManager_GetPotentialCandidates(t *testing.T) {
	t.Parallel()

	ctx := testContext()
	keeper := keeperTestGenesis(ctx)
	validatorManager := NewValidatorManager(keeper)

	t.Run("success_only_1_candidate", func(t *testing.T) {
		t.Parallel()

		candidate, err := sdk.AccAddressFromBech32("cosmos1g64vzyutdjfdvw5kyae73fc39sksg3r7gzmrzy")
		require.NoError(t, err)
		require.NoError(t, keeper.IncBalance(ctx, candidate, 100))

		consensusKey := []byte("0x1")
		validatorManager.AddNode(ctx, &types.Node{
			ConsensusKey: &types.Pubkey{
				Type:  "ed25519",
				Bytes: consensusKey,
			},
			AccAddress:  candidate.String(),
			IsValidator: false,
			Status:      types.NodeStatus_Candidate,
		})

		got := validatorManager.GetPotentialCandidates(ctx, 1)
		require.Len(t, got, 1)
		require.Equal(t, got[0].AccAddress, candidate.String())
	})
}
