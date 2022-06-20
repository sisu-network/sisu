package sisu

import (
	"testing"

	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/stretchr/testify/require"
)

func TestTxOutProducerContractChangeOwnership(t *testing.T) {
	appKeys := common.NewMockAppKeys()
	ctx := testContext()
	keeper := keeperTestAfterContractDeployed(ctx)
	deyesClient := &tssclients.MockDeyesClient{}
	worldState := defaultWorldStateTest(ctx, keeper, deyesClient)

	txOutputProducer := DefaultTxOutputProducer{
		worldState: worldState,
		keeper:     keeper,
		appKeys:    appKeys,
	}

	chain := "ganache1"
	contractHash := "contractHash"
	newOwner := "newOwner"

	txOutWithSigner, err := txOutputProducer.ContractChangeOwnership(ctx, chain, contractHash, newOwner)
	require.NoError(t, err)
	require.Equal(t, "cosmos1qhktedg5njrjc8xy97m9y9vwnvg9atrk3sru7y", txOutWithSigner.Signer)
	require.Equal(t, chain, txOutWithSigner.Data.OutChain)
	require.Equal(t, contractHash, txOutWithSigner.Data.ContractHash)
}
