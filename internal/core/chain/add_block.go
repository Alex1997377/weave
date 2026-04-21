package chain

import (
	"errors"
	"fmt"
	"runtime"
	"sync"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/transaction"
)

// AddBlock добавляет новый блок в цепочку
func (bc *Blockchain) AddBlock(transactions []transaction.Transaction) error {
	if len(bc.Blocks) == 0 {
		return errors.New("cannot add block to empty blockchain")
	}

	if len(transactions) > 100 {
		if err := bc.validateTransactionsParallel(transactions); err != nil {
			return err
		}
	} else {
		// Валидация транзакций
		for i, tx := range transactions {
			if tx == nil {
				return fmt.Errorf("transaction at index %d is nil", i)
			}

			if err := tx.TransactionValidate(); err != nil {
				return fmt.Errorf("transaction validation failed at index %d: %w", i, err)
			}
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

	if size > maxBlockSize {
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

func (bc *Blockchain) validateTransactionsParallel(txs []transaction.Transaction) error {
	numWorkers := runtime.NumCPU()
	if len(txs) < numWorkers {
		numWorkers = len(txs)
	}

	tasks := make(chan int, len(txs))
	errCh := make(chan error, len(txs))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func ()  {
			defer wg.Done()
			for idx := range tasks {
				tx := txs[idx]
				if tx == nil {
					errCh <- fmt.Errorf("transaction at index %d is nil", idx)
					continue
				}
				if err := tx.TransactionValidate(); err != nil {
					errCh <- fmt.Errorf("transaction validation failed at index %d: %w", idx, err)
				}
			}
		}()
	}

	for i := range txs {
		tasks <- i
	}
	close(tasks)

	go func ()  {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}