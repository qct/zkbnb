// Code generated by MockGen. DO NOT EDIT.
// Source: api.go

// Package failtx is a generated GoMock package.
package failtx

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	tx "github.com/bnb-chain/zkbas/common/model/tx"
)

// MockModel is a mock of Model interface.
type MockModel struct {
	ctrl     *gomock.Controller
	recorder *MockModelMockRecorder
}

// MockModelMockRecorder is the mock recorder for MockModel.
type MockModelMockRecorder struct {
	mock *MockModel
}

// NewMockModel creates a new mock instance.
func NewMockModel(ctrl *gomock.Controller) *MockModel {
	mock := &MockModel{ctrl: ctrl}
	mock.recorder = &MockModelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModel) EXPECT() *MockModelMockRecorder {
	return m.recorder
}

// CreateFailTx mocks base method.
func (m *MockModel) CreateFailTx(failTx *tx.FailTx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFailTx", failTx)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateFailTx indicates an expected call of CreateFailTx.
func (mr *MockModelMockRecorder) CreateFailTx(failTx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFailTx", reflect.TypeOf((*MockModel)(nil).CreateFailTx), failTx)
}
