package sisu

import (
	"testing"

	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func TestDefaultValidatorManager_GetSlashValidators(t *testing.T) {
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
		validatorManager.AddValidator(ctx, &types.Node{
			ConsensusKey: &types.Pubkey{
				Type:  "ed25519",
				Bytes: pk,
			},
			AccAddress:  "0x1",
			IsValidator: true,
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
