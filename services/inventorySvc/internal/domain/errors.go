package domain

import "errors"

// Common domain errors
var (
	ErrNotFound = errors.New("entity not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInsufficientReservation = errors.New("insufficient reservation")
	ErrDuplicateEntity = errors.New("entity already exists")
	ErrInvalidOperation = errors.New("invalid operation")
)
