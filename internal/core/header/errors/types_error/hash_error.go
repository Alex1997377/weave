package typeserror

import "fmt"

type HashError struct {
	Field    string
	Hash     []byte
	Op       string
	Required int
}

func (e *HashError) Error() string {
	if e.Hash == nil {
		return fmt.Sprintf("header error [%s]: field '%s' is nil",
			e.Op, e.Field)
	}

	if len(e.Hash) == 0 {
		return fmt.Sprintf("header error [%s]: field '%s' is empty",
			e.Op, e.Field)
	}

	return fmt.Sprintf("header error [%s]: field '%s' has invalid length %d (expected %d)",
		e.Op, e.Field, len(e.Hash), e.Required)
}
