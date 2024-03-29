/*
 * Copyright (C) 2020-2024 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ARM-software/embedded-development-services-client-utils/utils/cache (interfaces: IServerCache)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	cache "github.com/ARM-software/embedded-development-services-client-utils/utils/cache"
)

// MockIServerCache is a mock of IServerCache interface.
type MockIServerCache struct {
	ctrl     *gomock.Controller
	recorder *MockIServerCacheMockRecorder
}

// MockIServerCacheMockRecorder is the mock recorder for MockIServerCache.
type MockIServerCacheMockRecorder struct {
	mock *MockIServerCache
}

// NewMockIServerCache creates a new mock instance.
func NewMockIServerCache(ctrl *gomock.Controller) *MockIServerCache {
	mock := &MockIServerCache{ctrl: ctrl}
	mock.recorder = &MockIServerCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIServerCache) EXPECT() *MockIServerCacheMockRecorder {
	return m.recorder
}

// GetCacheControl mocks base method.
func (m *MockIServerCache) GetCacheControl() cache.Control {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCacheControl")
	ret0, _ := ret[0].(cache.Control)
	return ret0
}

// GetCacheControl indicates an expected call of GetCacheControl.
func (mr *MockIServerCacheMockRecorder) GetCacheControl() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCacheControl", reflect.TypeOf((*MockIServerCache)(nil).GetCacheControl))
}

// GetKey mocks base method.
func (m *MockIServerCache) GetKey() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKey")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetKey indicates an expected call of GetKey.
func (mr *MockIServerCacheMockRecorder) GetKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKey", reflect.TypeOf((*MockIServerCache)(nil).GetKey))
}

// Invalidate mocks base method.
func (m *MockIServerCache) Invalidate(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Invalidate", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Invalidate indicates an expected call of Invalidate.
func (mr *MockIServerCacheMockRecorder) Invalidate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Invalidate", reflect.TypeOf((*MockIServerCache)(nil).Invalidate), arg0)
}

// SetCacheControl mocks base method.
func (m *MockIServerCache) SetCacheControl(arg0 cache.Control) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCacheControl", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCacheControl indicates an expected call of SetCacheControl.
func (mr *MockIServerCacheMockRecorder) SetCacheControl(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCacheControl", reflect.TypeOf((*MockIServerCache)(nil).SetCacheControl), arg0)
}

// SetKey mocks base method.
func (m *MockIServerCache) SetKey(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetKey", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetKey indicates an expected call of SetKey.
func (mr *MockIServerCacheMockRecorder) SetKey(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetKey", reflect.TypeOf((*MockIServerCache)(nil).SetKey), arg0)
}
