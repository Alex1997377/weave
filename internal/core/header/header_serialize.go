package header

import (
	"bytes"
	"encoding/binary"

	"github.com/Alex1997377/weave/internal/core/header/errors"
	"github.com/Alex1997377/weave/internal/core/header/errors/constants"
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

	return nil
}

func (h *Header) Serialize() ([]byte, error) {
	// Получаем сериализованный заголовок с нулевым nonce и смещение
	data, nonceOffset, err := h.SerializeWithoutNonce()
	if err != nil {
		return nil, err
	}
	// Записываем реальный nonce как uint64 (8 байт) в нужное место
	binary.LittleEndian.PutUint64(data[nonceOffset:], uint64(h.Nonce))
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
