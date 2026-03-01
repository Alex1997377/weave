package core

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"time"

	"github.com/Alex1997377/weave/internal/crypto"
)

type Block struct {
	Header      Header
	Transaction []Transaction
	Hash        Hash // The unique "fingerprint" of the current block, calculated based on all the fields above.
	Size        int  `json:"size"`
}

func (b *Block) CalculateHash() []byte {
	data := b.Header.Serialize()
	hash := sha256.Sum256(data)
	return hash[:]
}

func (b *Block) SetMerkleRoot() {
	var ids [][]byte
	for _, tx := range b.Transaction {
		ids = append(ids, tx.GetID())
	}

	b.Header.MerkleRoot = crypto.CalculateMerkleRoot(ids)
}

func (b *Block) Serialize() []byte {
	buf := new(bytes.Buffer)

	buf.Write(b.Header.Serialize())

	for _, tx := range b.Transaction {
		buf.Write(tx.Serialize())
	}

	return buf.Bytes()
}

func (b *Block) CalculateSize() int {
	headerSize := len(b.Header.Serialize())

	transactionsSize := 0
	for _, tx := range b.Transaction {
		transactionsSize += len(tx.Serialize())
	}

	hashSize := len(b.Hash)

	return headerSize + transactionsSize + hashSize
}

func (b *Block) Validate() error {
	if !b.Hash.IsValidForDifficulty(b.Header.Difficulty) {
		return errors.New("invalid proof of work")
	}

	calculatedHash := b.CalculateHash()
	if !bytes.Equal(b.Hash[:], calculatedHash[:]) {
		return errors.New("block hash doesn`t match content")
	}

	return nil
}

func NewBlock(
	transaction []Transaction,
	PreviousHash []byte,
	index int,
	difficulty int) *Block {

	// Fills in the time, data, and the reference to the previous block.
	block := &Block{
		Header: Header{
			Index:        index,
			Timestamp:    time.Now().Unix(),
			PreviousHash: PreviousHash,
			Difficulty:   difficulty,
			Nonce:        0,
		},
		Transaction: transaction,
	}

	block.SetMerkleRoot()

	// block.Mine()

	// Crucially, it calls CalculateHash at the end so the block gets its own unique ID
	block.Hash = block.CalculateHash()

	return block
}
