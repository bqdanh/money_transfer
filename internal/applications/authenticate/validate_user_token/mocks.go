// Code generated by MockGen. DO NOT EDIT.
// Source: ./validate_user_token.go

// Package validate_user_token is a generated GoMock package.
package validate_user_token

import (
	reflect "reflect"

	authenticate "github.com/bqdanh/money_transfer/internal/entities/authenticate"
	gomock "github.com/golang/mock/gomock"
)

// MocktokenValidator is a mock of tokenValidator interface.
type MocktokenValidator struct {
	ctrl     *gomock.Controller
	recorder *MocktokenValidatorMockRecorder
}

// MocktokenValidatorMockRecorder is the mock recorder for MocktokenValidator.
type MocktokenValidatorMockRecorder struct {
	mock *MocktokenValidator
}

// NewMocktokenValidator creates a new mock instance.
func NewMocktokenValidator(ctrl *gomock.Controller) *MocktokenValidator {
	mock := &MocktokenValidator{ctrl: ctrl}
	mock.recorder = &MocktokenValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocktokenValidator) EXPECT() *MocktokenValidatorMockRecorder {
	return m.recorder
}

// ValidateToken mocks base method.
func (m *MocktokenValidator) ValidateToken(token string) (authenticate.UserAuthenticateData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", token)
	ret0, _ := ret[0].(authenticate.UserAuthenticateData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateToken indicates an expected call of ValidateToken.
func (mr *MocktokenValidatorMockRecorder) ValidateToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MocktokenValidator)(nil).ValidateToken), token)
}