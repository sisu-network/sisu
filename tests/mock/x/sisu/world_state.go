// Code generated by MockGen. DO NOT EDIT.
// Source: x/sisu/world_state.go

// Package mock is a generated GoMock package.
package mock

import (
	big "math/big"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	types "github.com/sisu-network/sisu/x/sisu/types"
)

// MockWorldState is a mock of WorldState interface.
type MockWorldState struct {
	ctrl     *gomock.Controller
	recorder *MockWorldStateMockRecorder
}

// MockWorldStateMockRecorder is the mock recorder for MockWorldState.
type MockWorldStateMockRecorder struct {
	mock *MockWorldState
}

// NewMockWorldState creates a new mock instance.
func NewMockWorldState(ctrl *gomock.Controller) *MockWorldState {
	mock := &MockWorldState{ctrl: ctrl}
	mock.recorder = &MockWorldStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWorldState) EXPECT() *MockWorldStateMockRecorder {
	return m.recorder
}

// GetGasPrice mocks base method.
func (m *MockWorldState) GetGasPrice(chain string) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGasPrice", chain)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGasPrice indicates an expected call of GetGasPrice.
func (mr *MockWorldStateMockRecorder) GetGasPrice(chain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGasPrice", reflect.TypeOf((*MockWorldState)(nil).GetGasPrice), chain)
}

// GetNativeTokenPriceForChain mocks base method.
func (m *MockWorldState) GetNativeTokenPriceForChain(chain string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNativeTokenPriceForChain", chain)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNativeTokenPriceForChain indicates an expected call of GetNativeTokenPriceForChain.
func (mr *MockWorldStateMockRecorder) GetNativeTokenPriceForChain(chain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNativeTokenPriceForChain", reflect.TypeOf((*MockWorldState)(nil).GetNativeTokenPriceForChain), chain)
}

// GetTokenFromAddress mocks base method.
func (m *MockWorldState) GetTokenFromAddress(chain, tokenAddr string) *types.Token {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenFromAddress", chain, tokenAddr)
	ret0, _ := ret[0].(*types.Token)
	return ret0
}

// GetTokenFromAddress indicates an expected call of GetTokenFromAddress.
func (mr *MockWorldStateMockRecorder) GetTokenFromAddress(chain, tokenAddr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenFromAddress", reflect.TypeOf((*MockWorldState)(nil).GetTokenFromAddress), chain, tokenAddr)
}

// GetTokenPrice mocks base method.
func (m *MockWorldState) GetTokenPrice(token string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenPrice", token)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenPrice indicates an expected call of GetTokenPrice.
func (mr *MockWorldStateMockRecorder) GetTokenPrice(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenPrice", reflect.TypeOf((*MockWorldState)(nil).GetTokenPrice), token)
}

// LoadData mocks base method.
func (m *MockWorldState) LoadData() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LoadData")
}

// LoadData indicates an expected call of LoadData.
func (mr *MockWorldStateMockRecorder) LoadData() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadData", reflect.TypeOf((*MockWorldState)(nil).LoadData))
}

// SetChain mocks base method.
func (m *MockWorldState) SetChain(chain *types.Chain) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetChain", chain)
}

// SetChain indicates an expected call of SetChain.
func (mr *MockWorldStateMockRecorder) SetChain(chain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetChain", reflect.TypeOf((*MockWorldState)(nil).SetChain), chain)
}

// SetTokens mocks base method.
func (m *MockWorldState) SetTokens(tokenPrices map[string]*types.Token) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTokens", tokenPrices)
}

// SetTokens indicates an expected call of SetTokens.
func (mr *MockWorldStateMockRecorder) SetTokens(tokenPrices interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTokens", reflect.TypeOf((*MockWorldState)(nil).SetTokens), tokenPrices)
}

// UseAndIncreaseNonce mocks base method.
func (m *MockWorldState) UseAndIncreaseNonce(chain string) int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UseAndIncreaseNonce", chain)
	ret0, _ := ret[0].(int64)
	return ret0
}

// UseAndIncreaseNonce indicates an expected call of UseAndIncreaseNonce.
func (mr *MockWorldStateMockRecorder) UseAndIncreaseNonce(chain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UseAndIncreaseNonce", reflect.TypeOf((*MockWorldState)(nil).UseAndIncreaseNonce), chain)
}