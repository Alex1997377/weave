package core

import (
	"bytes"
	"encoding/hex"
	"fmt"
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
func (bc *Blockchain) AddBlock(transactions []Transaction) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	newBlock := NewBlock(transactions, prevBlock.Hash, prevBlock.Header.Index+1, DIFFICULTY)

	bc.Blocks = append(bc.Blocks, newBlock)
}

// Verifies that the data has not been tampered with.
func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		current := bc.Blocks[i]
		previous := bc.Blocks[i-1]

		// Re-calculates the current block's hash (CalculateHash) and compares it with the one stored inside.
		// If the Data was modified, the hashes won't match.
		if !bytes.Equal(current.Hash, current.CalculateHash()) {
			return false
		}

		// Verifies that the PrevBlockHash of the current block matches the actual Hash of the previous one.
		// This ensures blocks haven't been swapped or removed.
		if !bytes.Equal(current.Header.PreviousHash, previous.Hash) {
			return false
		}
	}
	return true
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
