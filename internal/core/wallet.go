package core

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"

	"github.com/Alex1997377/weave/internal/crypto"
)

type Wallet struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

func CreateWallet() (*Wallet, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, NewCreateWalletError("generate wallet error or something went wrong", err)
	}

	if pub == nil || priv == nil {
		return nil, errors.New("generated keys are nil")
	}
	return &Wallet{
		PrivateKey: priv,
		PublicKey:  pub,
	}, nil
}

func (w *Wallet) GetAddres() ([]byte, error) {
	if w.PublicKey == nil {
		return nil, errors.New("wallet public key is nil")
	}

	address := crypto.HashPublicKey(w.PublicKey)
	if address == nil {
		return nil, errors.New("failed to generate address from public key")
	}

	return address, nil
}

func (w *Wallet) GetAddressString() (string, error) {
	address, err := w.GetAddres()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(address), nil
}

func (w *Wallet) SignTransaction(tx *BankTransaction) error {
	if tx == nil {
		return errors.New("transaction is nil")
	}

	if w.PrivateKey == nil {
		return errors.New("wallet private key is nil")
	}

	if w.PublicKey == nil {
		return errors.New("wallet public key is nil")
	}

	data, err := tx.TransactionSerialize()
	if err != nil {
		return errors.New("failed to serialize transaction")
	}

	signature := ed25519.Sign(w.PrivateKey, data)
	if signature == nil {
		return errors.New("failed to generate signature")
	}

	tx.Signature = signature
	tx.Sender = w.PublicKey

	hash := sha256.Sum256(data)
	tx.ID = hash[:]

	return nil
}

func (w *Wallet) SaveToFile(filename string) error {
	if filename == "" {
		return errors.New("filename cannot be empty")
	}

	if w.PrivateKey == nil {
		return errors.New("wallet private key is nil")
	}

	if len(w.PrivateKey) != ed25519.PrivateKeySize {
		return errors.New("invalid private key size")
	}

	err := os.WriteFile(filename, w.PrivateKey, 0600)
	if err != nil {
		return &os.PathError{
			Op:   "write",
			Path: filename,
			Err:  err,
		}
	}

	return nil
}

func LoadFromFile(filename string) (*Wallet, error) {
	if filename == "" {
		return nil, errors.New("filename cannot be empty")
	}

	privBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, &os.PathError{
			Op:   "read",
			Path: filename,
			Err:  err,
		}
	}

	if len(privBytes) != ed25519.PrivateKeySize {
		return nil, errors.New("invalid private key size")
	}

	priv := ed25519.PrivateKey(privBytes)
	if priv == nil {
		return nil, errors.New("failed to create private key from bytes")
	}

	pub, ok := priv.Public().(ed25519.PublicKey)
	if !ok || pub == nil {
		return nil, errors.New("failed to extract public key from private key")
	}

	return &Wallet{
		PrivateKey: priv,
		PublicKey:  pub,
	}, nil
}

func (w *Wallet) IsValid() error {
	if w == nil {
		return errors.New("wallet is nil")
	}

	if w.PrivateKey == nil {
		return errors.New("private key is nil")
	}

	if w.PublicKey == nil {
		return errors.New("public key is nil")
	}

	if len(w.PrivateKey) != ed25519.PrivateKeySize {
		return errors.New("private key has invalid size")
	}

	if len(w.PublicKey) != ed25519.PublicKeySize {
		return errors.New("public key has invalid size")
	}

	return nil
}
