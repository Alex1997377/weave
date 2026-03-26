package tests

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/block/tests/helpers"
	"github.com/Alex1997377/weave/internal/core/block/tests/mocks"
	"github.com/Alex1997377/weave/internal/core/transaction"
)

func CreateBenchmarkBlockData(txCount uint32, transactions []transaction.Transaction) *bytes.Buffer {
	buf := &bytes.Buffer{}

	buf.Write(bytes.Repeat([]byte{0xAA}, 32))
	binary.Write(buf, binary.LittleEndian, txCount)

	for _, tx := range transactions {
		txData, err := tx.TransactionSerialize()
		if err != nil {
			panic(err)
		}
		buf.Write(txData)
	}
	buf.Write(bytes.Repeat([]byte{0xCC}, 32))
	binary.Write(buf, binary.LittleEndian, uint32(12345))

	return buf
}

func BenchmarkDeserializeBlock(b *testing.B) {
	txCounts := []uint32{0, 1, 10, 100, 1000, 10000}

	for _, txCount := range txCounts {
		b.Run(fmt.Sprintf("tx=%d", txCount), func(b *testing.B) {
			transactions := make([]transaction.Transaction, txCount)

			for i := uint32(0); i < txCount; i++ {
				transactions[i] = helpers.CreateTestTransaction(byte(i % 256))
			}
			data := CreateBenchmarkBlockData(txCount, transactions)
			dataBytes := data.Bytes()

			opts := block.DeserializeOptions{
				Header: &mocks.MockHeaderDeserializer{},
				Tx:     &mocks.MockTransactionDeserializer{},
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err := block.DeserializeBlockWithparallelPooled(dataBytes, opts)

				if err != nil {
					b.Fatalf("deserialize failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkDeserializedBlockParallel(b *testing.B) {
	txCount := uint32(1000)
	transactions := make([]transaction.Transaction, txCount)

	for i := uint32(0); i < txCount; i++ {
		transactions[i] = helpers.CreateTestTransaction(byte(i * 256))
	}
	data := CreateBenchmarkBlockData(txCount, transactions)
	databytes := data.Bytes()

	opts := block.DeserializeOptions{
		Header: &mocks.MockHeaderDeserializer{},
		Tx:     &mocks.MockTransactionDeserializer{},
	}

	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := block.DeserializeBlockWithparallelPooled(databytes, opts)
			if err != nil {
				b.Fatalf("deserialize failed: %v", err)
			}
		}
	})
}
