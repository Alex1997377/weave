package core

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/Alex1997377/weave/internal/crypto"
)

const DIFFICULTY int = 4

type Blockchain struct {
	Blocks []*Block
}

func NewBlockchain() *Blockchain {
	emptyHash := make([]byte, 32)

	genesisBlock := NewBlock([]Transaction{}, emptyHash, 0, DIFFICULTY)

	return &Blockchain{Blocks: []*Block{genesisBlock}}
}

// Adds a new link to the end of the chain
func (bc *Blockchain) AddBlock(transactions []Transaction) error {
	// check empty blockchain
	if len(bc.Blocks) == 0 {
		return errors.New("cannot add block to empty blockchain")
	}

	// validation tratransactions
	for _, tx := range transactions {
		if err := tx.Validate(); err != nil {
			return NewInvalidBlockError("transaction validation failed", err)
		}
	}

	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(transactions, prevBlock.Hash, prevBlock.Header.Index+1, DIFFICULTY)

	// check size new block`s`
	if newBlock.CalculateSize() > 1024*1024 {
		return errors.New("block size exceeds limit")
	}

	if err := newBlock.Validate(); err != nil {
		return NewInvalidBlockError("new block validation failed", err)
	}

	bc.Blocks = append(bc.Blocks, newBlock)
	return nil
}

func (b *Block) CalculateMerkleRoot() []byte {
	var txIDs [][]byte

	for _, tx := range b.Transaction {
		txIDs = append(txIDs, tx.GetID())
	}

	return crypto.CalculateMerkleRoot(txIDs)
}

// Verifies that the data has not been tampered with.
func (bc *Blockchain) IsValid() error {
	if !bytes.Equal(bc.Blocks[0].Hash, bc.Blocks[0].CalculateHash()) {
		return NewInvalidHashError("Genesis block hash mismatch", nil)
	}

	for i := 1; i < len(bc.Blocks); i++ {
		current := bc.Blocks[i]
		previous := bc.Blocks[i-1]

		// Re-calculates the current block's hash (CalculateHash) and compares it with the one stored inside.
		// If the Data was modified, the hashes won't match.
		if !bytes.Equal(current.Hash, current.CalculateHash()) {
			return NewInvalidHashError(
				fmt.Sprintf("block %d hash mismatch: data has been tampered with", i),
				nil,
			)
		}

		// Verifies that the PrevBlockHash of the current block matches the actual Hash of the previous one.
		// This ensures blocks haven't been swapped or removed.
		if !bytes.Equal(current.Header.PreviousHash, previous.Hash) {
			return NewChainCorruptedError(
				fmt.Sprintf("block %d: PreviousHash does not match hash of block %d", i, i-1),
				nil,
			)
		}

		if !bytes.Equal(current.Header.MerkleRoot, current.CalculateMerkleRoot()) {
			return NewInvalidBlockError(
				fmt.Sprintf("block %d: Merkle Root mismatch (transactions modifies)", i),
				nil,
			)
		}

		hashStr := hex.EncodeToString(current.Hash)
		target := strings.Repeat("0", current.Header.Difficulty)
		if !strings.HasPrefix(hashStr, target) {
			return NewInvalidBlockError(
				fmt.Sprintf("block %d: hash does not satisfy difficulty %d", i, current.Header.Difficulty),
				nil,
			)
		}
	}
	return nil
}

// Display displays all blocks in the blockchain
func (bc *Blockchain) Display() {
	for i, block := range bc.Blocks {
		fmt.Printf("--- Block ID: %d ---\n", i)
		fmt.Printf("Timestamp: 	%d\n", block.Header.Timestamp)
		fmt.Printf("Data: 	   	%s\n", block.Transaction)
		fmt.Printf("Prev Hash:  %s\n", hex.EncodeToString(block.Header.PreviousHash))
		fmt.Printf("Hash: 		%s\n", hex.EncodeToString(block.Hash))
		fmt.Println("  --- ฿ ---  ")
	}
}
