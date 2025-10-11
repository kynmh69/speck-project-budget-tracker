package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewAppError creates a new AppError
func NewAppError(code, message string, statusCode int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// Common errors
var (
	// BadRequest errors
	ErrInvalidInput = func(err error) *AppError {
		return NewAppError("INVALID_INPUT", "Invalid input provided", http.StatusBadRequest, err)
	}

	ErrValidationFailed = func(message string) *AppError {
		return NewAppError("VALIDATION_FAILED", message, http.StatusBadRequest, nil)
	}

	// Unauthorized errors
	ErrUnauthorized = func() *AppError {
		return NewAppError("UNAUTHORIZED", "Authentication required", http.StatusUnauthorized, nil)
	}

	ErrInvalidCredentials = func() *AppError {
		return NewAppError("INVALID_CREDENTIALS", "Invalid email or password", http.StatusUnauthorized, nil)
	}

	ErrInvalidToken = func() *AppError {
		return NewAppError("INVALID_TOKEN", "Invalid or expired token", http.StatusUnauthorized, nil)
	}

	// Forbidden errors
	ErrForbidden = func() *AppError {
		return NewAppError("FORBIDDEN", "You don't have permission to access this resource", http.StatusForbidden, nil)
	}

	// NotFound errors
	ErrNotFound = func(resource string) *AppError {
		return NewAppError("NOT_FOUND", fmt.Sprintf("%s not found", resource), http.StatusNotFound, nil)
	}

	// Conflict errors
	ErrAlreadyExists = func(resource string) *AppError {
		return NewAppError("ALREADY_EXISTS", fmt.Sprintf("%s already exists", resource), http.StatusConflict, nil)
	}

	// Internal errors
	ErrInternal = func(err error) *AppError {
		return NewAppError("INTERNAL_ERROR", "An internal error occurred", http.StatusInternalServerError, err)
	}

	ErrDatabaseError = func(err error) *AppError {
		return NewAppError("DATABASE_ERROR", "Database operation failed", http.StatusInternalServerError, err)
	}
)
