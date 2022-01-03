package tss

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"

	"github.com/golang/mock/gomock"
	"github.com/sisu-network/cosmos-sdk/crypto/keys/ed25519"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/tests/mock"
)

func getMockAppKey(ctrl *gomock.Controller) common.AppKeys {
	priv := ed25519.GenPrivKey()
	addr := sdk.AccAddress(priv.PubKey().Address())
	appKeysMock := mock.NewMockAppKeys(ctrl)
	appKeysMock.EXPECT().GetSignerAddress().Return(addr).MinTimes(1)

	return appKeysMock
}
