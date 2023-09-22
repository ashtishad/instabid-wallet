package lib

import (
	"fmt"
	"net/http"
)

// APIError represents an custom api error,
// Code() returns http status code as integer.
// Error() returns customized string with hiding internal error.
// WithCauses includes internal actual error as causes.
// Wrap() is for manually wrapping actual error to api error, which not included in Error() method.
type APIError interface {
	Error() string
	WithCauses() string
	Wrap(err error) APIError
	Code() int
}

// apiError is a concrete implementation of the APIError interface.
type apiError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status"`
	Causes     string `json:"causes"`
}

// Code returns http status code
func (e *apiError) Code() int {
	return e.StatusCode
}

// Error returns error message and code. But hides internal server/db related errors
func (e *apiError) Error() string {
	return e.Message
}

// WithCauses returns error message, code and internal server/db related errors
func (e *apiError) WithCauses() string {
	return fmt.Sprintf("message: %s - status: %d - causes: %s",
		e.Message, e.StatusCode, e.Causes)
}

// Wrap wraps internal causes of error into an APIError
// use-cases: adding internal causes into error.
func (e *apiError) Wrap(err error) APIError {
	if err != nil {
		e.Causes = err.Error()
	}

	return e
}

// InternalServerError creates a new APIError for internal server errors.
// returns http.StatusInternalServerError 500.
// Example usage:
//
//	err := InternalServerError("internal server error", innerErr)
func InternalServerError(message string, err error) APIError {
	result := &apiError{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}

	return result.Wrap(err)
}

// BadRequestError creates a new APIError for bad requests.
// returns http.StatusBadRequest 400.
// Example usage:
//
//	err := BadRequestError("invalid input").Wrap(innerErr)
func BadRequestError(message string) APIError {
	return &apiError{
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NotFoundError creates a new APIError for not found errors.
// returns http.StatusNotFound 404.
// Example usage:
//
//	err := NotFoundError("resource not found")
func NotFoundError(message string) APIError {
	return &apiError{
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

// UnauthorizedError creates a new APIError for unauthorized requests.
// returns http.StatusUnauthorized 401.
// Example usage:
//
//	err := UnauthorizedError("unauthorized")
func UnauthorizedError(message string) APIError {
	return &apiError{
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// RateLimitError creates a new APIError for rate limit error.
// returns http.StatusTooManyRequests 429.
// Example usage:
//
//	err := RateLimitError("api request limited")
func RateLimitError(message string) APIError {
	return &apiError{
		Message:    message,
		StatusCode: http.StatusTooManyRequests,
	}
}

// ConflictError creates a new APIError for duplicate fields,
// returns http.StatusConflict 409.
// Example usage:
//
//	err := ConflictError("user name already exists")
func ConflictError(message string) APIError {
	return &apiError{
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}
