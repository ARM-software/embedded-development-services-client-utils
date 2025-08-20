/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Package api provides common helpers related to API calls.
package api

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/perimeterx/marshmallow"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/errors"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/parallelisation"
	"github.com/ARM-software/golang-utils/utils/reflection"
	"github.com/ARM-software/golang-utils/utils/safeio"
)

const requiredFieldError = "no value given for required property"

// IsCallSuccessful determines whether an API response is successful or not
func IsCallSuccessful(r *http.Response) bool {
	if r == nil {
		return false
	}
	return r.StatusCode >= http.StatusOK && r.StatusCode < http.StatusMultipleChoices
}

// CheckAPICallSuccess verifies whether an API response is successful or not and if not, populates an error with all the information needed.
// errorContext corresponds to the description of what led to the error if error there is e.g. `Failed adding a user`.
// resp corresponds to the HTTP response from a certain endpoint. The body of such response is not closed by this function.
// apiErr corresponds to the error which may be returned by the HTTP client when calling the endpoint.
func CheckAPICallSuccess(ctx context.Context, errorContext string, resp *http.Response, apiErr error) (err error) {
	err = parallelisation.DetermineContextError(ctx)
	if err != nil {
		return
	}
	if !IsCallSuccessful(resp) {
		statusCode := 0
		errorMessage := strings.Builder{}
		respErr := commonerrors.ErrUnexpected
		if resp != nil {
			statusCode = resp.StatusCode
			respErr = errors.MapErrorToHttpResponseCode(statusCode)
			if respErr == nil {
				respErr = commonerrors.ErrUnexpected
			}
			errorDetails, subErr := errors.FetchAPIErrorDescriptionWithContext(ctx, resp)
			if commonerrors.Ignore(subErr, commonerrors.ErrMarshalling) != nil {
				err = commonerrors.Join(commonerrors.New(respErr, errorContext), subErr)
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
		err = commonerrors.Newf(respErr, "%v (%d): %v%v", errorContext, statusCode, errorMessage.String(), extra)
	}
	return
}

// CallAndCheckSuccess is a wrapper for making an API call and then checking success with `CheckAPICallSuccess`
// errorContext corresponds to the description of what led to the error if error there is e.g. `Failed adding a user`.
// apiCallFunc corresponds to a generic function that will be called to make the API call
func CallAndCheckSuccess[T any](ctx context.Context, errorContext string, apiCallFunc func(ctx context.Context) (*T, *http.Response, error)) (result *T, err error) {
	if err = parallelisation.DetermineContextError(ctx); err != nil {
		return
	}

	result, resp, apiErr := apiCallFunc(ctx)
	if resp != nil && resp.Body != nil {
		defer func() {
			if resp != nil && resp.Body != nil {
				_ = resp.Body.Close()
			}
		}()
	}

	err = checkResponse(ctx, apiErr, resp, result, errorContext)
	return
}

func checkResponse(ctx context.Context, apiErr error, resp *http.Response, result any, errorContext string) (err error) {
	err = CheckAPICallSuccess(ctx, errorContext, resp, apiErr)
	if err != nil {
		return
	}

	if apiErr != nil {
		err = commonerrors.WrapError(commonerrors.ErrMarshalling, apiErr, "API call was successful but an error occurred during response marshalling")
		if commonerrors.CorrespondTo(apiErr, requiredFieldError) {
			return
		}
		if resp == nil || resp.Body == nil {
			return
		}
		// At this point, the marshalling problem may be due to the present of unknown fields in the response due to an API extension.
		// See https://github.com/OpenAPITools/openapi-generator/issues/21446
		var respB []byte
		respB, err = safeio.ReadAll(ctx, resp.Body)
		if err != nil {
			return
		}
		_, err = marshmallow.Unmarshal(respB, result, marshmallow.WithSkipPopulateStruct(false), marshmallow.WithExcludeKnownFieldsFromMap(true))
		if err != nil {
			err = commonerrors.WrapError(commonerrors.ErrMarshalling, err, "API call was successful but an error occurred during response marshalling")
			return
		}
	}
	if reflection.IsEmpty(result) {
		err = commonerrors.New(commonerrors.ErrMarshalling, "unmarshalled response is empty")
		return
	}
	return
}

// GenericCallAndCheckSuccess is similar to CallAndCheckSuccess but for function returning interfaces rather than concrete types.
// T must be an interface.
// errorContext corresponds to the description of what led to the error if error there is e.g. `Failed adding a user`.
// apiCallFunc corresponds to a generic function that will be called to make the API call
func GenericCallAndCheckSuccess[T any](ctx context.Context, errorContext string, apiCallFunc func(ctx context.Context) (T, *http.Response, error)) (result T, err error) {
	if err = parallelisation.DetermineContextError(ctx); err != nil {
		return
	}

	result, resp, apiErr := apiCallFunc(ctx)
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}

	err = checkResponse(ctx, apiErr, resp, result, errorContext)
	if err != nil {
		return
	}

	if reflect.ValueOf(result).Kind() != reflect.Ptr {
		err = commonerrors.Newf(commonerrors.ErrConflict, "result of the call is of type [%T] and so, not a pointer as expected", result)
		return
	}

	return
}
