package block

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Alex1997377/weave/internal/core/header"
	"github.com/Alex1997377/weave/internal/core/transaction"
	"github.com/Alex1997377/weave/internal/crypto/hash"
)

type Block struct {
	Header      header.Header
	Transaction []transaction.Transaction
	Hash        hash.Hash
	Size        uint32
}

func NewBlock(
	transactions []transaction.Transaction,
	PreviousHash []byte,
	index int,
	difficulty int) (*Block, error) {

	if PreviousHash == nil {
		return nil, errors.New("previous hash cannot be nil")
	}

	if index < 0 {
		return nil, fmt.Errorf("block index cannot be negative: %d", index)
	}

	if difficulty < 0 {
		return nil, fmt.Errorf("difficulty cannot be negative: %d", difficulty)
	}

	if index > 0 && len(transactions) == 0 {
		return nil, errors.New("non-genesis block must have at least one transaction")
	}

	for i, tx := range transactions {
		if tx == nil {
			return nil, fmt.Errorf("transaction at index %d is nil", i)
		}
	}

	block := &Block{
		Header: header.Header{
			Index:        index,
			Timestamp:    time.Now().Unix(),
			PreviousHash: PreviousHash,
			Difficulty:   difficulty,
			Nonce:        0,
			MerkleRoot:   nil,
		},
		Transaction: transactions,
	}

	if err := block.SetMerkleRoot(); err != nil {
		return nil, fmt.Errorf("failed to set merkle root: %w", err)
	}

	if err := block.Mine(context.Background(), MineConfig{}); err != nil {
		return nil, fmt.Errorf("failed to mine block: %w", err)
	}

	size, err := block.CalculateSize()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate block size: %w", err)
	}
	block.Size = size

	if err := block.Validate(); err != nil {
		return nil, fmt.Errorf("created block is invalid: %w", err)
	}

	return block, nil
}
