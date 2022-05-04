// nolint:gomnd
package apierrors

import (
	"net/http"

	status2 "github.com/alkaid/goerrors/apierrors/http/status"
)

// BadRequest new BadRequest error that is mapped to a 400 response.
func BadRequest(reason, message string, pretty string) *Error {
	return New(http.StatusBadRequest, reason, message, "")
}

// IsBadRequest determines if err is an error which indicates a BadRequest error.
// It supports wrapped errors.
func IsBadRequest(err error) bool {
	return Code(err) == http.StatusBadRequest
}

// Unauthorized new Unauthorized error that is mapped to a 401 response.
func Unauthorized(reason, message string, pretty string) *Error {
	return New(http.StatusUnauthorized, reason, message, "")
}

// IsUnauthorized determines if err is an error which indicates a Unauthorized error.
// It supports wrapped errors.
func IsUnauthorized(err error) bool {
	return Code(err) == http.StatusUnauthorized
}

// Forbidden new Forbidden error that is mapped to a 403 response.
func Forbidden(reason, message string, pretty string) *Error {
	return New(http.StatusForbidden, reason, message, "")
}

// IsForbidden determines if err is an error which indicates a Forbidden error.
// It supports wrapped errors.
func IsForbidden(err error) bool {
	return Code(err) == http.StatusForbidden
}

// NotFound new NotFound error that is mapped to a 404 response.
func NotFound(reason, message string, pretty string) *Error {
	return New(http.StatusNotFound, reason, message, "")
}

// IsNotFound determines if err is an error which indicates an NotFound error.
// It supports wrapped errors.
func IsNotFound(err error) bool {
	return Code(err) == http.StatusNotFound
}

// Conflict new Conflict error that is mapped to a 409 response.
func Conflict(reason, message string, pretty string) *Error {
	return New(409, reason, message, "")
}

// IsConflict determines if err is an error which indicates a Conflict error.
// It supports wrapped errors.
func IsConflict(err error) bool {
	return Code(err) == 409
}

// InternalServer new InternalServer error that is mapped to a 500 response.
func InternalServer(reason, message string, pretty string) *Error {
	return New(500, reason, message, "")
}

// IsInternalServer determines if err is an error which indicates an Internal error.
// It supports wrapped errors.
func IsInternalServer(err error) bool {
	return Code(err) == 500
}

// ServiceUnavailable new ServiceUnavailable error that is mapped to a HTTP 503 response.
func ServiceUnavailable(reason, message string, pretty string) *Error {
	return New(503, reason, message, "")
}

// IsServiceUnavailable determines if err is an error which indicates a Unavailable error.
// It supports wrapped errors.
func IsServiceUnavailable(err error) bool {
	return Code(err) == 503
}

// GatewayTimeout new GatewayTimeout error that is mapped to a HTTP 504 response.
func GatewayTimeout(reason, message string, pretty string) *Error {
	return New(504, reason, message, "")
}

// IsGatewayTimeout determines if err is an error which indicates a GatewayTimeout error.
// It supports wrapped errors.
func IsGatewayTimeout(err error) bool {
	return Code(err) == 504
}

// ClientClosed new ClientClosed error that is mapped to a HTTP 499 response.
func ClientClosed(reason, message string, pretty string) *Error {
	return New(status2.ClientClosed, reason, message, "")
}

// IsClientClosed determines if err is an error which indicates a IsClientClosed error.
// It supports wrapped errors.
func IsClientClosed(err error) bool {
	return Code(err) == status2.ClientClosed
}
