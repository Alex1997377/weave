package core

import (
	"errors"
)

type Transaction interface {
	GetID() []byte
	GetSender() []byte
	GetRecipient() []byte
	GetAmount() float64
	Validate() error
	Sign(privateKey []byte) error
	Verify(publicKey []byte) bool
}

type BankTransaction struct {
	ID        []byte  `json:"id"`
	Sender    []byte  `json:"sender"`
	Recipient []byte  `json:"recipient"`
	Amount    float64 `json:"amount"`
	Signature []byte  `json:"signature"`
}

func (bt *BankTransaction) GetID() []byte {
	return bt.ID
}

func (bt *BankTransaction) GetSender() []byte {
	return bt.Sender
}

func (bt *BankTransaction) GetRecipient() []byte {
	return bt.Recipient
}

func (bt *BankTransaction) GetAmount() float64 {
	return bt.Amount
}

func (bt *BankTransaction) Validate() error {
	if bt.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	if len(bt.Sender) == 0 || len(bt.Recipient) == 0 {
		return errors.New("sender and resipient cannot be empty")
	}
	return nil
}

func (bt *BankTransaction) Sign(privateKey []byte) error {
	// TODO signing implementation
	return nil
}

func (bt *BankTransaction) Verify(publicKey []byte) bool {
	// TODO varification implementation
	return true
}
