// internal/core/block/block_validate.go
// Валидация блока: проверка корректности хеша, proof-of-work, целостности и транзакций.
package block

import (
	"bytes"
	"errors"
	"fmt"
)

// Validate проверяет блок на корректность.
// Выполняет следующие проверки:
//  1. Блок не nil.
//  2. Хеш блока не nil.
//  3. Сложность не отрицательная.
//  4. Хеш удовлетворяет заданной сложности (proof-of-work).
//  5. Хеш блока совпадает с вычисленным на основе заголовка.
//  6. Все транзакции не nil и валидны (TransactionValidate).
//
// Возвращает:
//
//	error – nil, если блок корректен, иначе описание ошибки.
//
// Время выполнения: O(N), где N – количество транзакций (каждая транзакция валидируется).
// Для блока с 1000 транзакций: ~1–2 мс (зависит от сложности валидации транзакций).
// Аллокации: незначительные (создание ошибок, возможно, копирование при сериализации заголовка внутри CalculateHash).
func (b *Block) Validate() error {
	if b == nil {
		return errors.New("block is nil")
	}
	if b.Hash == nil {
		return errors.New("block hash is nil")
	}
	// Дополнительная проверка длины хеша (должен быть 32 байта)
	if len(b.Hash) != 32 {
		return fmt.Errorf("invalid hash length: expected 32, got %d", len(b.Hash))
	}
	if b.Header.Difficulty < 0 {
		return errors.New("block difficulty cannot be negative")
	}
	// Проверка proof-of-work
	if !b.Hash.IsValidForDifficulty(b.Header.Difficulty) {
		return errors.New("invalid proof of work")
	}
	// Проверка соответствия хеша содержимому блока (хеш должен быть вычислен из заголовка)
	calculatedHash, err := b.CalculateHash()
	if err != nil {
		return fmt.Errorf("failed to calculate hash for validation: %w", err)
	}
	if !bytes.Equal(b.Hash, calculatedHash) {
		return errors.New("block hash doesn't match content")
	}
	// Валидация транзакций
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
