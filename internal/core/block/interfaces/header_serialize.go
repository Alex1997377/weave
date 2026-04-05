package interfaces

type HeaderSerializer interface {
	SerializeWithoutNonce() ([]byte, int, error)
}
