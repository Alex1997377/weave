package chain

import (
	"errors"
	"fmt"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/transaction"
)

// AddBlock добавляет новый блок в цепочку
func (bc *Blockchain) AddBlock(transactions []transaction.Transaction) error {
	if len(bc.Blocks) == 0 {
		return errors.New("cannot add block to empty blockchain")
	}

	// Валидация транзакций
	for i, tx := range transactions {
		if tx == nil {
			return fmt.Errorf("transaction at index %d is nil", i)
		}
		if err := tx.TransactionValidate(); err != nil {
			return fmt.Errorf("transaction validation failed at index %d: %w", i, err)
		}
	}

	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock, err := block.NewBlock(transactions, prevBlock.Hash, prevBlock.Header.Index+1, DIFFICULTY)
	if err != nil {
		return fmt.Errorf("failed to create new block: %w", err)
	}

	// Проверка размера блока
	size, err := newBlock.CalculateSize()
	if err != nil {
		return fmt.Errorf("failed to calculate block size: %w", err)
	}

	if size > 1024*1024 {
		return fmt.Errorf("block size %d exceeds limit of 1MB", size)
	}

	if err := newBlock.Validate(); err != nil {
		return fmt.Errorf("new block validation failed: %w", err)
	}

	// Сохраняем в хранилище
	if err := bc.store.SaveBlock(newBlock); err != nil {
		return fmt.Errorf("failed to save block to store: %w", err)
	}

	bc.Blocks = append(bc.Blocks, newBlock)
	bc.Tip = newBlock.Hash
	return nil
}
