package store

import (
	"errors"
	"fmt"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/dgraph-io/badger/v4"
)

// SaveBlock сохраняет блок в БД (реализация block.BlockStore)
func (r *Repository) SaveBlock(b *block.Block) error {
	if b == nil {
		return ErrNilBlock
	}

	if b.Hash == nil {
		return errors.New("block hash is nil")
	}

	return r.db.Update(func(txn *badger.Txn) error {
		// Сериализуем блок
		blockData, err := b.Serialize()
		if err != nil {
			return fmt.Errorf("failed to serialize block: %w", err)
		}

		// Сохраняем блок по ключу b + hash
		key := append([]byte("b"), b.Hash...)
		if err := txn.Set(key, blockData); err != nil {
			return fmt.Errorf("failed to set block data: %w", err)
		}

		// Обновляем указатель на последний блок
		if err := txn.Set([]byte("l"), b.Hash); err != nil {
			return fmt.Errorf("failed to update last hash: %w", err)
		}

		return nil
	})
}

// GetBlock получает блок по хешу (реализация block.BlockStore)
func (r *Repository) GetBlock(hash []byte) (*block.Block, error) {
	if hash == nil {
		return nil, ErrNilHash
	}

	var resultBlock *block.Block

	err := r.db.View(func(txn *badger.Txn) error {
		key := append([]byte("b"), hash...)
		item, err := txn.Get(key)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return ErrBlockNotFound
			}
			return fmt.Errorf("failed to get block from db: %w", err)
		}

		return item.Value(func(val []byte) error {
			var err error
			resultBlock, err = block.DeserializeBlock(val)
			if err != nil {
				return fmt.Errorf("failed to deserialize block: %w", err)
			}
			return nil
		})
	})

	return resultBlock, err
}

// GetLastHash получает хеш последнего блока (реализация block.BlockStore)
func (r *Repository) GetLastHash() ([]byte, error) {
	var lastHash []byte
	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("l"))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil // нет последнего хеша - это нормально для новой БД
			}
			return fmt.Errorf("failed to get last hash: %w", err)
		}

		return item.Value(func(val []byte) error {
			lastHash = make([]byte, len(val))
			copy(lastHash, val)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return lastHash, nil
}

// Close закрывает соединение с БД (реализация block.BlockStore)
func (r *Repository) Close() error {
	if r.db == nil {
		return errors.New("database connection is nil")
	}
	return r.db.Close()
}
