package chain

import (
	"github.com/Alex1997377/weave/internal/core/block"
)

// GetBlockByHash возвращает блок по хешу
func (bc *Blockchain) GetBlockByHash(hash []byte) (*block.Block, error) {
	if hash == nil {
		return nil, NewInvalidHashError("hash cannot be nil", nil)
	}

	blk, err := bc.store.GetBlock(hash)
	if err != nil {
		return nil, NewChainCorruptedError("failed to get block from store", err)
	}
	if blk == nil {
		return nil, NewBlockNotFoundError("block not found", nil)
	}

	return blk, nil

}

// GetBlockByIndex возвращает блок по индексу
func (bc *Blockchain) GetBlockByIndex(index int) (*block.Block, error) {
	if index < 0 || index >= len(bc.Blocks) {
		return nil, NewBlockNotFoundError("block index out of range", nil)
	}
	return bc.Blocks[index], nil
}

func (bc *Blockchain) GetLastBlock() (*block.Block, error) {
	if len(bc.Blocks) == 0 {
		return nil, NewBlockNotFoundError("no blocks in chain", nil)
	}
	return bc.Blocks[len(bc.Blocks)-1], nil
}
