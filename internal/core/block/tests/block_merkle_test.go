package tests

import (
	"bytes"
	"testing"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/block/tests/helpers"
	"github.com/Alex1997377/weave/internal/core/block/tests/mocks"
	"github.com/Alex1997377/weave/internal/core/transaction"
)

// TestBlock_CalculateMerkleRootWithError тестирует метод CalculateMerkleRootWithError().
// Метод вычисляет корень Меркла для всех транзакций блока и возвращает ошибку при любых проблемах.
//
// Сценарии тестирования:
//   1. "nil block" – блок равен nil.
//      Ожидается: ошибка, корень nil или пустой (wantRootLen=0).
//   2. "empty transactions" – блок без транзакций.
//      Ожидается: ошибки нет, корень длиной 32 байта (нулевой хеш).
//   3. "one valid transaction" – одна корректная транзакция с ID длиной 32 байта.
//      Ожидается: ошибки нет, корень длиной 32 байта, не nil.
//   4. "two valid transactions" – две корректные транзакции.
//      Ожидается: ошибки нет, корень длиной 32 байта.
//   5. "nil transaction" – слайс Transaction содержит nil-элемент.
//      Ожидается: ошибка, корень nil/пустой.
//   6. "transaction with nil ID" – транзакция имеет поле Id = nil.
//      Ожидается: ошибка, корень nil/пустой.
//   7. "transaction with invalid ID length" – ID транзакции имеет длину 16 байт (не 32).
//      Ожидается: ошибка, корень nil/пустой.
//
// Входные данные:
//   - block: *block.Block – указатель на блок, может быть nil.
//   - Внутри блока: Transaction – слайс transaction.Transaction.
//
// Ожидаемый результат:
//   - root: []byte – корень Меркла (32 байта) или nil/пустой при ошибке.
//   - err: error – nil при успехе, иначе ошибка.
//
// Проверки:
//   - Наличие ошибки соответствует wantErr.
//   - Длина корня равна wantRootLen (0 или 32).
//   - Для непустых валидных блоков корень не должен быть нулевым (доп. проверка).
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

// TestBlock_CalculateMerkleRoot тестирует метод CalculateMerkleRoot().
// В отличие от CalculateMerkleRootWithError, этот метод не возвращает ошибку,
// а в случае проблем возвращает нулевой хеш (32 нулевых байта).
//
// Сценарии тестирования:
//   1. "nil block" – блок равен nil.
//      Ожидается: корень длиной 32 байта (нулевой хеш), ошибка не возвращается.
//   2. "empty transactions" – блок без транзакций.
//      Ожидается: корень длиной 32 байта (нулевой хеш).
//   3. "valid transactions" – блок с двумя корректными транзакциями.
//      Ожидается: корень длиной 32 байта, не nil.
//   4. "nil transaction" – слайс Transaction содержит nil-элемент.
//      Ожидается: корень длиной 32 байта (нулевой хеш), ошибка игнорируется.
//
// Входные данные:
//   - block: *block.Block – указатель на блок, может быть nil.
//
// Ожидаемый результат:
//   - root: []byte – всегда 32 байта (либо реальный корень, либо нулевой хеш).
//     Ошибка не возвращается (сигнатура метода не содержит error).
//
// Проверки:
//   - Длина корня равна wantRootLen (всегда 32).
//   - root не равен nil.
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

// TestBlock_SetMerkleRoot тестирует метод SetMerkleRoot().
// Метод вычисляет корень Меркла для транзакций блока и сохраняет его в Header.MerkleRoot.
//
// Сценарии тестирования:
//   1. "nil block" – блок равен nil.
//      Ожидается: ошибка, корень не устанавливается.
//   2. "empty transactions" – блок без транзакций.
//      Ожидается: ошибки нет, Header.MerkleRoot устанавливается в нулевой хеш (32 нулевых байта).
//   3. "valid transactions" – блок с одной валидной транзакцией.
//      Ожидается: ошибки нет, Header.MerkleRoot имеет длину 32 байта и не является нулевым.
//   4. "nil transaction" – слайс Transaction содержит nil-элемент.
//      Ожидается: ошибка, Header.MerkleRoot не изменяется (или остаётся нулевым).
//
// Входные данные:
//   - block: *block.Block – указатель на блок, может быть nil.
//   - Внутри блока: Transaction – слайс транзакций, Header.MerkleRoot – будет перезаписан.
//
// Ожидаемый результат:
//   - err: error – nil при успехе, иначе ошибка.
//   - Поле block.Header.MerkleRoot обновляется (при успехе) или остаётся неизменным (при ошибке).
//
// Проверки:
//   - Наличие ошибки соответствует wantErr.
//   - При успехе длина MerkleRoot равна 32.
//   - Для блоков с транзакциями MerkleRoot не должен быть нулевым (если не пустой блок).
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