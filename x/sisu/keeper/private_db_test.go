package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTxHashIndex(t *testing.T) {
	db := GetTestStorage()

	key := "test_key"
	index := db.GetTxHashIndex(key)
	require.Equal(t, uint32(0), index)

	db.SetTxHashIndex(key, index+1)
	index = db.GetTxHashIndex(key)
	require.Equal(t, uint32(1), index)
}
