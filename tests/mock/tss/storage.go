// Code generated by MockGen. DO NOT EDIT.
// Source: x/sisu/keeper/storage.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	types "github.com/sisu-network/sisu/x/sisu/types"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// CreateContractAddress mocks base method.
func (m *MockStorage) CreateContractAddress(chain, txOutHash, address string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateContractAddress", chain, txOutHash, address)
}

// CreateContractAddress indicates an expected call of CreateContractAddress.
func (mr *MockStorageMockRecorder) CreateContractAddress(chain, txOutHash, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateContractAddress", reflect.TypeOf((*MockStorage)(nil).CreateContractAddress), chain, txOutHash, address)
}

// GetAllKeygenPubkeys mocks base method.
func (m *MockStorage) GetAllKeygenPubkeys() map[string][]byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllKeygenPubkeys")
	ret0, _ := ret[0].(map[string][]byte)
	return ret0
}

// GetAllKeygenPubkeys indicates an expected call of GetAllKeygenPubkeys.
func (mr *MockStorageMockRecorder) GetAllKeygenPubkeys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllKeygenPubkeys", reflect.TypeOf((*MockStorage)(nil).GetAllKeygenPubkeys))
}

// GetAllKeygenResult mocks base method.
func (m *MockStorage) GetAllKeygenResult(keygenType string, index int32) []*types.KeygenResultWithSigner {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllKeygenResult", keygenType, index)
	ret0, _ := ret[0].([]*types.KeygenResultWithSigner)
	return ret0
}

// GetAllKeygenResult indicates an expected call of GetAllKeygenResult.
func (mr *MockStorageMockRecorder) GetAllKeygenResult(keygenType, index interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllKeygenResult", reflect.TypeOf((*MockStorage)(nil).GetAllKeygenResult), keygenType, index)
}

// GetAllTokenPricesRecord mocks base method.
func (m *MockStorage) GetAllTokenPricesRecord() map[string]*types.TokenPriceRecord {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTokenPricesRecord")
	ret0, _ := ret[0].(map[string]*types.TokenPriceRecord)
	return ret0
}

// GetAllTokenPricesRecord indicates an expected call of GetAllTokenPricesRecord.
func (mr *MockStorageMockRecorder) GetAllTokenPricesRecord() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTokenPricesRecord", reflect.TypeOf((*MockStorage)(nil).GetAllTokenPricesRecord))
}

// GetContract mocks base method.
func (m *MockStorage) GetContract(chain, hash string, includeByteCode bool) *types.Contract {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContract", chain, hash, includeByteCode)
	ret0, _ := ret[0].(*types.Contract)
	return ret0
}

// GetContract indicates an expected call of GetContract.
func (mr *MockStorageMockRecorder) GetContract(chain, hash, includeByteCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContract", reflect.TypeOf((*MockStorage)(nil).GetContract), chain, hash, includeByteCode)
}

// GetGasPriceRecord mocks base method.
func (m *MockStorage) GetGasPriceRecord(chain string, height int64) *types.GasPriceRecord {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGasPriceRecord", chain, height)
	ret0, _ := ret[0].(*types.GasPriceRecord)
	return ret0
}

// GetGasPriceRecord indicates an expected call of GetGasPriceRecord.
func (mr *MockStorageMockRecorder) GetGasPriceRecord(chain, height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGasPriceRecord", reflect.TypeOf((*MockStorage)(nil).GetGasPriceRecord), chain, height)
}

// GetKeygenPubkey mocks base method.
func (m *MockStorage) GetKeygenPubkey(keyType string) []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKeygenPubkey", keyType)
	ret0, _ := ret[0].([]byte)
	return ret0
}

// GetKeygenPubkey indicates an expected call of GetKeygenPubkey.
func (mr *MockStorageMockRecorder) GetKeygenPubkey(keyType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKeygenPubkey", reflect.TypeOf((*MockStorage)(nil).GetKeygenPubkey), keyType)
}

// GetLatestContractAddressByName mocks base method.
func (m *MockStorage) GetLatestContractAddressByName(chain, name string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestContractAddressByName", chain, name)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetLatestContractAddressByName indicates an expected call of GetLatestContractAddressByName.
func (mr *MockStorageMockRecorder) GetLatestContractAddressByName(chain, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestContractAddressByName", reflect.TypeOf((*MockStorage)(nil).GetLatestContractAddressByName), chain, name)
}

// GetNetworkGasPrice mocks base method.
func (m *MockStorage) GetNetworkGasPrice(chain string) int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNetworkGasPrice", chain)
	ret0, _ := ret[0].(int64)
	return ret0
}

// GetNetworkGasPrice indicates an expected call of GetNetworkGasPrice.
func (mr *MockStorageMockRecorder) GetNetworkGasPrice(chain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNetworkGasPrice", reflect.TypeOf((*MockStorage)(nil).GetNetworkGasPrice), chain)
}

// GetPendingContracts mocks base method.
func (m *MockStorage) GetPendingContracts(chain string) []*types.Contract {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPendingContracts", chain)
	ret0, _ := ret[0].([]*types.Contract)
	return ret0
}

// GetPendingContracts indicates an expected call of GetPendingContracts.
func (mr *MockStorageMockRecorder) GetPendingContracts(chain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPendingContracts", reflect.TypeOf((*MockStorage)(nil).GetPendingContracts), chain)
}

// GetTxOut mocks base method.
func (m *MockStorage) GetTxOut(outChain, hash string) *types.TxOut {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxOut", outChain, hash)
	ret0, _ := ret[0].(*types.TxOut)
	return ret0
}

// GetTxOut indicates an expected call of GetTxOut.
func (mr *MockStorageMockRecorder) GetTxOut(outChain, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxOut", reflect.TypeOf((*MockStorage)(nil).GetTxOut), outChain, hash)
}

// GetTxOutSig mocks base method.
func (m *MockStorage) GetTxOutSig(outChain, hashWithSig string) *types.TxOutSig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxOutSig", outChain, hashWithSig)
	ret0, _ := ret[0].(*types.TxOutSig)
	return ret0
}

// GetTxOutSig indicates an expected call of GetTxOutSig.
func (mr *MockStorageMockRecorder) GetTxOutSig(outChain, hashWithSig interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxOutSig", reflect.TypeOf((*MockStorage)(nil).GetTxOutSig), outChain, hashWithSig)
}

// IsContractExisted mocks base method.
func (m *MockStorage) IsContractExisted(msg *types.Contract) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsContractExisted", msg)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsContractExisted indicates an expected call of IsContractExisted.
func (mr *MockStorageMockRecorder) IsContractExisted(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsContractExisted", reflect.TypeOf((*MockStorage)(nil).IsContractExisted), msg)
}

// IsContractExistedAtAddress mocks base method.
func (m *MockStorage) IsContractExistedAtAddress(chain, address string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsContractExistedAtAddress", chain, address)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsContractExistedAtAddress indicates an expected call of IsContractExistedAtAddress.
func (mr *MockStorageMockRecorder) IsContractExistedAtAddress(chain, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsContractExistedAtAddress", reflect.TypeOf((*MockStorage)(nil).IsContractExistedAtAddress), chain, address)
}

// IsKeygenAddress mocks base method.
func (m *MockStorage) IsKeygenAddress(keyType, address string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsKeygenAddress", keyType, address)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsKeygenAddress indicates an expected call of IsKeygenAddress.
func (mr *MockStorageMockRecorder) IsKeygenAddress(keyType, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsKeygenAddress", reflect.TypeOf((*MockStorage)(nil).IsKeygenAddress), keyType, address)
}

// IsKeygenExisted mocks base method.
func (m *MockStorage) IsKeygenExisted(keyType string, index int) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsKeygenExisted", keyType, index)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsKeygenExisted indicates an expected call of IsKeygenExisted.
func (mr *MockStorageMockRecorder) IsKeygenExisted(keyType, index interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsKeygenExisted", reflect.TypeOf((*MockStorage)(nil).IsKeygenExisted), keyType, index)
}

// IsTxInExisted mocks base method.
func (m *MockStorage) IsTxInExisted(msg *types.TxIn) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsTxInExisted", msg)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsTxInExisted indicates an expected call of IsTxInExisted.
func (mr *MockStorageMockRecorder) IsTxInExisted(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsTxInExisted", reflect.TypeOf((*MockStorage)(nil).IsTxInExisted), msg)
}

// IsTxOutConfirmExisted mocks base method.
func (m *MockStorage) IsTxOutConfirmExisted(outChain, hash string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsTxOutConfirmExisted", outChain, hash)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsTxOutConfirmExisted indicates an expected call of IsTxOutConfirmExisted.
func (mr *MockStorageMockRecorder) IsTxOutConfirmExisted(outChain, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsTxOutConfirmExisted", reflect.TypeOf((*MockStorage)(nil).IsTxOutConfirmExisted), outChain, hash)
}

// IsTxOutExisted mocks base method.
func (m *MockStorage) IsTxOutExisted(msg *types.TxOut) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsTxOutExisted", msg)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsTxOutExisted indicates an expected call of IsTxOutExisted.
func (mr *MockStorageMockRecorder) IsTxOutExisted(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsTxOutExisted", reflect.TypeOf((*MockStorage)(nil).IsTxOutExisted), msg)
}

// IsTxRecordProcessed mocks base method.
func (m *MockStorage) IsTxRecordProcessed(hash []byte) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsTxRecordProcessed", hash)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsTxRecordProcessed indicates an expected call of IsTxRecordProcessed.
func (mr *MockStorageMockRecorder) IsTxRecordProcessed(hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsTxRecordProcessed", reflect.TypeOf((*MockStorage)(nil).IsTxRecordProcessed), hash)
}

// LoadValidators mocks base method.
func (m *MockStorage) LoadValidators() []*types.Node {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadValidators")
	ret0, _ := ret[0].([]*types.Node)
	return ret0
}

// LoadValidators indicates an expected call of LoadValidators.
func (mr *MockStorageMockRecorder) LoadValidators() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadValidators", reflect.TypeOf((*MockStorage)(nil).LoadValidators))
}

// PrintStore mocks base method.
func (m *MockStorage) PrintStore(name string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PrintStore", name)
}

// PrintStore indicates an expected call of PrintStore.
func (mr *MockStorageMockRecorder) PrintStore(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrintStore", reflect.TypeOf((*MockStorage)(nil).PrintStore), name)
}

// PrintStoreKeys mocks base method.
func (m *MockStorage) PrintStoreKeys(name string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PrintStoreKeys", name)
}

// PrintStoreKeys indicates an expected call of PrintStoreKeys.
func (mr *MockStorageMockRecorder) PrintStoreKeys(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrintStoreKeys", reflect.TypeOf((*MockStorage)(nil).PrintStoreKeys), name)
}

// ProcessTxRecord mocks base method.
func (m *MockStorage) ProcessTxRecord(hash []byte) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ProcessTxRecord", hash)
}

// ProcessTxRecord indicates an expected call of ProcessTxRecord.
func (mr *MockStorageMockRecorder) ProcessTxRecord(hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessTxRecord", reflect.TypeOf((*MockStorage)(nil).ProcessTxRecord), hash)
}

// SaveContract mocks base method.
func (m *MockStorage) SaveContract(msg *types.Contract, saveByteCode bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveContract", msg, saveByteCode)
}

// SaveContract indicates an expected call of SaveContract.
func (mr *MockStorageMockRecorder) SaveContract(msg, saveByteCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveContract", reflect.TypeOf((*MockStorage)(nil).SaveContract), msg, saveByteCode)
}

// SaveContracts mocks base method.
func (m *MockStorage) SaveContracts(msgs []*types.Contract, saveByteCode bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveContracts", msgs, saveByteCode)
}

// SaveContracts indicates an expected call of SaveContracts.
func (mr *MockStorageMockRecorder) SaveContracts(msgs, saveByteCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveContracts", reflect.TypeOf((*MockStorage)(nil).SaveContracts), msgs, saveByteCode)
}

// SaveKeygen mocks base method.
func (m *MockStorage) SaveKeygen(msg *types.Keygen) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveKeygen", msg)
}

// SaveKeygen indicates an expected call of SaveKeygen.
func (mr *MockStorageMockRecorder) SaveKeygen(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveKeygen", reflect.TypeOf((*MockStorage)(nil).SaveKeygen), msg)
}

// SaveKeygenResult mocks base method.
func (m *MockStorage) SaveKeygenResult(signerMsg *types.KeygenResultWithSigner) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveKeygenResult", signerMsg)
}

// SaveKeygenResult indicates an expected call of SaveKeygenResult.
func (mr *MockStorageMockRecorder) SaveKeygenResult(signerMsg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveKeygenResult", reflect.TypeOf((*MockStorage)(nil).SaveKeygenResult), signerMsg)
}

// SaveNetworkGasPrice mocks base method.
func (m *MockStorage) SaveNetworkGasPrice(chain string, gasPrice int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveNetworkGasPrice", chain, gasPrice)
}

// SaveNetworkGasPrice indicates an expected call of SaveNetworkGasPrice.
func (mr *MockStorageMockRecorder) SaveNetworkGasPrice(chain, gasPrice interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveNetworkGasPrice", reflect.TypeOf((*MockStorage)(nil).SaveNetworkGasPrice), chain, gasPrice)
}

// SaveNode mocks base method.
func (m *MockStorage) SaveNode(node *types.Node) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveNode", node)
}

// SaveNode indicates an expected call of SaveNode.
func (mr *MockStorageMockRecorder) SaveNode(node interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveNode", reflect.TypeOf((*MockStorage)(nil).SaveNode), node)
}

// SaveTxIn mocks base method.
func (m *MockStorage) SaveTxIn(msg *types.TxIn) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveTxIn", msg)
}

// SaveTxIn indicates an expected call of SaveTxIn.
func (mr *MockStorageMockRecorder) SaveTxIn(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTxIn", reflect.TypeOf((*MockStorage)(nil).SaveTxIn), msg)
}

// SaveTxOut mocks base method.
func (m *MockStorage) SaveTxOut(msg *types.TxOut) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveTxOut", msg)
}

// SaveTxOut indicates an expected call of SaveTxOut.
func (mr *MockStorageMockRecorder) SaveTxOut(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTxOut", reflect.TypeOf((*MockStorage)(nil).SaveTxOut), msg)
}

// SaveTxOutConfirm mocks base method.
func (m *MockStorage) SaveTxOutConfirm(msg *types.TxOutContractConfirm) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveTxOutConfirm", msg)
}

// SaveTxOutConfirm indicates an expected call of SaveTxOutConfirm.
func (mr *MockStorageMockRecorder) SaveTxOutConfirm(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTxOutConfirm", reflect.TypeOf((*MockStorage)(nil).SaveTxOutConfirm), msg)
}

// SaveTxOutSig mocks base method.
func (m *MockStorage) SaveTxOutSig(msg *types.TxOutSig) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveTxOutSig", msg)
}

// SaveTxOutSig indicates an expected call of SaveTxOutSig.
func (mr *MockStorageMockRecorder) SaveTxOutSig(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTxOutSig", reflect.TypeOf((*MockStorage)(nil).SaveTxOutSig), msg)
}

// SaveTxRecord mocks base method.
func (m *MockStorage) SaveTxRecord(hash []byte, signer string) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveTxRecord", hash, signer)
	ret0, _ := ret[0].(int)
	return ret0
}

// SaveTxRecord indicates an expected call of SaveTxRecord.
func (mr *MockStorageMockRecorder) SaveTxRecord(hash, signer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTxRecord", reflect.TypeOf((*MockStorage)(nil).SaveTxRecord), hash, signer)
}

// SetCalculatedTokenPrice mocks base method.
func (m *MockStorage) SetCalculatedTokenPrice(arg0 map[string]float32) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetCalculatedTokenPrice", arg0)
}

// SetCalculatedTokenPrice indicates an expected call of SetCalculatedTokenPrice.
func (mr *MockStorageMockRecorder) SetCalculatedTokenPrice(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCalculatedTokenPrice", reflect.TypeOf((*MockStorage)(nil).SetCalculatedTokenPrice), arg0)
}

// SetGasPrice mocks base method.
func (m *MockStorage) SetGasPrice(msg *types.GasPriceMsg) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetGasPrice", msg)
}

// SetGasPrice indicates an expected call of SetGasPrice.
func (mr *MockStorageMockRecorder) SetGasPrice(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetGasPrice", reflect.TypeOf((*MockStorage)(nil).SetGasPrice), msg)
}

// SetTokenPrices mocks base method.
func (m *MockStorage) SetTokenPrices(blockHeight uint64, msg *types.UpdateTokenPrice) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTokenPrices", blockHeight, msg)
}

// SetTokenPrices indicates an expected call of SetTokenPrices.
func (mr *MockStorageMockRecorder) SetTokenPrices(blockHeight, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTokenPrices", reflect.TypeOf((*MockStorage)(nil).SetTokenPrices), blockHeight, msg)
}

// UpdateContractAddress mocks base method.
func (m *MockStorage) UpdateContractAddress(chain, hash, address string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateContractAddress", chain, hash, address)
}

// UpdateContractAddress indicates an expected call of UpdateContractAddress.
func (mr *MockStorageMockRecorder) UpdateContractAddress(chain, hash, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateContractAddress", reflect.TypeOf((*MockStorage)(nil).UpdateContractAddress), chain, hash, address)
}

// UpdateContractsStatus mocks base method.
func (m *MockStorage) UpdateContractsStatus(chain, contractHash, status string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateContractsStatus", chain, contractHash, status)
}

// UpdateContractsStatus indicates an expected call of UpdateContractsStatus.
func (mr *MockStorageMockRecorder) UpdateContractsStatus(chain, contractHash, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateContractsStatus", reflect.TypeOf((*MockStorage)(nil).UpdateContractsStatus), chain, contractHash, status)
}
