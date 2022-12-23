/*
 * Copyright (C) 2020-2022 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Package errors provides common helpers related to errors coming from the web services.
package errors

import (
	"encoding/json"
	"fmt"
	"io"
	_http "net/http"
	"strings"

	"github.com/ARM-software/embedded-development-services-client/client"
	"github.com/ARM-software/golang-utils/utils/reflection"
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
	apiErrorMessage.WriteString("API call error")
	if code, _ := errorResponse.GetHttpStatusCodeOk(); !reflection.IsEmpty(code) {
		apiErrorMessage.WriteString(fmt.Sprintf(" (%v)", *code))
	}
	apiErrorMessage.WriteString(fmt.Sprintf(" [request-id: %v] ", errorResponse.GetRequestId()))
	apiErrorMessage.WriteString(errorResponse.GetMessage())
	if fields, has := errorResponse.GetFieldsOk(); has {
		apiErrorMessage.WriteString(" [")
		start := true
		for i := range fields {
			if !start {
				apiErrorMessage.WriteString(",")
			}
			field := fields[i]
			apiErrorMessage.WriteString(fmt.Sprintf("%v: %v (%v)", field.GetFieldName(), field.GetMessage(), field.GetFieldPath()))
			start = false
		}
		apiErrorMessage.WriteString("]")
	}
	message = apiErrorMessage.String()
	return
}
