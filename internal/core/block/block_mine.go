package block

import (
	"errors"
	"fmt"
	"time"

	"github.com/Alex1997377/weave/internal/crypto"
)

func (b *Block) Mine() error {
	if b == nil {
		return errors.New("block is nil")
	}

	if b.Header.Difficulty < 0 {
		return errors.New("block difficulty cannot be negative")
	}

	fmt.Printf("Mining block %d with difficulty %d...\n",
		b.Header.Index, b.Header.Difficulty)

	startTime := time.Now()
	hashAttempts := 0

	for {
		hashBytes, err := b.CalculateHash()
		if err != nil {
			return fmt.Errorf("failed to calculate hash during mining: %w", err)
		}

		hash := crypto.Hash(hashBytes)

		if hash.IsValidForDifficulty(b.Header.Difficulty) {
			b.Hash = hash
			miningDuration := time.Since(startTime)
			fmt.Printf("Mined! Hash: %s (attempts: %d, time: %v)\n",
				hash.String(), hashAttempts, miningDuration)
			break
		}

		b.Header.Nonce++
		hashAttempts++

		if b.Header.Nonce < 0 {
			return errors.New("nonce overflow during mining")
		}
	}

	return nil
}
