package errors

import (
	"errors"

	typeserror "github.com/Alex1997377/weave/internal/core/header/errors/types_error"
)

func NewNilHeaderError(op string) *typeserror.NilHeaderError {
	return &typeserror.NilHeaderError{Op: op}
}

func NewIndexError(op string, index int) *typeserror.IndexError {
	return &typeserror.IndexError{Index: index, Op: op}
}

func NewTimestampError(op string, timestamp int64, reason string) *typeserror.TimestampError {
	return &typeserror.TimestampError{Timestamp: timestamp, Op: op, Reason: reason}
}

func NewHashError(op, field string, hash []byte, requiredLen int) *typeserror.HashError {
	return &typeserror.HashError{
		Field:    field,
		Hash:     hash,
		Op:       op,
		Required: requiredLen,
	}
}

func NewDifficultyError(op string, difficulty, min, max int) *typeserror.DifficultyError {
	return &typeserror.DifficultyError{
		Difficulty: difficulty,
		Op:         op,
		Min:        min,
		Max:        max,
	}
}

func NewNonceError(op string, nonce int64) *typeserror.NonceError {
	return &typeserror.NonceError{Nonce: nonce, Op: op}
}

func NewSerializationError(op, stage string, err error) *typeserror.SerializationError {
	return &typeserror.SerializationError{Op: op, Stage: stage, Err: err}
}

func NewDeserializationError(op, stage string, err error) *typeserror.DeserializationError {
	return &typeserror.DeserializationError{Op: op, Stage: stage, Err: err}
}

func NewValidationError(op, field string, value interface{}, rule string) *typeserror.ValidationError {
	return &typeserror.ValidationError{
		Field: field,
		Value: value,
		Op:    op,
		Rule:  rule,
	}
}

func NewValidationErrors(op string, errs []error) *typeserror.ValidationErrors {
	return &typeserror.ValidationErrors{Op: op, Errors: errs}
}

// Check types errors

func IsNilHeaderError(err error) bool {
	var e *typeserror.NilHeaderError
	return errors.As(err, &e)
}

func IsIndexError(err error) bool {
	var e *typeserror.IndexError
	return errors.As(err, &e)
}

func IsTimestampError(err error) bool {
	var e *typeserror.TimestampError
	return errors.As(err, &e)
}

func IsHashError(err error) bool {
	var e *typeserror.HashError
	return errors.As(err, &e)
}

func IsDifficultyError(err error) bool {
	var e *typeserror.DifficultyError
	return errors.As(err, &e)
}

func IsNonceError(err error) bool {
	var e *typeserror.NonceError
	return errors.As(err, &e)
}

func IsSerializationError(err error) bool {
	var e *typeserror.SerializationError
	return errors.As(err, &e)
}

func IsDeserializationError(err error) bool {
	var e *typeserror.DeserializationError
	return errors.As(err, &e)
}

func IsValidationError(err error) bool {
	var e *typeserror.ValidationError
	return errors.As(err, &e)
}

func IsValidationErrors(err error) bool {
	var e *typeserror.ValidationErrors
	return errors.As(err, &e)
}
