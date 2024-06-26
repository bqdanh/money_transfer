// Code generated by MockGen. DO NOT EDIT.
// Source: ./login.go

// Package login is a generated GoMock package.
package login

import (
	context "context"
	reflect "reflect"

	user "github.com/bqdanh/money_transfer/internal/entities/user"
	gomock "github.com/golang/mock/gomock"
)

// MockuserNamePasswordValidator is a mock of userNamePasswordValidator interface.
type MockuserNamePasswordValidator struct {
	ctrl     *gomock.Controller
	recorder *MockuserNamePasswordValidatorMockRecorder
}

// MockuserNamePasswordValidatorMockRecorder is the mock recorder for MockuserNamePasswordValidator.
type MockuserNamePasswordValidatorMockRecorder struct {
	mock *MockuserNamePasswordValidator
}

// NewMockuserNamePasswordValidator creates a new mock instance.
func NewMockuserNamePasswordValidator(ctrl *gomock.Controller) *MockuserNamePasswordValidator {
	mock := &MockuserNamePasswordValidator{ctrl: ctrl}
	mock.recorder = &MockuserNamePasswordValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockuserNamePasswordValidator) EXPECT() *MockuserNamePasswordValidatorMockRecorder {
	return m.recorder
}

// ValidateUserNamePassword mocks base method.
func (m *MockuserNamePasswordValidator) ValidateUserNamePassword(ctx context.Context, username, password string) (user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateUserNamePassword", ctx, username, password)
	ret0, _ := ret[0].(user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateUserNamePassword indicates an expected call of ValidateUserNamePassword.
func (mr *MockuserNamePasswordValidatorMockRecorder) ValidateUserNamePassword(ctx, username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateUserNamePassword", reflect.TypeOf((*MockuserNamePasswordValidator)(nil).ValidateUserNamePassword), ctx, username, password)
}
