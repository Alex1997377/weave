package mocks

import (
	"github.com/Alex1997377/weave/internal/core/block/interfaces"
)

// MockHash реализует interfaces.Hash
type MockHash struct {
	Valid     bool
	BytesHash []byte
}

func (m MockHash) IsValidForDifficulty(difficulty int) bool {
	return m.Valid
}

// Bytes возвращает байтовое представление хеша
func (m MockHash) Bytes() []byte { // ← исправлено: было BytesHash
	return m.BytesHash
}

// MockHashCalculator реализует interfaces.HashCalculator
type MockHashCalculator struct {
	Valid bool // экспортируем для удобства
}

func (m *MockHashCalculator) Hash(data []byte) interfaces.Hash {
	return MockHash{
		Valid:     m.Valid,
		BytesHash: make([]byte, 32), // можно также копировать data, если нужно
	}
}
