package domain

import "errors"

// Common domain errors
var (
	// ErrNotFound is returned when a resource is not found
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidID is returned when an invalid ID is provided
	ErrInvalidID = errors.New("invalid ID format")


	// ErrValidation is returned when validation fails
	ErrValidation = errors.New("validation error")


	// ErrAlreadyExists is returned when a resource already exists
	ErrAlreadyExists = errors.New("resource already exists")


	// ErrInternal is returned when an internal error occurs
	ErrInternal = errors.New("internal server error")
)
