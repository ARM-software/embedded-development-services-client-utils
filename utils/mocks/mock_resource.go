/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ARM-software/embedded-development-services-client-utils/utils/resource (interfaces: IResource)
//
// Generated by this command:
//
//	mockgen -destination=../mocks/mock_resource.go -package=mocks github.com/ARM-software/embedded-development-services-client-utils/utils/resource IResource
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockIResource is a mock of IResource interface.
type MockIResource struct {
	ctrl     *gomock.Controller
	recorder *MockIResourceMockRecorder
	isgomock struct{}
}

// MockIResourceMockRecorder is the mock recorder for MockIResource.
type MockIResourceMockRecorder struct {
	mock *MockIResource
}

// NewMockIResource creates a new mock instance.
func NewMockIResource(ctrl *gomock.Controller) *MockIResource {
	mock := &MockIResource{ctrl: ctrl}
	mock.recorder = &MockIResourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIResource) EXPECT() *MockIResourceMockRecorder {
	return m.recorder
}

// FetchLinks mocks base method.
func (m *MockIResource) FetchLinks() (any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchLinks")
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchLinks indicates an expected call of FetchLinks.
func (mr *MockIResourceMockRecorder) FetchLinks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchLinks", reflect.TypeOf((*MockIResource)(nil).FetchLinks))
}

// FetchName mocks base method.
func (m *MockIResource) FetchName() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchName")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchName indicates an expected call of FetchName.
func (mr *MockIResourceMockRecorder) FetchName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchName", reflect.TypeOf((*MockIResource)(nil).FetchName))
}

// FetchTitle mocks base method.
func (m *MockIResource) FetchTitle() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchTitle")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchTitle indicates an expected call of FetchTitle.
func (mr *MockIResourceMockRecorder) FetchTitle() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchTitle", reflect.TypeOf((*MockIResource)(nil).FetchTitle))
}

// FetchType mocks base method.
func (m *MockIResource) FetchType() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchType")
	ret0, _ := ret[0].(string)
	return ret0
}

// FetchType indicates an expected call of FetchType.
func (mr *MockIResourceMockRecorder) FetchType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchType", reflect.TypeOf((*MockIResource)(nil).FetchType))
}
