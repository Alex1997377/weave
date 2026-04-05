package typeserror

import "fmt"

type DifficultyError struct {
	Difficulty int
	Op         string
	Min        int
	Max        int
}

func (e *DifficultyError) Error() string {
	if e.Difficulty < e.Min {
		return fmt.Sprintf("header error [%s]: difficulty %d is below minimum %d",
			e.Op, e.Difficulty, e.Min)
	}

	return fmt.Sprintf("header error [%s]: difficulty %d exceeds maximum %d",
		e.Op, e.Difficulty, e.Max)
}
