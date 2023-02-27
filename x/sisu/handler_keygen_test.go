package sisu

import (
	"testing"

	"github.com/sisu-network/sisu/x/sisu/background"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForHandlerKeygen() (sdk.Context, background.ManagerContainer) {
	ctx := testmock.TestContext()
	k := testmock.KeeperTestGenesis(ctx)
	globalData := &components.MockGlobalData{}
	pmm := components.NewPostedMessageManager(k)

	valsMag := &components.MockValidatorManager{}
	valsMag.GetValidatorPubkeysFunc = func(ctx sdk.Context) []ctypes.PubKey {
		return []ctypes.PubKey{}
	}

	dheartClient := &external.MockDheartClient{}

	mc := background.MockManagerContainer(k, pmm, globalData, valsMag, dheartClient)

	return ctx, mc
}

func TestHandlerKeygen_normal(t *testing.T) {
	submitCount := 0

	ctx, mc := mockForHandlerKeygen()
	dheartClient := mc.DheartClient().(*external.MockDheartClient)
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
	submitCount := 0
	ctx, mc := mockForHandlerKeygen()

	globalData := mc.GlobalData().(*components.MockGlobalData)
	globalData.IsCatchingUpFunc = func() bool {
		return true
	}
	dheartClient := mc.DheartClient().(*external.MockDheartClient)
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
