/*
 * Copyright (C) 2020-2022 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package api

import (
	"context"
	"fmt"
	_http "net/http"
	"strings"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/errors"
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
	if err = parallelisation.DetermineContextError(ctx); err != nil {
		return err
	}
	if !IsCallSuccessful(resp) {
		statusCode := 0
		errorMessage := strings.Builder{}
		if resp != nil {
			statusCode = resp.StatusCode
			errorDetails := errors.FetchAPIErrorDescription(resp)
			if !reflection.IsEmpty(errorDetails) {
				errorMessage.WriteString("further details: ")
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
