package errors

import (
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	ErrorCodeValidation    ErrorCode = "VALIDATION_ERROR"
	ErrorCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrorCodeConflict      ErrorCode = "CONFLICT"
	ErrorCodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden     ErrorCode = "FORBIDDEN"
	ErrorCodeInternal      ErrorCode = "INTERNAL_ERROR"
	ErrorCodeBadRequest    ErrorCode = "BAD_REQUEST"
	ErrorCodeUnprocessable ErrorCode = "UNPROCESSABLE_ENTITY"
)

type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	StatusCode int       `json:"-"`
}

func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewAppError(code ErrorCode, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

func NewValidationError(message string) *AppError {
	return NewAppError(ErrorCodeValidation, message, http.StatusBadRequest)
}

func NewNotFoundError(message string) *AppError {
	return NewAppError(ErrorCodeNotFound, message, http.StatusNotFound)
}

func NewConflictError(message string) *AppError {
	return NewAppError(ErrorCodeConflict, message, http.StatusConflict)
}

func NewBadRequestError(message string) *AppError {
	return NewAppError(ErrorCodeBadRequest, message, http.StatusBadRequest)
}

func NewInternalError(message string) *AppError {
	return NewAppError(ErrorCodeInternal, message, http.StatusInternalServerError)
}

var (
	ErrInvalidID         = NewBadRequestError("Invalid ID")
	ErrInvalidJSON       = NewBadRequestError("Invalid JSON")
	ErrInvalidInput      = NewBadRequestError("Invalid input")
	ErrResourceNotFound  = NewNotFoundError("Resource not found")
	ErrDuplicateResource = NewConflictError("Resource already exists")
	ErrInternalServer    = NewInternalError("Internal server error")
	ErrUnauthorized      = NewAppError(ErrorCodeUnauthorized, "Unauthorized", http.StatusUnauthorized)
	ErrForbidden         = NewAppError(ErrorCodeForbidden, "Forbidden", http.StatusForbidden)
)

func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return nil
}
