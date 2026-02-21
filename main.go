package main

import (
	"fmt"
	"time"
)

type Block struct {
	Index     int
	Timestamp time.Time
	Data      string
	Hash      string
}

type Blockchain struct {
	Blocks []*Block
}

func NewBlock(index int, data string) *Block {
	return &Block{
		Index:     index,
		Timestamp: time.Now(),
		Data:      data,
		Hash:      fmt.Sprintf("hash-%d", index),
	}
}

func (bc *Blockchain) AddBlock(data string) {
	index := len(bc.Blocks)
	block := NewBlock(index, data)
	bc.Blocks = append(bc.Blocks, block)
}

// Display displays all blocks in the blockchain
func (bc *Blockchain) Display() {
	fmt.Println("=== Blockchain ===")
	for _, block := range bc.Blocks {
		fmt.Printf("Block #%d\n", block.Index)
		fmt.Printf("Timestamp: %s\n", block.Timestamp.Format(time.RFC3339))
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Println("---")
	}
}

func main() {
	fmt.Println("🚀 Hello Blockchain!")
	fmt.Println("Creating your first blockchain...")

	// Create a new blockchain
	blockchain := &Blockchain{}

	// Add some blocks
	blockchain.AddBlock("Genesis Block - The beginning of our blockchain!")
	blockchain.AddBlock("Second Block - Learning Go and blockchain!")
	blockchain.AddBlock("Third Block - Building something amazing!")

	// Display the blockchain
	blockchain.Display()

	fmt.Println("✅ Your first blockchain is ready!")
	fmt.Println("In the next sections, we'll build a real blockchain with:")
	fmt.Println("- Cryptographic hashing")
	fmt.Println("- Proof of work mining")
	fmt.Println("- Transaction processing")
	fmt.Println("- Advanced consensus algorithms")
}
