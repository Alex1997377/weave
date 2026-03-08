package typeserror

import "fmt"

type NonceError struct {
	Nonce int64
	Op    string
}

func (e *NonceError) Error() string {
	return fmt.Sprintf("header error [%s]: invalid nonce %d (must be >= 0)",
		e.Op, e.Nonce)
}
