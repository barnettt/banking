// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/barnettt/banking/service (interfaces: CustomerService)

// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"

	exceptions "github.com/barnettt/banking-lib/exceptions"
	dto "github.com/barnettt/banking/dto"
	gomock "github.com/golang/mock/gomock"
)

// MockCustomerService is a mock of CustomerService interface.
type MockCustomerService struct {
	ctrl     *gomock.Controller
	recorder *MockCustomerServiceMockRecorder
}

// MockCustomerServiceMockRecorder is the mock recorder for MockCustomerService.
type MockCustomerServiceMockRecorder struct {
	mock *MockCustomerService
}

// NewMockCustomerService creates a new mock instance.
func NewMockCustomerService(ctrl *gomock.Controller) *MockCustomerService {
	mock := &MockCustomerService{ctrl: ctrl}
	mock.recorder = &MockCustomerServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCustomerService) EXPECT() *MockCustomerServiceMockRecorder {
	return m.recorder
}

// GetAllCustomers mocks base method.
func (m *MockCustomerService) GetAllCustomers() ([]dto.CustomerResponse, *exceptions.AppError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllCustomers")
	ret0, _ := ret[0].([]dto.CustomerResponse)
	ret1, _ := ret[1].(*exceptions.AppError)
	return ret0, ret1
}

// GetAllCustomers indicates an expected call of GetAllCustomers.
func (mr *MockCustomerServiceMockRecorder) GetAllCustomers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllCustomers", reflect.TypeOf((*MockCustomerService)(nil).GetAllCustomers))
}

// GetCustomer mocks base method.
func (m *MockCustomerService) GetCustomer(arg0 string) (*dto.CustomerResponse, *exceptions.AppError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCustomer", arg0)
	ret0, _ := ret[0].(*dto.CustomerResponse)
	ret1, _ := ret[1].(*exceptions.AppError)
	return ret0, ret1
}

// GetCustomer indicates an expected call of GetCustomer.
func (mr *MockCustomerServiceMockRecorder) GetCustomer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCustomer", reflect.TypeOf((*MockCustomerService)(nil).GetCustomer), arg0)
}

// GetCustomersByStatus mocks base method.
func (m *MockCustomerService) GetCustomersByStatus(arg0 string) ([]dto.CustomerResponse, *exceptions.AppError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCustomersByStatus", arg0)
	ret0, _ := ret[0].([]dto.CustomerResponse)
	ret1, _ := ret[1].(*exceptions.AppError)
	return ret0, ret1
}

// GetCustomersByStatus indicates an expected call of GetCustomersByStatus.
func (mr *MockCustomerServiceMockRecorder) GetCustomersByStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCustomersByStatus", reflect.TypeOf((*MockCustomerService)(nil).GetCustomersByStatus), arg0)
}
