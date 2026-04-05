package helpers

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/Alex1997377/weave/internal/core/transaction"
)

type TestTransaction struct {
	Id        []byte
	Sender    []byte
	Recipient []byte
	Amount    float64
	Signature []byte
}

func CreateTestTransaction(id byte) transaction.Transaction {
	signature := bytes.Repeat([]byte{0xC0 + id}, 64)
	return &TestTransaction{
		Id:        bytes.Repeat([]byte{id}, 32),
		Sender:    bytes.Repeat([]byte{0xA0 + id}, 32),
		Recipient: bytes.Repeat([]byte{0xB0 + id}, 32),
		Amount:    float64(id) * 100.5,
		Signature: signature,
	}
}

func (tt *TestTransaction) TransactionSerialize() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.Write(tt.Sender)
	buf.Write(tt.Recipient)
	buf.Write(tt.Id)
	binary.Write(buf, binary.LittleEndian, tt.Amount)
	sigLen := uint32(len(tt.Signature))
	binary.Write(buf, binary.LittleEndian, sigLen)
	buf.Write(tt.Signature)
	return buf.Bytes(), nil
}

func (tt *TestTransaction) TransactionGetID() []byte        { return tt.Id }
func (tt *TestTransaction) TransactionGetSender() []byte    { return tt.Sender }
func (tt *TestTransaction) TransactionGetRecipient() []byte { return tt.Recipient }
func (tt *TestTransaction) TransactionGetAmount() float64   { return tt.Amount }
func (tt *TestTransaction) TransactionValidate() error      { return nil }
func (tt *TestTransaction) TransactionSign([]byte) error    { return nil }
func (tt *TestTransaction) TransactionVerify([]byte) bool   { return true }

func CreateValidBlockData(t *testing.T, txCount uint32, transactions []transaction.Transaction) *bytes.Buffer {
	t.Helper()
	buf := &bytes.Buffer{}
	buf.Write(bytes.Repeat([]byte{0xAA}, 32))
	t.Logf("After header: %d bytes", buf.Len())
	binary.Write(buf, binary.LittleEndian, txCount)
	t.Logf("After txCount: %d bytes", buf.Len())
	for i, tx := range transactions {
		txData, err := tx.TransactionSerialize()
		if err != nil {
			t.Fatalf("failed to serialize transaction %d: %v", i, err)
		}
		buf.Write(txData)
		t.Logf("After tx %d: %d bytes", i, buf.Len())
	}
	buf.Write(bytes.Repeat([]byte{0xCC}, 32))
	t.Logf("After hash: %d bytes", buf.Len())
	if err := binary.Write(buf, binary.LittleEndian, uint32(12345)); err != nil {
		t.Fatalf("failed to write block size: %v", err)
	}
	t.Logf("After size: %d bytes", buf.Len())
	return buf
}
