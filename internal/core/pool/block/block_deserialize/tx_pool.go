package blockdeserialize

import (
	"sync"

	"github.com/Alex1997377/weave/internal/core/transaction"
)

var (
	TxPool = sync.Pool{
		New: func() interface{} {
			return &transaction.BankTransaction{}
		},
	}

	BytePool = sync.Pool{
		New: func() interface{} {
			buf := make([]byte, 0, 4096)
			return &buf
		},
	}

	TaskPool = sync.Pool{
		New: func() interface{} {
			return &TxTask{}
		},
	}
)

type TxTask struct {
	Index  uint32
	Data   []byte
	Result chan *TxResult
}

type TxResult struct {
	Index uint32
	Tx    transaction.Transaction
	Err   error
}
