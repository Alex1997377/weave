package benchmark

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/block/tests/helpers"
	"github.com/Alex1997377/weave/internal/core/block/tests/mocks"
	"github.com/Alex1997377/weave/internal/core/header"
	"github.com/Alex1997377/weave/internal/core/transaction"
)

func createBenchmarkBlockData(txCount uint32, transactions []transaction.Transaction) *bytes.Buffer {
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
			data := createBenchmarkBlockData(txCount, transactions)
			dataBytes := data.Bytes()

			headerMock := &mocks.MockHeaderDeserializer{
				MockFunc: func(r *bytes.Reader) (*header.Header, error) {
					r.Read(make([]byte, 32))
					return &header.Header{Index: 1}, nil
				},
			}

			txMock := &mocks.MockTransactionDeserializer{
				MockFunc: func(r *bytes.Reader) (transaction.Transaction, error) {
					// Читаем ровно 172 байта – размер одной тестовой транзакции
					buf := make([]byte, 172)
					if _, err := r.Read(buf); err != nil {
						return nil, err
					}
					// Возвращаем тестовую транзакцию (данные не важны, главное – сдвиг указателя)
					return helpers.CreateTestTransaction(1), nil
				},
			}

			opts := block.DeserializeOptions{
				Header: headerMock,
				Tx:     txMock,
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
	txCount := uint32(1000) // можно выбрать любое количество транзакций
	transactions := make([]transaction.Transaction, txCount)
	for i := uint32(0); i < txCount; i++ {
		transactions[i] = helpers.CreateTestTransaction(byte(i % 256))
	}
	data := createBenchmarkBlockData(txCount, transactions)
	dataBytes := data.Bytes()

	headerMock := &mocks.MockHeaderDeserializer{
		MockFunc: func(r *bytes.Reader) (*header.Header, error) {
			r.Read(make([]byte, 32))
			return &header.Header{Index: 1}, nil
		},
	}

	txMock := &mocks.MockTransactionDeserializer{
		MockFunc: func(r *bytes.Reader) (transaction.Transaction, error) {
			// Читаем ровно 172 байта – размер одной тестовой транзакции
			buf := make([]byte, 172)
			if _, err := r.Read(buf); err != nil {
				return nil, err
			}
			// Возвращаем тестовую транзакцию (можно любую)
			return helpers.CreateTestTransaction(1), nil
		},
	}

	opts := block.DeserializeOptions{
		Header: headerMock,
		Tx:     txMock,
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := block.DeserializeBlockWithparallelPooled(dataBytes, opts)
			if err != nil {
				b.Fatalf("deserialize failed: %v", err)
			}
		}
	})
}
