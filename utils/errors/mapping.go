/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package errors

import (
	"github.com/ARM-software/golang-utils/utils/http/errors"
)

// Deprecated: MapErrorToHTTPResponseCode maps a response status code to a common error.
// Use github.com/ARM-software/golang-utils/utils/http/errors instead
func MapErrorToHTTPResponseCode(statusCode int) error {
	return errors.MapErrorToHTTPResponseCode(statusCode)
}
