package chain

import (
	"fmt"

	"github.com/Alex1997377/weave/internal/core/block"
)

// loadBlocks загружает все блоки из хранилища
func (bc *Blockchain) loadBlocks() error {
	var blocks []*block.Block
	currentHash := bc.Tip

	for currentHash != nil {
		b, err := bc.store.GetBlock(currentHash)
		if err != nil {
			return fmt.Errorf("failed to get block %x: %w", currentHash, err)
		}

		// Вставляем в начало среза (обратный порядок)
		blocks = append([]*block.Block{b}, blocks...)

		// Проверяем, дошли ли до генезис блока
		isGenesis := true
		for _, bVal := range b.Header.PreviousHash {
			if bVal != 0 {
				isGenesis = false
				break
			}
		}

		if isGenesis {
			break
		}
		currentHash = b.Header.PreviousHash
	}

	bc.Blocks = blocks
	return nil
}
