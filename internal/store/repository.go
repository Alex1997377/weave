package store

import "github.com/dgraph-io/badger/v4"

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

// Close закрывает базу при выходе из программы
func (r *Repository) Close() {
	r.db.Close()
}
