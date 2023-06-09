// Code generated by MockGen. DO NOT EDIT.
// Source: jwt_manager_interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	types "passvault/pkg/types"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockJWTManagerInterface is a mock of JWTManagerInterface interface.
type MockJWTManagerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockJWTManagerInterfaceMockRecorder
}

// MockJWTManagerInterfaceMockRecorder is the mock recorder for MockJWTManagerInterface.
type MockJWTManagerInterfaceMockRecorder struct {
	mock *MockJWTManagerInterface
}

// NewMockJWTManagerInterface creates a new mock instance.
func NewMockJWTManagerInterface(ctrl *gomock.Controller) *MockJWTManagerInterface {
	mock := &MockJWTManagerInterface{ctrl: ctrl}
	mock.recorder = &MockJWTManagerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJWTManagerInterface) EXPECT() *MockJWTManagerInterfaceMockRecorder {
	return m.recorder
}

// GenerateToken mocks base method.
func (m *MockJWTManagerInterface) GenerateToken(arg0 time.Duration) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockJWTManagerInterfaceMockRecorder) GenerateToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockJWTManagerInterface)(nil).GenerateToken), arg0)
}

// VerifyToken mocks base method.
func (m *MockJWTManagerInterface) VerifyToken(arg0 string) (*types.Payload, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyToken", arg0)
	ret0, _ := ret[0].(*types.Payload)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyToken indicates an expected call of VerifyToken.
func (mr *MockJWTManagerInterfaceMockRecorder) VerifyToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyToken", reflect.TypeOf((*MockJWTManagerInterface)(nil).VerifyToken), arg0)
}
