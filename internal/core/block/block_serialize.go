// internal/core/block/block_serialize.go
// Сериализация блока, вычисление хеша и размера.
package block

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/Alex1997377/weave/internal/crypto/hash"
)

// Serialize преобразует блок в бинарный формат.
// Порядок записи:
//  1. Сериализованный заголовок (Header.Serialize())
//  2. Все транзакции подряд (TransactionSerialize() для каждой)
//  3. Хеш блока (32 байта)
//
// Возвращает:
//
//	[]byte – сериализованные данные.
//	error  – ошибка, если блок nil, транзакция nil, или ошибка сериализации.
//
// Выходные величины:
//
//	Размер блока (без учета самого размера, т.к. он не хранится) равен сумме длин заголовка,
//	всех транзакций и хеша. Для типового блока с 1000 транзакций (по 172 байта) + заголовок (~100)
//	+ хеш (32) ≈ 172 000 + 132 ≈ 172 КБ.
//
// Время выполнения: O(N), где N – количество транзакций. Для 10000 транзакций ~2–3 мс.
// Аллокации: создаётся один буфер, который растёт по мере записи.
func (b *Block) Serialize() ([]byte, error) {
	if b == nil {
		return nil, errors.New("block is nil")
	}

	buf := new(bytes.Buffer)

	// 1. Заголовок
	headerBytes, err := b.Header.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize header: %w", err)
	}
	if _, err := buf.Write(headerBytes); err != nil {
		return nil, fmt.Errorf("failed to write header to buffer: %w", err)
	}

	// 2. Транзакции
	for i, tx := range b.Transaction {
		if tx == nil {
			return nil, fmt.Errorf("transaction at index %d is nil", i)
		}
		txBytes, err := tx.TransactionSerialize()
		if err != nil {
			return nil, fmt.Errorf("failed to serialize transaction %d: %w", i, err)
		}
		if _, err := buf.Write(txBytes); err != nil {
			return nil, fmt.Errorf("failed to write transaction %d to buffer: %w", i, err)
		}
	}

	// 3. Хеш блока
	if _, err := buf.Write(b.Hash); err != nil {
		return nil, fmt.Errorf("failed to write block hash: %w", err)
	}

	return buf.Bytes(), nil
}

// CalculateHash вычисляет хеш блока как SHA-256 от сериализованного заголовка.
// Хеш зависит только от заголовка, так как Merkle root уже включён в заголовок,
// а транзакции не влияют на хеш напрямую.
// Возвращает:
//
//	[]byte – 32-байтовый хеш.
//	error  – ошибка, если блок nil или не удалось сериализовать заголовок.
//
// Время выполнения: ~1–2 мкс (зависит от размера заголовка).
// Аллокации: внутри hash.HashBytes (копия данных) и возвращаемый срез.
func (b *Block) CalculateHash() ([]byte, error) {
	if b == nil {
		return nil, errors.New("block is nil")
	}
	data, err := b.Header.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize header: %w", err)
	}
	// hash.HashBytes вычисляет SHA-256 и возвращает тип hash.Hash (который является []byte)
	return hash.HashBytes(data).Bytes(), nil
}

// CalculateSize вычисляет полный размер блока в байтах при сериализации.
// Включает:
//   - размер сериализованного заголовка
//   - сумму размеров всех сериализованных транзакций
//   - размер хеша (32 байта)
//
// Возвращаемое значение совпадает с длиной результата Serialize().
// Используется для заполнения поля Block.Size при создании блока.
// Время выполнения: O(N) – проход по всем транзакциям (каждая сериализуется).
// Для 10000 транзакций ~2–3 мс (аналогично Serialize).
// Аллокации: каждая транзакция сериализуется в отдельный байтовый срез, что даёт нагрузку на GC.
// Оптимизация: можно было бы вычислить размер без полной сериализации, зная структуру транзакции,
// но текущая реализация проще и надёжнее.
func (b *Block) CalculateSize() (uint32, error) {
	if b == nil {
		return 0, errors.New("block is nil")
	}

	// Заголовок
	headerBytes, err := b.Header.Serialize()
	if err != nil {
		return 0, fmt.Errorf("failed to serialize header for size calculation: %w", err)
	}
	headerSize := uint32(len(headerBytes))

	// Транзакции
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

	// Хеш
	hashSize := uint32(len(b.Hash))

	return headerSize + transactionsSize + hashSize, nil
}
