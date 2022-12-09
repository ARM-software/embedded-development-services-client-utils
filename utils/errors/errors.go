/*
 * Copyright (C) 2020-2022 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package errors

import (
	"encoding/json"
	"fmt"
	"io"
	_http "net/http"
	"strings"

	"github.com/ARM-software/embedded-development-services-client/client"
)

// FetchAPIErrorDescription returns the error message from an API response.
// This function does not close the response body.
func FetchAPIErrorDescription(resp *_http.Response) (message string) {
	if resp == nil {
		return
	}
	errorResponse := client.ErrorResponse{}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &errorResponse)
	if err != nil {
		message = string(bytes)
		return
	}
	apiErrorMessage := strings.Builder{}
	apiErrorMessage.WriteString(errorResponse.GetMessage())
	apiErrorMessage.WriteString(fmt.Sprintf("(request-id: %v)", errorResponse.GetRequestId()))
	if fields, has := errorResponse.GetFieldsOk(); has {
		apiErrorMessage.WriteString(" [")
		for i := range fields {
			field := fields[i]
			apiErrorMessage.WriteString(fmt.Sprintf("%v: %v (%v),", field.GetFieldName(), field.GetMessage(), field.GetFieldPath()))
		}
		apiErrorMessage.WriteString("]")
	}
	message = apiErrorMessage.String()
	return
}
