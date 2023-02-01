package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandSerialization(t *testing.T) {
	pauseVaultCmd := &PauseResumeVault{
		Paused: true,
	}

	cmd := &Command{
		Type: &Command_PauseResume{
			PauseResume: pauseVaultCmd,
		},
		Chain: "ganache1",
	}

	bz, err := cmd.Marshal()
	require.Nil(t, err)

	decoded := &Command{}
	err = decoded.Unmarshal(bz)
	require.Nil(t, err)

	cast, ok := decoded.Type.(*Command_PauseResume)
	require.True(t, ok)
	require.Equal(t, pauseVaultCmd, cast.PauseResume)
}
