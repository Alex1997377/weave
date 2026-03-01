package core

import "github.com/Alex1997377/weave/internal/crypto"

type Header struct {
	Index        int    `json:"index"`
	Timestamp    int64  `json:"timestamp"`
	PreviousHash []byte `json:"previoues_hash"`
	MerkleRoot   []byte `json:"merkle_root"`
	Nonce        int    `json:"nonce"`
	Difficulty   int    `json:"difficulty"`
}

func (h *Header) Serialize() []byte {
	return crypto.SerializeHeader(
		h.Index,
		h.Timestamp,
		h.PreviousHash,
		h.MerkleRoot,
		h.Nonce,
		h.Difficulty,
	)
}
