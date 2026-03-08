package transaction

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

func (bt *BankTransaction) TransactionSerialize() ([]byte, error) {
	buf := new(bytes.Buffer)

	if len(bt.Sender) != 32 {
		return nil, fmt.Errorf("invalid sender length: expected 32, got %d", len(bt.Sender))
	}
	n, err := buf.Write(bt.Sender)
	if err != nil {
		return nil, fmt.Errorf("failed to write sender: %w", err)
	}
	if n != 32 {
		return nil, fmt.Errorf("incomplete sender write: %d bytes written", n)
	}

	if len(bt.Recipient) != 32 {
		return nil, fmt.Errorf("invalid resipient length: expected 32, got %d", len(bt.Recipient))
	}
	n, err = buf.Write(bt.Recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to write recipient: %w", err)
	}
	if n != 32 {
		return nil, fmt.Errorf("incomplete recipient write: %d bytes written", n)
	}

	amountBits := math.Float64bits(bt.Amount)
	err = binary.Write(buf, binary.LittleEndian, amountBits)
	if err != nil {
		return nil, fmt.Errorf("failed to write amount: %w", err)
	}

	sigLen := uint32(len(bt.Signature))
	err = binary.Write(buf, binary.LittleEndian, sigLen)
	if err != nil {
		return nil, fmt.Errorf("failed to write signature length: %w", err)
	}

	if bt.Signature != nil {
		n, err = buf.Write(bt.Signature)
		if err != nil {
			return nil, fmt.Errorf("failed to write signature: %w", err)
		}
		if n != len(bt.Signature) {
			return nil, fmt.Errorf("incomplete signature write: %d/%d bytes written", n, len(bt.Signature))
		}
	}

	if len(bt.ID) != 32 {
		return nil, fmt.Errorf("invalid ID length: expected 32, got %d", len(bt.ID))
	}
	n, err = buf.Write(bt.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to write transaction ID: %w", err)
	}
	if n != 32 {
		return nil, fmt.Errorf("incomplete ID write: %d bytes written", n)
	}

	return buf.Bytes(), nil
}
