package eth

import (
	"testing"

	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func Test_getTokenOnChain(t *testing.T) {
	allTokens := map[string]*types.Token{
		"t1": {
			Id:        "t1",
			Chains:    []string{"ganache1", "ganache2"},
			Addresses: []string{"t1_addr1", "t1_addr2"},
		},
		"t2": {
			Id:        "t2",
			Chains:    []string{"ganache1", "ganache2"},
			Addresses: []string{"t2_addr1", "t2_addr2"},
		},
	}

	token := getTokenOnChain(allTokens, "t1_addr1", "ganache1")
	require.Equal(t, "t1", token.Id)

	token = getTokenOnChain(allTokens, "t1_addr2", "ganache1")
	require.Nil(t, token)

	token = getTokenOnChain(allTokens, "t2_addr2", "ganache2")
	require.Equal(t, "t2", token.Id)
}
