package serialize

import (
	"errors"
	"fmt"
)

var (
	ErrNilPreviousHash    = errors.New("previous hash cannot be nil")
	ErrEmptyPreviousHash  = errors.New("previous hash cannot be empty")
	ErrInvalidHashLength  = errors.New("invalid hash length")
	ErrNegativeIndex      = errors.New("index cannot be negative")
	ErrNegativeTimestamp  = errors.New("timestamp cannot be negative")
	ErrNegativeNonce      = errors.New("nonce cannot be negative")
	ErrNegativeDifficulty = errors.New("difficulty cannot be negative")
)

type SerializeHeaderError struct {
	Field string
	Value interface{}
	Err   error
}

func (e *SerializeHeaderError) Error() string {
	if e.Value != nil {
		return fmt.Sprintf("failed to serialize header field '%s' (value: %v): %v",
			e.Field, e.Value, e.Err)
	}

	return fmt.Sprintf("failed to serialize header field '%s': %v", e.Field, e.Err)
}

func (e *SerializeHeaderError) Unwrap() error {
	return e.Err
}

type HeaderValidationError struct {
	Field string
	Value interface{}
	Err   error
}

func (e *HeaderValidationError) Error() string {
	return fmt.Sprintf("header validation failed for field '%s' (value: %v): %v",
		e.Field, e.Value, e.Err)
}

func (e *HeaderValidationError) Unwrap() error {
	return e.Err
}
