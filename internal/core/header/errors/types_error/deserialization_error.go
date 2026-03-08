package typeserror

import "fmt"

type DeserializationError struct {
	Op    string
	Stage string
	Err   error
}

func (e *DeserializationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("header error [%s]: deserialization failed at stage '%s': %v",
			e.Op, e.Stage, e.Err)
	}

	return fmt.Sprintf("header error [%s]: deserialization failed at stage '%s'",
		e.Op, e.Stage)
}
