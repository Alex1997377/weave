package chain

const (
	ErrInvalidBlock     = "INVALID_BLOCK"
	ErrInvalidHash      = "INVALID_HASH"
	ErrInvalidSignature = "INVALID_SIGNATURE"
	ErrBlockNotFound    = "BLOCK_NOT_FOUND"
	ErrChainCorrupted   = "CHAIN_CORRUPTED"
	ErrCreateWallet     = "GENERATE_KEY_PAIR_ERROR"
	ErrInvalidAddress   = "INVALID_ADDRESS"
)

const (
	DIFFICULTY         = 4
	MAX_BLOCK_SIZE     = 1 << 20
	PARALLEL_THRESHOLD = 100
)
