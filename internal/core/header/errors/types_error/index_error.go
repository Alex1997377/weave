package typeserror

import "fmt"

type IndexError struct {
	Index int
	Op    string
}

func (e *IndexError) Error() string {
	return fmt.Sprintf("header error [%s]: invalid index %d (must be >= 0)",
		e.Op, e.Index)
}
