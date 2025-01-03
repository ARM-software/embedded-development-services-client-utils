/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package errors

import (
	"bytes"
	"context"
	"io"
	_http "net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/commonerrors/errortest"
)

func TestFetchAPIErrorDescription(t *testing.T) {
	t.Run("error response exists", func(t *testing.T) {
		resp := _http.Response{Body: io.NopCloser(bytes.NewReader([]byte("{\"message\": \"client error\",\"requestId\": \"761761721\"}")))}
		actualMessage := FetchAPIErrorDescription(&resp)
		expectedMessage := "API call error [request-id: 761761721] client error"
		assert.Equal(t, expectedMessage, actualMessage)
	})

	t.Run("error response exists with status", func(t *testing.T) {
		resp := _http.Response{Body: io.NopCloser(bytes.NewReader([]byte("{\"message\": \"client error\",\"requestId\": \"761761721\", \"httpStatusCode\": 406}")))}
		actualMessage := FetchAPIErrorDescription(&resp)
		expectedMessage := "API call error (406) [request-id: 761761721] client error"
		assert.Equal(t, expectedMessage, actualMessage)
	})

	t.Run("error response has fields", func(t *testing.T) {
		resp := _http.Response{Body: io.NopCloser(bytes.NewReader([]byte("{\"message\":\"client error\",\"requestId\":\"761761721\",\"fields\":[{\"fieldName\":\"client request error\",\"fieldPath\":\"https://foo.bar\",\"message\":\"client error\"}]}")))}
		actualMessage := FetchAPIErrorDescription(&resp)
		expectedMessage := "API call error [request-id: 761761721] client error [client request error: client error (https://foo.bar)]"
		assert.Equal(t, expectedMessage, actualMessage)
	})

	t.Run("error response empty", func(t *testing.T) {
		respBody := _http.Response{Body: io.NopCloser(bytes.NewReader(nil))}
		actualMessage := FetchAPIErrorDescription(&respBody)
		expectedMessage := ""
		assert.Equal(t, expectedMessage, actualMessage)
	})
}

func TestFetchAPIErrorDescriptionInterrupt(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp := _http.Response{Body: io.NopCloser(bytes.NewReader([]byte("{\"message\": \"client error\",\"requestId\": \"761761721\"}")))}
	_, err := FetchAPIErrorDescriptionWithContext(ctx, &resp)
	errortest.AssertError(t, err, commonerrors.ErrTimeout, commonerrors.ErrCancelled)
}
