package lisk

import (
	"testing"

	"github.com/sisu-network/deyes/config"
	"github.com/stretchr/testify/require"
)

func TestIntegration_CreateNewClient(t *testing.T) {
	t.Skip()
	config := config.Chain{Chain: "lisk-testnet", Rpcs: []string{"https://testnet-service.lisk.com/api/v2"}}
	client := NewLiskRPC(config).(*defaultLiskRPC)

	require.Equal(t, client.rpc, config.Rpcs[0])
	require.Equal(t, client.chain, config.Chain)
}
