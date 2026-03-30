// internal/core/block/block_hash.go
package block

import (
	"errors"

	"github.com/Alex1997377/weave/internal/crypto/hash"
)

func (b *Block) HashString() (string, error) {
	if b == nil {
		return "", errors.New("block is nil")
	}
	return hash.HashToString(b.Hash) // ← теперь crypto содержит HashToString
}

func (b *Block) ShortHash() (string, error) {
	if b == nil {
		return "", errors.New("block is nil")
	}
	return hash.ShortHashString(b.Hash)
}

func (b *Block) FormatHash(prefix string) (string, error) {
	if b == nil {
		return "", errors.New("block is nil")
	}
	return hash.FormatHashWithPrefix(prefix, b.Hash)
}
