// Code generated by MockGen. DO NOT EDIT.
// Source: C:/Users/lette/OneDrive/Документы/Projects/ygo/go-musthave-diploma-tpl/internal/app/interfaces/user.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/MrTomSawyer/loyalty-system/internal/app/models"
	gomock "github.com/golang/mock/gomock"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUserService) CreateUser(user models.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserServiceMockRecorder) CreateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserService)(nil).CreateUser), user)
}

// GetUserBalance mocks base method.
func (m *MockUserService) GetUserBalance(userID int) (models.Balance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBalance", userID)
	ret0, _ := ret[0].(models.Balance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserBalance indicates an expected call of GetUserBalance.
func (mr *MockUserServiceMockRecorder) GetUserBalance(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBalance", reflect.TypeOf((*MockUserService)(nil).GetUserBalance), userID)
}

// Login mocks base method.
func (m *MockUserService) Login(user models.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserServiceMockRecorder) Login(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserService)(nil).Login), user)
}