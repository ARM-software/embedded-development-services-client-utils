/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package job

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/job/jobtest"
	"github.com/ARM-software/embedded-development-services-client-utils/utils/logging"
	"github.com/ARM-software/embedded-development-services-client-utils/utils/messages"
	pagination2 "github.com/ARM-software/embedded-development-services-client-utils/utils/pagination"
	"github.com/ARM-software/golang-utils/utils/collection/pagination"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/commonerrors/errortest"
	"github.com/ARM-software/golang-utils/utils/field"
)

func TestManager_HasJobCompleted(t *testing.T) {
	logger, err := logging.NewStandardClientLogger("test", nil)
	require.NoError(t, err)
	loggerF := messages.NewMessageLoggerFactory(logger, false, time.Nanosecond)

	tests := []struct {
		jobFunc         func() (IAsynchronousJob, error)
		expectCompleted bool
		expectedError   error
	}{
		{
			jobFunc:         func() (IAsynchronousJob, error) { return nil, nil },
			expectCompleted: false,
			expectedError:   commonerrors.ErrUndefined,
		},
		{
			jobFunc:         mapFunc(jobtest.NewMockFailedAsynchronousJob),
			expectCompleted: true,
			expectedError:   commonerrors.ErrInvalid,
		},
		{
			jobFunc:         mapFunc(jobtest.NewMockUndoneAsynchronousJob),
			expectCompleted: false,
			expectedError:   nil,
		},
		{
			jobFunc:         mapFunc(jobtest.NewMockQueuedAsynchronousJob),
			expectCompleted: false,
			expectedError:   nil,
		},
		{
			jobFunc:         mapFunc(jobtest.NewMockSuccessfulAsynchronousJob),
			expectCompleted: true,
			expectedError:   nil,
		},
	}
	for i := range tests {
		test := tests[i]

		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			defer goleak.VerifyNone(t)
			job, err := test.jobFunc()
			factory, err := newMockJobManager(loggerF, 100*time.Millisecond, nil, job, err)
			require.NoError(t, err)
			require.NotNil(t, factory)
			completed, err := factory.HasJobCompleted(context.TODO(), job)
			if test.expectedError == nil {
				assert.NoError(t, err)
				assert.Equal(t, test.expectCompleted, completed)
			} else {
				assert.Error(t, err)
				errortest.AssertError(t, err, test.expectedError)
				assert.Equal(t, test.expectCompleted, completed)
			}
		})
	}
}

func TestManager_checkForMessageStreamExhaustion(t *testing.T) {
	logger, err := logging.NewStandardClientLogger("test", nil)
	require.NoError(t, err)
	loggerF := messages.NewMessageLoggerFactory(logger, false, time.Nanosecond)

	tests := []struct {
		jobFunc       func() (IAsynchronousJob, error)
		expectedError error
	}{
		{
			jobFunc:       func() (IAsynchronousJob, error) { return nil, nil },
			expectedError: commonerrors.ErrUndefined,
		},
		{
			jobFunc:       mapFunc(jobtest.NewMockFailedAsynchronousJob),
			expectedError: nil,
		},
		{
			jobFunc:       mapFunc(jobtest.NewMockSuccessfulAsynchronousJob),
			expectedError: nil,
		},
	}
	for i := range tests {
		test := tests[i]

		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			defer goleak.VerifyNone(t)
			job, err := test.jobFunc()
			factory, err := newMockJobManager(loggerF, 100*time.Millisecond, nil, job, err)
			require.NoError(t, err)
			require.NotNil(t, factory)

			messagePaginator, err := factory.createMessagePaginator(context.TODO(), job)
			if test.expectedError == nil {
				require.NoError(t, err)
				assert.NotNil(t, messagePaginator)
				assert.False(t, messagePaginator.IsRunningDry())
			} else {
				require.Error(t, err)
				assert.Nil(t, messagePaginator)
			}

			err = factory.checkForMessageStreamExhaustion(context.TODO(), messagePaginator, job)
			if test.expectedError == nil {
				assert.NoError(t, err)
				assert.True(t, messagePaginator.IsRunningDry())
			} else {
				assert.Error(t, err)
				errortest.AssertError(t, err, test.expectedError)
			}
		})
	}
}

func mapFunc(f func() (*jobtest.MockAsynchronousJob, error)) func() (IAsynchronousJob, error) {
	return func() (IAsynchronousJob, error) {
		job, err := f()
		return job, err
	}
}

func TestManager_logMessages(t *testing.T) {
	defer goleak.VerifyNone(t)
	tests := []struct {
		jobFunc       func() (IAsynchronousJob, error)
		expectedError []error
		timeout       *time.Duration
	}{
		{
			jobFunc:       mapFunc(jobtest.NewMockFailedAsynchronousJob),
			expectedError: nil,
		},
		{
			jobFunc:       mapFunc(jobtest.NewMockQueuedAsynchronousJob),
			expectedError: []error{commonerrors.ErrCondition, commonerrors.ErrTimeout, commonerrors.ErrCancelled},
			timeout:       field.ToOptionalDuration(500 * time.Millisecond),
		},
		{
			jobFunc:       mapFunc(jobtest.NewMockSuccessfulAsynchronousJob),
			expectedError: nil,
		},
	}
	for i := range tests {
		test := tests[i]

		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			// t.Parallel()
			logger, err := logging.NewStandardClientLogger(fmt.Sprintf("test #%v", i), nil)
			require.NoError(t, err)
			loggerF := messages.NewMessageLoggerFactory(logger, false, time.Nanosecond)
			job, err := test.jobFunc()
			runOut := time.Nanosecond
			factory, err := newMockJobManager(loggerF, time.Nanosecond, &runOut, job, err)

			require.NoError(t, err)
			require.NotNil(t, factory)
			if test.timeout == nil {
				err = factory.LogJobMessagesUntilNow(context.TODO(), job, 5*time.Minute)
			} else {
				err = factory.LogJobMessagesUntilNow(context.TODO(), job, *test.timeout)
			}
			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				errortest.AssertError(t, err, test.expectedError...)
			}
		})
	}
}

func TestManager_WaitForJobCompletion(t *testing.T) {
	defer goleak.VerifyNone(t)
	tests := []struct {
		jobFunc       func() (IAsynchronousJob, error)
		expectedError []error
		timeout       *time.Duration
	}{
		{
			jobFunc:       mapFunc(jobtest.NewMockFailedAsynchronousJob),
			expectedError: []error{commonerrors.ErrInvalid},
		},
		{
			jobFunc:       mapFunc(jobtest.NewMockQueuedAsynchronousJob),
			expectedError: []error{commonerrors.ErrCondition, commonerrors.ErrTimeout, commonerrors.ErrCancelled},
			timeout:       field.ToOptionalDuration(500 * time.Millisecond),
		},
		{
			jobFunc:       mapFunc(jobtest.NewMockSuccessfulAsynchronousJob),
			expectedError: nil,
		},
	}
	for i := range tests {
		test := tests[i]

		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			// t.Parallel()
			logger, err := logging.NewStandardClientLogger(fmt.Sprintf("test #%v", i), nil)
			require.NoError(t, err)
			loggerF := messages.NewMessageLoggerFactory(logger, false, time.Nanosecond)
			job, err := test.jobFunc()
			runOut := time.Nanosecond
			factory, err := newMockJobManager(loggerF, time.Nanosecond, &runOut, job, err)

			require.NoError(t, err)
			require.NotNil(t, factory)
			if test.timeout == nil {
				err = factory.WaitForJobCompletion(context.TODO(), job)
			} else {
				err = factory.WaitForJobCompletionWithTimeout(context.TODO(), job, *test.timeout)
			}
			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				errortest.AssertError(t, err, test.expectedError...)
			}
		})
	}
}

func TestManager_WaitForJobCompletionTimeout(t *testing.T) {
	defer goleak.VerifyNone(t)
	tests := []struct {
		jobFunc func() (IAsynchronousJob, error)
	}{
		{
			jobFunc: mapFunc(jobtest.NewMockFailedAsynchronousJob),
		},
		{
			jobFunc: mapFunc(jobtest.NewMockQueuedAsynchronousJob),
		},
		{
			jobFunc: mapFunc(jobtest.NewMockSuccessfulAsynchronousJob),
		},
	}
	for i := range tests {
		test := tests[i]

		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			// t.Parallel()
			logger, err := logging.NewStandardClientLogger(fmt.Sprintf("test #%v", i), nil)
			require.NoError(t, err)
			loggerF := messages.NewMessageLoggerFactory(logger, false, time.Nanosecond)
			job, err := test.jobFunc()
			runOut := time.Nanosecond
			factory, err := newMockJobManager(loggerF, time.Nanosecond, &runOut, job, err)

			require.NoError(t, err)
			require.NotNil(t, factory)

			err = factory.WaitForJobCompletionWithTimeout(context.TODO(), job, time.Nanosecond)
			assert.Error(t, err)
			errortest.AssertError(t, err, commonerrors.ErrInvalid, commonerrors.ErrTimeout, commonerrors.ErrCancelled, commonerrors.ErrCondition)
		})
	}
}

func newMockJobManager(logger *messages.MessageLoggerFactory, backOffPeriod time.Duration, messagePaginatorRunOutTimeout *time.Duration, job IAsynchronousJob, errToReturn error) (*Manager, error) {
	n, err := faker.RandomInt(1, 50)
	if err != nil {
		return nil, err
	}
	pageNumber := n[0]
	messageStream := messages.NewMockMessagePaginatorFactory(pageNumber)
	if messagePaginatorRunOutTimeout != nil {
		messageStream = messageStream.UpdateRunOutTimeout(*messagePaginatorRunOutTimeout)
	}

	return newJobManagerFromMessageFactory(logger, backOffPeriod, func(context.Context, string) (IAsynchronousJob, *http.Response, error) {
		return job, httptest.NewRecorder().Result(), errToReturn
	}, func(fctx context.Context, _ string) (pagination.IStaticPageStream, *http.Response, error) {
		firstPage, err := messages.NewMockNotificationFeedPage(fctx, pageNumber > 0, false)
		if err != nil {
			return nil, httptest.NewRecorder().Result(), err
		}
		return pagination2.ToStream(firstPage), httptest.NewRecorder().Result(), nil
	}, messageStream)
}
