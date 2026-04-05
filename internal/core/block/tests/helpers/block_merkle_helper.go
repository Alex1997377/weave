package helpers

import (
	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/block/tests/mocks"
	"github.com/Alex1997377/weave/internal/core/header"
	"github.com/Alex1997377/weave/internal/core/transaction"
)

func CreateTestBlockWithTxIDs(txIDs [][]byte) *block.Block {
	transactions := make([]transaction.Transaction, len(txIDs))
	for i, id := range txIDs {
		transactions[i] = &mocks.MockTransaction{Id: id}
	}
	return &block.Block{
		Transaction: transactions,
		Header:      header.Header{},
	}
}
