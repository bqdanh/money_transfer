// Code generated by MockGen. DO NOT EDIT.
// Source: link_bank_account.go

// Package link_account is a generated GoMock package.
package link_account

import (
	context "context"
	reflect "reflect"

	account "github.com/bqdanh/money_transfer/internal/entities/account"
	gomock "github.com/golang/mock/gomock"
)

// MockaccountRepository is a mock of accountRepository interface.
type MockaccountRepository struct {
	ctrl     *gomock.Controller
	recorder *MockaccountRepositoryMockRecorder
}

// MockaccountRepositoryMockRecorder is the mock recorder for MockaccountRepository.
type MockaccountRepositoryMockRecorder struct {
	mock *MockaccountRepository
}

// NewMockaccountRepository creates a new mock instance.
func NewMockaccountRepository(ctrl *gomock.Controller) *MockaccountRepository {
	mock := &MockaccountRepository{ctrl: ctrl}
	mock.recorder = &MockaccountRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockaccountRepository) EXPECT() *MockaccountRepositoryMockRecorder {
	return m.recorder
}

// CreateAccount mocks base method.
func (m *MockaccountRepository) CreateAccount(ctx context.Context, a account.Account) (account.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", ctx, a)
	ret0, _ := ret[0].(account.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockaccountRepositoryMockRecorder) CreateAccount(ctx, a interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockaccountRepository)(nil).CreateAccount), ctx, a)
}

// GetAccountsByUserID mocks base method.
func (m *MockaccountRepository) GetAccountsByUserID(ctx context.Context, userID int64) ([]account.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountsByUserID", ctx, userID)
	ret0, _ := ret[0].([]account.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountsByUserID indicates an expected call of GetAccountsByUserID.
func (mr *MockaccountRepositoryMockRecorder) GetAccountsByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountsByUserID", reflect.TypeOf((*MockaccountRepository)(nil).GetAccountsByUserID), ctx, userID)
}

// MockdistributeLock is a mock of distributeLock interface.
type MockdistributeLock struct {
	ctrl     *gomock.Controller
	recorder *MockdistributeLockMockRecorder
}

// MockdistributeLockMockRecorder is the mock recorder for MockdistributeLock.
type MockdistributeLockMockRecorder struct {
	mock *MockdistributeLock
}

// NewMockdistributeLock creates a new mock instance.
func NewMockdistributeLock(ctrl *gomock.Controller) *MockdistributeLock {
	mock := &MockdistributeLock{ctrl: ctrl}
	mock.recorder = &MockdistributeLockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockdistributeLock) EXPECT() *MockdistributeLockMockRecorder {
	return m.recorder
}

// AcquireCreateAccountLockByUserID mocks base method.
func (m *MockdistributeLock) AcquireCreateAccountLockByUserID(ctx context.Context, userID int64) (func(), error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcquireCreateAccountLockByUserID", ctx, userID)
	ret0, _ := ret[0].(func())
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AcquireCreateAccountLockByUserID indicates an expected call of AcquireCreateAccountLockByUserID.
func (mr *MockdistributeLockMockRecorder) AcquireCreateAccountLockByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcquireCreateAccountLockByUserID", reflect.TypeOf((*MockdistributeLock)(nil).AcquireCreateAccountLockByUserID), ctx, userID)
}