package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHoldProcessing(t *testing.T) {
	db := GetTestStorage()

	var hold bool
	hold = db.GetHoldProcessing("job", "ganache1")
	require.Equal(t, false, hold)

	db.SetHoldProcessing("job", "ganache1", true)
	hold = db.GetHoldProcessing("job", "ganache1")
	require.Equal(t, true, hold)

	hold = db.GetHoldProcessing("job", "ganache2")
	require.Equal(t, false, hold)
}
