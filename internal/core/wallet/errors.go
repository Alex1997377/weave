package wallet

import (
	"errors"
	"fmt"
)

type WalletError struct {
	Op      string
	Err     error
	Message string
}

func (e *WalletError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("wallet: %s: %s", e.Op, e.Message)
	}
	if e.Err != nil {
		return fmt.Sprintf("wallet: %s: %v", e.Op, e.Err)
	}
	return fmt.Sprintf("wallet: %s", e.Op)
}

func (e *WalletError) Unwrap() error {
	return e.Err
}

func NewWalletError(op string, msg string, err error) *WalletError {
	return &WalletError{
		Op:      op,
		Message: msg,
		Err:     err,
	}
}

type CreateWalletError struct {
	WalletError
}

func NewCreateWalletError(msg string, err error) *CreateWalletError {
	return &CreateWalletError{
		WalletError: *NewWalletError("create", msg, err),
	}
}

type SaveWalletError struct {
	WalletError
	Filename string
}

func NewSaveWalletError(filename, msg string, err error) *SaveWalletError {
	return &SaveWalletError{
		WalletError: *NewWalletError("save", msg, err),
		Filename:    filename,
	}
}

func (e *SaveWalletError) Error() string {
	return fmt.Sprintf("wallet: save: failed to save to '%s': %s", e.Filename, e.Message)
}

type LoadWalletError struct {
	WalletError
	Filename string
}

func NewLoadWalletError(filename, msg string, err error) *LoadWalletError {
	return &LoadWalletError{
		WalletError: *NewWalletError("load", msg, err),
		Filename:    filename,
	}
}

func (e *LoadWalletError) Error() string {
	return fmt.Sprintf("wallet: load: failed to load from '%s': %s", e.Filename, e.Message)
}

type SignTransactionError struct {
	WalletError
}

func NewSignTransactionError(msg string, err error) *SignTransactionError {
	return &SignTransactionError{
		WalletError: *NewWalletError("sign", msg, err),
	}
}

type InvalidWalletError struct {
	WalletError
	Reason string
}

func NewInvalidWalletError(reason string) *InvalidWalletError {
	return &InvalidWalletError{
		WalletError: *NewWalletError("validate", "invalid wallet", nil),
		Reason:      reason,
	}
}

func (e *InvalidWalletError) Error() string {
	return fmt.Sprintf("wallet: validate: %s", e.Reason)
}

func IsWalletError(err error) bool {
	var walletErr *WalletError
	return err != nil && errors.As(err, &walletErr)
}

func IsCreateWalletError(err error) bool {
	var createErr *CreateWalletError
	return err != nil && errors.As(err, &createErr)
}
