package store

import (
	"errors"

	"github.com/Alex1997377/weave/internal/core"
	"github.com/dgraph-io/badger/v4"
)

// SaveBlock записывает блок в БД и обновляет указатель на последний хеш

func (r *Repository) SaveBlock(block *core.Block) error {
	if block == nil {
		return ErrNilBlock
	}

	if block.Hash == nil {
		return errors.New("block hash is nil")
	}

	return r.db.Update(func(txn *badger.Txn) error {
		blockData, err := block.Serialize()
		if err != nil {
			return err
		}

		key := append([]byte("b"), block.Hash...)
		err = txn.Set(key, blockData)
		if err != nil {
			return err
		}

		return txn.Set([]byte("l"), block.Hash)
	})
}

func (r *Repository) GetBlock(hash []byte) (*core.Block, error) {
	if hash == nil {
		return nil, ErrNilHash
	}

	var block *core.Block
	err := r.db.View(func(txn *badger.Txn) error {
		key := append([]byte("b"), hash...)
		item, err := txn.Get(key)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return ErrBlockNotFound
			}
			return err
		}

		return item.Value(func(val []byte) error {
			// TODO
			block, err = core.DeserializationBlock(val)
			return err
		})
	})

	return block, err
}

// GetLastHash возвращает хеш последнего блока из БД
func (r *Repository) GetLastHash() ([]byte, error) {
	var lastHash []byte
	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("l"))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})
	})

	if err == badger.ErrKeyNotFound {
		return nil, nil
	}

	return lastHash, err
}

// Close закрывает базу при выходе из программы
func (r *Repository) Close() {
	r.db.Close()
}
