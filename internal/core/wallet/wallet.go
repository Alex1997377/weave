package wallet

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"os"

	"github.com/Alex1997377/weave/internal/core/transaction"
	"github.com/Alex1997377/weave/internal/crypto/hash"
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
		return nil, NewCreateWalletError("generated keys are nil", nil)
	}

	return &Wallet{
		PrivateKey: priv,
		PublicKey:  pub,
	}, nil
}

func (w *Wallet) GetAddres() ([]byte, error) {
	if w.PublicKey == nil {
		return nil, NewInvalidWalletError("public key is nil")
	}

	address := hash.HashPublicKey(w.PublicKey)
	if address == nil {
		return nil, NewInvalidWalletError("failed to generate address from public key")
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

func (w *Wallet) SignTransaction(tx *transaction.BankTransaction) error {
	if tx == nil {
		return NewSignTransactionError("transaction is nil", nil)
	}

	if err := w.IsValid(); err != nil {
		return NewSignTransactionError("invalid wallet", err)
	}

	data, err := tx.TransactionSerialize()
	if err != nil {
		return NewSignTransactionError("failed to serialize transaction", err)
	}

	signature := ed25519.Sign(w.PrivateKey, data)
	if signature == nil {
		return NewSignTransactionError("failed to generate signature", nil)
	}

	tx.Signature = signature
	tx.Sender = w.PublicKey

	hash := sha256.Sum256(data)
	tx.ID = hash[:]

	return nil
}

func (w *Wallet) SaveToFile(filename string) error {
	if filename == "" {
		return NewSaveWalletError(filename, "filename cannot be empty", nil)
	}

	if err := w.IsValid(); err != nil {
		return NewSaveWalletError(filename, "invalid wallet", err)
	}

	err := os.WriteFile(filename, w.PrivateKey, 0600)
	if err != nil {
		return NewSaveWalletError(filename, "failed to write file", err)
	}

	return nil
}

func LoadFromFile(filename string) (*Wallet, error) {
	if filename == "" {
		return nil, NewLoadWalletError(filename, "filename cannot be empty", nil)
	}

	privBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, NewLoadWalletError(filename, "failed to read file", err)
	}

	if len(privBytes) != ed25519.PrivateKeySize {
		return nil, NewLoadWalletError(filename, "invalid private key size", nil)
	}

	priv := ed25519.PrivateKey(privBytes)
	if priv == nil {
		return nil, NewLoadWalletError(filename, "failed to create private key from bytes", nil)
	}

	pub, ok := priv.Public().(ed25519.PublicKey)
	if !ok || pub == nil {
		return nil, NewLoadWalletError(filename, "failed to extract public key from private key", nil)
	}

	return &Wallet{
		PrivateKey: priv,
		PublicKey:  pub,
	}, nil
}

func (w *Wallet) IsValid() error {
	if w == nil {
		return NewInvalidWalletError("wallet is nil")
	}

	if w.PrivateKey == nil {
		return NewInvalidWalletError("private key is nil")
	}

	if w.PublicKey == nil {
		return NewInvalidWalletError("public key is nil")
	}

	if len(w.PrivateKey) != ed25519.PrivateKeySize {
		return NewInvalidWalletError("private key has invalid size")
	}

	if len(w.PublicKey) != ed25519.PublicKeySize {
		return NewInvalidWalletError("public key has invalid size")
	}

	return nil
}
