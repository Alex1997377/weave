package mocks

import (
	"bytes"
	"sync"
	"testing"

	"github.com/Alex1997377/weave/internal/core/header"
	"github.com/Alex1997377/weave/internal/core/transaction"
)

type MockHeaderDeserializer struct {
	MockFunc func(r *bytes.Reader) (*header.Header, error)
}

func (m MockHeaderDeserializer) DeserializeHeader(r *bytes.Reader) (*header.Header, error) {
	if m.MockFunc != nil {
		return m.MockFunc(r)
	}
	r.Seek(32, 0)
	return &header.Header{
		Index:        0,
		Timestamp:    1234567890,
		PreviousHash: bytes.Repeat([]byte{0xAA}, 32),
		MerkleRoot:   bytes.Repeat([]byte{0xBB}, 32),
		Nonce:        0,
		Difficulty:   1,
	}, nil
}

type MockTransactionDeserializer struct {
	mu           sync.Mutex
	MockFunc     func(r *bytes.Reader) (transaction.Transaction, error)
	CallCount    int
	Transactions []transaction.Transaction
}

func (m *MockTransactionDeserializer) DeserializeTransaction(r *bytes.Reader) (transaction.Transaction, error) {
	m.mu.Lock()
	m.CallCount++
	callCount := m.CallCount
	mockFunc := m.MockFunc
	transactions := m.Transactions
	m.mu.Unlock()

	if mockFunc != nil {
		return mockFunc(r)
	}
	if transactions != nil && callCount-1 < len(transactions) {
		return transactions[callCount-1], nil
	}
	return nil, nil // fallback (тесты должны задать поведение)
}

func (m *MockTransactionDeserializer) GetCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.CallCount
}

func (m *MockTransactionDeserializer) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CallCount = 0
}

func (m *MockTransactionDeserializer) AssertCalled(t *testing.T, expected int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.CallCount != expected {
		t.Errorf("expected %d calls, got %d", expected, m.CallCount)
	}
}
