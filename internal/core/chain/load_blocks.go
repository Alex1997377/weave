package chain

import (
	"github.com/Alex1997377/weave/internal/core/block"
)

// loadBlocks загружает все блоки из хранилища
func (bc *Blockchain) loadBlocks() error {
	var blocks []*block.Block
	currentHash := bc.Tip

	for currentHash != nil {
		b, err := bc.store.GetBlock(currentHash)
		if err != nil {
			return NewChainCorruptedError("failed to get block", err)
		}
		if b == nil {
			return NewBlockNotFoundError("block not found in store", nil)
		}
		blocks = append(blocks, b)

		if isGenesis(b) {
			break
		}
		currentHash = b.Header.PreviousHash
	}

	for i, j := 0, len(blocks)-1; i < j; i, j = i+1, j-1 {
		blocks[i], blocks[j] = blocks[j], blocks[i]
	}

	bc.Blocks = blocks
	return nil
}

func isGenesis(b *block.Block) bool {
	for _, v := range b.Header.PreviousHash {
		if v != 0 {
			return false
		}
	}
	return true
}
