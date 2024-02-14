/*
 * Copyright (C) 2020-2024 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"bytes"
	"context"
	"errors"
	"io"
	_http "net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ARM-software/golang-utils/utils/commonerrors"
)

func TestIsAPICallSuccessful(t *testing.T) {
	t.Run("api call successful", func(*testing.T) {
		resp := _http.Response{StatusCode: 200}
		isSuccessful := IsCallSuccessful(&resp)
		assert.True(t, isSuccessful)
	})

	t.Run("api call unsuccessful", func(*testing.T) {
		resp := _http.Response{StatusCode: 400}
		isSuccessful := IsCallSuccessful(&resp)
		assert.False(t, isSuccessful)
	})

	t.Run("api call returns nothing", func(*testing.T) {
		resp := _http.Response{}
		isSuccessful := IsCallSuccessful(&resp)
		assert.False(t, isSuccessful)
	})
}

func TestCheckAPICallSuccess(t *testing.T) {
	t.Run("context cancelled", func(*testing.T) {
		errMessage := "context cancelled"
		parentCtx := context.Background()
		ctx, cancelCtx := context.WithCancel(parentCtx)
		cancelCtx()
		respBody := _http.Response{Body: io.NopCloser(bytes.NewReader(nil))}
		actualErr := CheckAPICallSuccess(ctx, errMessage, &respBody, errors.New(errMessage))
		assert.True(t, commonerrors.Any(actualErr, commonerrors.ErrCancelled))
	})

	t.Run("api call not successful", func(t *testing.T) {
		errMessage := "client error"
		parentCtx := context.Background()
		resp := _http.Response{StatusCode: 400, Body: io.NopCloser(bytes.NewReader([]byte("{\"message\": \"client error\",\"requestId\": \"761761721\"}")))}
		actualErr := CheckAPICallSuccess(parentCtx, errMessage, &resp, errors.New(errMessage))
		expectedErr := "client error (400): API call error [request-id: 761761721] client error; client error"
		assert.Equal(t, actualErr.Error(), expectedErr)
	})

	t.Run("no context error, api call successful", func(t *testing.T) {
		errMessage := "no error"
		parentCtx := context.Background()
		resp := _http.Response{StatusCode: 200}
		err := CheckAPICallSuccess(parentCtx, errMessage, &resp, errors.New(errMessage))
		assert.NoError(t, err)
	})
}
