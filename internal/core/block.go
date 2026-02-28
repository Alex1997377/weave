package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"time"

	"github.com/Alex1997377/weave/internal/crypto"
)

type Header struct {
	Index        int    `json:"index"`
	Timestamp    int64  `json:"timestamp"`
	PreviousHash []byte `json:"previoues_hash"`
	MerkleRoot   []byte `json:"merkle_root"`
	Nonce        int    `json:"nonce"`
	Difficulty   int    `json:"difficulty"`
}

type Block struct {
	Header      Header
	Transaction []Transaction
	Hash        []byte `json:"hash"` // The unique "fingerprint" of the current block, calculated based on all the fields above.
}

func (h *Header) Serialize() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, int64(h.Index))
	binary.Write(buf, binary.LittleEndian, h.Timestamp)
	buf.Write([]byte(h.PreviousHash))
	buf.Write([]byte(h.MerkleRoot))
	binary.Write(buf, binary.LittleEndian, int64(h.Nonce))
	binary.Write(buf, binary.LittleEndian, int64(h.Difficulty))

	return buf.Bytes()
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
