package errors

import (
	"errors"
	"fmt"
)

// AppError is a custom error type for the application
type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Error codes
const (
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeBadRequest    = "BAD_REQUEST"
	ErrCodeInternal      = "INTERNAL_ERROR"
	ErrCodeInvalidID     = "INVALID_ID"
	ErrCodeFileTooLarge  = "FILE_TOO_LARGE"
	ErrCodeUnsupported   = "UNSUPPORTED_TYPE"
)

// Constructors
func NotFound(message string) *AppError {
	return &AppError{Code: ErrCodeNotFound, Message: message}
}

func BadRequest(message string) *AppError {
	return &AppError{Code: ErrCodeBadRequest, Message: message}
}

func Internal(message string, err error) *AppError {
	return &AppError{Code: ErrCodeInternal, Message: message, Err: err}
}

func InvalidID(message string) *AppError {
	return &AppError{Code: ErrCodeInvalidID, Message: message}
}

func FileTooLarge(message string) *AppError {
	return &AppError{Code: ErrCodeFileTooLarge, Message: message}
}

func UnsupportedType(message string) *AppError {
	return &AppError{Code: ErrCodeUnsupported, Message: message}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	return errors.As(err, &AppError{})
}

// GetAppError extracts AppError from wrapped errors
func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}
