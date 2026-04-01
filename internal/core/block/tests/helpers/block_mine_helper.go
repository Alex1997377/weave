package helpers

import (
	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/header"
)

func CreateTestBlock(index, difficulty int) *block.Block {
	return &block.Block{
		Header: header.Header{
			Index:        index,
			Difficulty:   difficulty,
			PreviousHash: make([]byte, 32),
			Timestamp:    1234567890,
			Nonce:        0,
			MerkleRoot:   make([]byte, 32),
		},
		Transaction: nil,
		Hash:        nil,
		Size:        0,
	}
}
