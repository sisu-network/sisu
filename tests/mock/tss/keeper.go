// Code generated by MockGen. DO NOT EDIT.
// Source: x/tss/keeper/keeper.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	types "github.com/sisu-network/cosmos-sdk/types"
	types0 "github.com/sisu-network/sisu/x/tss/types"
)

// MockKeeper is a mock of Keeper interface.
type MockKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockKeeperMockRecorder
}

// MockKeeperMockRecorder is the mock recorder for MockKeeper.
type MockKeeperMockRecorder struct {
	mock *MockKeeper
}

// NewMockKeeper creates a new mock instance.
func NewMockKeeper(ctrl *gomock.Controller) *MockKeeper {
	mock := &MockKeeper{ctrl: ctrl}
	mock.recorder = &MockKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeeper) EXPECT() *MockKeeperMockRecorder {
	return m.recorder
}

// GetAllPubKeys mocks base method.
func (m *MockKeeper) GetAllPubKeys(ctx types.Context) map[string][]byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllPubKeys", ctx)
	ret0, _ := ret[0].(map[string][]byte)
	return ret0
}

// GetAllPubKeys indicates an expected call of GetAllPubKeys.
func (mr *MockKeeperMockRecorder) GetAllPubKeys(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllPubKeys", reflect.TypeOf((*MockKeeper)(nil).GetAllPubKeys), ctx)
}

// GetPendingContracts mocks base method.
func (m *MockKeeper) GetPendingContracts(ctx types.Context, chain string) []*types0.Contract {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPendingContracts", ctx, chain)
	ret0, _ := ret[0].([]*types0.Contract)
	return ret0
}

// GetPendingContracts indicates an expected call of GetPendingContracts.
func (mr *MockKeeperMockRecorder) GetPendingContracts(ctx, chain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPendingContracts", reflect.TypeOf((*MockKeeper)(nil).GetPendingContracts), ctx, chain)
}

// IsContractExisted mocks base method.
func (m *MockKeeper) IsContractExisted(ctx types.Context, msg *types0.Contract) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsContractExisted", ctx, msg)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsContractExisted indicates an expected call of IsContractExisted.
func (mr *MockKeeperMockRecorder) IsContractExisted(ctx, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsContractExisted", reflect.TypeOf((*MockKeeper)(nil).IsContractExisted), ctx, msg)
}

// IsKeygenAddress mocks base method.
func (m *MockKeeper) IsKeygenAddress(ctx types.Context, keyType, address string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsKeygenAddress", ctx, keyType, address)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsKeygenAddress indicates an expected call of IsKeygenAddress.
func (mr *MockKeeperMockRecorder) IsKeygenAddress(ctx, keyType, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsKeygenAddress", reflect.TypeOf((*MockKeeper)(nil).IsKeygenAddress), ctx, keyType, address)
}

// IsKeygenExisted mocks base method.
func (m *MockKeeper) IsKeygenExisted(ctx types.Context, keyType string, index int) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsKeygenExisted", ctx, keyType, index)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsKeygenExisted indicates an expected call of IsKeygenExisted.
func (mr *MockKeeperMockRecorder) IsKeygenExisted(ctx, keyType, index interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsKeygenExisted", reflect.TypeOf((*MockKeeper)(nil).IsKeygenExisted), ctx, keyType, index)
}

// IsKeygenResultSuccess mocks base method.
func (m *MockKeeper) IsKeygenResultSuccess(ctx types.Context, signerMsg *types0.KeygenResultWithSigner) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsKeygenResultSuccess", ctx, signerMsg)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsKeygenResultSuccess indicates an expected call of IsKeygenResultSuccess.
func (mr *MockKeeperMockRecorder) IsKeygenResultSuccess(ctx, signerMsg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsKeygenResultSuccess", reflect.TypeOf((*MockKeeper)(nil).IsKeygenResultSuccess), ctx, signerMsg)
}

// IsTxInExisted mocks base method.
func (m *MockKeeper) IsTxInExisted(ctx types.Context, msg *types0.TxIn) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsTxInExisted", ctx, msg)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsTxInExisted indicates an expected call of IsTxInExisted.
func (mr *MockKeeperMockRecorder) IsTxInExisted(ctx, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsTxInExisted", reflect.TypeOf((*MockKeeper)(nil).IsTxInExisted), ctx, msg)
}

// IsTxOutExisted mocks base method.
func (m *MockKeeper) IsTxOutExisted(ctx types.Context, msg *types0.TxOut) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsTxOutExisted", ctx, msg)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsTxOutExisted indicates an expected call of IsTxOutExisted.
func (mr *MockKeeperMockRecorder) IsTxOutExisted(ctx, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsTxOutExisted", reflect.TypeOf((*MockKeeper)(nil).IsTxOutExisted), ctx, msg)
}

// PrintStore mocks base method.
func (m *MockKeeper) PrintStore(ctx types.Context, name string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PrintStore", ctx, name)
}

// PrintStore indicates an expected call of PrintStore.
func (mr *MockKeeperMockRecorder) PrintStore(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrintStore", reflect.TypeOf((*MockKeeper)(nil).PrintStore), ctx, name)
}

// SaveContracts mocks base method.
func (m *MockKeeper) SaveContracts(ctx types.Context, msgs []*types0.Contract, saveByteCode bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveContracts", ctx, msgs, saveByteCode)
}

// SaveContracts indicates an expected call of SaveContracts.
func (mr *MockKeeperMockRecorder) SaveContracts(ctx, msgs, saveByteCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveContracts", reflect.TypeOf((*MockKeeper)(nil).SaveContracts), ctx, msgs, saveByteCode)
}

// SaveKeygen mocks base method.
func (m *MockKeeper) SaveKeygen(ctx types.Context, msg *types0.Keygen) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveKeygen", ctx, msg)
}

// SaveKeygen indicates an expected call of SaveKeygen.
func (mr *MockKeeperMockRecorder) SaveKeygen(ctx, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveKeygen", reflect.TypeOf((*MockKeeper)(nil).SaveKeygen), ctx, msg)
}

// SaveKeygenResult mocks base method.
func (m *MockKeeper) SaveKeygenResult(ctx types.Context, signerMsg *types0.KeygenResultWithSigner) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveKeygenResult", ctx, signerMsg)
}

// SaveKeygenResult indicates an expected call of SaveKeygenResult.
func (mr *MockKeeperMockRecorder) SaveKeygenResult(ctx, signerMsg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveKeygenResult", reflect.TypeOf((*MockKeeper)(nil).SaveKeygenResult), ctx, signerMsg)
}

// SaveTxIn mocks base method.
func (m *MockKeeper) SaveTxIn(ctx types.Context, msg *types0.TxIn) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveTxIn", ctx, msg)
}

// SaveTxIn indicates an expected call of SaveTxIn.
func (mr *MockKeeperMockRecorder) SaveTxIn(ctx, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTxIn", reflect.TypeOf((*MockKeeper)(nil).SaveTxIn), ctx, msg)
}

// SaveTxOut mocks base method.
func (m *MockKeeper) SaveTxOut(ctx types.Context, msg *types0.TxOut) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveTxOut", ctx, msg)
}

// SaveTxOut indicates an expected call of SaveTxOut.
func (mr *MockKeeperMockRecorder) SaveTxOut(ctx, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTxOut", reflect.TypeOf((*MockKeeper)(nil).SaveTxOut), ctx, msg)
}

// UpdateContractsStatus mocks base method.
func (m *MockKeeper) UpdateContractsStatus(ctx types.Context, msgs []*types0.Contract, status string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateContractsStatus", ctx, msgs, status)
}

// UpdateContractsStatus indicates an expected call of UpdateContractsStatus.
func (mr *MockKeeperMockRecorder) UpdateContractsStatus(ctx, msgs, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateContractsStatus", reflect.TypeOf((*MockKeeper)(nil).UpdateContractsStatus), ctx, msgs, status)
}
