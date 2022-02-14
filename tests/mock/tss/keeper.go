// Code generated by MockGen. DO NOT EDIT.
// Source: x/sisu/keeper/keeper.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	types "github.com/cosmos/cosmos-sdk/types"
	gomock "github.com/golang/mock/gomock"
	types0 "github.com/sisu-network/sisu/x/sisu/types"
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

// CreateContractAddress mocks base method.
func (m *MockKeeper) CreateContractAddress(ctx types.Context, chain, txOutHash, address string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateContractAddress", ctx, chain, txOutHash, address)
}

// CreateContractAddress indicates an expected call of CreateContractAddress.
func (mr *MockKeeperMockRecorder) CreateContractAddress(ctx, chain, txOutHash, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateContractAddress", reflect.TypeOf((*MockKeeper)(nil).CreateContractAddress), ctx, chain, txOutHash, address)
}

// GetAllKeygenResult mocks base method.
func (m *MockKeeper) GetAllKeygenResult(ctx types.Context, keygenType string, index int32) []*types0.KeygenResultWithSigner {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllKeygenResult", ctx, keygenType, index)
	ret0, _ := ret[0].([]*types0.KeygenResultWithSigner)
	return ret0
}

// GetAllKeygenResult indicates an expected call of GetAllKeygenResult.
func (mr *MockKeeperMockRecorder) GetAllKeygenResult(ctx, keygenType, index interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllKeygenResult", reflect.TypeOf((*MockKeeper)(nil).GetAllKeygenResult), ctx, keygenType, index)
}

// GetContract mocks base method.
func (m *MockKeeper) GetContract(ctx types.Context, chain, hash string, includeByteCode bool) *types0.Contract {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContract", ctx, chain, hash, includeByteCode)
	ret0, _ := ret[0].(*types0.Contract)
	return ret0
}

// GetContract indicates an expected call of GetContract.
func (mr *MockKeeperMockRecorder) GetContract(ctx, chain, hash, includeByteCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContract", reflect.TypeOf((*MockKeeper)(nil).GetContract), ctx, chain, hash, includeByteCode)
}

// GetLatestContractAddressByName mocks base method.
func (m *MockKeeper) GetLatestContractAddressByName(ctx types.Context, chain, name string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestContractAddressByName", ctx, chain, name)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetLatestContractAddressByName indicates an expected call of GetLatestContractAddressByName.
func (mr *MockKeeperMockRecorder) GetLatestContractAddressByName(ctx, chain, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestContractAddressByName", reflect.TypeOf((*MockKeeper)(nil).GetLatestContractAddressByName), ctx, chain, name)
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

// GetTxOut mocks base method.
func (m *MockKeeper) GetTxOut(ctx types.Context, chain, hash string) *types0.TxOut {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxOut", ctx, chain, hash)
	ret0, _ := ret[0].(*types0.TxOut)
	return ret0
}

// GetTxOut indicates an expected call of GetTxOut.
func (mr *MockKeeperMockRecorder) GetTxOut(ctx, chain, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxOut", reflect.TypeOf((*MockKeeper)(nil).GetTxOut), ctx, chain, hash)
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

// IsContractExistedAtAddress mocks base method.
func (m *MockKeeper) IsContractExistedAtAddress(ctx types.Context, chain, address string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsContractExistedAtAddress", ctx, chain, address)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsContractExistedAtAddress indicates an expected call of IsContractExistedAtAddress.
func (mr *MockKeeperMockRecorder) IsContractExistedAtAddress(ctx, chain, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsContractExistedAtAddress", reflect.TypeOf((*MockKeeper)(nil).IsContractExistedAtAddress), ctx, chain, address)
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

// IsTxOutConfirmExisted mocks base method.
func (m *MockKeeper) IsTxOutConfirmExisted(ctx types.Context, outChain, hash string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsTxOutConfirmExisted", ctx, outChain, hash)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsTxOutConfirmExisted indicates an expected call of IsTxOutConfirmExisted.
func (mr *MockKeeperMockRecorder) IsTxOutConfirmExisted(ctx, outChain, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsTxOutConfirmExisted", reflect.TypeOf((*MockKeeper)(nil).IsTxOutConfirmExisted), ctx, outChain, hash)
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

// IsTxRecordProcessed mocks base method.
func (m *MockKeeper) IsTxRecordProcessed(ctx types.Context, hash []byte) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsTxRecordProcessed", ctx, hash)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsTxRecordProcessed indicates an expected call of IsTxRecordProcessed.
func (mr *MockKeeperMockRecorder) IsTxRecordProcessed(ctx, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsTxRecordProcessed", reflect.TypeOf((*MockKeeper)(nil).IsTxRecordProcessed), ctx, hash)
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

// ProcessTxRecord mocks base method.
func (m *MockKeeper) ProcessTxRecord(ctx types.Context, hash []byte) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ProcessTxRecord", ctx, hash)
}

// ProcessTxRecord indicates an expected call of ProcessTxRecord.
func (mr *MockKeeperMockRecorder) ProcessTxRecord(ctx, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessTxRecord", reflect.TypeOf((*MockKeeper)(nil).ProcessTxRecord), ctx, hash)
}

// SaveContract mocks base method.
func (m *MockKeeper) SaveContract(ctx types.Context, msg *types0.Contract, saveByteCode bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveContract", ctx, msg, saveByteCode)
}

// SaveContract indicates an expected call of SaveContract.
func (mr *MockKeeperMockRecorder) SaveContract(ctx, msg, saveByteCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveContract", reflect.TypeOf((*MockKeeper)(nil).SaveContract), ctx, msg, saveByteCode)
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

// SaveTxOutConfirm mocks base method.
func (m *MockKeeper) SaveTxOutConfirm(ctx types.Context, msg *types0.TxOutContractConfirm) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveTxOutConfirm", ctx, msg)
}

// SaveTxOutConfirm indicates an expected call of SaveTxOutConfirm.
func (mr *MockKeeperMockRecorder) SaveTxOutConfirm(ctx, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTxOutConfirm", reflect.TypeOf((*MockKeeper)(nil).SaveTxOutConfirm), ctx, msg)
}

// SaveTxRecord mocks base method.
func (m *MockKeeper) SaveTxRecord(ctx types.Context, hash []byte, val string) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveTxRecord", ctx, hash, val)
	ret0, _ := ret[0].(int)
	return ret0
}

// SaveTxRecord indicates an expected call of SaveTxRecord.
func (mr *MockKeeperMockRecorder) SaveTxRecord(ctx, hash, val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTxRecord", reflect.TypeOf((*MockKeeper)(nil).SaveTxRecord), ctx, hash, val)
}

// UpdateContractAddress mocks base method.
func (m *MockKeeper) UpdateContractAddress(ctx types.Context, chain, hash, address string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateContractAddress", ctx, chain, hash, address)
}

// UpdateContractAddress indicates an expected call of UpdateContractAddress.
func (mr *MockKeeperMockRecorder) UpdateContractAddress(ctx, chain, hash, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateContractAddress", reflect.TypeOf((*MockKeeper)(nil).UpdateContractAddress), ctx, chain, hash, address)
}

// UpdateContractsStatus mocks base method.
func (m *MockKeeper) UpdateContractsStatus(ctx types.Context, chain, contractHash, status string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateContractsStatus", ctx, chain, contractHash, status)
}

// UpdateContractsStatus indicates an expected call of UpdateContractsStatus.
func (mr *MockKeeperMockRecorder) UpdateContractsStatus(ctx, chain, contractHash, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateContractsStatus", reflect.TypeOf((*MockKeeper)(nil).UpdateContractsStatus), ctx, chain, contractHash, status)
}
