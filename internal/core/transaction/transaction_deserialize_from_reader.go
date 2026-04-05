package transaction

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

func DeserializeTransactionFromReader(buf *bytes.Reader) (*BankTransaction, error) {
	if buf == nil {
		return nil, errors.New("buffer is nil")
	}

	tx := &BankTransaction{}

	tx.Sender = make([]byte, 32)
	n, err := buf.Read(tx.Sender)
	if err != nil {
		return nil, fmt.Errorf("failed to read sender: %w", err)
	}
	if n != 32 {
		return nil, fmt.Errorf("invalid sender length: expected 32, got %d", n)
	}

	tx.Recipient = make([]byte, 32)
	n, err = buf.Read(tx.Recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to read recipient: %w", err)
	}
	if n != 32 {
		return nil, fmt.Errorf("invalid recipient length: expected 32, got %d", n)
	}

	var amountBits uint64
	if err := binary.Read(buf, binary.LittleEndian, &amountBits); err != nil {
		return nil, fmt.Errorf("failed to read amount: %w", err)
	}
	tx.Amount = math.Float64frombits(amountBits)

	var sigLen uint32
	if err := binary.Read(buf, binary.LittleEndian, &sigLen); err != nil {
		return nil, fmt.Errorf("failed to read signature length: %w", err)
	}

	if sigLen > 1024 {
		return nil, fmt.Errorf("signature length too large: %d", sigLen)
	}

	tx.Signature = make([]byte, sigLen)
	n, err = buf.Read(tx.Signature)
	if err != nil {
		return nil, fmt.Errorf("failed to read signature: %w", err)
	}
	if n != int(sigLen) {
		return nil, fmt.Errorf("invalid signature length: expected %d, got %d", sigLen, n)
	}

	tx.ID = make([]byte, 32)
	n, err = buf.Read(tx.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to read transaction ID: %w", err)
	}
	if n != 32 {
		return nil, fmt.Errorf("invalid ID length: expected 32, got %d", n)
	}

	return tx, nil
}
