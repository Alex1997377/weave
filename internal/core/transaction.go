package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

type Transaction interface {
	TransactionGetID() []byte
	TransactionGetSender() []byte
	TransactionGetRecipient() []byte
	TransactionGetAmount() float64
	TransactionValidate() error
	TransactionSign(privateKey []byte) error
	TransactionVerify(publicKey []byte) bool
	TransactionSerialize() []byte
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

	data := bt.TransactionSerialize()
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
	if len(bt.Signature) != ed25519.SignatureSize {
		return false
	}

	data := bt.TransactionSerialize()

	return ed25519.Verify(publicKey, data, bt.Signature)
}

func (bt *BankTransaction) TransactionSerialize() []byte {
	buf := new(bytes.Buffer)

	buf.Write(bt.Sender)
	buf.Write(bt.Recipient)

	binary.Write(buf, binary.LittleEndian, math.Float64bits(bt.Amount))

	return buf.Bytes()
}
