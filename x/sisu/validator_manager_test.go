package sisu

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func TestDefaultValidatorManager_GetExceedSlashThresholdValidators(t *testing.T) {
	t.Parallel()

	t.Run("emtpy", func(t *testing.T) {
		t.Parallel()

		ctx := testContext()
		keeper := keeperTestGenesis(ctx)
		validatorManager := NewValidatorManager(keeper)

		slashValidators, err := validatorManager.GetExceedSlashThresholdValidators(ctx)
		require.NoError(t, err)
		require.Empty(t, slashValidators)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := testContext()
		keeper := keeperTestGenesis(ctx)
		validatorManager := NewValidatorManager(keeper)

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

	t.Run("from_validator_to_candidate", func(t *testing.T) {
		t.Parallel()

		ctx := testContext()
		keeper := keeperTestGenesis(ctx)
		validatorManager := NewValidatorManager(keeper)

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

		ctx := testContext()
		keeper := keeperTestGenesis(ctx)
		validatorManager := NewValidatorManager(keeper)

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

	t.Run("success_only_1_candidate", func(t *testing.T) {
		t.Parallel()

		ctx := testContext()
		keeper := keeperTestGenesis(ctx)
		validatorManager := NewValidatorManager(keeper)

		candidate, err := sdk.AccAddressFromBech32("cosmos1g64vzyutdjfdvw5kyae73fc39sksg3r7gzmrzy")
		require.NoError(t, err)
		require.NoError(t, keeper.IncBondBalance(ctx, candidate, 100))

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

func TestValidatorManager_HasConsensus(t *testing.T) {
	t.Parallel()

	t.Run("Tx record reached consensus by 2 signer", func(t *testing.T) {
		t.Parallel()

		ctx := testContext()
		keeper := keeperTestGenesis(ctx)
		val1, err := sdk.AccAddressFromBech32("cosmos1g64vzyutdjfdvw5kyae73fc39sksg3r7gzmrzy")
		require.NoError(t, err)

		val2, err := sdk.AccAddressFromBech32("cosmos19ucgagq35rqhnj4xev5cvwszs05875rlc3k9rq")
		require.NoError(t, err)

		rcHash := []byte("record_hash")
		keeper.SaveTxRecord(ctx, rcHash, val1.String())
		keeper.SaveTxRecord(ctx, rcHash, val2.String())
		keeper.SaveParams(ctx, &types.Params{MajorityThreshold: 2})

		validatorManager := NewValidatorManager(keeper)
		validatorManager.AddNode(ctx, &types.Node{
			ConsensusKey: &types.Pubkey{
				Type:  "ed25519",
				Bytes: []byte("consensus_key_1"),
			},
			AccAddress:  val1.String(),
			IsValidator: true,
			Status:      types.NodeStatus_Validator,
		})
		validatorManager.AddNode(ctx, &types.Node{
			ConsensusKey: &types.Pubkey{
				Type:  "ed25519",
				Bytes: []byte("consensus_key_2"),
			},
			AccAddress:  val2.String(),
			IsValidator: true,
			Status:      types.NodeStatus_Validator,
		})

		got := validatorManager.HasConsensus(ctx, rcHash)
		require.True(t, got)
	})
}
