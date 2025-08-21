/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package messages

import (
	"context"
	"time"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/logging"
	"github.com/ARM-software/golang-utils/utils/collection/pagination"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/logs"
	"github.com/ARM-software/golang-utils/utils/parallelisation"
	"github.com/ARM-software/golang-utils/utils/reflection"
)

const (
	messageBufferSize = 1000
)

type logger struct {
	rawLogger    logging.ILogger
	printer      logs.WriterWithSource
	msgFormatter MessageFormatter
}

func newLogger(rawLogger logging.ILogger, printer logs.WriterWithSource, msgFormatter *MessageFormatter) (IMessageLogger, error) {
	if rawLogger == nil {
		return nil, commonerrors.ErrNoLogger
	}
	if msgFormatter == nil {
		return nil, commonerrors.UndefinedVariable("message formatter")
	}
	return &logger{
		rawLogger:    rawLogger,
		printer:      printer,
		msgFormatter: *msgFormatter,
	}, nil
}

func (l *logger) Check() error {
	return l.rawLogger.Check()
}

func (l *logger) SetLogSource(source string) error {
	return l.rawLogger.SetLogSource(source)
}

func (l *logger) SetLoggerSource(source string) error {
	return l.rawLogger.SetLoggerSource(source)
}

func (l *logger) Log(output ...interface{}) {
	l.rawLogger.Log(output...)
}

func (l *logger) LogError(err ...interface{}) {
	l.rawLogger.LogError(err...)
}

func (l *logger) Close() error {
	return l.printer.Close()
}

func (l *logger) SetSource(source string) error {
	return l.printer.SetSource(source)
}

func (l *logger) LogMessage(msg IMessage) {
	m, err := l.msgFormatter.FormatMessage(msg)
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

func (l *logger) LogEmptyMessageError() {
	l.rawLogger.LogErrorAndMessage(commonerrors.ErrEmpty, "empty message")
}

func (l *logger) LogMarshallingError(rawMessage *any) {
	if reflection.IsEmpty(rawMessage) {
		l.LogEmptyMessageError()
	} else {
		l.rawLogger.LogErrorAndMessage(commonerrors.ErrMarshalling, "message [%v] of type [%T] could not be marshalled", *rawMessage, *rawMessage)
	}
}

func (l *logger) LogMessagesCollection(ctx context.Context, messagePaginator pagination.IGenericPaginator) error {
	for {
		err := parallelisation.DetermineContextError(ctx)
		if err != nil {
			return err
		}
		if messagePaginator == nil {
			return commonerrors.UndefinedVariable("paginator")
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
			l.LogMarshallingError(&m)
		} else {
			l.LogMessage(messageItem)
		}
	}
}

// NewBasicAsynchronousMessageLogger creates an asynchronous logger for messages which prints them as they come. It uses the default message formatter.
func NewBasicAsynchronousMessageLogger(ctx context.Context, rawLogger logging.ILogger) (IMessageLogger, error) {
	return NewBasicAsynchronousMessageLoggerWithFormatter(ctx, rawLogger, DefaultMessageFormatter())
}

// NewBasicAsynchronousMessageLoggerWithFormatter creates an asynchronous logger for messages which prints them as they come.
func NewBasicAsynchronousMessageLoggerWithFormatter(ctx context.Context, rawLogger logging.ILogger, msgFormatter *MessageFormatter) (IMessageLogger, error) {
	if rawLogger == nil {
		return nil, commonerrors.ErrNoLogger
	}
	return newLogger(rawLogger, newBasicAsynchronousMessagePrinter(ctx, rawLogger), msgFormatter)
}

// NewPeriodicAsynchronousMessageLogger creates an asynchronous logger for messages which prints them at regular intervals. It uses the default message formatter.
func NewPeriodicAsynchronousMessageLogger(ctx context.Context, rawLogger logging.ILogger, printPeriod time.Duration) (IMessageLogger, error) {
	return NewPeriodicAsynchronousMessageLoggerWithFormatter(ctx, rawLogger, printPeriod, DefaultMessageFormatter())
}

// NewPeriodicAsynchronousMessageLoggerWithFormatter creates an asynchronous logger for messages which prints them at regular intervals.
func NewPeriodicAsynchronousMessageLoggerWithFormatter(ctx context.Context, rawLogger logging.ILogger, printPeriod time.Duration, msgFormatter *MessageFormatter) (IMessageLogger, error) {
	if rawLogger == nil {
		return nil, commonerrors.ErrNoLogger
	}
	return newLogger(rawLogger, newPeriodicAsynchronousMessagePrinter(ctx, rawLogger, printPeriod), msgFormatter)
}

// NewBasicSynchronousMessageLogger creates a synchronous logger for messages which prints them as they come. It uses the default message formatter.
func NewBasicSynchronousMessageLogger(ctx context.Context, rawLogger logging.ILogger) (IMessageLogger, error) {
	return NewBasicSynchronousMessageLoggerWithFormatter(ctx, rawLogger, DefaultMessageFormatter())
}

// NewBasicSynchronousMessageLoggerWithFormatter creates a synchronous logger for messages which prints them as they come.
func NewBasicSynchronousMessageLoggerWithFormatter(ctx context.Context, rawLogger logging.ILogger, msgFormatter *MessageFormatter) (IMessageLogger, error) {
	if rawLogger == nil {
		return nil, commonerrors.ErrNoLogger
	}
	return newLogger(rawLogger, newBasicSynchronousMessagePrinter(ctx, rawLogger), msgFormatter)
}

// NewPeriodicSynchronousMessageLogger creates a synchronous logger for messages which prints them at regular intervals. It uses the default message formatter.
func NewPeriodicSynchronousMessageLogger(ctx context.Context, rawLogger logging.ILogger, printPeriod time.Duration) (IMessageLogger, error) {
	if rawLogger == nil {
		return nil, commonerrors.ErrNoLogger
	}
	return newLogger(rawLogger, newPeriodicSynchronousMessagePrinter(ctx, rawLogger, printPeriod), DefaultMessageFormatter())
}

// NewPeriodicSynchronousMessageLoggerWithFormatter creates a synchronous logger for messages which prints them at regular intervals.
func NewPeriodicSynchronousMessageLoggerWithFormatter(ctx context.Context, rawLogger logging.ILogger, printPeriod time.Duration, msgFormatter *MessageFormatter) (IMessageLogger, error) {
	if rawLogger == nil {
		return nil, commonerrors.ErrNoLogger
	}
	return newLogger(rawLogger, newPeriodicSynchronousMessagePrinter(ctx, rawLogger, printPeriod), msgFormatter)
}

// MessageLoggerFactory defines a message logger factory
type MessageLoggerFactory struct {
	asynchronous bool
	period       time.Duration
	rawLogger    logging.ILogger
	msgFormatter *MessageFormatter
}

// Create returns a message logger.
func (f *MessageLoggerFactory) Create(ctx context.Context) (IMessageLogger, error) {
	if f.rawLogger == nil {
		return nil, commonerrors.ErrNoLogger
	}
	if f.asynchronous {
		if f.period > 0 {
			return NewPeriodicAsynchronousMessageLoggerWithFormatter(ctx, f.rawLogger, f.period, f.msgFormatter)
		}
		return NewBasicAsynchronousMessageLoggerWithFormatter(ctx, f.rawLogger, f.msgFormatter)
	}
	if f.period > 0 {
		return NewPeriodicSynchronousMessageLoggerWithFormatter(ctx, f.rawLogger, f.period, f.msgFormatter)
	}
	return NewBasicSynchronousMessageLoggerWithFormatter(ctx, f.rawLogger, f.msgFormatter)
}

// NewMessageLoggerFactory returns a message logger factory.
func NewMessageLoggerFactory(logger logging.ILogger, asynchronous bool, printingPeriod time.Duration) *MessageLoggerFactory {
	return NewMessageLoggerFactoryWithFormatter(logger, asynchronous, printingPeriod, DefaultMessageFormatter())
}

// NewMessageLoggerFactoryWithFormatter returns a message logger factory.
func NewMessageLoggerFactoryWithFormatter(logger logging.ILogger, asynchronous bool, printingPeriod time.Duration, formatter *MessageFormatter) *MessageLoggerFactory {
	return &MessageLoggerFactory{
		asynchronous: asynchronous,
		period:       printingPeriod,
		rawLogger:    logger,
		msgFormatter: formatter,
	}
}

// NewMessageLoggerFactoryWithFormattingOptions returns a message logger factory.
func NewMessageLoggerFactoryWithFormattingOptions(logger logging.ILogger, asynchronous bool, printingPeriod time.Duration, option ...FormatterOption) *MessageLoggerFactory {
	return NewMessageLoggerFactoryWithFormatter(logger, asynchronous, printingPeriod, NewMessageFormatter(option...))
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
