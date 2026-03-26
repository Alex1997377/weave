package block

import (
	"errors"
	"fmt"

	"github.com/Alex1997377/weave/internal/crypto/merkle"
)

func (b *Block) SetMerkleRoot() error {
	if b == nil {
		return errors.New("block is nil")
	}

	root, err := b.CalculateMerkleRootWithError()
	if err != nil {
		return fmt.Errorf("failed to calculate merkle root: %w", err)
	}

	b.Header.MerkleRoot = root
	return nil
}

// CalculateMerkleRoot вычисляет Merkle root из транзакций блока
func (b *Block) CalculateMerkleRoot() []byte {
	if b == nil || len(b.Transaction) == 0 {
		// Для пустого блока возвращаем нулевой хеш
		return make([]byte, 32)
	}

	// Собираем ID всех транзакций
	var txIDs [][]byte
	for _, tx := range b.Transaction {
		if tx == nil {
			continue
		}
		id := tx.TransactionGetID()
		if id == nil {
			continue
		}
		// Создаем копию, чтобы избежать проблем с ссылками
		idCopy := make([]byte, len(id))
		copy(idCopy, id)
		txIDs = append(txIDs, idCopy)
	}

	// Если нет валидных транзакций, возвращаем нулевой хеш
	if len(txIDs) == 0 {
		return make([]byte, 32)
	}

	// Вычисляем Merkle root через функцию из пакета crypto
	root, err := merkle.CalculateMerkleRoot(txIDs)
	if err != nil {
		// В случае ошибки возвращаем нулевой хеш
		// (в реальном приложении лучше логировать ошибку)
		return make([]byte, 32)
	}

	return root
}

// CalculateMerkleRootWithError вычисляет Merkle root и возвращает ошибку
func (b *Block) CalculateMerkleRootWithError() ([]byte, error) {
	if b == nil {
		return nil, errors.New("block is nil")
	}

	if len(b.Transaction) == 0 {
		return make([]byte, 32), nil
	}

	var txIDs [][]byte
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

		// Создаем копию
		idCopy := make([]byte, len(id))
		copy(idCopy, id)
		txIDs = append(txIDs, idCopy)
	}

	return merkle.CalculateMerkleRoot(txIDs)
}
