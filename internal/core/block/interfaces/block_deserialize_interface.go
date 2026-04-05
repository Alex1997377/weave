package interfaces

import (
	"bytes"

	"github.com/Alex1997377/weave/internal/core/header"
	"github.com/Alex1997377/weave/internal/core/transaction"
)

// HeaderDeserializer определяет контракт для десериализации заголовка блока
type HeaderDeserializer interface {
	DeserializeHeader(r *bytes.Reader) (*header.Header, error)
}

// TransactionDeserializer определяет контракт для десериализации транзакции
type TransactionDeserializer interface {
	DeserializeTransaction(r *bytes.Reader) (transaction.Transaction, error)
}

// RealHeaderDeserializer - реальная реализация, использующая header.DeserializeHeader
type RealHeaderDeserializer struct{}

func (RealHeaderDeserializer) DeserializeHeader(r *bytes.Reader) (*header.Header, error) {
	return header.DeserializeHeader(r)
}

// RealTransactionDeserializer - реальная реализация,
// использующая transaction.DeserializeTransactionFromReader
type RealTransactionDeserializer struct{}

func (RealTransactionDeserializer) DeserializeTransaction(r *bytes.Reader) (transaction.Transaction, error) {
	return transaction.DeserializeTransactionFromReader(r)
}
