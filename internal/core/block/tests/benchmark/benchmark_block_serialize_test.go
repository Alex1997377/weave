package benchmark

import (
	"testing"

	"github.com/Alex1997377/weave/internal/core/block/tests/helpers"
)

func BenchmarkSerialize(b *testing.B) {
	blk := helpers.CreateTestBlockForSerialize()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := blk.Serialize()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCalculateHash(b *testing.B) {
	blk := helpers.CreateTestBlockForSerialize()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := blk.CalculateHash()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCalculateSize(b *testing.B) {
	blk := helpers.CreateTestBlockForSerialize()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := blk.CalculateSize()
		if err != nil {
			b.Fatal(err)
		}
	}
}
