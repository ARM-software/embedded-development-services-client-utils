/*
 * Copyright (C) 2020-2024 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package logging

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/resource/resourcetests"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/filesystem"
	"github.com/ARM-software/golang-utils/utils/logs"
	"github.com/ARM-software/golang-utils/utils/logs/logstest"
)

func TestClientLogger(t *testing.T) {
	logger, err := NewClientLogger("test client Logger")
	require.NoError(t, err)

	defer func() { _ = logger.Close() }()

	require.NoError(t, logger.AppendLogger(logstest.NewNullTestLogger(), logstest.NewTestLogger(t), logstest.NewStdTestLogger()))
	require.NoError(t, logger.Check())
	sl, err := logs.NewStdLogger(faker.Name())
	require.NoError(t, err)
	require.NoError(t, logger.Append(sl))

	require.NoError(t, logger.SetLogSource(faker.Name()))
	logger.LogRawError(commonerrors.ErrUndefined)
	logger.LogErrorMessage("%v ....", faker.Sentence())
	logger.LogErrorAndMessage(commonerrors.ErrUnexpected, "some error %v (%v): %v", "no idea", 123, faker.Sentence())
	logger.LogInfo("information %v", faker.Sentence())
	resource, err := resourcetests.NewMockResource()
	require.NoError(t, err)
	logger.LogResource(resource)

	require.NoError(t, logger.Close())
}

func TestStandardClientLogger(t *testing.T) {
	filepath, err := filesystem.TempFileInTempDir("test-standard-logger-*.log")
	require.NoError(t, err)
	require.NoError(t, filepath.Close())
	logPath := filepath.Name()
	defer func() { _ = filesystem.Rm(logPath) }()

	empty, err := filesystem.IsEmpty(logPath)
	require.NoError(t, err)
	assert.True(t, empty)

	logger, err := NewStandardClientLogger("test client Logger", &logPath)
	require.NoError(t, err)

	defer func() { _ = logger.Close() }()

	require.NoError(t, logger.SetLogSource(faker.Name()))
	logger.LogRawError(commonerrors.ErrUndefined)
	logger.LogErrorMessage("%v ....", faker.Sentence())
	logger.LogErrorAndMessage(commonerrors.ErrUnexpected, "some error %v (%v): %v", "no idea", 123, faker.Sentence())
	logger.LogInfo("information %v", faker.Sentence())
	resource, err := resourcetests.NewMockResource()
	require.NoError(t, err)
	logger.LogResource(resource)

	require.NoError(t, logger.Close())

	empty, err = filesystem.IsEmpty(logPath)
	require.NoError(t, err)
	assert.False(t, empty)
}
