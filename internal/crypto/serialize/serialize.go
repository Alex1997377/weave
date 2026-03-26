package serialize

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

func SerializeHeader(index int, timestamp int64, prevHash, merkleRoot []byte, nonce int, difficulty int) ([]byte, error) {
	if index < 0 {
		return nil, &HeaderValidationError{
			Field: "index",
			Value: index,
			Err:   ErrNegativeIndex,
		}
	}

	if timestamp < 0 {
		return nil, &HeaderValidationError{
			Field: "timestamp",
			Value: timestamp,
			Err:   ErrNegativeTimestamp,
		}
	}

	if len(prevHash) == 0 {
		return nil, &HeaderValidationError{
			Field: "prevHash",
			Value: prevHash,
			Err:   ErrEmptyPreviousHash,
		}
	}

	if len(prevHash) != 32 {
		return nil, &HeaderValidationError{
			Field: "prevHash",
			Value: fmt.Sprintf("length: %d", len(prevHash)),
			Err: fmt.Errorf(
				"%w: expected 32 bytes, got %d",
				ErrInvalidHashLength, len(prevHash)),
		}
	}

	if merkleRoot != nil && len(merkleRoot) != 32 {
		return nil, &HeaderValidationError{
			Field: "merkleRoot",
			Value: fmt.Sprintf("length: %d", len(merkleRoot)),
			Err: fmt.Errorf(
				"%w: expected 32 bytes, got %d",
				ErrInvalidHashLength, len(merkleRoot)),
		}
	}

	if nonce < 0 {
		return nil, &HeaderValidationError{
			Field: "nonce",
			Value: nonce,
			Err:   ErrNegativeNonce,
		}
	}

	if difficulty < 0 {
		return nil, &HeaderValidationError{
			Field: "difficulty",
			Value: difficulty,
			Err:   ErrNegativeDifficulty,
		}
	}

	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, int64(index)); err != nil {
		return nil, &SerializeHeaderError{
			Field: "index",
			Value: index,
			Err:   fmt.Errorf("binary write failed: %w", err),
		}
	}

	if err := binary.Write(buf, binary.LittleEndian, timestamp); err != nil {
		return nil, &SerializeHeaderError{
			Field: "timestamp",
			Value: timestamp,
			Err:   fmt.Errorf("binary write failed: %w", err),
		}
	}

	n, err := buf.Write(prevHash)
	if err != nil {
		return nil, &SerializeHeaderError{
			Field: "prevHash",
			Value: fmt.Sprintf("length: %d", len(prevHash)),
			Err:   fmt.Errorf("buffer write failed: %w", err),
		}
	}
	if n != len(prevHash) {
		return nil, &SerializeHeaderError{
			Field: "prevHash",
			Value: fmt.Sprintf("written: %d, expected: %d", n, len(prevHash)),
			Err:   errors.New("incomplete write to buffer"),
		}
	}

	if merkleRoot != nil {
		n, err = buf.Write(merkleRoot)
		if err != nil {
			return nil, &SerializeHeaderError{
				Field: "merkleRoot",
				Value: fmt.Sprintf("length: %d", len(merkleRoot)),
				Err:   fmt.Errorf("buffer write failed: %w", err),
			}
		}
		if n != len(merkleRoot) {
			return nil, &SerializeHeaderError{
				Field: "merkleRoot",
				Value: fmt.Sprintf("written: %d, expected: %d", n, len(merkleRoot)),
				Err:   errors.New("incomplete write to buffer"),
			}
		}
	}

	if err := binary.Write(buf, binary.LittleEndian, int64(nonce)); err != nil {
		return nil, &SerializeHeaderError{
			Field: "nonce",
			Value: nonce,
			Err:   fmt.Errorf("binary write failed: %w", err),
		}
	}

	if err := binary.Write(buf, binary.LittleEndian, int64(difficulty)); err != nil {
		return nil, &SerializeHeaderError{
			Field: "difficulty",
			Value: difficulty,
			Err:   fmt.Errorf("binary write failed: %w", err),
		}
	}

	result := buf.Bytes()
	if len(result) == 0 {
		return nil, &SerializeHeaderError{
			Field: "result",
			Value: nil,
			Err:   errors.New("serialization produced empty result"),
		}
	}

	return result, nil
}
