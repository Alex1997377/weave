package crypto

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/Alex1997377/weave/internal/core/block/interfaces"
)

type Hash []byte

func (h Hash) String() string {
	if h == nil {
		return "<nil>"
	}
	return hex.EncodeToString(h)
}

func (h Hash) Bytes() []byte {
	if h == nil {
		return nil
	}
	result := make([]byte, len(h))
	copy(result, h)
	return result
}

func (h Hash) IsValidForDifficulty(difficulty int) bool {
	// difficulty is not been unpositive and empty hash is not been validate
	if difficulty < 0 || len(h) == 0 {
		return false
	}

	// protected from exit without slice range
	// Hash: [0, 0, 45, 12, ...]
	//     ↑  ↑
	//   byte0 byte1 (first 2 bytes must been 0)
	fullBytes := difficulty / 8
	if fullBytes > len(h) {
		fullBytes = len(h)
	}

	// check the first N - bytes is 0
	// ✅ [0, 0, 45, 12, ...] - valid (first 2 bytes = 0)
	// ❌ [0, 5, 45, 12, ...] - invalid (2-й byte ≠ 0)
	// ❌ [1, 0, 45, 12, ...] - invalid (1-й byte ≠ 0)
	for i := 0; i < fullBytes; i++ {
		if h[i] != 0 {
			return false
		}
	}

	// check remaining bits
	remainingBits := difficulty % 8
	if remainingBits > 0 && fullBytes < len(h) {
		mask := byte(0xFF << (8 - remainingBits))
		if h[fullBytes]&mask != 0 {
			return false
		}
	}

	return true
}

func HashPublicKey(pubKey []byte) []byte {
	hash := sha256.Sum256(pubKey)
	return hash[:20]
}

func HashBytes(data []byte) Hash {
	sum := sha256.Sum256(data)
	return Hash(sum[:])
}

type HashCalculatorImpl struct{}

func (HashCalculatorImpl) Hash(data []byte) interfaces.Hash {
	hash := sha256.Sum256(data)
	return Hash(hash[:])
}
