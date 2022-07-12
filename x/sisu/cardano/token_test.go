package cardano

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_wordToByteString(t *testing.T) {
	ret := wordToByteString("WRAP_ADA")
	require.Equal(t, "575241505f414441", ret)
}
