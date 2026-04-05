package typeserror

import "fmt"

type SerializationError struct {
	Op    string
	Stage string
	Err   error
}

func (e *SerializationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("headet error [%s]: serialization failed at stage '%s': %v",
			e.Op, e.Stage, e.Err)
	}

	return fmt.Sprintf("header error [%s]: serialization failed at stage '%s'",
		e.Op, e.Stage)
}

func (e *SerializationError) Unwrap() error {
	return e.Err
}
