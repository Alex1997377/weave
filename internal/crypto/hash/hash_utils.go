package hash

import (
	"encoding/hex"
	"errors"
	"fmt"
)

// HashToString конвертирует хеш в строку hex
func HashToString(hash Hash) (string, error) {
	if hash == nil {
		return "", errors.New("hash is nil")
	}
	if len(hash) == 0 {
		return "", errors.New("hash is empty")
	}
	return hex.EncodeToString(hash), nil
}

// HashFromString создаёт хеш из строки hex
func HashFromString(hashStr string) (Hash, error) {
	if hashStr == "" {
		return nil, errors.New("hash string is empty")
	}
	if len(hashStr) != 64 {
		return nil, fmt.Errorf("invalid hash length: expected 64 characters, got %d", len(hashStr))
	}
	bytes, err := hex.DecodeString(hashStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hash: %w", err)
	}
	return Hash(bytes), nil
}

// MustHashToString паникует при ошибке
func MustHashToString(hash Hash) string {
	str, err := HashToString(hash)
	if err != nil {
		panic(fmt.Sprintf("failed to convert hash to string: %v", err))
	}
	return str
}

// ShortHashString возвращает первые 8 символов хеша
func ShortHashString(hash Hash) (string, error) {
	full, err := HashToString(hash)
	if err != nil {
		return "", err
	}
	if len(full) > 8 {
		return full[:8], nil
	}
	return full, nil
}

// FormatHashWithPrefix добавляет префикс к строковому представлению хеша
func FormatHashWithPrefix(prefix string, hash Hash) (string, error) {
	hashStr, err := HashToString(hash)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s: %s", prefix, hashStr), nil
}

// HashEquals сравнивает два хеша
func HashEquals(h1, h2 Hash) bool {
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

// HashFromBytes создаёт хеш из сырых байтов (проверяет длину)
func HashFromBytes(data []byte) (Hash, error) {
	if data == nil {
		return nil, errors.New("data is nil")
	}
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	if len(data) != 32 {
		return nil, fmt.Errorf("invalid hash length: expected 32 bytes, got %d", len(data))
	}
	return Hash(data), nil
}
