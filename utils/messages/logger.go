/*
 * Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package messages

import (
	"context"
	"fmt"
	"time"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/logging"
	"github.com/ARM-software/golang-utils/utils/collection/pagination"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/logs"
	"github.com/ARM-software/golang-utils/utils/parallelisation"
	"github.com/ARM-software/golang-utils/utils/reflection"
)

const (
	// SleepBetweenMessages describes the frequency of message printing
	SleepBetweenMessages = 100 * time.Millisecond
	// SleepAtEnd describes the grace period which should happen when expecting message stream exhaustion.
	SleepAtEnd = 10 * SleepBetweenMessages

	MessageBackOff = 10 * time.Millisecond

	messageBufferSize = 1000
)

type Logger struct {
	rawLogger logging.ILogger
	printer   logs.WriterWithSource
}

func (l *Logger) Check() error {
	return l.rawLogger.Check()
}

func (l *Logger) SetLogSource(source string) error {
	return l.rawLogger.SetLogSource(source)
}

func (l *Logger) SetLoggerSource(source string) error {
	return l.rawLogger.SetLoggerSource(source)
}

func (l *Logger) Log(output ...interface{}) {
	l.rawLogger.Log(output...)
}

func (l *Logger) LogError(err ...interface{}) {
	l.rawLogger.LogError(err...)
}

func (l *Logger) Close() error {
	return l.printer.Close()
}

func (l *Logger) SetSource(source string) error {
	return l.printer.SetSource(source)
}

func (l *Logger) LogMessage(msg IMessage) {
	m, err := FormatMessage(msg)
	if err != nil {
		l.rawLogger.LogErrorAndMessage(err, "failed logging message")
		return
	}
	_, err = l.printer.Write([]byte(m))
	if err != nil {
		l.rawLogger.LogErrorAndMessage(err, "failed logging message")
		return
	}
}

func (l *Logger) LogEmptyMessageError() {
	l.rawLogger.LogErrorAndMessage(commonerrors.ErrEmpty, "empty message")
}

func (l *Logger) LogMarshallingError(rawMessage *any) {
	if reflection.IsEmpty(rawMessage) {
		l.LogEmptyMessageError()
	} else {
		l.rawLogger.LogErrorAndMessage(commonerrors.ErrMarshalling, "message [%v] of type [%T] could not be marshalled", *rawMessage, *rawMessage)
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
		messageItem, err := convertRawMessageIntoIMessage(m)
		if err != nil {
			l.LogMarshallingError(m)
		} else {
			l.LogMessage(messageItem)
		}
	}
}

// NewBasicAsynchronousMessageLogger creates an asynchronous logger for messages which prints them as they come.
func NewBasicAsynchronousMessageLogger(ctx context.Context, rawLogger logging.ILogger) (IMessageLogger, error) {
	if rawLogger == nil {
		return nil, commonerrors.ErrNoLogger
	}
	return &Logger{
		rawLogger: rawLogger,
		printer:   newBasicAsynchronousMessagePrinter(ctx, rawLogger),
	}, nil
}

// NewPeriodicAsynchronousMessageLogger creates an asynchronous logger for messages which prints them at regular intervals.
func NewPeriodicAsynchronousMessageLogger(ctx context.Context, rawLogger logging.ILogger, printPeriod time.Duration) (IMessageLogger, error) {
	if rawLogger == nil {
		return nil, commonerrors.ErrNoLogger
	}
	return &Logger{
		rawLogger: rawLogger,
		printer:   newPeriodicAsynchronousMessagePrinter(ctx, rawLogger, printPeriod),
	}, nil
}

// NewBasicSynchronousMessageLogger creates a synchronous logger for messages which prints them as they come.
func NewBasicSynchronousMessageLogger(ctx context.Context, rawLogger logging.ILogger) (IMessageLogger, error) {
	if rawLogger == nil {
		return nil, commonerrors.ErrNoLogger
	}
	return &Logger{
		rawLogger: rawLogger,
		printer:   newBasicSynchronousMessagePrinter(ctx, rawLogger),
	}, nil
}

// NewPeriodicSynchronousMessageLogger creates a synchronous logger for messages which prints them at regular intervals.
func NewPeriodicSynchronousMessageLogger(ctx context.Context, rawLogger logging.ILogger, printPeriod time.Duration) (IMessageLogger, error) {
	if rawLogger == nil {
		return nil, commonerrors.ErrNoLogger
	}
	return &Logger{
		rawLogger: rawLogger,
		printer:   newPeriodicSynchronousMessagePrinter(ctx, rawLogger, printPeriod),
	}, nil
}

// basicMessagePrinter will print messages as they come.
type basicMessagePrinter struct {
	logger logging.ILogger
	ctx    context.Context
}

func (m *basicMessagePrinter) Write(p []byte) (n int, err error) {
	err = parallelisation.DetermineContextError(m.ctx)
	if err != nil {
		return
	}
	m.logger.Log(string(p))
	return len(p), err
}

func (m *basicMessagePrinter) Close() error {
	return nil
}

func (m *basicMessagePrinter) SetSource(source string) error {
	return m.logger.SetLoggerSource(source)
}

// periodicMessagePrinter will print messages at a regular intervals.
type periodicMessagePrinter struct {
	basicMessagePrinter
	messagePrintPeriod time.Duration
	lastPrint          time.Time
}

func (m *periodicMessagePrinter) Write(p []byte) (int, error) {
	delay := time.Since(m.lastPrint)
	if delay < m.messagePrintPeriod {
		parallelisation.SleepWithContext(m.ctx, m.messagePrintPeriod-delay)
	}
	m.lastPrint = time.Now()
	return m.basicMessagePrinter.Write(p)
}

// newBasicAsynchronousMessagePrinter will print messages asynchronously in a disconnected manner from the paging mechanism.
func newBasicAsynchronousMessagePrinter(ctx context.Context, logger logging.ILogger) logs.WriterWithSource {
	return logs.NewDiodeWriterForSlowWriter(newBasicSynchronousMessagePrinter(ctx, logger), messageBufferSize, 0, logger)
}

// newBasicSynchronousMessagePrinter will print messages synchronously.
func newBasicSynchronousMessagePrinter(ctx context.Context, logger logging.ILogger) logs.WriterWithSource {
	return &basicMessagePrinter{
		logger: logger,
		ctx:    ctx,
	}
}

// newPeriodicAsynchronousMessagePrinter will print messages at regular intervals asynchronously in a disconnected manner from the paging mechanism.
func newPeriodicAsynchronousMessagePrinter(ctx context.Context, logger logging.ILogger, period time.Duration) logs.WriterWithSource {
	return logs.NewDiodeWriterForSlowWriter(newPeriodicSynchronousMessagePrinter(ctx, logger, period), messageBufferSize, 0, logger)
}

// newPeriodicSynchronousMessagePrinter will print messages synchronously but at regular intervals.
func newPeriodicSynchronousMessagePrinter(ctx context.Context, logger logging.ILogger, period time.Duration) logs.WriterWithSource {
	return &periodicMessagePrinter{
		basicMessagePrinter: basicMessagePrinter{logger: logger,
			ctx: ctx},
		messagePrintPeriod: period,
		lastPrint:          time.Now(),
	}
}
