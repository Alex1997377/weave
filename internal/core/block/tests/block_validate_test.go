package tests

import (
	"errors"
	"testing"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/block/tests/helpers"
	"github.com/Alex1997377/weave/internal/core/header"
	"github.com/Alex1997377/weave/internal/core/transaction"
	"github.com/Alex1997377/weave/internal/crypto/hash"
)

// Тесты
func TestBlock_Validate(t *testing.T) {
	tests := []struct {
		name      string
		block     *block.Block
		wantErr   bool
		errString string
	}{
		{
			name:      "nil block",
			block:     nil,
			wantErr:   true,
			errString: "block is nil",
		},
		{
			name: "nil hash",
			block: &block.Block{
				Header:      header.Header{Difficulty: 1},
				Hash:        nil,
				Transaction: nil,
			},
			wantErr:   true,
			errString: "block hash is nil",
		},
		{
			name: "negative difficulty",
			block: &block.Block{
				Header:      header.Header{Difficulty: -1},
				Hash:        hash.HashBytes([]byte{0}),
				Transaction: nil,
			},
			wantErr:   true,
			errString: "block difficulty cannot be negative",
		},
		{
			name: "invalid proof of work (hash too high)",
			block: &block.Block{
				Header:      header.Header{Difficulty: 5}, // требует 5 нулевых бит, а хеш не подходит
				Hash:        hash.HashBytes([]byte{0xFF}),
				Transaction: nil,
			},
			wantErr:   true,
			errString: "invalid proof of work",
		},
		{
			name: "hash mismatch",
			block: func() *block.Block {
				b := helpers.CreateValidBlockForValidate()
				b.Hash = hash.HashBytes([]byte{0x01}) // подменяем на другой
				return b
			}(),
			wantErr:   true,
			errString: "block hash doesn`t match content",
		},
		{
			name: "nil transaction in slice",
			block: func() *block.Block {
				b := helpers.CreateValidBlockForValidate()
				b.Transaction = append(b.Transaction, nil)
				return b
			}(),
			wantErr:   true,
			errString: "transaction at index 2 is nil",
		},
		{
			name: "invalid transaction",
			block: func() *block.Block {
				b := helpers.CreateValidBlockForValidate()
				// подменяем транзакцию на ту, что возвращает ошибку
				invalidTx := &helpers.TestTransactionWithValidate{
					Id:          []byte("bad"),
					Sender:      []byte("alice"),
					Recipient:   []byte("bob"),
					Amount:      0,
					Signature:   []byte("sig"),
					ValidateErr: errors.New("invalid tx"),
				}
				b.Transaction = []transaction.Transaction{invalidTx}
				return b
			}(),
			wantErr:   true,
			errString: "invalid transaction at index 0: invalid tx",
		},
		{
			name:      "valid block",
			block:     helpers.CreateValidBlockForValidate(),
			wantErr:   false,
			errString: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.block.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errString {
				t.Errorf("Validate() error = %v, want %v", err, tt.errString)
			}
		})
	}
}
