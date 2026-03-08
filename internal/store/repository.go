package store

import (
	"errors"

	"github.com/dgraph-io/badger/v4"
)

var (
	ErrBlockNotFound = errors.New("block not found")
	ErrNilBlock      = errors.New("block is nil")
	ErrNilHash       = errors.New("hash is nil")
)

type Repository struct {
	db *badger.DB
}

// NewRepository открывает базу данных в указанной папке
func NewRepository(dbPath string) (*Repository, error) {
	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil // Отключаем лишние логи БД

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}
