package chain

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// IsValid проверяет целостность цепочки
func (bc *Blockchain) IsValid() error {
	if len(bc.Blocks) == 0 {
		return errors.New("blockchain is empty")
	}

	for i := 0; i < len(bc.Blocks); i++ {
		current := bc.Blocks[i]

		// Проверяем хеш текущего блока
		currentHash, err := current.CalculateHash()
		if err != nil {
			return fmt.Errorf("failed to calculate hash for block %d: %w", i, err)
		}

		if !bytes.Equal(current.Hash, currentHash) {
			return fmt.Errorf("block %d hash mismatch: data has been tampered with", i)
		}

		// Проверяем ссылку на предыдущий блок (кроме генезиса)
		if i > 0 {
			previous := bc.Blocks[i-1]
			if !bytes.Equal(current.Header.PreviousHash, previous.Hash) {
				return fmt.Errorf("block %d: PreviousHash does not match hash of block %d", i, i-1)
			}
		}

		// Проверяем Merkle root
		expectedMerkleRoot := current.CalculateMerkleRoot()
		if !bytes.Equal(current.Header.MerkleRoot, expectedMerkleRoot) {
			return fmt.Errorf("block %d: Merkle Root mismatch (transactions modified)", i)
		}

		// Проверяем proof of work
		hashStr := hex.EncodeToString(current.Hash)
		target := strings.Repeat("0", current.Header.Difficulty)
		if !strings.HasPrefix(hashStr, target) {
			return fmt.Errorf("block %d: hash does not satisfy difficulty %d", i, current.Header.Difficulty)
		}
	}
	return nil
}

