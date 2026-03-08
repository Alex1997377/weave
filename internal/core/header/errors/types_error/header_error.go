package typeserror

import "fmt"

type HeaderError struct {
	Op      string
	Field   string
	Value   interface{}
	Message string
	Err     error
}

func (e *HeaderError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("header error [%s]: field '%s' = %v - %s: %v",
			e.Op, e.Field, e.Value, e.Message, e.Err)
	}
	return fmt.Sprintf("header error [%s]: field '%s' = %v - %s",
		e.Op, e.Field, e.Value, e.Message)
}

func (e *HeaderError) Unwrap() error {
	return e.Err
}
