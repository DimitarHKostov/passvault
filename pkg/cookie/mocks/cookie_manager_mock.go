// Code generated by MockGen. DO NOT EDIT.
// Source: cookie_manager_interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCookieManagerInterface is a mock of CookieManagerInterface interface.
type MockCookieManagerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockCookieManagerInterfaceMockRecorder
}

// MockCookieManagerInterfaceMockRecorder is the mock recorder for MockCookieManagerInterface.
type MockCookieManagerInterfaceMockRecorder struct {
	mock *MockCookieManagerInterface
}

// NewMockCookieManagerInterface creates a new mock instance.
func NewMockCookieManagerInterface(ctrl *gomock.Controller) *MockCookieManagerInterface {
	mock := &MockCookieManagerInterface{ctrl: ctrl}
	mock.recorder = &MockCookieManagerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCookieManagerInterface) EXPECT() *MockCookieManagerInterfaceMockRecorder {
	return m.recorder
}

// ProduceCookie mocks base method.
func (m *MockCookieManagerInterface) ProduceCookie() (*http.Cookie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProduceCookie")
	ret0, _ := ret[0].(*http.Cookie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProduceCookie indicates an expected call of ProduceCookie.
func (mr *MockCookieManagerInterfaceMockRecorder) ProduceCookie() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProduceCookie", reflect.TypeOf((*MockCookieManagerInterface)(nil).ProduceCookie))
}