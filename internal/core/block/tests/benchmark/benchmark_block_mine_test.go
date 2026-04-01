package benchmark

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/header"
	"github.com/Alex1997377/weave/internal/crypto/hash"
)

// createTestBlock создаёт блок для тестов (без транзакций).
func createTestBlock(index, difficulty int) *block.Block {
	return &block.Block{
		Header: header.Header{
			Index:        index,
			Difficulty:   difficulty,
			PreviousHash: make([]byte, 32),
			Timestamp:    time.Now().Unix(),
			MerkleRoot:   make([]byte, 32),
			Nonce:        0,
		},
		Transaction: nil,
		Hash:        nil,
		Size:        0,
	}
}

// BenchmarkMine измеряет время майнинга с разным количеством воркеров.
func BenchmarkMine(b *testing.B) {
	difficulty := 10 // умеренная сложность, чтобы каждая итерация была измеримой

	for _, workers := range []int{1, runtime.NumCPU(), 2 * runtime.NumCPU()} {
		b.Run(fmt.Sprintf("workers=%d", workers), func(b *testing.B) {
			hasher := &hash.HashCalculatorImpl{}
			config := block.MineConfig{
				NumWorkers: workers,
				Verbose:    false,
				Hasher:     hasher,
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				blk := createTestBlock(1, difficulty)
				err := blk.Mine(context.Background(), config)
				if err != nil {
					b.Fatalf("mine failed: %v", err)
				}
			}
		})
	}
}
