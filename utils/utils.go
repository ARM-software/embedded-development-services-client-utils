package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	_http "net/http"
	"strings"

	"github.com/ARM-software/golang-utils/utils/parallelisation"
	"github.com/ARM-software/golang-utils/utils/reflection"

	"github.com/Arm-Debug/solar-services-client/client"
)

func IsAPICallSuccessful(r *_http.Response) bool {
	if r == nil {
		return false
	}
	return r.StatusCode >= _http.StatusOK && r.StatusCode < _http.StatusMultipleChoices
}

func fetchAPIErrorDescription(resp *_http.Response) (message string) {
	apiErrorMessage := strings.Builder{}
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
		return
	}
	apiErrorMessage.WriteString(errorResponse.GetMessage())
	apiErrorMessage.WriteString(fmt.Sprintf("(request-id:%v)", errorResponse.GetRequestId()))
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

func CheckAPICallSuccess(ctx context.Context, errorDescription string, resp *_http.Response, apierr error) (err error) {
	if err = parallelisation.DetermineContextError(ctx); err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}
	if !IsAPICallSuccessful(resp) {
		statusCode := 0
		errorMessage := strings.Builder{}
		if resp != nil {
			statusCode = resp.StatusCode
			errorDetails := fetchAPIErrorDescription(resp)
			if !reflection.IsEmpty(errorDetails) {
				errorMessage.WriteString("further details: ")
				errorMessage.WriteString(errorDetails)
			}
			_ = resp.Body.Close()
		}
		err = fmt.Errorf("%v (%d): %v", errorDescription, statusCode, errorMessage.String())
	}
	return
}
