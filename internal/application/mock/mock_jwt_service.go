// Code generated by MockGen. DO NOT EDIT.
// Source: internal/application/authentication/jwt_service.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIJWTService is a mock of IJWTService interface.
type MockIJWTService struct {
	ctrl     *gomock.Controller
	recorder *MockIJWTServiceMockRecorder
}

// MockIJWTServiceMockRecorder is the mock recorder for MockIJWTService.
type MockIJWTServiceMockRecorder struct {
	mock *MockIJWTService
}

// NewMockIJWTService creates a new mock instance.
func NewMockIJWTService(ctrl *gomock.Controller) *MockIJWTService {
	mock := &MockIJWTService{ctrl: ctrl}
	mock.recorder = &MockIJWTServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIJWTService) EXPECT() *MockIJWTServiceMockRecorder {
	return m.recorder
}

// GenerateAccessToken mocks base method.
func (m *MockIJWTService) GenerateAccessToken(userID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateAccessToken", userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateAccessToken indicates an expected call of GenerateAccessToken.
func (mr *MockIJWTServiceMockRecorder) GenerateAccessToken(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateAccessToken", reflect.TypeOf((*MockIJWTService)(nil).GenerateAccessToken), userID)
}

// GetUserIDFromAccessToken mocks base method.
func (m *MockIJWTService) GetUserIDFromAccessToken(accessToken string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserIDFromAccessToken", accessToken)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIDFromAccessToken indicates an expected call of GetUserIDFromAccessToken.
func (mr *MockIJWTServiceMockRecorder) GetUserIDFromAccessToken(accessToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIDFromAccessToken", reflect.TypeOf((*MockIJWTService)(nil).GetUserIDFromAccessToken), accessToken)
}
