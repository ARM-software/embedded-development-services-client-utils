/*
 * Copyright (C) 2020-2024 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package job

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/logging"
	"github.com/ARM-software/embedded-development-services-client-utils/utils/messages"
	pagination2 "github.com/ARM-software/embedded-development-services-client-utils/utils/pagination"
	"github.com/ARM-software/golang-utils/utils/collection/pagination"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
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
			jobFunc:         NewMockFailedAsynchronousJob,
			expectCompleted: true,
			expectedError:   commonerrors.ErrInvalid,
		},
		{
			jobFunc:         NewMockUndoneAsynchronousJob,
			expectCompleted: false,
			expectedError:   nil,
		},
		{
			jobFunc:         NewMockSuccessfulAsynchronousJob,
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
				assert.True(t, commonerrors.Any(err, test.expectedError))
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
			jobFunc:       NewMockFailedAsynchronousJob,
			expectedError: nil,
		},
		{
			jobFunc:       NewMockSuccessfulAsynchronousJob,
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
				assert.True(t, commonerrors.Any(err, test.expectedError))
			}
		})
	}
}

func TestManager_WaitForJobCompletion(t *testing.T) {
	defer goleak.VerifyNone(t)
	tests := []struct {
		jobFunc       func() (IAsynchronousJob, error)
		expectedError error
	}{
		{
			jobFunc:       NewMockFailedAsynchronousJob,
			expectedError: commonerrors.ErrInvalid,
		},
		{
			jobFunc:       NewMockSuccessfulAsynchronousJob,
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
			err = factory.WaitForJobCompletion(context.TODO(), job)
			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.True(t, commonerrors.Any(err, test.expectedError))
			}
		})
	}
}

func newMockJobManager(logger *messages.MessageLoggerFactory, backOffPeriod time.Duration, messagePaginatorRunOutTimeout *time.Duration, job IAsynchronousJob, errToReturn error) (*Manager, error) {
	pageNumber := rand.Intn(50) //nolint:gosec //causes G404: Use of weak random number generator
	messageStream := messages.NewMockMessagePaginatorFactory(pageNumber)
	if messagePaginatorRunOutTimeout != nil {
		messageStream = messageStream.UpdateRunOutTimeout(*messagePaginatorRunOutTimeout)
	}

	return newJobManagerFromMessageFactory(logger, backOffPeriod, func(context.Context, string) (IAsynchronousJob, *http.Response, error) {
		return job, httptest.NewRecorder().Result(), errToReturn
	}, func(fctx context.Context, jobName string) (pagination.IStaticPageStream, *http.Response, error) {
		firstPage, err := messages.NewMockNotificationFeedPage(fctx, pageNumber > 0, false)
		if err != nil {
			return nil, httptest.NewRecorder().Result(), err
		}
		return pagination2.ToStream(firstPage), httptest.NewRecorder().Result(), nil
	}, messageStream)
}
