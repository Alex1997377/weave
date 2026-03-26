package chain

import (
	"errors"
	"fmt"

	"github.com/Alex1997377/weave/internal/core/block"
)

// GetBlockByHash возвращает блок по хешу
func (bc *Blockchain) GetBlockByHash(hash []byte) (*block.Block, error) {
	if hash == nil {
		return nil, errors.New("hash cannot be nil")
	}
	return bc.store.GetBlock(hash)
}

// GetBlockByIndex возвращает блок по индексу
func (bc *Blockchain) GetBlockByIndex(index int) (*block.Block, error) {
	if index < 0 || index >= len(bc.Blocks) {
		return nil, fmt.Errorf("block index %d out of range", index)
	}
	return bc.Blocks[index], nil
}
