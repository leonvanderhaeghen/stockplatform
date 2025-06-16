package errors

import (
	"errors"
	"fmt"
)

// Common error types following Uber Go style guide
var (
	// ErrNotFound indicates a resource was not found
	ErrNotFound = errors.New("resource not found")
	
	// ErrAlreadyExists indicates a resource already exists
	ErrAlreadyExists = errors.New("resource already exists")
	
	// ErrInvalidInput indicates invalid input parameters
	ErrInvalidInput = errors.New("invalid input")
	
	// ErrUnauthorized indicates unauthorized access
	ErrUnauthorized = errors.New("unauthorized")
	
	// ErrForbidden indicates forbidden access
	ErrForbidden = errors.New("forbidden")
	
	// ErrInternal indicates an internal server error
	ErrInternal = errors.New("internal server error")
	
	// ErrTimeout indicates a timeout occurred
	ErrTimeout = errors.New("operation timed out")
	
	// ErrUnavailable indicates service is unavailable
	ErrUnavailable = errors.New("service unavailable")
)

// Error represents a structured error with context
type Error struct {
	Code    string
	Message string
	Cause   error
	Context map[string]interface{}
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying cause
func (e *Error) Unwrap() error {
	return e.Cause
}

// Is checks if the error matches the target
func (e *Error) Is(target error) bool {
	if target == nil {
		return false
	}
	
	if e.Cause != nil && errors.Is(e.Cause, target) {
		return true
	}
	
	return e.Message == target.Error()
}

// New creates a new error with the given message
func New(message string) error {
	return &Error{
		Message: message,
	}
}

// Newf creates a new error with formatted message
func Newf(format string, args ...interface{}) error {
	return &Error{
		Message: fmt.Sprintf(format, args...),
	}
}

// Wrap wraps an error with additional context
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	
	return &Error{
		Message: message,
		Cause:   err,
	}
}

// Wrapf wraps an error with formatted message
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	
	return &Error{
		Message: fmt.Sprintf(format, args...),
		Cause:   err,
	}
}

// WithCode adds an error code
func WithCode(err error, code string) error {
	if err == nil {
		return nil
	}
	
	if e, ok := err.(*Error); ok {
		e.Code = code
		return e
	}
	
	return &Error{
		Code:    code,
		Message: err.Error(),
		Cause:   err,
	}
}

// WithContext adds context to an error
func WithContext(err error, key string, value interface{}) error {
	if err == nil {
		return nil
	}
	
	var e *Error
	if errors.As(err, &e) {
		if e.Context == nil {
			e.Context = make(map[string]interface{})
		}
		e.Context[key] = value
		return e
	}
	
	return &Error{
		Message: err.Error(),
		Cause:   err,
		Context: map[string]interface{}{
			key: value,
		},
	}
}

// GetCode extracts error code from error
func GetCode(err error) string {
	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}
	return ""
}

// GetContext extracts context from error
func GetContext(err error) map[string]interface{} {
	var e *Error
	if errors.As(err, &e) {
		return e.Context
	}
	return nil
}
