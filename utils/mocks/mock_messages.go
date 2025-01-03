// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ARM-software/embedded-development-services-client-utils/utils/messages (interfaces: IMessage,IMessageLogger)
//
// Generated by this command:
//
//	mockgen -destination=../mocks/mock_messages.go -package=mocks github.com/ARM-software/embedded-development-services-client-utils/utils/messages IMessage,IMessageLogger
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	messages "github.com/ARM-software/embedded-development-services-client-utils/utils/messages"
	pagination "github.com/ARM-software/golang-utils/utils/collection/pagination"
	gomock "go.uber.org/mock/gomock"
)

// MockIMessage is a mock of IMessage interface.
type MockIMessage struct {
	ctrl     *gomock.Controller
	recorder *MockIMessageMockRecorder
	isgomock struct{}
}

// MockIMessageMockRecorder is the mock recorder for MockIMessage.
type MockIMessageMockRecorder struct {
	mock *MockIMessage
}

// NewMockIMessage creates a new mock instance.
func NewMockIMessage(ctrl *gomock.Controller) *MockIMessage {
	mock := &MockIMessage{ctrl: ctrl}
	mock.recorder = &MockIMessageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIMessage) EXPECT() *MockIMessageMockRecorder {
	return m.recorder
}

// GetCtimeOk mocks base method.
func (m *MockIMessage) GetCtimeOk() (*time.Time, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCtimeOk")
	ret0, _ := ret[0].(*time.Time)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetCtimeOk indicates an expected call of GetCtimeOk.
func (mr *MockIMessageMockRecorder) GetCtimeOk() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCtimeOk", reflect.TypeOf((*MockIMessage)(nil).GetCtimeOk))
}

// GetMessageOk mocks base method.
func (m *MockIMessage) GetMessageOk() (*string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessageOk")
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetMessageOk indicates an expected call of GetMessageOk.
func (mr *MockIMessageMockRecorder) GetMessageOk() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessageOk", reflect.TypeOf((*MockIMessage)(nil).GetMessageOk))
}

// GetSeverityOk mocks base method.
func (m *MockIMessage) GetSeverityOk() (*string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSeverityOk")
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetSeverityOk indicates an expected call of GetSeverityOk.
func (mr *MockIMessageMockRecorder) GetSeverityOk() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSeverityOk", reflect.TypeOf((*MockIMessage)(nil).GetSeverityOk))
}

// GetSourceOk mocks base method.
func (m *MockIMessage) GetSourceOk() (*string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSourceOk")
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetSourceOk indicates an expected call of GetSourceOk.
func (mr *MockIMessageMockRecorder) GetSourceOk() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSourceOk", reflect.TypeOf((*MockIMessage)(nil).GetSourceOk))
}

// MockIMessageLogger is a mock of IMessageLogger interface.
type MockIMessageLogger struct {
	ctrl     *gomock.Controller
	recorder *MockIMessageLoggerMockRecorder
	isgomock struct{}
}

// MockIMessageLoggerMockRecorder is the mock recorder for MockIMessageLogger.
type MockIMessageLoggerMockRecorder struct {
	mock *MockIMessageLogger
}

// NewMockIMessageLogger creates a new mock instance.
func NewMockIMessageLogger(ctrl *gomock.Controller) *MockIMessageLogger {
	mock := &MockIMessageLogger{ctrl: ctrl}
	mock.recorder = &MockIMessageLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIMessageLogger) EXPECT() *MockIMessageLoggerMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockIMessageLogger) Check() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check")
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check.
func (mr *MockIMessageLoggerMockRecorder) Check() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockIMessageLogger)(nil).Check))
}

// Close mocks base method.
func (m *MockIMessageLogger) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockIMessageLoggerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockIMessageLogger)(nil).Close))
}

// Log mocks base method.
func (m *MockIMessageLogger) Log(output ...any) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range output {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Log", varargs...)
}

// Log indicates an expected call of Log.
func (mr *MockIMessageLoggerMockRecorder) Log(output ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Log", reflect.TypeOf((*MockIMessageLogger)(nil).Log), output...)
}

// LogEmptyMessageError mocks base method.
func (m *MockIMessageLogger) LogEmptyMessageError() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LogEmptyMessageError")
}

// LogEmptyMessageError indicates an expected call of LogEmptyMessageError.
func (mr *MockIMessageLoggerMockRecorder) LogEmptyMessageError() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogEmptyMessageError", reflect.TypeOf((*MockIMessageLogger)(nil).LogEmptyMessageError))
}

// LogError mocks base method.
func (m *MockIMessageLogger) LogError(err ...any) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range err {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "LogError", varargs...)
}

// LogError indicates an expected call of LogError.
func (mr *MockIMessageLoggerMockRecorder) LogError(err ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogError", reflect.TypeOf((*MockIMessageLogger)(nil).LogError), err...)
}

// LogMarshallingError mocks base method.
func (m *MockIMessageLogger) LogMarshallingError(rawMessage *any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LogMarshallingError", rawMessage)
}

// LogMarshallingError indicates an expected call of LogMarshallingError.
func (mr *MockIMessageLoggerMockRecorder) LogMarshallingError(rawMessage any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogMarshallingError", reflect.TypeOf((*MockIMessageLogger)(nil).LogMarshallingError), rawMessage)
}

// LogMessage mocks base method.
func (m *MockIMessageLogger) LogMessage(msg messages.IMessage) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LogMessage", msg)
}

// LogMessage indicates an expected call of LogMessage.
func (mr *MockIMessageLoggerMockRecorder) LogMessage(msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogMessage", reflect.TypeOf((*MockIMessageLogger)(nil).LogMessage), msg)
}

// LogMessagesCollection mocks base method.
func (m *MockIMessageLogger) LogMessagesCollection(ctx context.Context, messagePaginator pagination.IGenericPaginator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogMessagesCollection", ctx, messagePaginator)
	ret0, _ := ret[0].(error)
	return ret0
}

// LogMessagesCollection indicates an expected call of LogMessagesCollection.
func (mr *MockIMessageLoggerMockRecorder) LogMessagesCollection(ctx, messagePaginator any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogMessagesCollection", reflect.TypeOf((*MockIMessageLogger)(nil).LogMessagesCollection), ctx, messagePaginator)
}

// SetLogSource mocks base method.
func (m *MockIMessageLogger) SetLogSource(source string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLogSource", source)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetLogSource indicates an expected call of SetLogSource.
func (mr *MockIMessageLoggerMockRecorder) SetLogSource(source any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLogSource", reflect.TypeOf((*MockIMessageLogger)(nil).SetLogSource), source)
}

// SetLoggerSource mocks base method.
func (m *MockIMessageLogger) SetLoggerSource(source string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLoggerSource", source)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetLoggerSource indicates an expected call of SetLoggerSource.
func (mr *MockIMessageLoggerMockRecorder) SetLoggerSource(source any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLoggerSource", reflect.TypeOf((*MockIMessageLogger)(nil).SetLoggerSource), source)
}
