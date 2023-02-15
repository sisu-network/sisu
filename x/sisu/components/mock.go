package components

import (
	"encoding/hex"
	"errors"

	"github.com/cosmos/cosmos-sdk/client/rpc"
	keyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"
	tcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

type MockAppKeys struct {
	privKey tcrypto.PrivKey
	addr    sdk.AccAddress

	GetSignerAddressFunc func() sdk.AccAddress
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
	GetMyPubkeyFunc         func() tcrypto.PubKey
	SetReadOnlyContextFunc  func(ctx sdk.Context)
	GetReadOnlyContextFunc  func() sdk.Context
	AppInitializedFunc      func() bool
	SetAppInitializedFunc   func()
	RecalculateGasFunc      func(chain string)
	GetRecalculateGasFunc   func() []string
	ResetGasCalculationFunc func()
	ConfirmTxInFunc         func(txIn *types.TxIn)
	GetTxInQueueFunc        func() []*types.TxIn
	ResetTxInQueueFunc      func()
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

func (m *MockGlobalData) SetReadOnlyContext(ctx sdk.Context) {
	if m.SetReadOnlyContextFunc != nil {
		m.SetReadOnlyContextFunc(ctx)
	}
}

func (m *MockGlobalData) GetReadOnlyContext() sdk.Context {
	if m.GetReadOnlyContextFunc != nil {
		return m.GetReadOnlyContextFunc()
	}

	return sdk.Context{}
}

func (m *MockGlobalData) AppInitialized() bool {
	if m.AppInitializedFunc != nil {
		return m.AppInitializedFunc()
	}
	return false
}

func (m *MockGlobalData) SetAppInitialized() {
	if m.SetAppInitializedFunc != nil {
		m.SetAppInitializedFunc()
	}
}

func (m *MockGlobalData) RecalculateGas(chain string) {
	if m.RecalculateGasFunc != nil {
		m.RecalculateGasFunc(chain)
	}
}

func (m *MockGlobalData) GetRecalculateGas() []string {
	if m.GetRecalculateGasFunc != nil {
		return m.GetRecalculateGasFunc()
	}

	return nil
}

func (m *MockGlobalData) ResetGasCalculation() {
	if m.ResetGasCalculationFunc != nil {
		m.ResetGasCalculationFunc()
	}
}

func (m *MockGlobalData) GetMyPubkey() tcrypto.PubKey {
	if m.GetMyPubkeyFunc != nil {
		return m.GetMyPubkey()
	}

	return nil
}

func (m *MockGlobalData) ConfirmTxIn(txIn *types.TxIn) {
	if m.ConfirmTxInFunc != nil {
		m.ConfirmTxInFunc(txIn)
	}
}

func (m *MockGlobalData) GetTxInQueue() []*types.TxIn {
	if m.GetTxInQueueFunc != nil {
		return m.GetTxInQueueFunc()
	}

	return nil
}

func (m *MockGlobalData) ResetTxInQueue() {
	if m.ResetTxInQueueFunc != nil {
		m.ResetTxInQueueFunc()
	}
}

///// Validator Manager
type MockValidatorManager struct {
	AddValidatorFunc         func(ctx sdk.Context, node *types.Node)
	IsValidatorFunc          func(ctx sdk.Context, signer string) bool
	GetValidatorLengthFunc   func(ctx sdk.Context) int
	GetValidatorsFunc        func(ctx sdk.Context) []*types.Node
	GetAssignedValidatorFunc func(ctx sdk.Context, hash string) (*types.Node, error)
}

func (m *MockValidatorManager) AddValidator(ctx sdk.Context, node *types.Node) {
	if m.AddValidatorFunc != nil {
		m.AddValidatorFunc(ctx, node)
	}
}

func (m *MockValidatorManager) IsValidator(ctx sdk.Context, signer string) bool {
	if m.IsValidatorFunc != nil {
		return m.IsValidatorFunc(ctx, signer)
	}

	return false
}

func (m *MockValidatorManager) GetValidatorLength(ctx sdk.Context) int {
	if m.GetValidatorLengthFunc != nil {
		return m.GetValidatorLength(ctx)
	}

	return 0
}

func (m *MockValidatorManager) GetValidators(ctx sdk.Context) []*types.Node {
	if m.GetValidatorsFunc != nil {
		return m.GetValidators(ctx)
	}

	return nil
}

func (m *MockValidatorManager) GetAssignedValidator(ctx sdk.Context, hash string) (*types.Node, error) {
	if m.GetAssignedValidatorFunc != nil {
		return m.GetAssignedValidatorFunc(ctx, hash)
	}

	return nil, errors.New("invalid function")
}
