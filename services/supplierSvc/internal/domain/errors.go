package domain

import "errors"

// Common domain errors
var (
	ErrNotFound      = errors.New("resource not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
)
