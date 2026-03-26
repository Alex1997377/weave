package block

import (
	"errors"

	"github.com/Alex1997377/weave/pkg/utils"
)

func (b *Block) HashString() (string, error) {
	if b == nil {
		return "", errors.New("block is nil")
	}

	return utils.HashToString(b.Hash)
}

func (b *Block) ShortHash() (string, error) {
	if b == nil {
		return "", errors.New("block is nil")
	}

	return utils.ShortHashString(b.Hash)
}

func (b *Block) FormatHash(prefix string) (string, error) {
	if b == nil {
		return "", errors.New("block is nil")
	}

	return utils.FormatHahsWithPrefix(prefix, b.Hash)
}
