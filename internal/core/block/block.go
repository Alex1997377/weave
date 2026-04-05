// internal/core/block/block.go
// Основное определение блока и конструктор NewBlock.
package block

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Alex1997377/weave/internal/core/header"
	"github.com/Alex1997377/weave/internal/core/transaction"
	"github.com/Alex1997377/weave/internal/crypto/hash"
)

// Block представляет блок в блокчейне.
// Содержит заголовок, список транзакций, хеш блока (результат майнинга) и размер в байтах.
// Поля:
//
//	Header      – структура заголовка (индекс, временная метка, предыдущий хеш, Merkle root, nonce, сложность).
//	Transaction – слайс транзакций (каждая реализует интерфейс transaction.Transaction).
//	Hash        – 32-байтовый хеш блока, удовлетворяющий условию сложности.
//	Size        – общий размер блока при сериализации (вычисляется CalculateSize).
type Block struct {
	Header      header.Header
	Transaction []transaction.Transaction
	Hash        hash.Hash
	Size        uint32
}

// NewBlock создаёт новый блок с заданными транзакциями, предыдущим хешем, индексом и сложностью.
// Выполняет следующие шаги:
//  1. Проверяет корректность входных данных (PreviousHash не nil, неотрицательные index и difficulty,
//     не-генезис блок должен иметь хотя бы одну транзакцию, все транзакции не nil).
//  2. Инициализирует структуру Block с заголовком (текущее время, nonce=0, MerkleRoot пока nil).
//  3. Вычисляет Merkle root через SetMerkleRoot().
//  4. Запускает майнинг с настройками по умолчанию (MineConfig{}).
//  5. Вычисляет и устанавливает размер блока (CalculateSize).
//  6. Проверяет созданный блок через Validate().
//
// Возвращает:
//
//	*Block – указатель на созданный блок (при успехе).
//	error  – ошибка, если какой-либо шаг не удался.
//
// Примерные величины (для блока с 1000 транзакций, сложность 1):
//
//	Время создания: ~2–3 мс (включая майнинг).
//	Размер блока: ~172 КБ.
//	Аллокации: зависят от количества транзакций (каждая транзакция может выделять память).
func NewBlock(
	transactions []transaction.Transaction,
	PreviousHash []byte,
	index int,
	difficulty int,
) (*Block, error) {
	// Проверка входных параметров
	if PreviousHash == nil {
		return nil, errors.New("previous hash cannot be nil")
	}
	if index < 0 {
		return nil, fmt.Errorf("block index cannot be negative: %d", index)
	}
	if difficulty < 0 {
		return nil, fmt.Errorf("difficulty cannot be negative: %d", difficulty)
	}
	if index > 0 && len(transactions) == 0 {
		return nil, errors.New("non-genesis block must have at least one transaction")
	}
	for i, tx := range transactions {
		if tx == nil {
			return nil, fmt.Errorf("transaction at index %d is nil", i)
		}
	}
	// Создание блока
	block := &Block{
		Header: header.Header{
			Index:        index,
			Timestamp:    time.Now().Unix(),
			PreviousHash: PreviousHash,
			Difficulty:   difficulty,
			Nonce:        0,
			MerkleRoot:   nil,
		},
		Transaction: transactions,
	}
	// Вычисление Merkle root
	if err := block.SetMerkleRoot(); err != nil {
		return nil, fmt.Errorf("failed to set merkle root: %w", err)
	}
	// Майнинг (с настройками по умолчанию)
	if err := block.Mine(context.Background(), MineConfig{}); err != nil {
		return nil, fmt.Errorf("failed to mine block: %w", err)
	}
	// Вычисление размера
	size, err := block.CalculateSize()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate block size: %w", err)
	}
	block.Size = size
	// Финальная валидация
	if err := block.Validate(); err != nil {
		return nil, fmt.Errorf("created block is invalid: %w", err)
	}
	return block, nil
}
