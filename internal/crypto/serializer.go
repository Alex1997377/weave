package crypto

import (
	"bytes"
	"encoding/binary"
)

func SerializeHeader(index int, timestamp int64, prevHash, merkleRoot []byte, nonce, difficulty int) []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, index)
	binary.Write(buf, binary.LittleEndian, timestamp)
	buf.Write(prevHash)
	buf.Write(merkleRoot)
	binary.Write(buf, binary.LittleEndian, nonce)
	binary.Write(buf, binary.LittleEndian, difficulty)

	return buf.Bytes()
}
