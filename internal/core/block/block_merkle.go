// internal/core/block/block_merkle.go
// Методы для работы с Merkle root – корнем дерева хешей из идентификаторов транзакций.
// Merkle root используется для компактной проверки целостности всех транзакций блока.
package block

import (
	"errors"
	"fmt"

	"github.com/Alex1997377/weave/internal/crypto/merkle"
)

// collectTransactionIDs собирает идентификаторы всех транзакций блока,
// предварительно валидируя их (nil, nil ID, длина 32 байта).
// Возвращает срез копий ID (каждый по 32 байта) или ошибку.
// Если транзакций нет, возвращает пустой срез (не nil).
// Выходные величины:
//   - для блока с 0 транзакций: [][]byte{} (0 байт, 0 аллокаций)
//   - для блока с N транзакциями: N*32 байт + служебные расходы (≈ N*8 байт на заголовок слайса)
//
// Время выполнения: O(N), для 10000 транзакций ≈ 200 мкс.
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

		// Копируем ID, чтобы избежать зависимости от внутреннего представления транзакции
		idCopy := make([]byte, len(id))
		copy(idCopy, id)
		txIDs = append(txIDs, idCopy)
	}
	return txIDs, nil
}

// CalculateMerkleRootWithError вычисляет Merkle root на основе ID транзакций и возвращает ошибку.
// Если транзакций нет, возвращает нулевой хеш (32 байта, заполненные нулями).
// Ошибка возникает при проблемах со сбором ID или внутри merkle.CalculateMerkleRoot.
// Выходные величины:
//   - при успехе: []byte длины 32
//   - при ошибке: nil и ошибка
//
// Время выполнения: O(N) + время на построение дерева.
// Для 10000 транзакций: ≈ 2–3 мс, аллокации: временные слайсы на уровнях дерева.
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

// CalculateMerkleRoot упрощённая версия, не возвращающая ошибку.
// При ошибке возвращает нулевой хеш (32 нуля). Удобна для случаев,
// когда ошибка не критична (например, для логирования).
// Однако в production-коде предпочтительнее CalculateMerkleRootWithError.
func (b *Block) CalculateMerkleRoot() []byte {
	root, err := b.CalculateMerkleRootWithError()
	if err != nil {
		return make([]byte, 32)
	}
	return root
}

// SetMerkleRoot вычисляет Merkle root для блока и сохраняет его в поле b.Header.MerkleRoot.
// Вызывается перед майнингом или при создании блока.
// Возвращает ошибку, если вычисление не удалось.
func (b *Block) SetMerkleRoot() error {
	root, err := b.CalculateMerkleRootWithError()
	if err != nil {
		return fmt.Errorf("failed to calculate merkle root: %w", err)
	}
	b.Header.MerkleRoot = root
	return nil
}
