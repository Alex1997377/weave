package typeserror

import "fmt"

type ValidationErrors struct {
	Op     string
	Errors []error
}

func (e *ValidationErrors) Error() string {
	return fmt.Sprintf("header error [%s]: %d validation errors occurred",
		e.Op, len(e.Errors))
}

func (e *ValidationErrors) Unwrap() []error {
	return e.Errors
}
