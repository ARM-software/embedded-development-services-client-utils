/*
 * Copyright (C) 2020-2022 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package errors

import (
	"bytes"
	"io"
	_http "net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchAPIErrorDescription(t *testing.T) {
	t.Run("error response exists", func(*testing.T) {
		resp := _http.Response{Body: io.NopCloser(bytes.NewReader([]byte("{\"message\": \"client error\",\"requestId\": \"761761721\"}")))}
		actualMessage := FetchAPIErrorDescription(&resp)
		expectedMessage := "client error(request-id: 761761721)"
		assert.Equal(t, actualMessage, expectedMessage)
	})

	t.Run("error response has fields", func(*testing.T) {
		resp := _http.Response{Body: io.NopCloser(bytes.NewReader([]byte("{\"message\":\"client error\",\"requestId\":\"761761721\",\"fields\":[{\"fieldName\":\"client request error\",\"fieldPath\":\"https://foo.bar\",\"message\":\"client error\"}]}")))}
		actualMessage := FetchAPIErrorDescription(&resp)
		expectedMessage := "client error(request-id: 761761721) [client request error: client error (https://foo.bar),]"
		assert.Equal(t, actualMessage, expectedMessage)
	})

	t.Run("error response empty", func(t *testing.T) {
		respBody := _http.Response{Body: io.NopCloser(bytes.NewReader(nil))}
		actualMessage := FetchAPIErrorDescription(&respBody)
		expectedMessage := ""
		assert.Equal(t, actualMessage, expectedMessage)
	})
}
