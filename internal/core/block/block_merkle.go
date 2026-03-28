package block

import (
	"errors"
	"fmt"

	"github.com/Alex1997377/weave/internal/crypto/merkle"
)

func (b *Block) collectTransactionIDs() ([][]byte, error) {
	if b == nil {
		return nil, errors.New("block is nil")
	}
	if len(b.Transaction) == 0 {
		return [][]byte{}, nil
	}

	txIDs := make([][]byte, 0, len(b.Transaction))
	for i, tx := range b.Transaction {
		if tx == nil {
			return nil, fmt.Errorf("transaction at index %d is nil", i)
		}

		id := tx.TransactionGetID()
		if id == nil {
			return nil, fmt.Errorf("transaction at index %d has nil ID", i)
		}
		if len(id) != 32 {
			return nil, fmt.Errorf("transaction at index %d has invalid ID length: %d", i, len(id))
		}

		idCopy := make([]byte, len(id))
		copy(idCopy, id)
		txIDs = append(txIDs, idCopy)
	}
	return txIDs, nil
}

// CalculateMerkleRootWithError вычисляет Merkle root и возвращает ошибку
func (b *Block) CalculateMerkleRootWithError() ([]byte, error) {
	txIDs, err := b.collectTransactionIDs()
	if err != nil {
		return nil, err
	}
	if len(txIDs) == 0 {
		return make([]byte, 32), nil
	}
	return merkle.CalculateMerkleRoot(txIDs)
}

// CalculateMerkleRoot вычисляет Merkle root из транзакций блока
func (b *Block) CalculateMerkleRoot() []byte {
	root, err := b.CalculateMerkleRootWithError()
	if err != nil {
		return make([]byte, 32)
	}
	return root
}

// SetMerkleRoot вычисляет Merkle root для блока и устанавливает его в заголовок.
func (b *Block) SetMerkleRoot() error {
	root, err := b.CalculateMerkleRootWithError()
	if err != nil {
		return fmt.Errorf("failed to calculate merkle root: %w", err)
	}

	b.Header.MerkleRoot = root
	return nil
}
