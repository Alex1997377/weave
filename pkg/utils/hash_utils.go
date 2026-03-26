package utils

import (
	"errors"
	"fmt"

	"github.com/Alex1997377/weave/internal/crypto"
)

func HashToString(hash crypto.Hash) (string, error) {
	if hash == nil {
		return "", errors.New("hash is nil")
	}

	if len(hash) == 0 {
		return "", errors.New("hash is empty")
	}

	return BytesToHex(hash), nil
}

func HashFromString(hashStr string) (crypto.Hash, error) {
	if hashStr == "" {
		return nil, errors.New("hash string is empty")
	}

	if len(hashStr) != 64 {
		return nil, fmt.Errorf("invalid hash length: expected 64 characters, got %d", len(hashStr))
	}

	bytes, err := HexToBytes(hashStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hash: %w", err)
	}

	return crypto.Hash(bytes), nil
}

func MustHashToString(hash crypto.Hash) string {
	str, err := HashToString(hash)
	if err != nil {
		panic(fmt.Sprintf("fialed to convert hash to string: %v", err))
	}
	return str
}

func BlockHashString(block interface{ GetHahs() crypto.Hash }) (string, error) {
	if block == nil {
		return "", errors.New("block is nil")
	}

	hash := block.GetHahs()
	if hash == nil {
		return "", errors.New("block hash is nil")
	}

	return HashToString(hash)
}

func IsValidHahs(hashStr string) bool {
	if len(hashStr) != 64 {
		return false
	}

	_, err := HexToBytes(hashStr)
	return err == nil
}

func HashEquals(h1, h2 crypto.Hash) bool {
	if h1 == nil || h2 == nil {
		return h1 == nil && h2 == nil
	}

	if len(h1) != len(h2) {
		return false
	}

	for i := range h1 {
		if h1[i] != h2[i] {
			return false
		}
	}

	return true
}

func HashBytes(hash crypto.Hash) []byte {
	if hash == nil {
		return nil
	}

	result := make([]byte, len(hash))
	copy(result, hash)
	return result
}

func HashFromBytes(data []byte) (crypto.Hash, error) {
	if data == nil {
		return nil, errors.New("data is nil")
	}

	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}

	if len(data) != 32 {
		return nil, fmt.Errorf("invalid hash length: expected 32 bytes, got %d", len(data))
	}

	return crypto.Hash(data), nil
}

func ShortHashString(hash crypto.Hash) (string, error) {
	full, err := HashToString(hash)
	if err != nil {
		return "", err
	}

	if len(full) > 8 {
		return full[:8], nil
	}

	return full, nil
}

func FormatHahsWithPrefix(prefix string, hash crypto.Hash) (string, error) {
	hashStr, err := HashToString(hash)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s: %s", prefix, hashStr), nil
}
