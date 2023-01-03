/*
 * Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ARM-software/embedded-development-services-client-utils/utils/job (interfaces: IAsynchronousJob,IJobManager)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	job "github.com/ARM-software/embedded-development-services-client-utils/utils/job"
)

// MockIAsynchronousJob is a mock of IAsynchronousJob interface.
type MockIAsynchronousJob struct {
	ctrl     *gomock.Controller
	recorder *MockIAsynchronousJobMockRecorder
}

// MockIAsynchronousJobMockRecorder is the mock recorder for MockIAsynchronousJob.
type MockIAsynchronousJobMockRecorder struct {
	mock *MockIAsynchronousJob
}

// NewMockIAsynchronousJob creates a new mock instance.
func NewMockIAsynchronousJob(ctrl *gomock.Controller) *MockIAsynchronousJob {
	mock := &MockIAsynchronousJob{ctrl: ctrl}
	mock.recorder = &MockIAsynchronousJobMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAsynchronousJob) EXPECT() *MockIAsynchronousJobMockRecorder {
	return m.recorder
}

// FetchLinks mocks base method.
func (m *MockIAsynchronousJob) FetchLinks() (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchLinks")
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchLinks indicates an expected call of FetchLinks.
func (mr *MockIAsynchronousJobMockRecorder) FetchLinks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchLinks", reflect.TypeOf((*MockIAsynchronousJob)(nil).FetchLinks))
}

// FetchName mocks base method.
func (m *MockIAsynchronousJob) FetchName() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchName")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchName indicates an expected call of FetchName.
func (mr *MockIAsynchronousJobMockRecorder) FetchName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchName", reflect.TypeOf((*MockIAsynchronousJob)(nil).FetchName))
}

// FetchTitle mocks base method.
func (m *MockIAsynchronousJob) FetchTitle() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchTitle")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchTitle indicates an expected call of FetchTitle.
func (mr *MockIAsynchronousJobMockRecorder) FetchTitle() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchTitle", reflect.TypeOf((*MockIAsynchronousJob)(nil).FetchTitle))
}

// FetchType mocks base method.
func (m *MockIAsynchronousJob) FetchType() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchType")
	ret0, _ := ret[0].(string)
	return ret0
}

// FetchType indicates an expected call of FetchType.
func (mr *MockIAsynchronousJobMockRecorder) FetchType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchType", reflect.TypeOf((*MockIAsynchronousJob)(nil).FetchType))
}

// GetDone mocks base method.
func (m *MockIAsynchronousJob) GetDone() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDone")
	ret0, _ := ret[0].(bool)
	return ret0
}

// GetDone indicates an expected call of GetDone.
func (mr *MockIAsynchronousJobMockRecorder) GetDone() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDone", reflect.TypeOf((*MockIAsynchronousJob)(nil).GetDone))
}

// GetError mocks base method.
func (m *MockIAsynchronousJob) GetError() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetError")
	ret0, _ := ret[0].(bool)
	return ret0
}

// GetError indicates an expected call of GetError.
func (mr *MockIAsynchronousJobMockRecorder) GetError() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetError", reflect.TypeOf((*MockIAsynchronousJob)(nil).GetError))
}

// GetFailure mocks base method.
func (m *MockIAsynchronousJob) GetFailure() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFailure")
	ret0, _ := ret[0].(bool)
	return ret0
}

// GetFailure indicates an expected call of GetFailure.
func (mr *MockIAsynchronousJobMockRecorder) GetFailure() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFailure", reflect.TypeOf((*MockIAsynchronousJob)(nil).GetFailure))
}

// GetStatus mocks base method.
func (m *MockIAsynchronousJob) GetStatus() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatus")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetStatus indicates an expected call of GetStatus.
func (mr *MockIAsynchronousJobMockRecorder) GetStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatus", reflect.TypeOf((*MockIAsynchronousJob)(nil).GetStatus))
}

// GetSuccess mocks base method.
func (m *MockIAsynchronousJob) GetSuccess() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSuccess")
	ret0, _ := ret[0].(bool)
	return ret0
}

// GetSuccess indicates an expected call of GetSuccess.
func (mr *MockIAsynchronousJobMockRecorder) GetSuccess() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSuccess", reflect.TypeOf((*MockIAsynchronousJob)(nil).GetSuccess))
}

// MockIJobManager is a mock of IJobManager interface.
type MockIJobManager struct {
	ctrl     *gomock.Controller
	recorder *MockIJobManagerMockRecorder
}

// MockIJobManagerMockRecorder is the mock recorder for MockIJobManager.
type MockIJobManagerMockRecorder struct {
	mock *MockIJobManager
}

// NewMockIJobManager creates a new mock instance.
func NewMockIJobManager(ctrl *gomock.Controller) *MockIJobManager {
	mock := &MockIJobManager{ctrl: ctrl}
	mock.recorder = &MockIJobManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIJobManager) EXPECT() *MockIJobManagerMockRecorder {
	return m.recorder
}

// HasJobCompleted mocks base method.
func (m *MockIJobManager) HasJobCompleted(arg0 context.Context, arg1 job.IAsynchronousJob) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasJobCompleted", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasJobCompleted indicates an expected call of HasJobCompleted.
func (mr *MockIJobManagerMockRecorder) HasJobCompleted(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasJobCompleted", reflect.TypeOf((*MockIJobManager)(nil).HasJobCompleted), arg0, arg1)
}

// WaitForJobCompletion mocks base method.
func (m *MockIJobManager) WaitForJobCompletion(arg0 context.Context, arg1 job.IAsynchronousJob) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitForJobCompletion", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitForJobCompletion indicates an expected call of WaitForJobCompletion.
func (mr *MockIJobManagerMockRecorder) WaitForJobCompletion(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForJobCompletion", reflect.TypeOf((*MockIJobManager)(nil).WaitForJobCompletion), arg0, arg1)
}
