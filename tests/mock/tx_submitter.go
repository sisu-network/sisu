// Code generated by MockGen. DO NOT EDIT.
// Source: common/tx_submitter.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	types "github.com/cosmos/cosmos-sdk/types"
	gomock "github.com/golang/mock/gomock"
)

// MockTxSubmit is a mock of TxSubmit interface.
type MockTxSubmit struct {
	ctrl     *gomock.Controller
	recorder *MockTxSubmitMockRecorder
}

// MockTxSubmitMockRecorder is the mock recorder for MockTxSubmit.
type MockTxSubmitMockRecorder struct {
	mock *MockTxSubmit
}

// NewMockTxSubmit creates a new mock instance.
func NewMockTxSubmit(ctrl *gomock.Controller) *MockTxSubmit {
	mock := &MockTxSubmit{ctrl: ctrl}
	mock.recorder = &MockTxSubmitMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTxSubmit) EXPECT() *MockTxSubmitMockRecorder {
	return m.recorder
}

// SubmitMessageAsync mocks base method.
func (m *MockTxSubmit) SubmitMessageAsync(msg types.Msg) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitMessageAsync", msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// SubmitMessageAsync indicates an expected call of SubmitMessageAsync.
func (mr *MockTxSubmitMockRecorder) SubmitMessageAsync(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitMessageAsync", reflect.TypeOf((*MockTxSubmit)(nil).SubmitMessageAsync), msg)
}

// SubmitMessageSync mocks base method.
func (m *MockTxSubmit) SubmitMessageSync(msg types.Msg) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitMessageSync", msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// SubmitMessageSync indicates an expected call of SubmitMessageSync.
func (mr *MockTxSubmitMockRecorder) SubmitMessageSync(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitMessageSync", reflect.TypeOf((*MockTxSubmit)(nil).SubmitMessageSync), msg)
}
