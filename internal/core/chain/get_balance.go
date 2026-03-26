package chain

import (
	"bytes"
	"errors"
)

// GetBalance вычисляет баланс адреса
func (bc *Blockchain) GetBalance(address []byte) (float64, error) {
	if address == nil {
		return 0, errors.New("address cannot be nil")
	}

	var balance float64

	for _, b := range bc.Blocks {
		for _, tx := range b.Transaction {
			if tx == nil {
				continue
			}

			senderAddr := tx.TransactionGetSender()
			recipientAddr := tx.TransactionGetRecipient()

			if bytes.Equal(senderAddr, address) {
				balance -= tx.TransactionGetAmount()
			}

			if bytes.Equal(recipientAddr, address) {
				balance += tx.TransactionGetAmount()
			}
		}
	}
	return balance, nil
}
