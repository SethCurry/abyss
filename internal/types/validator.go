package types

import "fmt"

func NewValidationError(from any, field string, reason string) *ValidationError {
	return &ValidationError{
		Type:   fmt.Sprintf("%T", from),
		Field:  field,
		Reason: reason,
	}
}

type ValidationError struct {
	Type   string
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid field in type %s: %s", e.Type, e.Reason)
}

type Validator interface {
	Validate() error
}
