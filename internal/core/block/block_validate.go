package block

import (
	"bytes"
	"errors"
	"fmt"
)

func (b *Block) Validate() error {
	if b == nil {
		return errors.New("block is nil")
	}

	if b.Hash == nil {
		return errors.New("block hash is nil")
	}

	if b.Header.Difficulty < 0 {
		return errors.New("block difficulty cannot be negative")
	}

	if !b.Hash.IsValidForDifficulty(b.Header.Difficulty) {
		return errors.New("invalid proof of work")
	}

	calculatedHash, err := b.CalculateHash()
	if err != nil {
		return fmt.Errorf("failed to calculate hash for validation: %w", err)
	}

	if !bytes.Equal(b.Hash[:], calculatedHash[:]) {
		return errors.New("block hash doesn`t match content")
	}

	for i, tx := range b.Transaction {
		if tx == nil {
			return fmt.Errorf("transaction at index %d is nil", i)
		}

		if err := tx.TransactionValidate(); err != nil {
			return fmt.Errorf("invalid transaction at index %d: %w", i, err)
		}
	}

	return nil
}
