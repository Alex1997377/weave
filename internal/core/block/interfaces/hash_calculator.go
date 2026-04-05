package interfaces

type Hash interface {
	IsValidForDifficulty(difficulty int) bool
	Bytes() []byte
}

type HashCalculator interface {
	Hash(data []byte) Hash
}
