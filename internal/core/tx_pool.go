package core

import (
	"errors"
	"sync"
)

type Mempool struct {
	mu           sync.RWMutex
	transactions []Transaction
	maxSize      int
}

func NewMempool(capacity int) *Mempool {
	return &Mempool{
		transactions: make([]Transaction, 0, capacity),
		maxSize:      capacity,
	}
}

func (mp *Mempool) Add(tx Transaction) error {
	if err := tx.TransactionValidate(); err != nil {
		return err
	}

	if !tx.TransactionVerify(tx.TransactionGetSender()) {
		return NewInvalidSignatureError("invalid tx signature", nil)
	}

	mp.mu.Lock()
	defer mp.mu.Unlock()

	if len(mp.transactions) >= mp.maxSize {
		return errors.New("mempool is full")
	}

	mp.transactions = append(mp.transactions, tx)
	return nil
}

func (mp *Mempool) GetForBlock(limit int) []Transaction {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	count := len(mp.transactions)
	if count > limit {
		count = limit
	}

	txsForBlock := mp.transactions[:count]

	mp.transactions = mp.transactions[count:]

	return txsForBlock
}
