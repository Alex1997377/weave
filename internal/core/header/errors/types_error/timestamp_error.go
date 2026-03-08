package typeserror

import "fmt"

type TimestampError struct {
	Timestamp int64
	Op        string
	Reason    string
}

func (e *TimestampError) Error() string {
	return fmt.Sprintf("header error [%s]: invalid timestamp %d - %s",
		e.Op, e.Timestamp, e.Reason)
}
