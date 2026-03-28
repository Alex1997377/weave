package mocks

type MockTransaction struct {
	Id []byte
}

func (m *MockTransaction) TransactionGetID() []byte { return m.Id }

func (m *MockTransaction) TransactionGetSender() []byte { return nil }

func (m *MockTransaction) TransactionGetRecipient() []byte { return nil }

func (m *MockTransaction) TransactionGetAmount() float64 { return 0 }

func (m *MockTransaction) TransactionValidate() error { return nil }

func (m *MockTransaction) TransactionSign([]byte) error { return nil }

func (m *MockTransaction) TransactionVerify([]byte) bool { return true }

func (m *MockTransaction) TransactionSerialize() ([]byte, error) { return nil, nil }
