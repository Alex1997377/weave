package hash

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/Alex1997377/weave/internal/core/block/interfaces"
)

// Hash представляет собой 32-байтовый хеш
type Hash []byte

// String возвращает hex-представление хеша
func (h Hash) String() string {
	if h == nil {
		return "<nil>"
	}
	return hex.EncodeToString(h)
}

// Bytes возвращает копию хеша в виде []byte
func (h Hash) Bytes() []byte {
	if h == nil {
		return nil
	}
	result := make([]byte, len(h))
	copy(result, h)
	return result
}

// IsValidForDifficulty проверяет, удовлетворяет ли хеш заданной сложности
func (h Hash) IsValidForDifficulty(difficulty int) bool {
	if difficulty < 0 || len(h) == 0 {
		return false
	}

	fullBytes := difficulty / 8
	if fullBytes > len(h) {
		fullBytes = len(h)
	}

	for i := 0; i < fullBytes; i++ {
		if h[i] != 0 {
			return false
		}
	}

	remainingBits := difficulty % 8
	if remainingBits > 0 && fullBytes < len(h) {
		mask := byte(0xFF << (8 - remainingBits))
		if h[fullBytes]&mask != 0 {
			return false
		}
	}

	return true
}

// HashPublicKey хеширует публичный ключ (возвращает первые 20 байт)
func HashPublicKey(pubKey []byte) []byte {
	hash := sha256.Sum256(pubKey)
	return hash[:20]
}

// HashBytes вычисляет SHA256 от данных и возвращает Hash
func HashBytes(data []byte) Hash {
	sum := sha256.Sum256(data)
	return Hash(sum[:])
}

// HashCalculatorImpl реализует интерфейс interfaces.HashCalculator
type HashCalculatorImpl struct{}

// Hash вычисляет хеш от данных и возвращает интерфейс interfaces.Hash
func (HashCalculatorImpl) Hash(data []byte) interfaces.Hash {
	return HashBytes(data)
}
