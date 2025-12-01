/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	_http "net/http"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ARM-software/embedded-development-services-client/client"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/commonerrors/errortest"
	"github.com/ARM-software/golang-utils/utils/field"
)

func TestCheckAPICallSuccess(t *testing.T) {
	t.Run("context cancelled", func(t *testing.T) {
		errMessage := "context cancelled"
		parentCtx := context.Background()
		ctx, cancelCtx := context.WithCancel(parentCtx)
		cancelCtx()
		respBody := _http.Response{Body: io.NopCloser(bytes.NewReader(nil))}
		actualErr := CheckAPICallSuccess(ctx, errMessage, &respBody, errors.New(errMessage))
		errortest.AssertError(t, actualErr, commonerrors.ErrCancelled)

	})

	t.Run("api call not successful", func(t *testing.T) {
		errMessage := "client error"
		parentCtx := context.Background()
		resp := _http.Response{StatusCode: _http.StatusBadRequest, Body: io.NopCloser(bytes.NewReader([]byte("{\"message\": \"client error\",\"requestId\": \"761761721\"}")))}
		actualErr := CheckAPICallSuccess(parentCtx, errMessage, &resp, errors.New(errMessage))
		expectedErr := "client error (400): API call error [request-id: 761761721] client error; client error"
		assert.Contains(t, actualErr.Error(), expectedErr)
		errortest.AssertError(t, actualErr, commonerrors.ErrInvalid)
	})

	t.Run("api call not successful", func(t *testing.T) {
		errMessage := "client error"
		parentCtx := context.Background()
		resp := _http.Response{StatusCode: _http.StatusServiceUnavailable, Body: io.NopCloser(bytes.NewReader([]byte("{\"message\": \"client error\",\"requestId\": \"761761721\"}")))}
		actualErr := CheckAPICallSuccess(parentCtx, errMessage, &resp, errors.New(errMessage))
		expectedErr := "client error (503): API call error [request-id: 761761721] client error; client error"
		assert.Contains(t, actualErr.Error(), expectedErr)
		errortest.AssertError(t, actualErr, commonerrors.ErrUnavailable)
	})

	t.Run("api call not successful", func(t *testing.T) {
		errMessage := "client error"
		parentCtx := context.Background()
		resp := _http.Response{StatusCode: _http.StatusUnauthorized, Body: io.NopCloser(bytes.NewReader([]byte("{\"message\": \"client error\",\"requestId\": \"761761721\"}")))}
		actualErr := CheckAPICallSuccess(parentCtx, errMessage, &resp, errors.New(errMessage))
		expectedErr := "client error (401): API call error [request-id: 761761721] client error; client error"
		assert.Contains(t, actualErr.Error(), expectedErr)
		errortest.AssertError(t, actualErr, commonerrors.ErrUnauthorised)
	})

	t.Run("api call not successful (no JSON response)", func(t *testing.T) {
		errMessage := "response error"
		parentCtx := context.Background()
		resp := _http.Response{StatusCode: _http.StatusForbidden, Body: io.NopCloser(bytes.NewReader([]byte("<html><head><title>403 Forbidden</title></head></html>")))}
		actualErr := CheckAPICallSuccess(parentCtx, errMessage, &resp, errors.New("403 Forbidden"))
		expectedErr := "response error (403): <html><head><title>403 Forbidden</title></head></html>; 403 Forbidden"
		assert.Contains(t, actualErr.Error(), expectedErr)
		errortest.AssertError(t, actualErr, commonerrors.ErrForbidden)
	})

	t.Run("no context error, api call successful", func(t *testing.T) {
		errMessage := "no error"
		parentCtx := context.Background()
		resp := _http.Response{StatusCode: 200}
		err := CheckAPICallSuccess(parentCtx, errMessage, &resp, errors.New(errMessage))
		assert.NoError(t, err)
	})
}

func TestCallAndCheckSuccess(t *testing.T) {
	t.Run("context cancelled", func(t *testing.T) {
		errMessage := "context cancelled"
		parentCtx := context.Background()
		ctx, cancelCtx := context.WithCancel(parentCtx)
		cancelCtx()
		_, actualErr := CallAndCheckSuccess(ctx, errMessage,
			func(ctx context.Context) (*struct{}, *_http.Response, error) {
				return nil, &_http.Response{Body: io.NopCloser(bytes.NewReader(nil))}, nil
			})
		errortest.AssertError(t, actualErr, commonerrors.ErrCancelled)
	})

	t.Run("api call not successful", func(t *testing.T) {
		errMessage := "client error"
		parentCtx := context.Background()
		_, actualErr := CallAndCheckSuccess(parentCtx, errMessage,
			func(ctx context.Context) (*struct{}, *_http.Response, error) {
				resp := _http.Response{StatusCode: 400, Body: io.NopCloser(bytes.NewReader([]byte("{\"message\": \"client error\",\"requestId\": \"761761721\"}")))}
				return nil, &resp, errors.New(errMessage)
			})
		expectedErr := "client error (400): API call error [request-id: 761761721] client error; client error"
		assert.Contains(t, actualErr.Error(), expectedErr)
		errortest.AssertError(t, actualErr, commonerrors.ErrInvalid)
	})

	t.Run("api call successful, marshalling failed due to missing required field in response", func(t *testing.T) {
		expectedErrorMessage := client.ErrorResponse{
			Fields: []client.FieldObject{{
				FieldName: faker.Name(),
				FieldPath: field.ToOptionalString(faker.Name()),
				Message:   faker.Sentence(),
			}},
			HttpStatusCode: 200,
			Message:        faker.Sentence(),
			RequestId:      faker.UUIDDigit(),
		}
		response, err := expectedErrorMessage.ToMap()
		require.NoError(t, err)
		delete(response, "message")

		reducedResponse, err := json.Marshal(response)
		require.NoError(t, err)

		errorResponse := client.ErrorResponse{}
		errM := errorResponse.UnmarshalJSON(reducedResponse)
		require.Error(t, errM)

		parentCtx := context.Background()
		_, err = CallAndCheckSuccess[client.ErrorResponse](parentCtx, "test",
			func(ctx context.Context) (*client.ErrorResponse, *_http.Response, error) {
				return &errorResponse, &_http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(reducedResponse))}, errM
			})
		require.Error(t, err)
		errortest.AssertError(t, err, commonerrors.ErrMarshalling)
	})

	t.Run("api call successful, strict marshalling failed but recovery", func(t *testing.T) {
		expectedErrorMessage := client.ErrorResponse{
			Fields: []client.FieldObject{{
				FieldName: faker.Name(),
				FieldPath: field.ToOptionalString(faker.Name()),
				Message:   faker.Sentence(),
			}},
			HttpStatusCode: 200,
			Message:        faker.Sentence(),
			RequestId:      faker.UUIDDigit(),
		}
		response, err := expectedErrorMessage.ToMap()
		require.NoError(t, err)
		response[faker.Word()] = faker.Name()
		response[faker.Word()] = faker.Sentence()
		response[faker.Word()] = faker.Paragraph()
		response[faker.Word()] = faker.UUIDDigit()
		extendedResponse, err := json.Marshal(response)
		require.NoError(t, err)
		errMessage := "no error"
		parentCtx := context.Background()
		result, err := CallAndCheckSuccess[client.ErrorResponse](parentCtx, errMessage,
			func(ctx context.Context) (*client.ErrorResponse, *_http.Response, error) {
				return &client.ErrorResponse{}, &_http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(extendedResponse))}, errors.New(errMessage)
			})
		require.NoError(t, err)
		assert.Equal(t, expectedErrorMessage, *result)
	})

	t.Run("api call successful, empty response", func(t *testing.T) {
		errMessage := "no error"
		parentCtx := context.Background()
		_, err := CallAndCheckSuccess(parentCtx, errMessage,
			func(ctx context.Context) (*struct{}, *_http.Response, error) {
				return &struct{}{}, &_http.Response{StatusCode: 200}, errors.New(errMessage)
			})
		errortest.AssertError(t, err, commonerrors.ErrMarshalling)
		errortest.AssertErrorDescription(t, err, "API call was successful but an error occurred during response marshalling")
	})

	t.Run("api call successful, broken response decode", func(t *testing.T) {
		errMessage := "no error"
		parentCtx := context.Background()
		_, err := CallAndCheckSuccess(parentCtx, errMessage,
			func(ctx context.Context) (*struct{}, *_http.Response, error) {
				return &struct{}{}, &_http.Response{StatusCode: 200}, nil
			})
		errortest.AssertError(t, err, commonerrors.ErrMarshalling)
		errortest.AssertErrorDescription(t, err, "unmarshalled response is empty")
	})
}

func TestGenericCallAndCheckSuccess(t *testing.T) {
	t.Run("context cancelled", func(t *testing.T) {
		errMessage := "context cancelled"
		parentCtx := context.Background()
		ctx, cancelCtx := context.WithCancel(parentCtx)
		cancelCtx()
		_, actualErr := GenericCallAndCheckSuccess(ctx, errMessage,
			func(ctx context.Context) (*struct{}, *_http.Response, error) {
				return nil, &_http.Response{Body: io.NopCloser(bytes.NewReader(nil))}, errors.New(errMessage)
			})
		errortest.AssertError(t, actualErr, commonerrors.ErrCancelled)
	})

	t.Run("api call not successful", func(t *testing.T) {
		errMessage := "client error"
		parentCtx := context.Background()
		_, actualErr := GenericCallAndCheckSuccess(parentCtx, errMessage,
			func(ctx context.Context) (*struct{}, *_http.Response, error) {
				resp := _http.Response{StatusCode: 400, Body: io.NopCloser(bytes.NewReader([]byte("{\"message\": \"client error\",\"requestId\": \"761761721\"}")))}
				return nil, &resp, errors.New(errMessage)
			})
		expectedErr := "client error (400): API call error [request-id: 761761721] client error; client error"
		assert.Contains(t, actualErr.Error(), expectedErr)
		errortest.AssertError(t, actualErr, commonerrors.ErrInvalid)
	})

	t.Run("api call successful but error marshalling", func(t *testing.T) {
		errMessage := "no error"
		parentCtx := context.Background()
		_, err := GenericCallAndCheckSuccess(parentCtx, errMessage,
			func(ctx context.Context) (any, *_http.Response, error) {
				tmp := struct {
					test string
				}{
					test: faker.Word(),
				}
				return &tmp, &_http.Response{StatusCode: 200}, errors.New(errMessage)
			})
		require.Error(t, err)
		errortest.AssertError(t, err, commonerrors.ErrMarshalling)
	})

	t.Run("api call successful, empty response", func(t *testing.T) {
		errMessage := "response error"
		parentCtx := context.Background()
		_, err := GenericCallAndCheckSuccess(parentCtx, errMessage,
			func(ctx context.Context) (*struct{}, *_http.Response, error) {
				return &struct{}{}, &_http.Response{StatusCode: 200}, errors.New(errMessage)
			})
		errortest.AssertError(t, err, commonerrors.ErrMarshalling)
	})

	t.Run("api call successful, incorrect response", func(t *testing.T) {
		parentCtx := context.Background()
		_, err := GenericCallAndCheckSuccess(parentCtx, "response error",
			func(ctx context.Context) (struct{ Blah string }, *_http.Response, error) {
				return struct{ Blah string }{Blah: "fsadsfs"}, &_http.Response{StatusCode: 200}, nil
			})
		errortest.AssertError(t, err, commonerrors.ErrConflict)
	})
}
