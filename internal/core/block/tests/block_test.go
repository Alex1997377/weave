package tests

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/block/tests/helpers"
	"github.com/Alex1997377/weave/internal/core/transaction"
)

// TestNewBlock тестирует функцию создания нового блока block.NewBlock.
// Функция NewBlock создаёт блок с заданными транзакциями, хешем предыдущего блока,
// индексом (высотой) и сложностью майнинга.
//
// Сценарии тестирования (таблица tests):
//  1. "nil previous hash" – prevHash = nil.
//     Ожидается: ошибка с текстом "previous hash cannot be nil".
//  2. "negative index" – index = -1.
//     Ожидается: ошибка с текстом "block index cannot be negative".
//  3. "negative difficulty" – difficulty = -1.
//     Ожидается: ошибка с текстом "difficulty cannot be negative".
//  4. "non-genesis with no transactions" – index = 1, transactions = [].
//     Ожидается: ошибка с текстом "non-genesis block must have at least one transaction".
//  5. "nil transaction in slice" – в слайсе транзакций есть nil-элемент.
//     Ожидается: ошибка с текстом "transaction at index 1 is nil".
//  6. "valid genesis block" – index = 0, transactions = [].
//     Ожидается: успех (ошибка nil), блок не nil.
//  7. "valid normal block" – index = 1, две корректные транзакции.
//     Ожидается: успех, блок не nil.
//
// Входные параметры функции NewBlock:
//   - transactions: []transaction.Transaction – слайс транзакций блока.
//   - prevHash: []byte – хеш предыдущего блока (должен быть 32 байта, но тест не проверяет длину).
//   - index: int – высота блока (номер в цепочке).
//   - difficulty: int – сложность майнинга (количество ведущих нулей в хеше).
//
// Ожидаемые выходные значения NewBlock:
//   - got: *block.Block – указатель на созданный блок (при успехе).
//   - err: error – nil при успехе, иначе ошибка валидации.
//
// Проверки теста:
//   - Наличие ошибки соответствует wantErr.
//   - При ошибке сообщение содержит ожидаемую подстроку errContains.
//   - При успехе блок не равен nil.
func TestNewBlock(t *testing.T) {
	prevHash := bytes.Repeat([]byte{0xAA}, 32) // правильная длина 32

	tests := []struct {
		name         string
		transactions []transaction.Transaction
		prevHash     []byte
		index        int
		difficulty   int
		wantErr      bool
		errContains  string
	}{
		{
			name:         "nil previous hash",
			transactions: nil,
			prevHash:     nil,
			index:        0,
			difficulty:   0,
			wantErr:      true,
			errContains:  "previous hash cannot be nil",
		},
		{
			name:         "negative index",
			transactions: nil,
			prevHash:     prevHash,
			index:        -1,
			difficulty:   0,
			wantErr:      true,
			errContains:  "block index cannot be negative",
		},
		{
			name:         "negative difficulty",
			transactions: nil,
			prevHash:     prevHash,
			index:        0,
			difficulty:   -1,
			wantErr:      true,
			errContains:  "difficulty cannot be negative",
		},
		{
			name:         "non-genesis with no transactions",
			transactions: []transaction.Transaction{},
			prevHash:     prevHash,
			index:        1,
			difficulty:   0,
			wantErr:      true,
			errContains:  "non-genesis block must have at least one transaction",
		},
		{
			name: "nil transaction in slice",
			transactions: []transaction.Transaction{
				helpers.CreateTestTransaction(1),
				nil,
			},
			prevHash:    prevHash,
			index:       1,
			difficulty:  0,
			wantErr:     true,
			errContains: "transaction at index 1 is nil",
		},
		{
			name:         "valid genesis block",
			transactions: []transaction.Transaction{},
			prevHash:     prevHash,
			index:        0,
			difficulty:   0,
			wantErr:      false,
		},
		{
			name: "valid normal block",
			transactions: []transaction.Transaction{
				helpers.CreateTestTransaction(1),
				helpers.CreateTestTransaction(2),
			},
			prevHash:   prevHash,
			index:      1,
			difficulty: 0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := block.NewBlock(tt.transactions, tt.prevHash, tt.index, tt.difficulty)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("error = %v, want contain %q", err, tt.errContains)
			}
			if !tt.wantErr && got == nil {
				t.Error("NewBlock() returned nil block")
			}
		})
	}
}
