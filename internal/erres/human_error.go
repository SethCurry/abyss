package erres

// NewHumanError wraps err with a human-readable message.
func NewHumanError(message string, err error) *BaseHumanError {
	return &BaseHumanError{
		Message: message,
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
