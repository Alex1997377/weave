package chain

import (
	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/transaction"
	"github.com/Alex1997377/weave/internal/store"
)

type Blockchain struct {
	store  store.BlockStore
	Tip    []byte
	Blocks []*block.Block
}

// NewBlockchain создает новую или восстанавливает существующую цепочку
func NewBlockchain(store store.BlockStore, genesis *block.Block) (*Blockchain, error) {
	if store == nil {
		return nil, NewInvalidBlockError("store cannot be nil", nil)
	}
	if genesis == nil {
		return nil, NewInvalidBlockError("genesis block is nil", nil)
	}

	lastHash, err := store.GetLastHash()
	if err != nil {
		return nil, NewChainCorruptedError("failed to get last hash", err)
	}

	// Если нет последнего хеша, создаем генезис блок
	if lastHash == nil {
		genesis, err := block.NewBlock([]transaction.Transaction{}, make([]byte, 32), 0, DIFFICULTY)
		if err != nil {
			return nil, NewInvalidBlockError("failed to create genesis block", err)
		}

		if err := store.SaveBlock(genesis); err != nil {
			return nil, NewInvalidBlockError("failed to save genesis block", err)
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
		return nil, err
	}

	return bc, nil
}

// Close закрывает хранилище
func (bc *Blockchain) Close() error {
	return bc.store.Close()
}
