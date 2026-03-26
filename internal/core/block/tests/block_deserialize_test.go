package tests

import (
	"bytes"
	"encoding/binary"
	"errors"
	"sync"
	"testing"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/block/tests/helpers"
	"github.com/Alex1997377/weave/internal/core/block/tests/mocks"
	"github.com/Alex1997377/weave/internal/core/header"
	"github.com/Alex1997377/weave/internal/core/transaction"
)

func TestDeserializeBlockWithDeps_Success(t *testing.T) {
	txCount := uint32(3)
	transactions := []transaction.Transaction{
		helpers.CreateTestTransaction(1),
		helpers.CreateTestTransaction(2),
		helpers.CreateTestTransaction(3),
	}
	data := helpers.CreateValidBlockData(t, txCount, transactions)
	dataBytes := data.Bytes()

	var mu sync.Mutex
	var callCount int

	txMock := &mocks.MockTransactionDeserializer{
		MockFunc: func(r *bytes.Reader) (transaction.Transaction, error) {
			mu.Lock()
			idx := callCount
			callCount++
			mu.Unlock()
			t.Logf("Mock called for tx %d", idx)
			return helpers.CreateTestTransaction(byte(idx + 1)), nil
		},
	}

	headerMock := &mocks.MockHeaderDeserializer{
		MockFunc: func(r *bytes.Reader) (*header.Header, error) {
			r.Read(make([]byte, 32))
			return &header.Header{Index: 42}, nil
		},
	}

	opts := block.DeserializeOptions{
		Header: headerMock,
		Tx:     txMock,
	}

	blk, err := block.DeserializeBlockWithparallelPooled(dataBytes, opts)
	if err != nil {
		t.Fatalf("failed to deserialize: %v", err)
	}

	if len(blk.Transaction) != int(txCount) {
		t.Errorf("got %d transactions, want %d", len(blk.Transaction), txCount)
	}

	if blk.Size != 12345 {
		t.Errorf("got size %d, want 12345", blk.Size)
	}

	if !bytes.Equal(blk.Hash, bytes.Repeat([]byte{0xCC}, 32)) {
		t.Error("hash mismatch")
	}

	mu.Lock()

	if callCount != int(txCount) {
		t.Errorf("tx mock called %d times, want %d", callCount, txCount)
	}

	mu.Unlock()
}

func TestDeserializeBlockWithDeps_ZeroTransactions(t *testing.T) {
	data := helpers.CreateValidBlockData(t, 0, []transaction.Transaction{})
	opts := block.DeserializeOptions{
		Header: &mocks.MockHeaderDeserializer{},
		Tx:     &mocks.MockTransactionDeserializer{},
	}

	blk, err := block.DeserializeBlockWithparallelPooled(data.Bytes(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blk.Transaction) != 0 {
		t.Errorf("expected 0 transactions, got %d", len(blk.Transaction))
	}
}

func TestDeserializeBlockWithDeps_TransactionCountTooHigh(t *testing.T) {
	buf := &bytes.Buffer{}

	buf.Write(bytes.Repeat([]byte{0xAA}, 32))
	binary.Write(buf, binary.LittleEndian, uint32(10001))
	buf.Write(bytes.Repeat([]byte{0xCC}, 32))
	binary.Write(buf, binary.LittleEndian, uint32(12345))

	opts := block.DeserializeOptions{
		Header: &mocks.MockHeaderDeserializer{},
		Tx:     &mocks.MockTransactionDeserializer{},
	}

	_, err := block.DeserializeBlockWithparallelPooled(buf.Bytes(), opts)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := "transaction count too high: 10001 (max: 10000)"
	if err.Error() != expected {
		t.Errorf("got %q, want %q", err.Error(), expected)
	}
}

func TestDeserializeBlockWithDeps_TransactionError(t *testing.T) {
	buf := &bytes.Buffer{}

	buf.Write(bytes.Repeat([]byte{0xAA}, 32))
	binary.Write(buf, binary.LittleEndian, uint32(2))

	tx1 := helpers.CreateTestTransaction(1)

	tx1Data, _ := tx1.TransactionSerialize()
	buf.Write(tx1Data)

	tx2 := helpers.CreateTestTransaction(2)

	tx2Data, _ := tx2.TransactionSerialize()
	buf.Write(tx2Data)

	buf.Write(bytes.Repeat([]byte{0xCC}, 32))
	binary.Write(buf, binary.LittleEndian, uint32(12345))

	txMock := &mocks.MockTransactionDeserializer{
		MockFunc: func(r *bytes.Reader) (transaction.Transaction, error) {
			sender := make([]byte, 32)
			r.Read(sender)

			recipient := make([]byte, 32)
			r.Read(recipient)

			id := make([]byte, 32)
			r.Read(id)

			var amount float64
			binary.Read(r, binary.LittleEndian, &amount)

			var sigLen uint32
			binary.Read(r, binary.LittleEndian, &sigLen)

			signature := make([]byte, sigLen)
			r.Read(signature)

			if id[0] == 2 {
				return nil, errors.New("simulated error")
			}

			return &helpers.TestTransaction{
				Id:        id,
				Sender:    sender,
				Recipient: recipient,
				Amount:    amount,
				Signature: signature,
			}, nil
		},
	}
	opts := block.DeserializeOptions{
		Header: &mocks.MockHeaderDeserializer{},
		Tx:     txMock,
	}

	_, err := block.DeserializeBlockWithparallelPooled(buf.Bytes(), opts)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "tx 1 missing" {
		t.Errorf("expected 'tx 1 missing', got %q", err.Error())
	}

	if txMock.GetCallCount() != 2 {
		t.Errorf("expected 2 calls, got %d", txMock.GetCallCount())
	}
}

func TestDeserializeBlockWithDeps_ExtraData(t *testing.T) {
	buf := &bytes.Buffer{}

	buf.Write(bytes.Repeat([]byte{0xAA}, 32))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	buf.Write(bytes.Repeat([]byte{0xCC}, 32))
	binary.Write(buf, binary.LittleEndian, uint32(100))
	buf.WriteByte(0xFF) // extra

	opts := block.DeserializeOptions{
		Header: &mocks.MockHeaderDeserializer{},
		Tx:     &mocks.MockTransactionDeserializer{},
	}

	blk, err := block.DeserializeBlockWithparallelPooled(buf.Bytes(), opts)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := "extra data after block deserialization: 1 bytes remaining"

	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}

	if blk != nil {
		t.Error("block should be nil")
	}
}

func TestDeserializeBlockWithDeps_InvalidBlockSize(t *testing.T) {
	buf := &bytes.Buffer{}

	buf.Write(bytes.Repeat([]byte{0xAA}, 32))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	buf.Write(bytes.Repeat([]byte{0xCC}, 32))
	binary.Write(buf, binary.LittleEndian, uint32(0)) // size 0

	opts := block.DeserializeOptions{
		Header: &mocks.MockHeaderDeserializer{},
		Tx:     &mocks.MockTransactionDeserializer{},
	}

	blk, err := block.DeserializeBlockWithparallelPooled(buf.Bytes(), opts)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := "invalid block size"
	if err.Error() != expected {
		t.Errorf("got %q, want %q", err.Error(), expected)
	}

	if blk != nil {
		t.Error("block should be nil")
	}
}

func TestDeserializeBlockWithDeps_MaxTransactions(t *testing.T) {
	txCount := uint32(10000)
	buf := &bytes.Buffer{}

	buf.Write(bytes.Repeat([]byte{0xAA}, 32))
	binary.Write(buf, binary.LittleEndian, txCount)

	for i := uint32(0); i < txCount; i++ {
		tx := helpers.CreateTestTransaction(byte(i % 256))

		txData, _ := tx.TransactionSerialize()

		buf.Write(txData)
	}

	buf.Write(bytes.Repeat([]byte{0xCC}, 32))
	binary.Write(buf, binary.LittleEndian, uint32(12345))

	txMock := &mocks.MockTransactionDeserializer{
		MockFunc: func(r *bytes.Reader) (transaction.Transaction, error) {
			sender := make([]byte, 32)
			r.Read(sender)

			recipient := make([]byte, 32)
			r.Read(recipient)

			id := make([]byte, 32)
			r.Read(id)

			var amount float64
			binary.Read(r, binary.LittleEndian, &amount)

			var sigLen uint32
			binary.Read(r, binary.LittleEndian, &sigLen)

			signature := make([]byte, sigLen)
			r.Read(signature)

			return &helpers.TestTransaction{
				Id:        id,
				Sender:    sender,
				Recipient: recipient,
				Amount:    amount,
				Signature: signature,
			}, nil
		},
	}

	opts := block.DeserializeOptions{
		Header: &mocks.MockHeaderDeserializer{},
		Tx:     txMock,
	}

	blk, err := block.DeserializeBlockWithparallelPooled(buf.Bytes(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blk.Transaction) != int(txCount) {
		t.Errorf("got %d transactions, want %d", len(blk.Transaction), txCount)
	}

	if txMock.GetCallCount() != int(txCount) {
		t.Errorf("expected %d calls, got %d", txCount, txMock.GetCallCount())
	}
}
