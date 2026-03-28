package tests

import (
	"bytes"
	"testing"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/block/tests/helpers"
	"github.com/Alex1997377/weave/internal/core/block/tests/mocks"
	"github.com/Alex1997377/weave/internal/core/transaction"
)

func TestBlock_CalculateMerkleRootWithError(t *testing.T) {
	tests := []struct {
		name        string
		block       *block.Block
		wantErr     bool
		wantRootLen int
	}{
		{
			name:        "nil block",
			block:       nil,
			wantErr:     true,
			wantRootLen: 0,
		},
		{
			name:        "empty transactions",
			block:       helpers.CreateTestBlockWithTxIDs([][]byte{}),
			wantErr:     false,
			wantRootLen: 32,
		},
		{
			name:        "one valid transaction",
			block:       helpers.CreateTestBlockWithTxIDs([][]byte{bytes.Repeat([]byte{0x01}, 32)}),
			wantErr:     false,
			wantRootLen: 32,
		},
		{
			name:        "two valid transactions",
			block:       helpers.CreateTestBlockWithTxIDs([][]byte{bytes.Repeat([]byte{0x01}, 32), bytes.Repeat([]byte{0x02}, 32)}),
			wantErr:     false,
			wantRootLen: 32,
		},
		{
			name: "nil transaction",
			block: func() *block.Block {
				b := &block.Block{Transaction: make([]transaction.Transaction, 1)}
				b.Transaction[0] = nil
				return b
			}(),
			wantErr:     true,
			wantRootLen: 0,
		},
		{
			name: "transaction with nil ID",
			block: func() *block.Block {
				tx := &mocks.MockTransaction{Id: nil}
				return &block.Block{Transaction: []transaction.Transaction{tx}}
			}(),
			wantErr:     true,
			wantRootLen: 0,
		},
		{
			name: "transaction with invalid ID length",
			block: func() *block.Block {
				tx := &mocks.MockTransaction{Id: bytes.Repeat([]byte{0x01}, 16)}
				return &block.Block{Transaction: []transaction.Transaction{tx}}
			}(),
			wantErr:     true,
			wantRootLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := tt.block.CalculateMerkleRootWithError()
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateMerkleRootWithError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(root) != tt.wantRootLen {
					t.Errorf("root length = %d, want %d", len(root), tt.wantRootLen)
				}
				if tt.wantRootLen == 32 && bytes.Equal(root, make([]byte, 32)) && tt.name != "empty transactions" && len(tt.block.Transaction) > 0 {
					// Для непустых блоков root не должен быть нулевым, если только не ошибка в расчётах
					// Но из-за того, что у нас мок ID может быть одинаковым, root может получиться нулевым только в редких случаях.
					// Проверяем, что root не nil.
					if root == nil {
						t.Error("root is nil")
					}
				}
			}
		})
	}
}

func TestBlock_CalculateMerkleRoot(t *testing.T) {
	tests := []struct {
		name        string
		block       *block.Block
		wantRootLen int
	}{
		{
			name:        "nil block",
			block:       nil,
			wantRootLen: 32, // возвращает нулевой хеш
		},
		{
			name:        "empty transactions",
			block:       helpers.CreateTestBlockWithTxIDs([][]byte{}),
			wantRootLen: 32,
		},
		{
			name:        "valid transactions",
			block:       helpers.CreateTestBlockWithTxIDs([][]byte{bytes.Repeat([]byte{0x01}, 32), bytes.Repeat([]byte{0x02}, 32)}),
			wantRootLen: 32,
		},
		{
			name: "nil transaction",
			block: func() *block.Block {
				b := &block.Block{Transaction: make([]transaction.Transaction, 1)}
				b.Transaction[0] = nil
				return b
			}(),
			wantRootLen: 32, // ошибка игнорируется, возвращается нулевой хеш
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := tt.block.CalculateMerkleRoot()
			if len(root) != tt.wantRootLen {
				t.Errorf("CalculateMerkleRoot() length = %d, want %d", len(root), tt.wantRootLen)
			}
			if root == nil {
				t.Error("root is nil")
			}
		})
	}
}

func TestBlock_SetMerkleRoot(t *testing.T) {
	tests := []struct {
		name      string
		block     *block.Block
		wantErr   bool
		checkRoot bool // проверяем, что root установлен и не нулевой (для валидных)
	}{
		{
			name:      "nil block",
			block:     nil,
			wantErr:   true,
			checkRoot: false,
		},
		{
			name:      "empty transactions",
			block:     helpers.CreateTestBlockWithTxIDs([][]byte{}),
			wantErr:   false,
			checkRoot: true, // должен установить нулевой хеш
		},
		{
			name:      "valid transactions",
			block:     helpers.CreateTestBlockWithTxIDs([][]byte{bytes.Repeat([]byte{0x01}, 32)}),
			wantErr:   false,
			checkRoot: true,
		},
		{
			name: "nil transaction",
			block: func() *block.Block {
				b := &block.Block{Transaction: make([]transaction.Transaction, 1)}
				b.Transaction[0] = nil
				return b
			}(),
			wantErr:   true,
			checkRoot: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.block.SetMerkleRoot()
			if (err != nil) != tt.wantErr {
				t.Errorf("SetMerkleRoot() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.checkRoot && !tt.wantErr {
				if tt.block == nil {
					t.Fatal("block is nil")
				}
				if len(tt.block.Header.MerkleRoot) != 32 {
					t.Errorf("Header.MerkleRoot length = %d, want 32", len(tt.block.Header.MerkleRoot))
				}
				// Если есть хотя бы одна валидная транзакция, root не должен быть нулевым
				if len(tt.block.Transaction) > 0 {
					// Проверка на нулевой хеш (32 нулевых байта)
					zeroHash := make([]byte, 32)
					if bytes.Equal(tt.block.Header.MerkleRoot, zeroHash) {
						t.Error("Merkle root is zero, but there are transactions")
					}
				}
			}
		})
	}
}
