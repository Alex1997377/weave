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

func NewRepository(db *badger.DB) *Repository {
	return &Repository{
		db: db,
	}
}
