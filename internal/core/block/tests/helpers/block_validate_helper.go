package helpers

import (
	"bytes"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/header"
	"github.com/Alex1997377/weave/internal/core/transaction"
	"github.com/Alex1997377/weave/internal/crypto/hash"
)

// testTransactionWithValidate – реализует transaction.Transaction с возможностью управлять ошибкой валидации.
type TestTransactionWithValidate struct {
	Id          []byte
	Sender      []byte
	Recipient   []byte
	Amount      float64
	Signature   []byte
	ValidateErr error
}

func (tt *TestTransactionWithValidate) TransactionGetID() []byte        { return tt.Id }
func (tt *TestTransactionWithValidate) TransactionGetSender() []byte    { return tt.Sender }
func (tt *TestTransactionWithValidate) TransactionGetRecipient() []byte { return tt.Recipient }
func (tt *TestTransactionWithValidate) TransactionGetAmount() float64   { return tt.Amount }
func (tt *TestTransactionWithValidate) TransactionValidate() error      { return tt.ValidateErr }
func (tt *TestTransactionWithValidate) TransactionSign([]byte) error    { return nil }
func (tt *TestTransactionWithValidate) TransactionVerify([]byte) bool   { return true }
func (tt *TestTransactionWithValidate) TransactionSerialize() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.Write(tt.Id)
	buf.Write(tt.Sender)
	buf.Write(tt.Recipient)
	return buf.Bytes(), nil
}

// createValidBlockForValidate создаёт блок, проходящий валидацию.
func CreateValidBlockForValidate() *block.Block {
	// Подготавливаем заголовок
	h := header.Header{
		Index:        1,
		Timestamp:    1234567890,
		PreviousHash: bytes.Repeat([]byte{0xAA}, 32),
		MerkleRoot:   bytes.Repeat([]byte{0xBB}, 32),
		Nonce:        0,
		Difficulty:   1,
	}
	// Вычисляем хеш блока (заголовка)
	data, _ := h.Serialize()
	hashBlock := hash.HashBytes(data)

	// Создаём транзакции
	tx1 := &TestTransactionWithValidate{
		Id:          []byte("tx1"),
		Sender:      []byte("alice"),
		Recipient:   []byte("bob"),
		Amount:      100,
		Signature:   []byte("sig1"),
		ValidateErr: nil,
	}
	tx2 := &TestTransactionWithValidate{
		Id:          []byte("tx2"),
		Sender:      []byte("bob"),
		Recipient:   []byte("charlie"),
		Amount:      50,
		Signature:   []byte("sig2"),
		ValidateErr: nil,
	}
	return &block.Block{
		Header:      h,
		Transaction: []transaction.Transaction{tx1, tx2},
		Hash:        hashBlock,
		Size:        0,
	}
}
