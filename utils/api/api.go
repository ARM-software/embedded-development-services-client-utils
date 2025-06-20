/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Package api provides common helpers related to API calls.
package api

import (
	"context"
	"fmt"
	_http "net/http"
	"reflect"
	"strings"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/errors"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/parallelisation"
	"github.com/ARM-software/golang-utils/utils/reflection"
)

// IsCallSuccessful determines whether an API response is successful or not
func IsCallSuccessful(r *_http.Response) bool {
	if r == nil {
		return false
	}
	return r.StatusCode >= _http.StatusOK && r.StatusCode < _http.StatusMultipleChoices
}

// CheckAPICallSuccess verifies whether an API response is successful or not and if not, populates an error with all the information needed.
// errorContext corresponds to the description of what led to the error if error there is e.g. `Failed adding a user`.
// resp corresponds to the HTTP response from a certain endpoint. The body of such response is not closed by this function.
// apiErr corresponds to the error which may be returned by the HTTP client when calling the endpoint.
func CheckAPICallSuccess(ctx context.Context, errorContext string, resp *_http.Response, apiErr error) (err error) {
	err = parallelisation.DetermineContextError(ctx)
	if err != nil {
		return
	}
	if !IsCallSuccessful(resp) {
		statusCode := 0
		errorMessage := strings.Builder{}
		if resp != nil {
			statusCode = resp.StatusCode
			errorDetails, subErr := errors.FetchAPIErrorDescriptionWithContext(ctx, resp)
			if commonerrors.Ignore(subErr, commonerrors.ErrMarshalling) != nil {
				err = subErr
				return
			}
			if !reflection.IsEmpty(errorDetails) {
				errorMessage.WriteString(errorDetails)
			}
			_ = resp.Body.Close()
		}
		extra := ""
		if apiErr != nil {
			extra = fmt.Sprintf("; %v", apiErr.Error())
		}
		err = fmt.Errorf("%v (%d): %v%v", errorContext, statusCode, errorMessage.String(), extra)
	}
	return
}

// CallAndCheckSuccess is a wrapper for making an API call and then checking success with `CheckAPICallSuccess`
// errorContext corresponds to the description of what led to the error if error there is e.g. `Failed adding a user`.
// apiCallFunc corresponds to a generic function that will be called to make the API call
func CallAndCheckSuccess[T any](ctx context.Context, errorContext string, apiCallFunc func(ctx context.Context) (*T, *_http.Response, error)) (result *T, err error) {
	if err = parallelisation.DetermineContextError(ctx); err != nil {
		return
	}

	result, resp, apiErr := apiCallFunc(ctx)
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}

	if err = CheckAPICallSuccess(ctx, errorContext, resp, apiErr); err != nil {
		return
	}

	if apiErr != nil {
		err = commonerrors.WrapError(commonerrors.ErrMarshalling, apiErr, "API call was successful but an error occured during response marshalling")
		return
	}

	if result != nil && reflection.IsEmpty(result) {
		err = commonerrors.New(commonerrors.ErrMarshalling, "unmarshalled response is empty")
		return
	}

	return
}

// GenericCallAndCheckSuccess is similar to CallAndCheckSuccess but for function returning interfaces rather than concrete types.
// T must be an interface.
// errorContext corresponds to the description of what led to the error if error there is e.g. `Failed adding a user`.
// apiCallFunc corresponds to a generic function that will be called to make the API call
func GenericCallAndCheckSuccess[T any](ctx context.Context, errorContext string, apiCallFunc func(ctx context.Context) (T, *_http.Response, error)) (result T, err error) {
	if err = parallelisation.DetermineContextError(ctx); err != nil {
		return
	}

	result, resp, apiErr := apiCallFunc(ctx)
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}

	if err = CheckAPICallSuccess(ctx, errorContext, resp, apiErr); err != nil {
		return
	}

	if reflect.ValueOf(result).Kind() != reflect.Ptr {
		err = commonerrors.Newf(commonerrors.ErrConflict, "result of the call is of type [%T] and so, not a pointer as expected", result)
		return
	}

	if reflection.IsEmpty(result) {
		err = commonerrors.New(commonerrors.ErrMarshalling, "unmarshalled response is empty")
		return
	}

	return
}
