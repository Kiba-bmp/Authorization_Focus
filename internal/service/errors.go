package service

import (
	"errors"
	"net/http"
)

type AppError struct {
	StatusCode int
	Message    string
}

func (e *AppError) Error() string {
	return e.Message
}

func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		return nil, false
	}
	return appErr, true
}

func badRequest(message string) error {
	return &AppError{StatusCode: http.StatusBadRequest, Message: message}
}

func conflict(message string) error {
	return &AppError{StatusCode: http.StatusConflict, Message: message}
}

func unauthorized(message string) error {
	return &AppError{StatusCode: http.StatusUnauthorized, Message: message}
}
