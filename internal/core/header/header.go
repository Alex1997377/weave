package header

import (
	"github.com/Alex1997377/weave/internal/core/header/errors"
	"github.com/Alex1997377/weave/internal/core/header/errors/constants"
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

func (h *Header) Serialize() ([]byte, error) {
	if h == nil {
		return nil, errors.NewNilHeaderError(constants.OpSerialize)
	}

	if h.Index < 0 {
		return nil, errors.NewIndexError(
			constants.OpSerialize,
			h.Index)
	}

	if h.Timestamp <= 0 {
		return nil, errors.NewTimestampError(
			constants.OpSerialize,
			h.Timestamp, "must be positive")
	}

	if h.PreviousHash == nil {
		return nil, errors.NewHashError(
			constants.OpSerialize,
			constants.FieldPreviousHash, nil, 32)
	}

	if len(h.PreviousHash) == 0 {
		return nil, errors.NewHashError(
			constants.OpSerialize,
			constants.FieldPreviousHash,
			h.PreviousHash, 32)
	}

	if len(h.PreviousHash) != 32 {
		return nil, errors.NewHashError(
			constants.OpSerialize,
			constants.FieldPreviousHash,
			h.PreviousHash, 32)
	}

	if h.MerkleRoot != nil && len(h.MerkleRoot) == 0 {
		return nil, errors.NewHashError(
			constants.OpSerialize,
			constants.FieldMerkleRoot,
			h.MerkleRoot, 32)
	}

	if h.Difficulty < 0 {
		return nil, errors.NewDifficultyError(
			constants.OpSerialize,
			h.Difficulty, 0, 255)
	}

	if h.Nonce < 0 {
		return nil, errors.NewNonceError(
			constants.OpSerialize, int64(h.Nonce))
	}

	data, err := crypto.SerializeHeader(
		h.Index,
		h.Timestamp,
		h.PreviousHash,
		h.MerkleRoot,
		h.Nonce,
		h.Difficulty,
	)

	if err != nil {
		return nil, errors.NewSerializationError(
			constants.OpSerialize,
			"crypto_serialize", err)
	}

	if data == nil {
		return nil, errors.NewSerializationError(
			constants.OpSerialize,
			"nil_result", nil)
	}

	if len(data) == 0 {
		return nil, errors.NewSerializationError(
			constants.OpSerialize,
			"empty_result", nil)
	}

	return data, nil
}
