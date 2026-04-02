package helpers

import (
	"bytes"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/header"
	"github.com/Alex1997377/weave/internal/core/transaction"
	"github.com/Alex1997377/weave/internal/crypto/hash"
)

type TestTransactionFromSerialize struct {
	id        []byte
	sender    []byte
	recipient []byte
	amount    float64
	signature []byte
}

func (tt *TestTransactionFromSerialize) TransactionGetID() []byte        { return tt.id }
func (tt *TestTransactionFromSerialize) TransactionGetSender() []byte    { return tt.sender }
func (tt *TestTransactionFromSerialize) TransactionGetRecipient() []byte { return tt.recipient }
func (tt *TestTransactionFromSerialize) TransactionGetAmount() float64   { return tt.amount }
func (tt *TestTransactionFromSerialize) TransactionValidate() error      { return nil }
func (tt *TestTransactionFromSerialize) TransactionSign([]byte) error    { return nil }
func (tt *TestTransactionFromSerialize) TransactionVerify([]byte) bool   { return true }
func (tt *TestTransactionFromSerialize) TransactionSerialize() ([]byte, error) {

	buf := &bytes.Buffer{}
	buf.Write(tt.id)
	buf.Write(tt.sender)
	buf.Write(tt.recipient)

	return buf.Bytes(), nil
}

func CreateTestBlockForSerialize() *block.Block {
	tx1 := &TestTransactionFromSerialize{
		id:        []byte("tx1"),
		sender:    []byte("alice"),
		recipient: []byte("bob"),
		amount:    100,
		signature: []byte("sig1"),
	}
	tx2 := &TestTransactionFromSerialize{
		id:        []byte("tx2"),
		sender:    []byte("bob"),
		recipient: []byte("charlie"),
		amount:    50,
		signature: []byte("sig2"),
	}

	return &block.Block{
		Header: header.Header{
			Index:        1,
			Timestamp:    1234567890,
			PreviousHash: bytes.Repeat([]byte{0xAA}, 32),
			MerkleRoot:   bytes.Repeat([]byte{0xBB}, 32),
			Nonce:        42,
			Difficulty:   2,
		},
		Transaction: []transaction.Transaction{tx1, tx2},
		Hash:        hash.HashBytes(bytes.Repeat([]byte{0xCC}, 32)),
		Size:        0,
	}
}
