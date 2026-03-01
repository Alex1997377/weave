package core

import "fmt"

type BlockchainError struct {
	Code    string
	Message string
	Err     error
}

func (e *BlockchainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *BlockchainError) Unwrap() error {
	return e.Err
}

const (
	ErrInvalidBlock     = "INVALID_BLOCK"
	ErrInvalidHash      = "INVALID_HASH"
	ErrInvalidSignature = "INVALID_SIGNATURE"
	ErrBlockNotFound    = "BLOCK_NOT_FOUND"
	ErrChainCorrupted   = "CHAIN_CORRUPTED"
)

func NewInvalidBlockError(message string, err error) *BlockchainError {
	return &BlockchainError{
		Code:    ErrInvalidBlock,
		Message: message,
		Err:     err,
	}
}

func NewInvalidHashError(message string, err error) *BlockchainError {
	return &BlockchainError{
		Code:    ErrInvalidHash,
		Message: message,
		Err:     err,
	}
}

func NewInvalidSignatureError(message string, err error) *BlockchainError {
	return &BlockchainError{
		Code:    ErrInvalidSignature,
		Message: message,
		Err:     err,
	}
}

func NewBlockNotFoundError(message string, err error) *BlockchainError {
	return &BlockchainError{
		Code:    ErrBlockNotFound,
		Message: message,
		Err:     err,
	}
}

func NewChainCorruptedError(message string, err error) *BlockchainError {
	return &BlockchainError{
		Code:    ErrChainCorrupted,
		Message: message,
		Err:     err,
	}
}
