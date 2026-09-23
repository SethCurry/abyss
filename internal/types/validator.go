// Package types defines validation error types used across the abyss project.
package types

import "fmt"

// NewValidationError constructs a ValidationError describing the given field
// and reason for the failure on the provided value.
func NewValidationError(from any, field string, reason string) *ValidationError {
	return &ValidationError{
		Type:   fmt.Sprintf("%T", from),
		Field:  field,
		Reason: reason,
	}
}

// ValidationError is a custom error type that embeds information about
// a failed validation, such as the field that was invalid.
type ValidationError struct {
	// The name of the type that failed validation.  Usually the name of a struct.
	Type string

	// The name of the field that failed to validate.
	Field string

	// A human-readable explanation of why validation failed, like "should be non-zero"
	Reason string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid field in type %s: %s", e.Type, e.Reason)
}

// Validator is an object that can be validated and returns errors
// if the data within the object is invalid.
type Validator interface {
	Validate() error
}
