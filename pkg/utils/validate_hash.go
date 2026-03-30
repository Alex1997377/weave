package utils

import "github.com/Alex1997377/weave/internal/core/header/errors"

func ValidateHash(op, field string, hash []byte, required bool) error {
	if hash == nil && required {
		return errors.NewHashError(op, field, hash, 32)
	}
	if len(hash) > 0 && len(hash) != 32 {
		return errors.NewHashError(op, field, hash, 32)
	}

	return nil
}
