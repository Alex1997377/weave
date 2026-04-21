package chain

import (
	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/transaction"
	"github.com/Alex1997377/weave/internal/store"
)

type Blockchain struct {
	store    store.BlockStore
	Tip      []byte
	Blocks   []*block.Block
	balances map[string]float64
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

		bc := &Blockchain{
			store:    store,
			Tip:      genesis.Hash,
			Blocks:   []*block.Block{genesis},
			balances: make(map[string]float64),
		}
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

func (bc *Blockchain) updateBalancesFromBlock(block *block.Block) {
	for _, tx := range block.Transaction {
		if tx == nil {
			continue
		}

		sender := string(tx.TransactionGetSender())
		recipient := string(tx.TransactionGetRecipient())
		amount := tx.TransactionGetAmount()

		bc.balances[sender] -= amount
		bc.balances[recipient] += amount
	}
}

// Close закрывает хранилище
func (bc *Blockchain) Close() error {
	return bc.store.Close()
}
