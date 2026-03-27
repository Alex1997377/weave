package block

import (
	"errors"

	"github.com/Alex1997377/weave/pkg/utils"
)

func (b *Block) hashBytes() ([]byte, error) {
	if b == nil {
		return nil, errors.New("block is nil")
	}
	return b.Hash, nil
}

func (b *Block) HashString() (string, error) {
	hash, err := b.hashBytes()
	if err != nil {
		return "", err
	}

	return utils.HashToString(hash)
}

func (b *Block) ShortHash() (string, error) {
	hash, err := b.hashBytes()
	if err != nil {
		return "", err
	}

	return utils.ShortHashString(hash)
}

func (b *Block) FormatHash(prefix string) (string, error) {
	hash, err := b.hashBytes()
	if err != nil {
		return "", err
	}

	return utils.FormatHashWithPrefix(prefix, hash)
}
