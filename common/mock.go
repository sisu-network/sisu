package common

import (
	"encoding/hex"

	keyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	tcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

type MockAppKeys struct {
	privKey tcrypto.PrivKey
	addr    sdk.AccAddress
}

func NewMockAppKeys() AppKeys {
	bz, err := hex.DecodeString("fd914bab512acb5d8fdd5537146339fdf8bbd141f9cba3c039b8d732188b0d6a62c5713923f2e087ce62554b07d738c07363d84edfc4a3e05c4e65e0")
	if err != nil {
		panic(err)
	}

	addr, err := sdk.AccAddressFromBech32("cosmos1qhktedg5njrjc8xy97m9y9vwnvg9atrk3sru7y")
	if err != nil {
		panic(err)
	}

	privKey := ed25519.PrivKey(bz)
	return &MockAppKeys{
		privKey: privKey,
		addr:    addr,
	}
}

func (ak *MockAppKeys) Init() {
}

func (ak *MockAppKeys) GetSignerInfo() keyring.Info {
	return nil
}

func (ak *MockAppKeys) GetSignerAddress() sdk.AccAddress {
	return ak.addr
}

func (ak *MockAppKeys) GetKeyring() keyring.Keyring {
	panic("Unsupported")
}

func (ak *MockAppKeys) GetEncryptedPrivKey() ([]byte, error) {
	return ak.privKey.Bytes(), nil
}

func (ak *MockAppKeys) GetAesEncrypted(msg []byte) ([]byte, error) {
	return utils.AESDEncrypt(msg, ak.privKey.Bytes())
}
