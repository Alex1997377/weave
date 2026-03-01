package crypto

import (
	"bytes"
	"crypto/sha256"
)

func CalculateMerkleRoot(hashes [][]byte) []byte {
	if len(hashes) == 0 {
		return make([]byte, 32)
	}

	tree := append([][]byte{}, hashes...)

	for len(tree) > 1 {
		if len(tree)%2 != 0 {
			tree = append(tree, tree[len(tree)-1])
		}

		var level [][]byte
		for i := 0; i < len(tree); i += 2 {
			pair := bytes.Join([][]byte{tree[i], tree[i+1]}, []byte{})
			hash := sha256.Sum256(pair)
			level = append(level, hash[:])
		}
		tree = level
	}
	return tree[0]
}
