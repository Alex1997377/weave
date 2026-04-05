package typeserror

import "fmt"

type NilHeaderError struct {
	Op string
}

func (e *NilHeaderError) Error() string {
	return fmt.Sprintf("header error [%s]: header is nil", e.Op)
}
