/*
 * Copyright (C) 2020-2024 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Package errors provides common helpers related to errors coming from the web services.
package errors

import (
	"context"
	"encoding/json"
	"fmt"
	_http "net/http"
	"strings"

	"github.com/ARM-software/embedded-development-services-client/client"
	"github.com/ARM-software/golang-utils/utils/parallelisation"
	"github.com/ARM-software/golang-utils/utils/reflection"
	"github.com/ARM-software/golang-utils/utils/safeio"
)

// FetchAPIErrorDescription returns the error message from an API response.
// This function does not close the response body.
func FetchAPIErrorDescription(resp *_http.Response) (message string) {
	message, _ = FetchAPIErrorDescriptionWithContext(context.Background(), resp)
	return
}

// FetchAPIErrorDescriptionWithContext returns the error message from an API response.
// This function does not close the response body.
func FetchAPIErrorDescriptionWithContext(ctx context.Context, resp *_http.Response) (message string, err error) {
	if resp == nil {
		return
	}
	err = parallelisation.DetermineContextError(ctx)
	if err != nil {
		return
	}
	errorResponse := client.ErrorResponse{}
	bytes, err := safeio.ReadAll(ctx, resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &errorResponse)
	if err != nil {
		alternativeErrorResponse := alternativeErrorMessage{}
		err = json.Unmarshal(bytes, &alternativeErrorResponse)
		if err == nil {
			alternativeErrorResponse.setErrorResponse(&errorResponse)
		}
	}
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
			subErr := parallelisation.DetermineContextError(ctx)
			if subErr != nil {
				err = subErr
				return
			}
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

type alternativeErrorMessage struct {
	Fields         []client.FieldObject `json:"fields,omitempty"`
	HTTPStatusCode *int32               `json:"httpStatusCode,omitempty"`
	Message        string               `json:"message"`
	RequestID      string               `json:"requestId"`
}

func (e *alternativeErrorMessage) setErrorResponse(err *client.ErrorResponse) {
	err.SetMessage(e.Message)
	err.SetRequestId(e.RequestID)
	err.SetFields(e.Fields)
	if e.HTTPStatusCode != nil {
		err.SetHttpStatusCode(*e.HTTPStatusCode)
	}
}
