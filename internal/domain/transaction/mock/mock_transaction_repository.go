// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/transaction/transaction_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	transaction "github.com/ucho456job/pocgo/internal/domain/transaction"
)

// MockITransactionRepository is a mock of ITransactionRepository interface.
type MockITransactionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockITransactionRepositoryMockRecorder
}

// MockITransactionRepositoryMockRecorder is the mock recorder for MockITransactionRepository.
type MockITransactionRepositoryMockRecorder struct {
	mock *MockITransactionRepository
}

// NewMockITransactionRepository creates a new mock instance.
func NewMockITransactionRepository(ctrl *gomock.Controller) *MockITransactionRepository {
	mock := &MockITransactionRepository{ctrl: ctrl}
	mock.recorder = &MockITransactionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockITransactionRepository) EXPECT() *MockITransactionRepositoryMockRecorder {
	return m.recorder
}

// ListByAccountID mocks base method.
func (m *MockITransactionRepository) ListByAccountID(accountID string, limit, offset *int) ([]*transaction.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByAccountID", accountID, limit, offset)
	ret0, _ := ret[0].([]*transaction.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByAccountID indicates an expected call of ListByAccountID.
func (mr *MockITransactionRepositoryMockRecorder) ListByAccountID(accountID, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByAccountID", reflect.TypeOf((*MockITransactionRepository)(nil).ListByAccountID), accountID, limit, offset)
}

// Save mocks base method.
func (m *MockITransactionRepository) Save(transaction *transaction.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", transaction)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockITransactionRepositoryMockRecorder) Save(transaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockITransactionRepository)(nil).Save), transaction)
}