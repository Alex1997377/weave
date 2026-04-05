// internal/core/block/block_hash.go
// Методы для получения строковых представлений хеша блока.
package block

import (
	"errors"

	"github.com/Alex1997377/weave/internal/crypto/hash"
)

// HashString возвращает полное шестнадцатеричное представление хеша блока.
// Если блок nil, возвращает ошибку.
// Результат: строка длиной 64 символа (32 байта → 64 hex-символа).
// Пример вывода: "a1b2c3d4e5f67890..."
func (b *Block) HashString() (string, error) {
	if b == nil {
		return "", errors.New("block is nil")
	}
	// hash.HashToString проверяет хеш на nil/пустоту и кодирует в hex
	return hash.HashToString(b.Hash)
}

// ShortHash возвращает сокращённое представление хеша (первые 8 символов).
// Удобно для быстрой идентификации блока в логах или UI.
// Если блок nil, возвращает ошибку.
// Пример вывода: "a1b2c3d4"
func (b *Block) ShortHash() (string, error) {
	if b == nil {
		return "", errors.New("block is nil")
	}
	return hash.ShortHashString(b.Hash)
}

// FormatHash возвращает строку с префиксом и полным хешем блока.
// Полезно для форматированного вывода, например: "Block hash: a1b2c3d4..."
// Если блок nil, возвращает ошибку.
// Пример вызова: block.FormatHash("Hash")
// Пример вывода: "Hash: a1b2c3d4e5f67890..."
func (b *Block) FormatHash(prefix string) (string, error) {
	if b == nil {
		return "", errors.New("block is nil")
	}
	return hash.FormatHashWithPrefix(prefix, b.Hash)
}
