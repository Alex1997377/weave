package block

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
)

func (b *Block) Serialize() ([]byte, error) {
	if b == nil {
		return nil, errors.New("block is nil")
	}

	buf := new(bytes.Buffer)

	headerBytes, err := b.Header.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize headerL %w", err)
	}

	_, err = buf.Write(headerBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to write header to buffer: %w", err)
	}

	for i, tx := range b.Transaction {
		if tx == nil {
			return nil, fmt.Errorf("transaction at index %d is nil", i)
		}

		txBytes, err := tx.TransactionSerialize()
		if err != nil {
			return nil, fmt.Errorf("failed to serialize transaction %d:", err)
		}

		_, err = buf.Write(txBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to write transaction %d to buffer: %w", i, err)
		}
	}

	return buf.Bytes(), nil
}

func (b *Block) CalculateHash() ([]byte, error) {
	if b == nil {
		return nil, errors.New("block is nil")
	}

	data, err := b.Header.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize header: %w", err)
	}

	hash := sha256.Sum256(data)
	return hash[:], nil
}

func (b *Block) CalculateSize() (uint32, error) {
	if b == nil {
		return 0, errors.New("block is nil")
	}

	headerBytes, err := b.Header.Serialize()
	if err != nil {
		return 0, fmt.Errorf("failed to serialize header for size calculation: %w", err)
	}
	headerSize := uint32(len(headerBytes))

	var transactionsSize uint32 = 0
	for i, tx := range b.Transaction {
		if tx == nil {
			return 0, fmt.Errorf("transaction at index %d is nil during size calculation", i)
		}

		txBytes, err := tx.TransactionSerialize()
		if err != nil {
			return 0, fmt.Errorf("failed to serialize transaction %d for size calculation: %w", i, err)
		}

		transactionsSize += uint32(len(txBytes))
	}

	hashSize := uint32(len(b.Hash))

	return headerSize + transactionsSize + hashSize, nil
}

