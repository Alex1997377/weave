package store

import "github.com/Alex1997377/weave/internal/core/block"

type BlockStore interface {
	SaveBlock(block *block.Block) error
	GetBlock(hash []byte) (*block.Block, error)
	GetLastHash() ([]byte, error)
	Close() error
}
