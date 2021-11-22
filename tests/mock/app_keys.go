package mock

import (
	"github.com/sisu-network/cosmos-sdk/crypto/keyring"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
)

// Make sure struct implement interface at compile-time
var _ common.AppKeys = (*AppKeys)(nil)

type AppKeys struct {
	InitFunc                func()
	GetSignerInfoFunc       func() keyring.Info
	GetSignerAddressFunc    func() sdk.AccAddress
	GetKeyringFunc          func() keyring.Keyring
	GetEncryptedPrivKeyFunc func() ([]byte, error)
	GetAesEncryptedFunc     func(msg []byte) ([]byte, error)
}

func (a AppKeys) Init() {
	if a.InitFunc == nil {
		panic("function is not defined")
	}

	a.InitFunc()
}

func (a AppKeys) GetSignerInfo() keyring.Info {
	if a.GetSignerInfoFunc == nil {
		panic("function is not defined")
	}

	return a.GetSignerInfoFunc()
}

func (a AppKeys) GetSignerAddress() sdk.AccAddress {
	if a.GetSignerAddressFunc == nil {
		panic("function is not defined")
	}

	return a.GetSignerAddressFunc()
}

func (a AppKeys) GetKeyring() keyring.Keyring {
	if a.GetKeyringFunc == nil {
		panic("function is not defined")
	}

	return a.GetKeyringFunc()
}

func (a AppKeys) GetEncryptedPrivKey() ([]byte, error) {
	if a.GetAesEncryptedFunc == nil {
		panic("function is not defined")
	}

	return a.GetEncryptedPrivKeyFunc()
}

func (a AppKeys) GetAesEncrypted(msg []byte) ([]byte, error) {
	if a.GetAesEncryptedFunc == nil {
		panic("function is not defined")
	}

	return a.GetAesEncryptedFunc(msg)
}
