// Code generated by MockGen. DO NOT EDIT.
// Source: x/sisu/posted_message_manager.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	types "github.com/cosmos/cosmos-sdk/types"
	gomock "github.com/golang/mock/gomock"
)

// MockPostedMessageManager is a mock of PostedMessageManager interface.
type MockPostedMessageManager struct {
	ctrl     *gomock.Controller
	recorder *MockPostedMessageManagerMockRecorder
}

// MockPostedMessageManagerMockRecorder is the mock recorder for MockPostedMessageManager.
type MockPostedMessageManagerMockRecorder struct {
	mock *MockPostedMessageManager
}

// NewMockPostedMessageManager creates a new mock instance.
func NewMockPostedMessageManager(ctrl *gomock.Controller) *MockPostedMessageManager {
	mock := &MockPostedMessageManager{ctrl: ctrl}
	mock.recorder = &MockPostedMessageManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostedMessageManager) EXPECT() *MockPostedMessageManagerMockRecorder {
	return m.recorder
}

// PreProcessingMsg mocks base method.
func (m *MockPostedMessageManager) PreProcessingMsg(ctx types.Context, msg types.Msg) (bool, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PreProcessingMsg", ctx, msg)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PreProcessingMsg indicates an expected call of PreProcessingMsg.
func (mr *MockPostedMessageManagerMockRecorder) PreProcessingMsg(ctx, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PreProcessingMsg", reflect.TypeOf((*MockPostedMessageManager)(nil).PreProcessingMsg), ctx, msg)
}
