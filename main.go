package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

func (b *Block) CalculateHash() []byte {
	headers := bytes.Join(
		[][]byte{
			b.PrevBlockHash,
			b.Data,
			[]byte(fmt.Sprintf("%d", b.Timestamp)),
		},
		[]byte{},
	)
	hash := sha256.Sum256(headers)
	return hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {

	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
	}

	block.Hash = block.CalculateHash()
	return block
}

type Blockchain struct {
	Blocks []*Block
}

func (bc *Blockchain) AddBlock(data string) {
	var prevHash []byte

	if len(bc.Blocks) > 0 {
		prevHash = bc.Blocks[len(bc.Blocks)-1].Hash
	} else {
		prevHash = []byte{}
	}

	newBlock := NewBlock(data, prevHash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		current := bc.Blocks[i]
		previous := bc.Blocks[i-1]

		if !bytes.Equal(current.Hash, current.CalculateHash()) {
			return false
		}

		if !bytes.Equal(current.PrevBlockHash, previous.Hash) {
			return false
		}
	}
	return true
}

// Display displays all blocks in the blockchain
func (bc *Blockchain) Display() {
	for i, block := range bc.Blocks {
		fmt.Printf("--- Block ID: %d ---\n", i)
		fmt.Printf("Timestamp: 	%d\n", block.Timestamp)
		fmt.Printf("Data: 	   	%s\n", block.Data)
		fmt.Printf("Prev Hash:  %s\n", hex.EncodeToString(block.PrevBlockHash))
		fmt.Printf("Hash: 		%s\n", hex.EncodeToString(block.Hash))
		fmt.Println("  --- ฿ ---  ")
	}
}

func main() {
	fmt.Println("🚀 Hello Blockchain!")
	fmt.Println("Creating your first blockchain...")

	genesisBlock := NewBlock("Check working blockchain", []byte{})
	// Create a new blockchain
	blockchain := &Blockchain{Blocks: []*Block{genesisBlock}}

	// Add some blocks
	blockchain.AddBlock("Genesis Block - The beginning of our blockchain!")
	blockchain.AddBlock("Second Block - Learning Go and blockchain!")

	fmt.Printf("Chain is valid: %v\n", blockchain.IsValid())
	// Display the blockchain
	blockchain.Display()

	fmt.Println("✅ Your first blockchain is ready!")
	fmt.Println("In the next sections, we'll build a real blockchain with:")
	fmt.Println("- Cryptographic hashing")
	fmt.Println("- Proof of work mining")
	fmt.Println("- Transaction processing")
	fmt.Println("- Advanced consensus algorithms")
}
