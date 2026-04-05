package transaction

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"fmt"
)

type Transaction interface {
	TransactionGetID() []byte
	TransactionGetSender() []byte
	TransactionGetRecipient() []byte
	TransactionGetAmount() float64
	TransactionValidate() error
	TransactionSign(privateKey []byte) error
	TransactionVerify(publicKey []byte) bool
	TransactionSerialize() ([]byte, error)
}

type BankTransaction struct {
	ID        []byte  `json:"id"`
	Sender    []byte  `json:"sender"`
	Recipient []byte  `json:"recipient"`
	Amount    float64 `json:"amount"`
	Signature []byte  `json:"signature"`
}

func (bt *BankTransaction) TransactionGetID() []byte {
	return bt.ID
}

func (bt *BankTransaction) TransactionGetSender() []byte {
	return bt.Sender
}

func (bt *BankTransaction) TransactionGetRecipient() []byte {
	return bt.Recipient
}

func (bt *BankTransaction) TransactionGetAmount() float64 {
	return bt.Amount
}

func (bt *BankTransaction) TransactionValidate() error {
	if bt.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	if len(bt.Sender) == 0 || len(bt.Recipient) == 0 {
		return errors.New("sender and resipient cannot be empty")
	}
	return nil
}

func (bt *BankTransaction) TransactionSign(privateKey []byte) error {
	privKey, err := x509.ParseECPrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	data, err := bt.TransactionSerialize()
	if err != nil {
		return fmt.Errorf("failed to serialize transaction: %w", err)
	}

	hash := sha256.Sum256(data)

	signature, err := ecdsa.SignASN1(rand.Reader, privKey, hash[:])
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	bt.Signature = signature

	bt.ID = hash[:]

	return nil
}

func (bt *BankTransaction) TransactionVerify(publicKey []byte) bool {
	if bt.Signature == nil {
		return false
	}

	data, err := bt.TransactionSerialize()
	if err != nil {
		return false
	}

	hash := sha256.Sum256(data)

	pubKey, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return false
	}

	ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return false
	}

	return ecdsa.VerifyASN1(ecdsaPubKey, hash[:], bt.Signature)
}
