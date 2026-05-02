package usecase

import "errors"

// Common use case errors.
var (
	ErrNotFound    = errors.New("resource not found")
	ErrValidation  = errors.New("validation failed")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden   = errors.New("forbidden")
)
