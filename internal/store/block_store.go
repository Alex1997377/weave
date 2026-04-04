package block

type BlockStore interface {
	SaveBlock(block *Block) error
	GetBlock(hash []byte) (*Block, error)
	GetLastHash() ([]byte, error)
	Close() error
}
