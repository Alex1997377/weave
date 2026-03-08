package typeserror

import "fmt"

type ValidationError struct {
	Field string
	Value interface{}
	Op    string
	Rule  string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("header error [%s]: validation failed for field '%s' = %v - rule: %s",
		e.Op, e.Field, e.Value, e.Rule)
}
