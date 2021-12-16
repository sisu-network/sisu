package tss

import (
	"testing"

	"github.com/golang/mock/gomock"
	ctypes "github.com/sisu-network/cosmos-sdk/crypto/types"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/sisu/tests/mock"
	mockcommon "github.com/sisu-network/sisu/tests/mock/common"
	mocktss "github.com/sisu-network/sisu/tests/mock/tss"
	mocktssclients "github.com/sisu-network/sisu/tests/mock/tss/tssclients"
	"github.com/sisu-network/sisu/x/tss/types"
)

// Test happy case for keygen proposal.
func TestProcessor_DeliverKeygenProposal(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	ctx := sdk.Context{}

	mockDb := mock.NewMockDatabase(ctrl)
	mockDb.EXPECT().GetKeyGen(libchain.KEY_TYPE_ECDSA).Return(nil, nil).Times(1)
	mockDb.EXPECT().CreateKeygen(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockKeeper := mocktss.NewMockKeeper(ctrl)
	mockKeeper.EXPECT().IsKeygenProposalExisted(gomock.Any(), gomock.Any()).Return(false).Times(1)
	mockKeeper.EXPECT().SaveKeygenProposal(ctx, gomock.Any()).Times(1)

	mockGlobalData := mockcommon.NewMockGlobalData(ctrl)
	mockGlobalData.EXPECT().IsCatchingUp().Return(false).Times(1)

	mockDheartClient := mocktssclients.NewMockDheartClient(ctrl)
	mockDheartClient.EXPECT().KeyGen(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockPartyManager := mocktss.NewMockPartyManager(ctrl)
	mockPartyManager.EXPECT().GetActivePartyPubkeys().Return([]ctypes.PubKey{}).Times(1)

	wrapper := &types.KeygenProposalWithSigner{
		Signer: "",
		Data: &types.KeygenProposal{
			KeyType: libchain.KEY_TYPE_ECDSA,
			Id:      "keygen",
		},
	}

	p := &Processor{
		db:           mockDb,
		keeper:       mockKeeper,
		globalData:   mockGlobalData,
		partyManager: mockPartyManager,
		dheartClient: mockDheartClient,
	}
	p.currentHeight.Store(int64(0))

	p.DeliverKeyGenProposal(ctx, wrapper)
}

// Test Deliver KeygenProposal while the node is catching up.
func TestProcessor_DeliverKeygenProposal_CatchingUp(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	ctx := sdk.Context{}

	mockDb := mock.NewMockDatabase(ctrl)
	mockDb.EXPECT().GetKeyGen(libchain.KEY_TYPE_ECDSA).Return(nil, nil).Times(0)

	mockKeeper := mocktss.NewMockKeeper(ctrl)
	mockKeeper.EXPECT().IsKeygenProposalExisted(gomock.Any(), gomock.Any()).Return(false).Times(1)
	mockKeeper.EXPECT().SaveKeygenProposal(ctx, gomock.Any()).Times(1)

	mockGlobalData := mockcommon.NewMockGlobalData(ctrl)
	mockGlobalData.EXPECT().IsCatchingUp().Return(true).Times(1) // block is catching up.

	wrapper := &types.KeygenProposalWithSigner{
		Signer: "",
		Data: &types.KeygenProposal{
			KeyType: libchain.KEY_TYPE_ECDSA,
			Id:      "keygen",
		},
	}

	p := &Processor{
		db:         mockDb,
		keeper:     mockKeeper,
		globalData: mockGlobalData,
	}
	p.currentHeight.Store(int64(0))

	p.DeliverKeyGenProposal(ctx, wrapper)
}
