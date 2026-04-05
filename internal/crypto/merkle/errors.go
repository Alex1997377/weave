package merkle

import (
	"errors"
	"fmt"
)

var (
	ErrNilHashIsSlice      = errors.New("nil hash in slice")
	ErrEmptyHash           = errors.New("empty hash")
	ErrInvalidHashLength   = errors.New("invalid hash length")
	ErrMerkleRootCalculate = errors.New("failed to calculate merkle root")
)

type MerkleRootError struct {
	Op      string
	Index   int
	Message string
	Err     error
}

func (e *MerkleRootError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("merkle root error [%s] at index %d: %s: %v",
			e.Op, e.Index, e.Message, e.Err)
	}
	return fmt.Sprintf("merkle root error [%s] at index %d: %s",
		e.Op, e.Index, e.Message)
}

func (e *MerkleRootError) Unwrap() error {
	return e.Err
}
