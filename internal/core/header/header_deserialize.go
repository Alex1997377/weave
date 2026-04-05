package header

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func DeserializeHeader(buf *bytes.Reader) (*Header, error) {
	header := &Header{}

	if err := binary.Read(buf, binary.LittleEndian, &header.Index); err != nil {
		return nil, fmt.Errorf("failed to read index: %w", err)
	}

	if err := binary.Read(buf, binary.LittleEndian, &header.Timestamp); err != nil {
		return nil, fmt.Errorf("failed to read timestamp: %w", err)
	}

	header.PreviousHash = make([]byte, 32)
	n, err := buf.Read(header.PreviousHash)
	if err != nil {
		return nil, fmt.Errorf("failed to read previous hash: %w", err)
	}
	if n != 32 {
		return nil, fmt.Errorf("invalid previous hash length: expected 32, got %d", n)
	}

	header.MerkleRoot = make([]byte, 32)
	n, err = buf.Read(header.MerkleRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to read merkle root: %w", err)
	}
	if n != 32 {
		return nil, fmt.Errorf("invalid previous hash length: expected 32, got %d", n)
	}

	if err := binary.Read(buf, binary.LittleEndian, &header.Nonce); err != nil {
		return nil, fmt.Errorf("failed to read nonce: %w", err)
	}

	if err := binary.Read(buf, binary.LittleEndian, &header.Difficulty); err != nil {
		return nil, fmt.Errorf("failed to read difficulty: %w", err)
	}

	return header, nil

}
