/*
 * Copyright (C) 2020-2022 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package messages

import (
	"context"
	"fmt"
	"time"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/logging"
	"github.com/ARM-software/embedded-development-services-client/client"
	"github.com/ARM-software/golang-utils/utils/collection/pagination"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/logs"
	"github.com/ARM-software/golang-utils/utils/parallelisation"
	"github.com/ARM-software/golang-utils/utils/reflection"
)

const (
	// SleepBetweenMessages describes the frequency of message printing
	SleepBetweenMessages = 100 * time.Millisecond
	// SleepAtEnd describes the grace period which should happen to wait for message stream exhaustion.
	SleepAtEnd = 10 * SleepBetweenMessages
	//
	MessageBackOff = 10 * time.Millisecond

	messageBufferSize = 200
)

type Logger struct {
	logger  logging.ILogger
	printer logs.WriterWithSource
}

func (l *Logger) Close() error {
	return l.printer.Close()
}

func (l *Logger) SetSource(source string) error {
	return l.printer.SetSource(source)
}

func (l *Logger) LogMessage(msg *client.MessageObject) {
	m, err := FormatMessage(msg)
	if err != nil {
		l.logger.LogErrorAndMessage(err, "failed logging message")
		return
	}
	_, err = l.printer.Write([]byte(m))
	if err != nil {
		l.logger.LogErrorAndMessage(err, "failed logging message")
		return
	}
}

func (l *Logger) LogEmptyMessageError() {
	l.logger.LogErrorAndMessage(commonerrors.ErrEmpty, "empty message")
}

func (l *Logger) LogMarshallingError(rawMessage *any) {
	if reflection.IsEmpty(rawMessage) {
		l.LogEmptyMessageError()
	} else {
		l.logger.LogErrorAndMessage(commonerrors.ErrMarshalling, "message [%v] could not be marshalled", *rawMessage)
	}
}

func (l *Logger) LogMessagesCollection(ctx context.Context, messagePaginator pagination.IGenericPaginator) error {
	for {
		err := parallelisation.DetermineContextError(ctx)
		if err != nil {
			return err
		}
		if messagePaginator == nil {
			return fmt.Errorf("%w: missing paginator", commonerrors.ErrUndefined)
		}
		if !messagePaginator.HasNext() {
			return nil
		}
		m, err := messagePaginator.GetNext()
		if err != nil {
			return err
		}
		if reflection.IsEmpty(m) {
			l.LogEmptyMessageError()
		} else {
			if messageItem, ok := (*m).(client.MessageObject); ok {
				l.LogMessage(&messageItem)
			} else {
				l.LogMarshallingError(m)
			}
		}
	}
}

// NewMessageLogger creates a logger for logging messages.
func NewMessageLogger(rawLogger logging.ILogger) (*Logger, error) {
	if rawLogger == nil {
		return nil, commonerrors.ErrNoLogger
	}
	return &Logger{
		logger:  rawLogger,
		printer: newMessagePrinter(rawLogger),
	}, nil
}

// messagePrinter will print messages at a regular time disconnected from the paging mechanism.
type messagePrinter struct {
	logger logging.ILogger
}

func (m *messagePrinter) Write(p []byte) (n int, err error) {
	m.logger.Log(string(p))
	return len(p), err
}

func (m *messagePrinter) Close() error {
	return nil
}

func (m *messagePrinter) SetSource(source string) error {
	return m.logger.SetLoggerSource(source)
}

func newMessagePrinter(logger logging.ILogger) logs.WriterWithSource {
	return logs.NewDiodeWriterForSlowWriter(&messagePrinter{
		logger: logger,
	}, messageBufferSize, SleepBetweenMessages, logger)
}
