package sisu

import (
	"testing"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForHandlerKeygen() (sdk.Context, ManagerContainer) {
	ctx := testContext()
	k := keeperTestGenesis(ctx)
	globalData := &common.MockGlobalData{}
	pmm := NewPostedMessageManager(k)

	partyManager := &MockPartyManager{}
	partyManager.GetActivePartyPubkeysFunc = func() []ctypes.PubKey {
		return []ctypes.PubKey{}
	}

	dheartClient := &tssclients.MockDheartClient{}

	mc := MockManagerContainer(k, pmm, globalData, partyManager, dheartClient)

	return ctx, mc
}

func TestHandlerKeygen_normal(t *testing.T) {
	t.Parallel()

	submitCount := 0

	ctx, mc := mockForHandlerKeygen()
	dheartClient := mc.DheartClient().(*tssclients.MockDheartClient)
	dheartClient.KeyGenFunc = func(keygenId, chain string, pubKeys []ctypes.PubKey) error {
		submitCount = 1
		return nil
	}

	msg := &types.KeygenWithSigner{
		Signer: "signer",
		Data: &types.Keygen{
			KeyType: libchain.KEY_TYPE_ECDSA,
			Index:   0,
		},
	}

	handler := NewHandlerKeygen(mc)
	_, err := handler.DeliverMsg(ctx, msg)

	require.NoError(t, err)
	require.Equal(t, 1, submitCount)
}

func TestHandlerKeygen_CatchingUp(t *testing.T) {
	t.Parallel()

	submitCount := 0
	ctx, mc := mockForHandlerKeygen()

	globalData := mc.GlobalData().(*common.MockGlobalData)
	globalData.IsCatchingUpFunc = func() bool {
		return true
	}
	dheartClient := mc.DheartClient().(*tssclients.MockDheartClient)
	dheartClient.KeyGenFunc = func(keygenId, chain string, pubKeys []ctypes.PubKey) error {
		submitCount = 1
		return nil
	}

	msg := &types.KeygenWithSigner{
		Signer: "signer",
		Data: &types.Keygen{
			KeyType: libchain.KEY_TYPE_ECDSA,
			Index:   0,
		},
	}

	handler := NewHandlerKeygen(mc)
	_, err := handler.DeliverMsg(ctx, msg)

	require.NoError(t, err)
	require.Equal(t, 0, submitCount)
}
