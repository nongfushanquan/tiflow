// Code generated by MockGen. DO NOT EDIT.
// Source: cdc/owner/owner.go

// Package mock_owner is a generated GoMock package.
package mock_owner

import (
	context "context"
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/pingcap/tiflow/cdc/model"
	owner "github.com/pingcap/tiflow/cdc/owner"
	orchestrator "github.com/pingcap/tiflow/pkg/orchestrator"
)

// MockOwner is a mock of Owner interface.
type MockOwner struct {
	ctrl     *gomock.Controller
	recorder *MockOwnerMockRecorder
}

// MockOwnerMockRecorder is the mock recorder for MockOwner.
type MockOwnerMockRecorder struct {
	mock *MockOwner
}

// NewMockOwner creates a new mock instance.
func NewMockOwner(ctrl *gomock.Controller) *MockOwner {
	mock := &MockOwner{ctrl: ctrl}
	mock.recorder = &MockOwnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOwner) EXPECT() *MockOwnerMockRecorder {
	return m.recorder
}

// AsyncStop mocks base method.
func (m *MockOwner) AsyncStop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AsyncStop")
}

// AsyncStop indicates an expected call of AsyncStop.
func (mr *MockOwnerMockRecorder) AsyncStop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AsyncStop", reflect.TypeOf((*MockOwner)(nil).AsyncStop))
}

// EnqueueJob mocks base method.
func (m *MockOwner) EnqueueJob(adminJob model.AdminJob, done chan<- error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "EnqueueJob", adminJob, done)
}

// EnqueueJob indicates an expected call of EnqueueJob.
func (mr *MockOwnerMockRecorder) EnqueueJob(adminJob, done interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnqueueJob", reflect.TypeOf((*MockOwner)(nil).EnqueueJob), adminJob, done)
}

// Query mocks base method.
func (m *MockOwner) Query(query *owner.Query, done chan<- error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Query", query, done)
}

// Query indicates an expected call of Query.
func (mr *MockOwnerMockRecorder) Query(query, done interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockOwner)(nil).Query), query, done)
}

// RebalanceTables mocks base method.
func (m *MockOwner) RebalanceTables(cfID model.ChangeFeedID, done chan<- error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RebalanceTables", cfID, done)
}

// RebalanceTables indicates an expected call of RebalanceTables.
func (mr *MockOwnerMockRecorder) RebalanceTables(cfID, done interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RebalanceTables", reflect.TypeOf((*MockOwner)(nil).RebalanceTables), cfID, done)
}

// ScheduleTable mocks base method.
func (m *MockOwner) ScheduleTable(cfID model.ChangeFeedID, toCapture model.CaptureID, tableID model.TableID, done chan<- error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ScheduleTable", cfID, toCapture, tableID, done)
}

// ScheduleTable indicates an expected call of ScheduleTable.
func (mr *MockOwnerMockRecorder) ScheduleTable(cfID, toCapture, tableID, done interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduleTable", reflect.TypeOf((*MockOwner)(nil).ScheduleTable), cfID, toCapture, tableID, done)
}

// Tick mocks base method.
func (m *MockOwner) Tick(ctx context.Context, state orchestrator.ReactorState) (orchestrator.ReactorState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tick", ctx, state)
	ret0, _ := ret[0].(orchestrator.ReactorState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Tick indicates an expected call of Tick.
func (mr *MockOwnerMockRecorder) Tick(ctx, state interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tick", reflect.TypeOf((*MockOwner)(nil).Tick), ctx, state)
}

// WriteDebugInfo mocks base method.
func (m *MockOwner) WriteDebugInfo(w io.Writer, done chan<- error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteDebugInfo", w, done)
}

// WriteDebugInfo indicates an expected call of WriteDebugInfo.
func (mr *MockOwnerMockRecorder) WriteDebugInfo(w, done interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteDebugInfo", reflect.TypeOf((*MockOwner)(nil).WriteDebugInfo), w, done)
}
