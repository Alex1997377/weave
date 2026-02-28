package utils

import (
	"encoding/hex"
	"fmt"
)

// Convert byte`s slice into string
func BytesToHex(data []byte) string {
	return hex.EncodeToString(data)
}

func HexToBytes(s string) ([]byte, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %w", err)
	}
	return b, nil
}
