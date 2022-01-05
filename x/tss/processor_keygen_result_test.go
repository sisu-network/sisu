package tss

import (
	"testing"

	"github.com/golang/mock/gomock"
	ctypes "github.com/sisu-network/cosmos-sdk/crypto/types"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"
	mockcommon "github.com/sisu-network/sisu/tests/mock/common"
	mocktss "github.com/sisu-network/sisu/tests/mock/tss"
	mocktssclients "github.com/sisu-network/sisu/tests/mock/tss/tssclients"
	"github.com/sisu-network/sisu/x/tss/types"
)

// Test happy case for keygen proposal.
func TestProcessor_deliverKeyGen_normal(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	ctx := sdk.Context{}

	mockKeeper := mocktss.NewMockKeeper(ctrl)
	mockKeeper.EXPECT().IsKeygenExisted(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).Times(1)
	mockKeeper.EXPECT().SaveKeygen(ctx, gomock.Any()).Times(1)

	mockPrivateDb := mocktss.NewMockPrivateDb(ctrl)
	mockPrivateDb.EXPECT().SaveKeygen(gomock.Any()).Times(1)

	mockGlobalData := mockcommon.NewMockGlobalData(ctrl)
	mockGlobalData.EXPECT().IsCatchingUp().Return(false).Times(1)

	mockDheartClient := mocktssclients.NewMockDheartClient(ctrl)
	mockDheartClient.EXPECT().KeyGen(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockPartyManager := mocktss.NewMockPartyManager(ctrl)
	mockPartyManager.EXPECT().GetActivePartyPubkeys().Return([]ctypes.PubKey{}).Times(1)

	wrapper := &types.KeygenWithSigner{
		Signer: "",
		Data: &types.Keygen{
			KeyType: libchain.KEY_TYPE_ECDSA,
			Index:   0,
		},
	}

	p := &Processor{
		keeper:       mockKeeper,
		privateDb:    mockPrivateDb,
		globalData:   mockGlobalData,
		partyManager: mockPartyManager,
		dheartClient: mockDheartClient,
	}

	p.deliverKeygen(ctx, wrapper)
}

// Test Deliver KeygenProposal while the node is catching up.
func TestProcessor_deliverKeyGen_CatchingUp(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	ctx := sdk.Context{}

	mockPrivateDb := mocktss.NewMockPrivateDb(ctrl)
	mockPrivateDb.EXPECT().SaveKeygen(gomock.Any()).Times(1)

	mockKeeper := mocktss.NewMockKeeper(ctrl)
	mockKeeper.EXPECT().IsKeygenExisted(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).Times(1)
	mockKeeper.EXPECT().SaveKeygen(ctx, gomock.Any()).Times(1)

	mockGlobalData := mockcommon.NewMockGlobalData(ctrl)
	mockGlobalData.EXPECT().IsCatchingUp().Return(true).Times(1) // block is catching up.

	wrapper := &types.KeygenWithSigner{
		Signer: "",
		Data: &types.Keygen{
			KeyType: libchain.KEY_TYPE_ECDSA,
			Index:   0,
		},
	}

	p := &Processor{
		keeper:     mockKeeper,
		privateDb:  mockPrivateDb,
		globalData: mockGlobalData,
	}

	p.deliverKeygen(ctx, wrapper)
}
