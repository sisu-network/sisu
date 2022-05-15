// Code generated by MockGen. DO NOT EDIT.
// Source: x/sisu/keeper/keeper.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	types "github.com/cosmos/cosmos-sdk/types"
	gomock "github.com/golang/mock/gomock"
	keeper "github.com/sisu-network/sisu/x/sisu/keeper"
	types0 "github.com/sisu-network/sisu/x/sisu/types"
	types1 "github.com/tendermint/tendermint/abci/types"
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

// DecBondBalance mocks base method.
func (m *MockKeeper) DecBondBalance(ctx types.Context, address types.AccAddress, amount int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecBondBalance", ctx, address, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// DecBondBalance indicates an expected call of DecBondBalance.
func (mr *MockKeeperMockRecorder) DecBondBalance(ctx, address, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecBondBalance", reflect.TypeOf((*MockKeeper)(nil).DecBondBalance), ctx, address, amount)
}

// DecSlashToken mocks base method.
func (m *MockKeeper) DecSlashToken(ctx types.Context, address types.AccAddress, amount int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecSlashToken", ctx, address, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// DecSlashToken indicates an expected call of DecSlashToken.
func (mr *MockKeeperMockRecorder) DecSlashToken(ctx, address, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecSlashToken", reflect.TypeOf((*MockKeeper)(nil).DecSlashToken), ctx, address, amount)
}

// GetAllChains mocks base method.
func (m *MockKeeper) GetAllChains(ctx types.Context) map[string]*types0.Chain {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllChains", ctx)
	ret0, _ := ret[0].(map[string]*types0.Chain)
	return ret0
}

// GetAllChains indicates an expected call of GetAllChains.
func (mr *MockKeeperMockRecorder) GetAllChains(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllChains", reflect.TypeOf((*MockKeeper)(nil).GetAllChains), ctx)
}

// GetAllDheartIPAddresses mocks base method.
func (m *MockKeeper) GetAllDheartIPAddresses(ctx types.Context) []keeper.AccAddressDheartIP {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllDheartIPAddresses", ctx)
	ret0, _ := ret[0].([]keeper.AccAddressDheartIP)
	return ret0
}

// GetAllDheartIPAddresses indicates an expected call of GetAllDheartIPAddresses.
func (mr *MockKeeperMockRecorder) GetAllDheartIPAddresses(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllDheartIPAddresses", reflect.TypeOf((*MockKeeper)(nil).GetAllDheartIPAddresses), ctx)
}

// GetAllKeygenPubkeys mocks base method.
func (m *MockKeeper) GetAllKeygenPubkeys(ctx types.Context) map[string][]byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllKeygenPubkeys", ctx)
	ret0, _ := ret[0].(map[string][]byte)
	return ret0
}

// GetAllKeygenPubkeys indicates an expected call of GetAllKeygenPubkeys.
func (mr *MockKeeperMockRecorder) GetAllKeygenPubkeys(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllKeygenPubkeys", reflect.TypeOf((*MockKeeper)(nil).GetAllKeygenPubkeys), ctx)
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

// GetAllLiquidities mocks base method.
func (m *MockKeeper) GetAllLiquidities(ctx types.Context) map[string]*types0.Liquidity {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllLiquidities", ctx)
	ret0, _ := ret[0].(map[string]*types0.Liquidity)
	return ret0
}

// GetAllLiquidities indicates an expected call of GetAllLiquidities.
func (mr *MockKeeperMockRecorder) GetAllLiquidities(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllLiquidities", reflect.TypeOf((*MockKeeper)(nil).GetAllLiquidities), ctx)
}

// GetAllTokenPricesRecord mocks base method.
func (m *MockKeeper) GetAllTokenPricesRecord(ctx types.Context) map[string]*types0.TokenPriceRecords {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTokenPricesRecord", ctx)
	ret0, _ := ret[0].(map[string]*types0.TokenPriceRecords)
	return ret0
}

// GetAllTokenPricesRecord indicates an expected call of GetAllTokenPricesRecord.
func (mr *MockKeeperMockRecorder) GetAllTokenPricesRecord(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTokenPricesRecord", reflect.TypeOf((*MockKeeper)(nil).GetAllTokenPricesRecord), ctx)
}

// GetAllTokens mocks base method.
func (m *MockKeeper) GetAllTokens(ctx types.Context) map[string]*types0.Token {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTokens", ctx)
	ret0, _ := ret[0].(map[string]*types0.Token)
	return ret0
}

// GetAllTokens indicates an expected call of GetAllTokens.
func (mr *MockKeeperMockRecorder) GetAllTokens(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTokens", reflect.TypeOf((*MockKeeper)(nil).GetAllTokens), ctx)
}

// GetBondBalance mocks base method.
func (m *MockKeeper) GetBondBalance(ctx types.Context, address types.AccAddress) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBondBalance", ctx, address)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBondBalance indicates an expected call of GetBondBalance.
func (mr *MockKeeperMockRecorder) GetBondBalance(ctx, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBondBalance", reflect.TypeOf((*MockKeeper)(nil).GetBondBalance), ctx, address)
}

// GetChain mocks base method.
func (m *MockKeeper) GetChain(ctx types.Context, chain string) *types0.Chain {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChain", ctx, chain)
	ret0, _ := ret[0].(*types0.Chain)
	return ret0
}

// GetChain indicates an expected call of GetChain.
func (mr *MockKeeperMockRecorder) GetChain(ctx, chain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChain", reflect.TypeOf((*MockKeeper)(nil).GetChain), ctx, chain)
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

// GetGasPriceRecord mocks base method.
func (m *MockKeeper) GetGasPriceRecord(ctx types.Context, chain string, height int64) *types0.GasPriceRecord {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGasPriceRecord", ctx, chain, height)
	ret0, _ := ret[0].(*types0.GasPriceRecord)
	return ret0
}

// GetGasPriceRecord indicates an expected call of GetGasPriceRecord.
func (mr *MockKeeperMockRecorder) GetGasPriceRecord(ctx, chain, height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGasPriceRecord", reflect.TypeOf((*MockKeeper)(nil).GetGasPriceRecord), ctx, chain, height)
}

// GetIncomingValidatorUpdates mocks base method.
func (m *MockKeeper) GetIncomingValidatorUpdates(ctx types.Context) types1.ValidatorUpdates {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIncomingValidatorUpdates", ctx)
	ret0, _ := ret[0].(types1.ValidatorUpdates)
	return ret0
}

// GetIncomingValidatorUpdates indicates an expected call of GetIncomingValidatorUpdates.
func (mr *MockKeeperMockRecorder) GetIncomingValidatorUpdates(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIncomingValidatorUpdates", reflect.TypeOf((*MockKeeper)(nil).GetIncomingValidatorUpdates), ctx)
}

// GetKeygenPubkey mocks base method.
func (m *MockKeeper) GetKeygenPubkey(ctx types.Context, keyType string) []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKeygenPubkey", ctx, keyType)
	ret0, _ := ret[0].([]byte)
	return ret0
}

// GetKeygenPubkey indicates an expected call of GetKeygenPubkey.
func (mr *MockKeeperMockRecorder) GetKeygenPubkey(ctx, keyType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKeygenPubkey", reflect.TypeOf((*MockKeeper)(nil).GetKeygenPubkey), ctx, keyType)
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

// GetLiquidity mocks base method.
func (m *MockKeeper) GetLiquidity(ctx types.Context, chain string) *types0.Liquidity {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLiquidity", ctx, chain)
	ret0, _ := ret[0].(*types0.Liquidity)
	return ret0
}

// GetLiquidity indicates an expected call of GetLiquidity.
func (mr *MockKeeperMockRecorder) GetLiquidity(ctx, chain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLiquidity", reflect.TypeOf((*MockKeeper)(nil).GetLiquidity), ctx, chain)
}

// GetParams mocks base method.
func (m *MockKeeper) GetParams(ctx types.Context) *types0.Params {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetParams", ctx)
	ret0, _ := ret[0].(*types0.Params)
	return ret0
}

// GetParams indicates an expected call of GetParams.
func (mr *MockKeeperMockRecorder) GetParams(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetParams", reflect.TypeOf((*MockKeeper)(nil).GetParams), ctx)
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

// GetSlashToken mocks base method.
func (m *MockKeeper) GetSlashToken(ctx types.Context, address types.AccAddress) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSlashToken", ctx, address)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSlashToken indicates an expected call of GetSlashToken.
func (mr *MockKeeperMockRecorder) GetSlashToken(ctx, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSlashToken", reflect.TypeOf((*MockKeeper)(nil).GetSlashToken), ctx, address)
}

// GetTokens mocks base method.
func (m *MockKeeper) GetTokens(ctx types.Context, tokens []string) map[string]*types0.Token {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokens", ctx, tokens)
	ret0, _ := ret[0].(map[string]*types0.Token)
	return ret0
}

// GetTokens indicates an expected call of GetTokens.
func (mr *MockKeeperMockRecorder) GetTokens(ctx, tokens interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokens", reflect.TypeOf((*MockKeeper)(nil).GetTokens), ctx, tokens)
}

// GetTopBondBalance mocks base method.
func (m *MockKeeper) GetTopBondBalance(ctx types.Context, n int) []types.AccAddress {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopBondBalance", ctx, n)
	ret0, _ := ret[0].([]types.AccAddress)
	return ret0
}

// GetTopBondBalance indicates an expected call of GetTopBondBalance.
func (mr *MockKeeperMockRecorder) GetTopBondBalance(ctx, n interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopBondBalance", reflect.TypeOf((*MockKeeper)(nil).GetTopBondBalance), ctx, n)
}

// GetTxOut mocks base method.
func (m *MockKeeper) GetTxOut(ctx types.Context, outChain, hash string) *types0.TxOut {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxOut", ctx, outChain, hash)
	ret0, _ := ret[0].(*types0.TxOut)
	return ret0
}

// GetTxOut indicates an expected call of GetTxOut.
func (mr *MockKeeperMockRecorder) GetTxOut(ctx, outChain, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxOut", reflect.TypeOf((*MockKeeper)(nil).GetTxOut), ctx, outChain, hash)
}

// GetTxOutSig mocks base method.
func (m *MockKeeper) GetTxOutSig(ctx types.Context, outChain, hashWithSig string) *types0.TxOutSig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxOutSig", ctx, outChain, hashWithSig)
	ret0, _ := ret[0].(*types0.TxOutSig)
	return ret0
}

// GetTxOutSig indicates an expected call of GetTxOutSig.
func (mr *MockKeeperMockRecorder) GetTxOutSig(ctx, outChain, hashWithSig interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxOutSig", reflect.TypeOf((*MockKeeper)(nil).GetTxOutSig), ctx, outChain, hashWithSig)
}

// IncBondBalance mocks base method.
func (m *MockKeeper) IncBondBalance(ctx types.Context, address types.AccAddress, amount int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncBondBalance", ctx, address, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncBondBalance indicates an expected call of IncBondBalance.
func (mr *MockKeeperMockRecorder) IncBondBalance(ctx, address, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncBondBalance", reflect.TypeOf((*MockKeeper)(nil).IncBondBalance), ctx, address, amount)
}

// IncSlashToken mocks base method.
func (m *MockKeeper) IncSlashToken(ctx types.Context, address types.AccAddress, amount int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncSlashToken", ctx, address, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncSlashToken indicates an expected call of IncSlashToken.
func (mr *MockKeeperMockRecorder) IncSlashToken(ctx, address, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncSlashToken", reflect.TypeOf((*MockKeeper)(nil).IncSlashToken), ctx, address, amount)
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

// LoadNodesByStatus mocks base method.
func (m *MockKeeper) LoadNodesByStatus(ctx types.Context, status types0.NodeStatus) []*types0.Node {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadNodesByStatus", ctx, status)
	ret0, _ := ret[0].([]*types0.Node)
	return ret0
}

// LoadNodesByStatus indicates an expected call of LoadNodesByStatus.
func (mr *MockKeeperMockRecorder) LoadNodesByStatus(ctx, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadNodesByStatus", reflect.TypeOf((*MockKeeper)(nil).LoadNodesByStatus), ctx, status)
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

// PrintStoreKeys mocks base method.
func (m *MockKeeper) PrintStoreKeys(ctx types.Context, name string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PrintStoreKeys", ctx, name)
}

// PrintStoreKeys indicates an expected call of PrintStoreKeys.
func (mr *MockKeeperMockRecorder) PrintStoreKeys(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrintStoreKeys", reflect.TypeOf((*MockKeeper)(nil).PrintStoreKeys), ctx, name)
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

// SaveChain mocks base method.
func (m *MockKeeper) SaveChain(ctx types.Context, chain *types0.Chain) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveChain", ctx, chain)
}

// SaveChain indicates an expected call of SaveChain.
func (mr *MockKeeperMockRecorder) SaveChain(ctx, chain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveChain", reflect.TypeOf((*MockKeeper)(nil).SaveChain), ctx, chain)
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

// SaveDheartIPAddress mocks base method.
func (m *MockKeeper) SaveDheartIPAddress(ctx types.Context, address types.AccAddress, ip string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveDheartIPAddress", ctx, address, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveDheartIPAddress indicates an expected call of SaveDheartIPAddress.
func (mr *MockKeeperMockRecorder) SaveDheartIPAddress(ctx, address, ip interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveDheartIPAddress", reflect.TypeOf((*MockKeeper)(nil).SaveDheartIPAddress), ctx, address, ip)
}

// SaveIncomingValidatorUpdates mocks base method.
func (m *MockKeeper) SaveIncomingValidatorUpdates(ctx types.Context, validatorUpdates types1.ValidatorUpdates) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveIncomingValidatorUpdates", ctx, validatorUpdates)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveIncomingValidatorUpdates indicates an expected call of SaveIncomingValidatorUpdates.
func (mr *MockKeeperMockRecorder) SaveIncomingValidatorUpdates(ctx, validatorUpdates interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveIncomingValidatorUpdates", reflect.TypeOf((*MockKeeper)(nil).SaveIncomingValidatorUpdates), ctx, validatorUpdates)
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

// SaveNode mocks base method.
func (m *MockKeeper) SaveNode(ctx types.Context, node *types0.Node) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveNode", ctx, node)
}

// SaveNode indicates an expected call of SaveNode.
func (mr *MockKeeperMockRecorder) SaveNode(ctx, node interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveNode", reflect.TypeOf((*MockKeeper)(nil).SaveNode), ctx, node)
}

// SaveParams mocks base method.
func (m *MockKeeper) SaveParams(ctx types.Context, params *types0.Params) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveParams", ctx, params)
}

// SaveParams indicates an expected call of SaveParams.
func (mr *MockKeeperMockRecorder) SaveParams(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveParams", reflect.TypeOf((*MockKeeper)(nil).SaveParams), ctx, params)
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

// SaveTxOutSig mocks base method.
func (m *MockKeeper) SaveTxOutSig(ctx types.Context, msg *types0.TxOutSig) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveTxOutSig", ctx, msg)
}

// SaveTxOutSig indicates an expected call of SaveTxOutSig.
func (mr *MockKeeperMockRecorder) SaveTxOutSig(ctx, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTxOutSig", reflect.TypeOf((*MockKeeper)(nil).SaveTxOutSig), ctx, msg)
}

// SaveTxRecord mocks base method.
func (m *MockKeeper) SaveTxRecord(ctx types.Context, hash []byte, signer string) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveTxRecord", ctx, hash, signer)
	ret0, _ := ret[0].(int)
	return ret0
}

// SaveTxRecord indicates an expected call of SaveTxRecord.
func (mr *MockKeeperMockRecorder) SaveTxRecord(ctx, hash, signer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTxRecord", reflect.TypeOf((*MockKeeper)(nil).SaveTxRecord), ctx, hash, signer)
}

// SetGasPrice mocks base method.
func (m *MockKeeper) SetGasPrice(ctx types.Context, msg *types0.GasPriceMsg) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetGasPrice", ctx, msg)
}

// SetGasPrice indicates an expected call of SetGasPrice.
func (mr *MockKeeperMockRecorder) SetGasPrice(ctx, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetGasPrice", reflect.TypeOf((*MockKeeper)(nil).SetGasPrice), ctx, msg)
}

// SetLiquidities mocks base method.
func (m *MockKeeper) SetLiquidities(ctx types.Context, liquidities map[string]*types0.Liquidity) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetLiquidities", ctx, liquidities)
}

// SetLiquidities indicates an expected call of SetLiquidities.
func (mr *MockKeeperMockRecorder) SetLiquidities(ctx, liquidities interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLiquidities", reflect.TypeOf((*MockKeeper)(nil).SetLiquidities), ctx, liquidities)
}

// SetTokenPrices mocks base method.
func (m *MockKeeper) SetTokenPrices(ctx types.Context, blockHeight uint64, msg *types0.UpdateTokenPrice) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTokenPrices", ctx, blockHeight, msg)
}

// SetTokenPrices indicates an expected call of SetTokenPrices.
func (mr *MockKeeperMockRecorder) SetTokenPrices(ctx, blockHeight, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTokenPrices", reflect.TypeOf((*MockKeeper)(nil).SetTokenPrices), ctx, blockHeight, msg)
}

// SetTokens mocks base method.
func (m *MockKeeper) SetTokens(ctx types.Context, tokens map[string]*types0.Token) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTokens", ctx, tokens)
}

// SetTokens indicates an expected call of SetTokens.
func (mr *MockKeeperMockRecorder) SetTokens(ctx, tokens interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTokens", reflect.TypeOf((*MockKeeper)(nil).SetTokens), ctx, tokens)
}

// SetValidators mocks base method.
func (m *MockKeeper) SetValidators(ctx types.Context, nodes []*types0.Node) ([]*types0.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetValidators", ctx, nodes)
	ret0, _ := ret[0].([]*types0.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetValidators indicates an expected call of SetValidators.
func (mr *MockKeeperMockRecorder) SetValidators(ctx, nodes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetValidators", reflect.TypeOf((*MockKeeper)(nil).SetValidators), ctx, nodes)
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

// UpdateNodeStatus mocks base method.
func (m *MockKeeper) UpdateNodeStatus(ctx types.Context, consKey []byte, status types0.NodeStatus) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateNodeStatus", ctx, consKey, status)
}

// UpdateNodeStatus indicates an expected call of UpdateNodeStatus.
func (mr *MockKeeperMockRecorder) UpdateNodeStatus(ctx, consKey, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNodeStatus", reflect.TypeOf((*MockKeeper)(nil).UpdateNodeStatus), ctx, consKey, status)
}
