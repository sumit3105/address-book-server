package error

import "net/http"

type AppError struct {
	StatusCode int `json:"-"`
	Code string `json:"error"`
	Message string `json:"message"`
	Details map[string]string `json:"details,omitempty"`
	Err error `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func NewError(statusCode int, code string, message string, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

func NewValidationError(details map[string]string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Code:       "VALIDATION_ERROR",
		Message:    "Validation failed",
		Details:    details,
	}
}

func BadRequest(message string, err error) *AppError {
	return NewError(
		http.StatusBadRequest,
		"BAD_REQUEST",
		message,
		err,
	)
}

func NotFound(message string, err error) *AppError {
	return NewError(
		http.StatusNotFound,
		"NOT_FOUND",
		message,
		err,
	)
}

func Unauthorized(message string, err error) *AppError {
	return NewError(
		http.StatusUnauthorized,
		"UNAUTHORIZED",
		message,
		err,
	)
}

func Forbidden(message string, err error) *AppError {
	return NewError(
		http.StatusForbidden,
		"FORBIDDEN",
		message,
		err,
	)
}

func Internal(message string, err error) *AppError {
	return NewError(
		http.StatusInternalServerError,
		"INTERNAL_ERROR",
		message,
		err,
	)
}
