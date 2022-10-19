// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/anytypeio/go-anytype-infrastructure-experiments/common/commonspace/synctree/updatelistener (interfaces: UpdateListener)

// Package mock_updatelistener is a generated GoMock package.
package mock_updatelistener

import (
	reflect "reflect"

	tree "github.com/anytypeio/go-anytype-infrastructure-experiments/common/pkg/acl/tree"
	gomock "github.com/golang/mock/gomock"
)

// MockUpdateListener is a mock of UpdateListener interface.
type MockUpdateListener struct {
	ctrl     *gomock.Controller
	recorder *MockUpdateListenerMockRecorder
}

// MockUpdateListenerMockRecorder is the mock recorder for MockUpdateListener.
type MockUpdateListenerMockRecorder struct {
	mock *MockUpdateListener
}

// NewMockUpdateListener creates a new mock instance.
func NewMockUpdateListener(ctrl *gomock.Controller) *MockUpdateListener {
	mock := &MockUpdateListener{ctrl: ctrl}
	mock.recorder = &MockUpdateListenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUpdateListener) EXPECT() *MockUpdateListenerMockRecorder {
	return m.recorder
}

// Rebuild mocks base method.
func (m *MockUpdateListener) Rebuild(arg0 tree.ObjectTree) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Rebuild", arg0)
}

// Rebuild indicates an expected call of Rebuild.
func (mr *MockUpdateListenerMockRecorder) Rebuild(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rebuild", reflect.TypeOf((*MockUpdateListener)(nil).Rebuild), arg0)
}

// Update mocks base method.
func (m *MockUpdateListener) Update(arg0 tree.ObjectTree) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Update", arg0)
}

// Update indicates an expected call of Update.
func (mr *MockUpdateListenerMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUpdateListener)(nil).Update), arg0)
}