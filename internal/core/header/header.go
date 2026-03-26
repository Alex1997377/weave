package header

type Header struct {
	Index        int    `json:"index"`
	Timestamp    int64  `json:"timestamp"`
	PreviousHash []byte `json:"previoues_hash"`
	MerkleRoot   []byte `json:"merkle_root"`
	Nonce        int    `json:"nonce"`
	Difficulty   int    `json:"difficulty"`
}
