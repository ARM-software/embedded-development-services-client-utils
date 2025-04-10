/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ARM-software/embedded-development-services-client-utils/utils/artefacts (interfaces: IArtefactManager)
//
// Generated by this command:
//
//	mockgen -destination=../mocks/mock_artefacts.go -package=mocks github.com/ARM-software/embedded-development-services-client-utils/utils/artefacts IArtefactManager
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	client "github.com/ARM-software/embedded-development-services-client/client"
	pagination "github.com/ARM-software/golang-utils/utils/collection/pagination"
	gomock "go.uber.org/mock/gomock"
)

// MockIArtefactManager is a mock of IArtefactManager interface.
type MockIArtefactManager struct {
	ctrl     *gomock.Controller
	recorder *MockIArtefactManagerMockRecorder
	isgomock struct{}
}

// MockIArtefactManagerMockRecorder is the mock recorder for MockIArtefactManager.
type MockIArtefactManagerMockRecorder struct {
	mock *MockIArtefactManager
}

// NewMockIArtefactManager creates a new mock instance.
func NewMockIArtefactManager(ctrl *gomock.Controller) *MockIArtefactManager {
	mock := &MockIArtefactManager{ctrl: ctrl}
	mock.recorder = &MockIArtefactManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIArtefactManager) EXPECT() *MockIArtefactManagerMockRecorder {
	return m.recorder
}

// DownloadAllJobArtefacts mocks base method.
func (m *MockIArtefactManager) DownloadAllJobArtefacts(ctx context.Context, jobName, outputDirectory string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadAllJobArtefacts", ctx, jobName, outputDirectory)
	ret0, _ := ret[0].(error)
	return ret0
}

// DownloadAllJobArtefacts indicates an expected call of DownloadAllJobArtefacts.
func (mr *MockIArtefactManagerMockRecorder) DownloadAllJobArtefacts(ctx, jobName, outputDirectory any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadAllJobArtefacts", reflect.TypeOf((*MockIArtefactManager)(nil).DownloadAllJobArtefacts), ctx, jobName, outputDirectory)
}

// DownloadAllJobArtefactsWithTree mocks base method.
func (m *MockIArtefactManager) DownloadAllJobArtefactsWithTree(ctx context.Context, jobName string, maintainTreeStructure bool, outputDirectory string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadAllJobArtefactsWithTree", ctx, jobName, maintainTreeStructure, outputDirectory)
	ret0, _ := ret[0].(error)
	return ret0
}

// DownloadAllJobArtefactsWithTree indicates an expected call of DownloadAllJobArtefactsWithTree.
func (mr *MockIArtefactManagerMockRecorder) DownloadAllJobArtefactsWithTree(ctx, jobName, maintainTreeStructure, outputDirectory any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadAllJobArtefactsWithTree", reflect.TypeOf((*MockIArtefactManager)(nil).DownloadAllJobArtefactsWithTree), ctx, jobName, maintainTreeStructure, outputDirectory)
}

// DownloadJobArtefact mocks base method.
func (m *MockIArtefactManager) DownloadJobArtefact(ctx context.Context, jobName, outputDirectory string, artefactManager *client.ArtefactManagerItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadJobArtefact", ctx, jobName, outputDirectory, artefactManager)
	ret0, _ := ret[0].(error)
	return ret0
}

// DownloadJobArtefact indicates an expected call of DownloadJobArtefact.
func (mr *MockIArtefactManagerMockRecorder) DownloadJobArtefact(ctx, jobName, outputDirectory, artefactManager any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadJobArtefact", reflect.TypeOf((*MockIArtefactManager)(nil).DownloadJobArtefact), ctx, jobName, outputDirectory, artefactManager)
}

// DownloadJobArtefactFromLink mocks base method.
func (m *MockIArtefactManager) DownloadJobArtefactFromLink(ctx context.Context, jobName, outputDirectory string, artefactManagerItemLink *client.HalLinkData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadJobArtefactFromLink", ctx, jobName, outputDirectory, artefactManagerItemLink)
	ret0, _ := ret[0].(error)
	return ret0
}

// DownloadJobArtefactFromLink indicates an expected call of DownloadJobArtefactFromLink.
func (mr *MockIArtefactManagerMockRecorder) DownloadJobArtefactFromLink(ctx, jobName, outputDirectory, artefactManagerItemLink any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadJobArtefactFromLink", reflect.TypeOf((*MockIArtefactManager)(nil).DownloadJobArtefactFromLink), ctx, jobName, outputDirectory, artefactManagerItemLink)
}

// DownloadJobArtefactFromLinkWithTree mocks base method.
func (m *MockIArtefactManager) DownloadJobArtefactFromLinkWithTree(ctx context.Context, jobName string, maintainTreeLocation bool, outputDirectory string, artefactManagerItemLink *client.HalLinkData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadJobArtefactFromLinkWithTree", ctx, jobName, maintainTreeLocation, outputDirectory, artefactManagerItemLink)
	ret0, _ := ret[0].(error)
	return ret0
}

// DownloadJobArtefactFromLinkWithTree indicates an expected call of DownloadJobArtefactFromLinkWithTree.
func (mr *MockIArtefactManagerMockRecorder) DownloadJobArtefactFromLinkWithTree(ctx, jobName, maintainTreeLocation, outputDirectory, artefactManagerItemLink any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadJobArtefactFromLinkWithTree", reflect.TypeOf((*MockIArtefactManager)(nil).DownloadJobArtefactFromLinkWithTree), ctx, jobName, maintainTreeLocation, outputDirectory, artefactManagerItemLink)
}

// DownloadJobArtefactWithTree mocks base method.
func (m *MockIArtefactManager) DownloadJobArtefactWithTree(ctx context.Context, jobName string, maintainTreeLocation bool, outputDirectory string, artefactManager *client.ArtefactManagerItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadJobArtefactWithTree", ctx, jobName, maintainTreeLocation, outputDirectory, artefactManager)
	ret0, _ := ret[0].(error)
	return ret0
}

// DownloadJobArtefactWithTree indicates an expected call of DownloadJobArtefactWithTree.
func (mr *MockIArtefactManagerMockRecorder) DownloadJobArtefactWithTree(ctx, jobName, maintainTreeLocation, outputDirectory, artefactManager any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadJobArtefactWithTree", reflect.TypeOf((*MockIArtefactManager)(nil).DownloadJobArtefactWithTree), ctx, jobName, maintainTreeLocation, outputDirectory, artefactManager)
}

// ListJobArtefacts mocks base method.
func (m *MockIArtefactManager) ListJobArtefacts(ctx context.Context, jobName string) (pagination.IPaginatorAndPageFetcher, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListJobArtefacts", ctx, jobName)
	ret0, _ := ret[0].(pagination.IPaginatorAndPageFetcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListJobArtefacts indicates an expected call of ListJobArtefacts.
func (mr *MockIArtefactManagerMockRecorder) ListJobArtefacts(ctx, jobName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListJobArtefacts", reflect.TypeOf((*MockIArtefactManager)(nil).ListJobArtefacts), ctx, jobName)
}
