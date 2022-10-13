package utils

import (
	"bytes"
	"context"
	"errors"
	"io"
	_http "net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAPICallSuccessful(t *testing.T) {
	t.Run("api call successful", func(*testing.T) {
		resp := _http.Response{StatusCode: 200}
		isSuccessful := IsAPICallSuccessful(&resp)
		assert.True(t, isSuccessful)
	})

	t.Run("api call unsuccessful", func(*testing.T) {
		resp := _http.Response{StatusCode: 400}
		isSuccessful := IsAPICallSuccessful(&resp)
		assert.False(t, isSuccessful)
	})

	t.Run("api call returns nothing", func(*testing.T) {
		resp := _http.Response{}
		isSuccessful := IsAPICallSuccessful(&resp)
		assert.False(t, isSuccessful)
	})
}

func TestFetchAPIErrorDescription(t *testing.T) {
	t.Run("error response exists", func(*testing.T) {
		resp := _http.Response{Body: io.NopCloser(bytes.NewReader([]byte("{\"message\": \"client error\",\"requestId\": \"761761721\"}")))}
		actualMessage := fetchAPIErrorDescription(&resp)
		expectedMessage := "client error(request-id:761761721)"
		assert.Equal(t, actualMessage, expectedMessage)
	})

	t.Run("error response has fields", func(*testing.T) {
		resp := _http.Response{Body: io.NopCloser(bytes.NewReader([]byte("{\"message\":\"client error\",\"requestId\":\"761761721\",\"fields\":[{\"fieldName\":\"client request error\",\"fieldPath\":\"https://foo.bar\",\"message\":\"client error\"}]}")))}
		actualMessage := fetchAPIErrorDescription(&resp)
		expectedMessage := "client error(request-id:761761721) [client request error: client error (https://foo.bar),]"
		assert.Equal(t, actualMessage, expectedMessage)
	})

	t.Run("error response empty", func(t *testing.T) {
		respBody := _http.Response{Body: io.NopCloser(bytes.NewReader(nil))}
		actualMessage := fetchAPIErrorDescription(&respBody)
		expectedMessage := ""
		assert.Equal(t, actualMessage, expectedMessage)
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
		expectedErr := "context cancelled: cancelled"
		assert.Equal(t, actualErr.Error(), expectedErr)
	})

	t.Run("api call not successful", func(t *testing.T) {
		errMessage := "client error"
		parentCtx := context.Background()
		resp := _http.Response{StatusCode: 400, Body: io.NopCloser(bytes.NewReader([]byte("{\"message\": \"client error\",\"requestId\": \"761761721\"}")))}
		actualErr := CheckAPICallSuccess(parentCtx, errMessage, &resp, errors.New(errMessage))
		expectedErr := "client error (400): further details: client error(request-id:761761721)"
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
