package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Block struct {
	Timestamp     int64  // A time marker. It is needed to keep each block unique, even if the Data is identical.
	Data          []byte // The payload (transactions, messages).
	PrevBlockHash []byte // This is the "glue" that stores the hash of the previous block, linking them into a chain.
	Hash          []byte // The unique "fingerprint" of the current block, calculated based on all the fields above.
}

// Turns the block's data into a unique fixed-length string
func (b *Block) CalculateHash() []byte {
	// Combines all block fields (PrevBlockHash, Data, and time) into a single long byte slice.
	headers := bytes.Join(
		[][]byte{
			b.PrevBlockHash,
			b.Data,
			[]byte(fmt.Sprintf("%d", b.Timestamp)),
		},
		[]byte{},
	)
	// A mathematical function. If even one bit of the source data is changed,
	// the result (hash) will change beyond recognition.
	hash := sha256.Sum256(headers)
	return hash[:]
}

// Creates a block object and "signs" it immediately.
func NewBlock(data string, prevBlockHash []byte) *Block {

	// Fills in the time, data, and the reference to the previous block.
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
	}

	// Crucially, it calls CalculateHash at the end so the block gets its own unique ID
	block.Hash = block.CalculateHash()
	return block
}

type Blockchain struct {
	Blocks []*Block
}

// Adds a new link to the end of the chain
func (bc *Blockchain) AddBlock(data string) {
	var prevHash []byte

	// Checks if there are already blocks in the chain.
	if len(bc.Blocks) > 0 {
		prevHash = bc.Blocks[len(bc.Blocks)-1].Hash
	} else {
		prevHash = []byte{}
	}

	newBlock := NewBlock(data, prevHash)
	// Adds the newly created block to the blockchain slice (array
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
