package common

import (
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/client/rpc"
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

///// TxSubmit
type MockTxSubmit struct {
	SubmitMessageAsyncFunc func(msg sdk.Msg) error
	SubmitMessageSyncFunc  func(msg sdk.Msg) error
}

func (m *MockTxSubmit) SubmitMessageAsync(msg sdk.Msg) error {
	if m.SubmitMessageAsyncFunc != nil {
		return m.SubmitMessageAsyncFunc(msg)
	}

	return nil
}

func (m *MockTxSubmit) SubmitMessageSync(msg sdk.Msg) error {
	if m.SubmitMessageSyncFunc != nil {
		return m.SubmitMessageSyncFunc(msg)
	}

	return nil
}

///// GlobalData

type MockGlobalData struct {
	InitFunc                func()
	UpdateCatchingUpFunc    func() bool
	UpdateValidatorSetsFunc func()
	IsCatchingUpFunc        func() bool
	GetValidatorSetFunc     func() []rpc.ValidatorOutput
	GetMyValidatorAddrFunc  func() string
}

func (m *MockGlobalData) Init() {
	if m.InitFunc != nil {
		m.InitFunc()
	}
}

func (m *MockGlobalData) UpdateCatchingUp() bool {
	if m.UpdateCatchingUpFunc != nil {
		return m.UpdateCatchingUpFunc()
	}

	return false
}

func (m *MockGlobalData) UpdateValidatorSets() {
	if m.UpdateValidatorSetsFunc != nil {
		m.UpdateValidatorSetsFunc()
	}
}

func (m *MockGlobalData) IsCatchingUp() bool {
	if m.IsCatchingUpFunc != nil {
		return m.IsCatchingUpFunc()
	}

	return false
}

func (m *MockGlobalData) GetValidatorSet() []rpc.ValidatorOutput {
	if m.GetValidatorSetFunc != nil {
		return m.GetValidatorSetFunc()
	}

	return nil
}

func (m *MockGlobalData) GetMyValidatorAddr() string {
	if m.GetMyValidatorAddrFunc != nil {
		return m.GetMyValidatorAddrFunc()
	}

	return ""
}
