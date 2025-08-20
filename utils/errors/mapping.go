package errors

import (
	"net/http"

	"github.com/ARM-software/golang-utils/utils/commonerrors"
)

// MapErrorToHTTPResponseCode maps a response status code to a common error.
func MapErrorToHTTPResponseCode(statusCode int) error {
	if statusCode < http.StatusBadRequest {
		return nil
	}
	switch statusCode {
	case http.StatusBadRequest:
		return commonerrors.ErrInvalid
	case http.StatusUnauthorized:
		return commonerrors.ErrUnauthorised
	case http.StatusPaymentRequired:
		return commonerrors.ErrUnknown
	case http.StatusForbidden:
		return commonerrors.ErrForbidden
	case http.StatusNotFound:
		return commonerrors.ErrNotFound
	case http.StatusMethodNotAllowed:
		return commonerrors.ErrNotFound
	case http.StatusNotAcceptable:
		return commonerrors.ErrUnsupported
	case http.StatusProxyAuthRequired:
		return commonerrors.ErrUnauthorised
	case http.StatusRequestTimeout:
		return commonerrors.ErrTimeout
	case http.StatusConflict:
		return commonerrors.ErrConflict
	case http.StatusGone:
		return commonerrors.ErrNotFound
	case http.StatusLengthRequired:
		return commonerrors.ErrInvalid
	case http.StatusPreconditionFailed:
		return commonerrors.ErrCondition
	case http.StatusRequestEntityTooLarge:
		return commonerrors.ErrTooLarge
	case http.StatusRequestURITooLong:
		return commonerrors.ErrTooLarge
	case http.StatusUnsupportedMediaType:
		return commonerrors.ErrUnsupported
	case http.StatusRequestedRangeNotSatisfiable:
		return commonerrors.ErrOutOfRange
	case http.StatusExpectationFailed:
		return commonerrors.ErrUnsupported
	case http.StatusTeapot:
		return commonerrors.ErrUnknown
	case http.StatusMisdirectedRequest:
		return commonerrors.ErrUnsupported
	case http.StatusUnprocessableEntity:
		return commonerrors.ErrMarshalling
	case http.StatusLocked:
		return commonerrors.ErrLocked
	case http.StatusFailedDependency:
		return commonerrors.ErrFailed
	case http.StatusTooEarly:
		return commonerrors.ErrUnexpected
	case http.StatusUpgradeRequired:
		return commonerrors.ErrUnsupported
	case http.StatusPreconditionRequired:
		return commonerrors.ErrCondition
	case http.StatusTooManyRequests:
		return commonerrors.ErrUnavailable
	case http.StatusRequestHeaderFieldsTooLarge:
		return commonerrors.ErrTooLarge
	case http.StatusUnavailableForLegalReasons:
		return commonerrors.ErrUnavailable

	case http.StatusInternalServerError:
		return commonerrors.ErrUnexpected
	case http.StatusNotImplemented:
		return commonerrors.ErrNotImplemented
	case http.StatusBadGateway:
		return commonerrors.ErrUnavailable
	case http.StatusServiceUnavailable:
		return commonerrors.ErrUnavailable
	case http.StatusGatewayTimeout:
		return commonerrors.ErrTimeout
	case http.StatusHTTPVersionNotSupported:
		return commonerrors.ErrUnsupported
	case http.StatusVariantAlsoNegotiates:
		return commonerrors.ErrUnexpected
	case http.StatusInsufficientStorage:
		return commonerrors.ErrUnexpected
	case http.StatusLoopDetected:
		return commonerrors.ErrUnexpected
	case http.StatusNotExtended:
		return commonerrors.ErrUnexpected
	case http.StatusNetworkAuthenticationRequired:
		return commonerrors.ErrUnauthorised
	default:
		return commonerrors.ErrUnexpected
	}
}
