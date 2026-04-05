package benchmark

import (
	"testing"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/block/tests/helpers"
	"github.com/Alex1997377/weave/internal/core/transaction"
)

func BenchmarkNewBlock(b *testing.B) {
	tx := helpers.CreateTestTransaction(1)
	prevHash := make([]byte, 32)
	for i := 0; i < b.N; i++ {
		_, err := block.NewBlock([]transaction.Transaction{tx}, prevHash, 1, 0)
		if err != nil {
			b.Fatal(err)
		}
	}
}
