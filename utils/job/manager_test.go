/*
 * Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/logging"
	"github.com/ARM-software/embedded-development-services-client-utils/utils/messages"
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
			factory, err := newJobManagerFromMessageFactory(loggerF, 100*time.Millisecond, func(context.Context, string) (IAsynchronousJob, *http.Response, error) {
				return job, httptest.NewRecorder().Result(), err
			}, messages.NewMockMessagePaginatorFactory())
			require.NoError(t, err)
			require.NotNil(t, factory)
			completed, err := factory.HasJobCompleted(context.TODO(), job)
			if test.expectedError == nil {
				assert.NoError(t, err)
				assert.Equal(t, test.expectCompleted, completed)
			} else {
				assert.Error(t, err)
				assert.True(t, commonerrors.Any(err, test.expectedError))
			}
		})
	}
}

// FIXME to enable
// func TestManager_WaitForJobCompletion(t *testing.T) {
//	logger, err := logging.NewStandardClientLogger("test", nil)
//	require.NoError(t, err)
//	loggerF := messages.NewMessageLoggerFactory(logger, false, time.Nanosecond)
//
//	tests := []struct {
//		jobFunc       func() (IAsynchronousJob, error)
//		expectedError error
//	}{
//		//{
//		//	jobFunc:       NewMockFailedAsynchronousJob,
//		//	expectedError: commonerrors.ErrInvalid,
//		//},
//		{
//			jobFunc:       NewMockSuccessfulAsynchronousJob,
//			expectedError: nil,
//		},
//	}
//	for i := range tests {
//		test := tests[i]
//
//		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
//         defer goleak.VerifyNone(t)
//			job, err := test.jobFunc()
//			factory, err := newJobManagerFromMessageFactory(loggerF, time.Nanosecond, func(context.Context, string) (IAsynchronousJob, *http.Response, error) {
//				return job, httptest.NewRecorder().Result(), err
//			}, messages.NewMockMessagePaginatorFactory().UpdateRunOutTimeout(time.Nanosecond))
//			require.NoError(t, err)
//			require.NotNil(t, factory)
//			err = factory.WaitForJobCompletion(context.TODO(), job)
//			if test.expectedError == nil {
//				assert.NoError(t, err)
//			} else {
//				assert.Error(t, err)
//				assert.True(t, commonerrors.Any(err, test.expectedError))
//			}
//		})
//	}
// }
