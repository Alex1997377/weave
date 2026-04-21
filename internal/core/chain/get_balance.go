package chain

// GetBalance вычисляет баланс адреса
func (bc *Blockchain) GetBalance(address []byte) (float64, error) {
	if address == nil {
		return 0, NewInvalidAddressError("address cannot be nil", nil)
	}

	key := string(address)
	balance, ok := bc.balances[key]
	if !ok {
		return 0, nil
	}

	return balance, nil
}
