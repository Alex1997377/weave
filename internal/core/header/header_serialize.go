package header

import (
	"bytes"
	"encoding/binary"

	"github.com/Alex1997377/weave/internal/core/header/errors"
	"github.com/Alex1997377/weave/internal/core/header/errors/constants"
	"github.com/Alex1997377/weave/internal/crypto/serialize"
	"github.com/Alex1997377/weave/pkg/utils"
)

func (h *Header) Validate(op string) error {
	if h == nil {
		return errors.NewNilHeaderError(op)
	}
	if h.Index < 0 {
		return errors.NewIndexError(op, h.Index)
	}
	if h.Timestamp <= 0 {
		return errors.NewTimestampError(op, h.Timestamp, "must be positive")
	}

	if err := utils.ValidateHash(op, constants.FieldPreviousHash, h.PreviousHash, true); err != nil {
		return err
	}
	if err := utils.ValidateHash(op, constants.FieldMerkleRoot, h.MerkleRoot, false); err != nil {
		return err
	}

	if h.Difficulty < 0 {
		return errors.NewDifficultyError(op, h.Difficulty, 0, 255)
	}
	if h.Nonce < 0 {
		return errors.NewNonceError(op, int64(h.Nonce))
	}

	return nil
}

func (h *Header) Serialize() ([]byte, error) {
	if err := h.Validate(constants.OpSerialize); err != nil {
		return nil, err
	}

	data, err := serialize.SerializeHeader(
		h.Index, h.Timestamp, h.PreviousHash,
		h.MerkleRoot, h.Nonce, h.Difficulty,
	)

	if err != nil || len(data) == 0 {
		reason := "crypto_serialize"
		if err == nil {
			reason = "empty result"
		}
		return nil, errors.NewSerializationError(constants.OpSerialize, reason, err)
	}

	return data, nil
}

func (h *Header) SerializeWithoutNonce() ([]byte, int, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 104))

	binary.Write(buf, binary.LittleEndian, h.Index)
	binary.Write(buf, binary.LittleEndian, h.Timestamp)

	binary.Write(buf, binary.LittleEndian, uint32(len(h.PreviousHash)))
	buf.Write(h.PreviousHash)

	binary.Write(buf, binary.LittleEndian, uint32(len(h.MerkleRoot)))
	buf.Write(h.MerkleRoot)

	binary.Write(buf, binary.LittleEndian, h.Difficulty)

	nonceOffset := buf.Len()
	binary.Write(buf, binary.LittleEndian, uint64(0))

	return buf.Bytes(), nonceOffset, nil
}
