package benchmark

import (
	"testing"

	"github.com/Alex1997377/weave/internal/core/block/tests/helpers"
)

// Бенчмарк для валидации корректного блока
func BenchmarkValidate(b *testing.B) {
	blk := helpers.CreateValidBlockForValidate()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := blk.Validate(); err != nil {
			b.Fatalf("validation failed: %v", err)
		}
	}
}
