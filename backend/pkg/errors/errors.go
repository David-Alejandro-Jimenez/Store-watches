package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
    Code    int
    Message string
}

func (e *AppError) Error() string {
    return fmt.Sprintf("(%d) %s", e.Code, e.Message)
}

func NewBadRequestError(message string) *AppError {
    return &AppError{
        Code:    400,
        Message: message,
    }
}

func NewConflictError(message string) *AppError {
    return &AppError{
        Code:    409,
        Message: message,
    }
}

func NewInternalError(message string) *AppError {
    return &AppError{
        Code:    500,
        Message: message,
    }
}

func NewAuthError(message string) *AppError {
    return &AppError{
        Code:    http.StatusUnauthorized,
        Message: message,
    }
}

func NewTooManyRequestsError(message string) *AppError {
    return &AppError{
        Code:    http.StatusTooManyRequests,
        Message: message,
    }
}