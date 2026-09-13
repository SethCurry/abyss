package erres

import "strings"

// NewHumanError wraps err with a human-readable message.
func NewHumanError(err error, messages ...string) *BaseHumanError {
	return &BaseHumanError{
		Message: strings.Join(messages, "\n"),
		Err:     err,
	}
}

// BaseHumanError is the most basic version of HumanError,
// and simply embeds a static human error message inside it.
type BaseHumanError struct {
	Message string
	Err     error
}

// Error returns the underlying error message.
func (e *BaseHumanError) Error() string {
	return e.Err.Error()
}

// HumanError returns the human-readable message.
func (e *BaseHumanError) HumanError() string {
	return e.Message
}

// HumanError is an extension of the built-in error that also bundles
// a more human-readable variant of the message.  I.e. where the normal
// error might say "failed to open file ...: $error message", HumanError would
// contain something more helpful like "failed while trying to read your
// configuration file.  Make sure it exists and has correct permissions."
type HumanError interface {
	error
	HumanError() string
}
