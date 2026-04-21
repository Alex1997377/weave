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

// TestBlock_Validate тестирует метод Validate() блока.
// Метод Validate проверяет целостность и корректность блока:
//   - блок не nil
//   - хеш блока не nil
//   - сложность не отрицательная
//   - proof of work (хеш блока соответствует сложности)
//   - хеш блока соответствует содержимому (хеш заголовка + транзакции)
//   - все транзакции не nil
//   - каждая транзакция валидна (вызов Transaction.Validate())
//
// Сценарии тестирования (таблица tests):
//   1. "nil block" – блок равен nil.
//      Ожидается: ошибка "block is nil".
//   2. "nil hash" – поле Hash блока равно nil.
//      Ожидается: ошибка "block hash is nil".
//   3. "negative difficulty" – сложность блока (Header.Difficulty) отрицательная.
//      Ожидается: ошибка "block difficulty cannot be negative".
//   4. "invalid proof of work (hash too high)" – хеш блока не удовлетворяет требуемой сложности
//      (например, difficulty=5, а хеш = 0xFF...).
//      Ожидается: ошибка "invalid proof of work".
//   5. "hash mismatch" – хеш блока не соответствует вычисленному хешу содержимого.
//      Ожидается: ошибка "block hash doesn`t match content".
//   6. "nil transaction in slice" – в слайсе Transaction есть nil-элемент.
//      Ожидается: ошибка "transaction at index 2 is nil" (индекс зависит от позиции).
//   7. "invalid transaction" – одна из транзакций возвращает ошибку при Validate().
//      Ожидается: ошибка "invalid transaction at index 0: invalid tx".
//   8. "valid block" – все проверки проходят успешно.
//      Ожидается: ошибка nil.
//
// Входные данные:
//   - block: *block.Block – указатель на блок для валидации.
//
// Выходные данные:
//   - err: error – nil если блок корректен, иначе ошибка с описанием проблемы.
//
// Проверки теста:
//   - Наличие ошибки соответствует wantErr.
//   - При ошибке текст сообщения совпадает с ожидаемым errString.
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