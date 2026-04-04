package header

type Header struct {
	Index        int
	Timestamp    int64
	PreviousHash []byte
	MerkleRoot   []byte
	Nonce        uint64
	Difficulty   int
}
