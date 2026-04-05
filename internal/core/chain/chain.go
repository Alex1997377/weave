package chain

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/transaction"
	"github.com/Alex1997377/weave/internal/store"
)

const DIFFICULTY int = 4

type Blockchain struct {
	store  store.BlockStore // Приватное поле
	Tip    []byte
	Blocks []*block.Block
}

// NewBlockchain создает новую или восстанавливает существующую цепочку
func NewBlockchain(store store.BlockStore) (*Blockchain, error) {
	if store == nil {
		return nil, errors.New("store cannot be nil")
	}

	lastHash, err := store.GetLastHash()
	if err != nil {
		return nil, fmt.Errorf("failed to get last hash: %w", err)
	}

	// Если нет последнего хеша, создаем генезис блок
	if lastHash == nil {
		genesis, err := block.NewBlock([]transaction.Transaction{}, make([]byte, 32), 0, DIFFICULTY)
		if err != nil {
			return nil, fmt.Errorf("failed to create genesis block: %w", err)
		}

		if err := store.SaveBlock(genesis); err != nil {
			return nil, fmt.Errorf("failed to save genesis block: %w", err)
		}

		return &Blockchain{
			store:  store,
			Tip:    genesis.Hash,
			Blocks: []*block.Block{genesis},
		}, nil
	}

	// Загружаем существующую цепочку
	bc := &Blockchain{
		store: store,
		Tip:   lastHash,
	}

	if err := bc.loadBlocks(); err != nil {
		return nil, fmt.Errorf("failed to load blocks: %w", err)
	}

	return bc, nil
}

// Display отображает все блоки
func (bc *Blockchain) Display() {
	for i, b := range bc.Blocks {
		fmt.Printf("--- Block ID: %d ---\n", i)
		fmt.Printf("Timestamp: 	%d\n", b.Header.Timestamp)
		fmt.Printf("Transactions: 	%d\n", len(b.Transaction))
		fmt.Printf("Prev Hash:  %s\n", hex.EncodeToString(b.Header.PreviousHash))
		fmt.Printf("Size: %d bytes\n", b.Size)
		fmt.Printf("Hash: 		%s\n", hex.EncodeToString(b.Hash))
		fmt.Println("  --- ฿ ---  ")
	}
}

// Close закрывает хранилище
func (bc *Blockchain) Close() error {
	return bc.store.Close()
}
