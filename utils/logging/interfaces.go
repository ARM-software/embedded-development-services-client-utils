/*
 * Copyright (C) 2020-2022 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Package logging provides helpers regarding logging.
package logging

import (
	"github.com/ARM-software/embedded-development-services-client-utils/utils/resource"
	"github.com/ARM-software/golang-utils/utils/logs"
)

// ILogger defines a generic Client logger.
type ILogger interface {
	logs.IMultipleLoggers
	// LogRawError logs an error.
	LogRawError(err error)
	LogErrorAndMessage(err error, format string, args ...interface{})
	LogErrorMessage(format string, args ...interface{})
	LogInfo(format string, args ...interface{})
	// LogResource logs the description of an API resource
	LogResource(r resource.IResource)
}
