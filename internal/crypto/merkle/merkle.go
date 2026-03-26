package merkle

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
)

func CalculateMerkleRoot(hashes [][]byte) ([]byte, error) {
	if len(hashes) == 0 {
		return make([]byte, 32), nil
	}

	for i, hash := range hashes {
		if err := validateHash(hash, i); err != nil {
			return nil, err
		}
	}

	tree := make([][]byte, len(hashes))
	for i, hash := range hashes {
		tree[i] = make([]byte, len(hash))
		copy(tree[i], hash)
	}

	for len(tree) > 1 {
		if len(tree)%2 != 0 {
			lastHash := make([]byte, len(tree[len(tree)-1]))
			copy(lastHash, tree[len(tree)-1])
			tree = append(tree, lastHash)
		}

		var level [][]byte
		for i := 0; i < len(tree); i += 2 {
			if i+1 >= len(tree) {
				return nil, &MerkleRootError{
					Op:      "build_level",
					Index:   i,
					Message: "unexpected end of tree",
					Err:     errors.New("incomplete pair"),
				}
			}

			pair := bytes.Join([][]byte{tree[i], tree[i+1]}, []byte{})
			if len(pair) == 0 {
				return nil, &MerkleRootError{
					Op:      "concatenate",
					Index:   i,
					Message: "empty pair after concatenation",
				}
			}

			hash := sha256.Sum256(pair)

			level = append(level, hash[:])
		}

		tree = level
	}

	if len(tree) == 0 {
		return nil, &MerkleRootError{
			Op:      "final",
			Index:   -1,
			Message: "tree became empty during calculation",
			Err:     ErrMerkleRootCalculate,
		}
	}

	if tree[0] == nil {
		return nil, &MerkleRootError{
			Op:      "final",
			Index:   0,
			Message: "calculated root is nil",
			Err:     ErrMerkleRootCalculate,
		}
	}

	if len(tree[0]) != 32 {
		return nil, &MerkleRootError{
			Op:      "final",
			Index:   0,
			Message: fmt.Sprintf("invalid root hash length: %d", len(tree[0])),
			Err:     ErrInvalidHashLength,
		}
	}

	return tree[0], nil
}

func validateHash(hash []byte, index int) error {
	if hash == nil {
		return &MerkleRootError{
			Op:      "validate",
			Index:   index,
			Message: "hash is nil",
			Err:     ErrNilHashIsSlice,
		}
	}

	if len(hash) == 0 {
		return &MerkleRootError{
			Op:      "validate",
			Index:   index,
			Message: "hash is empty",
			Err:     ErrEmptyHash,
		}
	}

	if len(hash) != 32 {
		return &MerkleRootError{
			Op:      "validate",
			Index:   index,
			Message: fmt.Sprintf("invalid hash length: expected 32, got %d", len(hash)),
			Err:     ErrInvalidHashLength,
		}
	}

	return nil
}

func CalculateMerkleRootWithContext(hashes [][]byte, context string) ([]byte, error) {
	root, err := CalculateMerkleRoot(hashes)
	if err != nil {
		return nil, fmt.Errorf("merkle root calculation failed for %s: %w", context, err)
	}
	return root, nil
}

func VerifyMerkleProof(hash []byte, proof [][]byte, root []byte, index int) (bool, error) {
	if err := validateHash(hash, -1); err != nil {
		return false, fmt.Errorf("invalid hash to verify: %w", err)
	}

	if err := validateHash(root, -1); err != nil {
		return false, fmt.Errorf("invalid root hash: %w", err)
	}

	if index < 0 {
		return false, errors.New("index cannot be negative")
	}

	currentHash := make([]byte, len(hash))
	copy(currentHash, hash)

	for i, proofHash := range proof {
		if err := validateHash(proofHash, i); err != nil {
			return false, fmt.Errorf("invalid proof hash at index %d: %w", i, err)
		}

		var pair []byte
		if index%2 == 0 {
			pair = bytes.Join([][]byte{currentHash, proofHash}, []byte{})
		} else {
			pair = bytes.Join([][]byte{proofHash, currentHash}, []byte{})
		}

		computedHash := sha256.Sum256(pair)
		currentHash = computedHash[:]
		index >>= 1
	}

	return bytes.Equal(currentHash, root), nil
}
