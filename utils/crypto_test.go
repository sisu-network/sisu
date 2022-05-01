package utils

import (
	"testing"

	"encoding/base64"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestGetCosmosPubKey(t *testing.T) {
	t.Parallel()

	bz, err := base64.StdEncoding.DecodeString("sk9Ab7wGydi2YEXJzOWAHlKBc/3un2i78hXEl0Mlohs=")
	require.NoError(t, err)
	abci.Ed25519ValidatorUpdate(bz, 100)
}
