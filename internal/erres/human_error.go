package erres

func NewHumanError(message string, err error) *BaseHumanError {
	return &BaseHumanError{
		Message: message,
		Err:     err,
	}
}

type BaseHumanError struct {
	Message string
	Err     error
}

func (e *BaseHumanError) Error() string {
	return e.Err.Error()
}

func (e *BaseHumanError) HumanError() string {
	return e.Message
}

type HumanError interface {
	error
	HumanError() string
}
