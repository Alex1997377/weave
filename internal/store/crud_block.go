package store

import (
	"github.com/Alex1997377/weave/internal/core"
	"github.com/dgraph-io/badger/v4"
)

// SaveBlock записывает блок в БД и обновляет указатель на последний хеш
func (r *Repository) SaveBlock(block *core.Block) error {
	return r.db.Update(func(txn *badger.Txn) error {
		// 1. Сериализуем блок (используем твой метод)
		blockData := block.Serialize()

		// 2. Сохраняем блок по ключу его хеша
		err := txn.Set(block.Hash, blockData)
		if err != nil {
			return err
		}

		// 3. Обновляем ключ "l" (последний хеш в сети)
		return txn.Set([]byte("l"), block.Hash)
	})
}

// GetLastHash возвращает хеш последнего блока из БД
func (r *Repository) GetLastHash() ([]byte, error) {
	var lastHash []byte
	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("l"))
		if err != nil {
			return err
		}
		lastHash, err = item.ValueCopy(nil)
		return err
	})
	return lastHash, err
}
