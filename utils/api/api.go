/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Package api provides common helpers related to API calls.
package api

import (
	"context"
	"net/http"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/errors"
	"github.com/ARM-software/golang-utils/utils/http/api"
)

// Deprecated: Use github.com/ARM-software/golang-utils/utils/http/api instead
// IsCallSuccessful determines whether an API response is successful or not
func IsCallSuccessful(r *http.Response) bool {
	return api.IsCallSuccessful(r)
}

// CheckAPICallSuccess verifies whether an API response is successful or not and if not, populates an error with all the information needed.
// errorContext corresponds to the description of what led to the error if error there is e.g. `Failed adding a user`.
// resp corresponds to the HTTP response from a certain endpoint. The body of such response is not closed by this function.
// apiErr corresponds to the error which may be returned by the HTTP client when calling the endpoint.
func CheckAPICallSuccess(ctx context.Context, errorContext string, resp *http.Response, apiErr error) error {
	return api.CheckAPICallSuccess(ctx, errorContext, errors.FetchAPIErrorDescriptionWithContext, resp, apiErr)
}

// CallAndCheckSuccess is a wrapper for making an API call and then checking success with `CheckAPICallSuccess`
// errorContext corresponds to the description of what led to the error if error there is e.g. `Failed adding a user`.
// apiCallFunc corresponds to a generic function that will be called to make the API call
func CallAndCheckSuccess[T any](ctx context.Context, errorContext string, apiCallFunc func(ctx context.Context) (*T, *http.Response, error)) (*T, error) {
	return api.CallAndCheckSuccess[T](ctx, errorContext, errors.FetchAPIErrorDescriptionWithContext, apiCallFunc)
}

// GenericCallAndCheckSuccess is similar to CallAndCheckSuccess but for function returning interfaces rather than concrete types.
// T must be an interface.
// errorContext corresponds to the description of what led to the error if error there is e.g. `Failed adding a user`.
// apiCallFunc corresponds to a generic function that will be called to make the API call
func GenericCallAndCheckSuccess[T any](ctx context.Context, errorContext string, apiCallFunc func(ctx context.Context) (T, *http.Response, error)) (T, error) {
	return api.GenericCallAndCheckSuccess[T](ctx, errorContext, errors.FetchAPIErrorDescriptionWithContext, apiCallFunc)
}
