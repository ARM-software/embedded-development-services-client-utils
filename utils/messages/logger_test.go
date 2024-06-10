/*
 * Copyright (C) 2020-2024 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package messages

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/logging"
	"github.com/ARM-software/embedded-development-services-client/client"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/field"
)

func TestNewMessageLogger(t *testing.T) {
	period := 100 * time.Nanosecond
	tests := []struct {
		messageLogger func(context.Context, logging.ILogger) (IMessageLogger, error)
	}{
		{
			messageLogger: NewBasicSynchronousMessageLogger,
		},
		{
			messageLogger: NewBasicAsynchronousMessageLogger,
		},
		{
			messageLogger: func(ctx context.Context, l logging.ILogger) (IMessageLogger, error) {
				return NewPeriodicSynchronousMessageLogger(ctx, l, period)
			},
		},
		{
			messageLogger: func(ctx context.Context, l logging.ILogger) (IMessageLogger, error) {
				return NewPeriodicAsynchronousMessageLogger(ctx, l, period)
			},
		},
		{
			messageLogger: func(ctx context.Context, l logging.ILogger) (IMessageLogger, error) {
				return NewMessageLoggerFactory(l, false, 0).Create(ctx)
			},
		},
		{
			messageLogger: func(ctx context.Context, l logging.ILogger) (IMessageLogger, error) {
				return NewMessageLoggerFactory(l, false, period).Create(ctx)
			},
		},
		{
			messageLogger: func(ctx context.Context, l logging.ILogger) (IMessageLogger, error) {
				return NewMessageLoggerFactory(l, true, 0).Create(ctx)
			},
		},
		{
			messageLogger: func(ctx context.Context, l logging.ILogger) (IMessageLogger, error) {
				return NewMessageLoggerFactory(l, true, period).Create(ctx)
			},
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			defer goleak.VerifyNone(t)
			stdLogger, err := logging.NewStandardClientLogger(faker.Name(), nil)
			require.NoError(t, err)
			logger, err := test.messageLogger(context.TODO(), stdLogger)
			require.NoError(t, err)
			logger.LogEmptyMessageError()
			logger.LogMarshallingError(nil)
			logger.LogMarshallingError(field.ToOptionalAny(*client.NewMessageObject(faker.Sentence())))
			logger.LogMessage(client.NewMessageObject(faker.Sentence()))
			logger.LogMessage(client.NewNotificationMessageObject(faker.Sentence()))
			require.NoError(t, logger.Close())
		})
	}
}

func TestLogMessageCollectionCancel(t *testing.T) {
	period := 100 * time.Nanosecond
	tests := []struct {
		messageLogger func(context.Context, logging.ILogger) (IMessageLogger, error)
	}{
		{
			messageLogger: NewBasicSynchronousMessageLogger,
		},
		{
			messageLogger: NewBasicAsynchronousMessageLogger,
		},
		{
			messageLogger: func(ctx context.Context, l logging.ILogger) (IMessageLogger, error) {
				return NewPeriodicSynchronousMessageLogger(ctx, l, period)
			},
		},
		{
			messageLogger: func(ctx context.Context, l logging.ILogger) (IMessageLogger, error) {
				return NewPeriodicAsynchronousMessageLogger(ctx, l, period)
			},
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			defer goleak.VerifyNone(t)
			stdLogger, err := logging.NewStandardClientLogger(faker.Name(), nil)
			require.NoError(t, err)
			logger, err := test.messageLogger(context.TODO(), stdLogger)
			require.NoError(t, err)
			defer func() { _ = logger.Close() }()
			messages, err := NewMockNotificationFeedPaginator(context.TODO())
			require.NoError(t, err)
			gtx, cancel := context.WithCancel(context.TODO())
			cancel()
			err = logger.LogMessagesCollection(gtx, messages)
			require.Error(t, err)
			assert.True(t, commonerrors.Any(err, commonerrors.ErrCancelled))
			require.NoError(t, logger.Close())
		})
	}
}

func TestLogMessageCollection(t *testing.T) {
	period := 100 * time.Nanosecond
	tests := []struct {
		messageLogger func(context.Context, logging.ILogger) (IMessageLogger, error)
	}{
		{
			messageLogger: NewBasicSynchronousMessageLogger,
		},
		{
			messageLogger: NewBasicAsynchronousMessageLogger,
		},
		{
			messageLogger: func(ctx context.Context, l logging.ILogger) (IMessageLogger, error) {
				return NewPeriodicSynchronousMessageLogger(ctx, l, period)
			},
		},
		{
			messageLogger: func(ctx context.Context, l logging.ILogger) (IMessageLogger, error) {
				return NewPeriodicAsynchronousMessageLogger(ctx, l, period)
			},
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			defer goleak.VerifyNone(t)
			stdLogger, err := logging.NewStandardClientLogger(faker.Name(), nil)
			require.NoError(t, err)
			logger, err := test.messageLogger(context.TODO(), stdLogger)
			require.NoError(t, err)
			defer func() { _ = logger.Close() }()
			messages, err := NewMockNotificationFeedPaginator(context.TODO())
			require.NoError(t, err)
			require.NoError(t, logger.LogMessagesCollection(context.TODO(), messages))
			require.NoError(t, logger.Close())
		})
	}
}

func TestLogMessageStream(t *testing.T) {
	period := 100 * time.Nanosecond
	tests := []struct {
		messageLogger func(context.Context, logging.ILogger) (IMessageLogger, error)
	}{
		{
			messageLogger: NewBasicSynchronousMessageLogger,
		},
		{
			messageLogger: NewBasicAsynchronousMessageLogger,
		},
		{
			messageLogger: func(ctx context.Context, l logging.ILogger) (IMessageLogger, error) {
				return NewPeriodicSynchronousMessageLogger(ctx, l, period)
			},
		},
		{
			messageLogger: func(ctx context.Context, l logging.ILogger) (IMessageLogger, error) {
				return NewPeriodicAsynchronousMessageLogger(ctx, l, period)
			},
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			defer goleak.VerifyNone(t)
			stdLogger, err := logging.NewStandardClientLogger(faker.Name(), nil)
			require.NoError(t, err)
			logger, err := test.messageLogger(context.TODO(), stdLogger)
			require.NoError(t, err)
			defer func() { _ = logger.Close() }()
			messages, err := NewMockNotificationFeedStreamPaginator(context.TODO())
			require.NoError(t, err)
			require.NoError(t, logger.LogMessagesCollection(context.TODO(), messages))
			require.NoError(t, logger.Close())
		})
	}
}
