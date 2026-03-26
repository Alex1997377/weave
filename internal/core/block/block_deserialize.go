// block_deserialize.go (полностью, с исправлениями)
package block

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"runtime"
	"sync"

	"github.com/Alex1997377/weave/internal/core/block/interfaces"
	blockdeserialize "github.com/Alex1997377/weave/internal/core/pool/block/block_deserialize"
	"github.com/Alex1997377/weave/internal/core/transaction"
)

const (
	MaxTransactions = 10000
	HashSize        = 32
	minBlockSize    = 32 + 4 + 32 + 4
)

type DeserializeOptions struct {
	Header interfaces.HeaderDeserializer
	Tx     interfaces.TransactionDeserializer
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func DeserializeBlockWithparallelPooled(data []byte, opts DeserializeOptions) (*Block, error) {
	if len(data) < minBlockSize {
		return nil, fmt.Errorf("data too short for block")
	}

	buf := bytes.NewReader(data)
	block := &Block{}

	header, err := opts.Header.DeserializeHeader(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize header: %w", err)
	}
	block.Header = *header

	var txCount uint32
	if err := binary.Read(buf, binary.LittleEndian, &txCount); err != nil {
		return nil, fmt.Errorf("failed to read transaction count: %w", err)
	}
	if txCount > MaxTransactions {
		return nil, fmt.Errorf("transaction count too high: %d (max: %d)", txCount, MaxTransactions)
	}

	txDataStart := 32 + 4
	txBoundaries, err := findTransactionBoundaries(data[txDataStart:], txCount)
	if err != nil {
		return nil, fmt.Errorf("failed to find tx boundaries: %w", err)
	}

	block.Transaction = make([]transaction.Transaction, txCount)
	if txCount > 0 {
		numWorkers := min(int(txCount), runtime.NumCPU())
		wp := blockdeserialize.NewWorkerPool(numWorkers, opts.Tx)

		var resultsWg sync.WaitGroup
		resultsWg.Add(1)
		go func() {
			defer resultsWg.Done()
			for result := range wp.Results {
				if result.Err == nil && result.Tx != nil {
					block.Transaction[result.Index] = result.Tx
				}
			}
		}()

		for i := uint32(0); i < txCount; i++ {
			task := blockdeserialize.TaskPool.Get().(*blockdeserialize.TxTask)
			task.Index = i
			task.Data = data[txDataStart+txBoundaries[i] : txDataStart+txBoundaries[i+1]]
			task.Result = wp.Results
			wp.Tasks <- task
		}

		close(wp.Tasks)
		wp.Wg.Wait()
		close(wp.Results)
		resultsWg.Wait()

		for i, tx := range block.Transaction {
			if tx == nil {
				return nil, fmt.Errorf("tx %d missing", i)
			}
		}
	}

	footerStart := txDataStart + txBoundaries[txCount]
	if len(data) < footerStart+HashSize+4 {
		return nil, errors.New("insufficient data for footer")
	}

	block.Hash = make([]byte, HashSize)
	copy(block.Hash, data[footerStart:footerStart+HashSize])
	block.Size = binary.LittleEndian.Uint32(data[footerStart+HashSize:])

	if block.Size == 0 {
		return nil, errors.New("invalid block size")
	}

	if len(data) > footerStart+HashSize+4 {
		return nil, fmt.Errorf("extra data after block deserialization: %d bytes remaining", len(data)-(footerStart+HashSize+4))
	}

	return block, nil
}

func DeserializeBlock(data []byte) (*Block, error) {
	return DeserializeBlockWithparallelPooled(data, DeserializeOptions{
		Header: interfaces.RealHeaderDeserializer{},
		Tx:     interfaces.RealTransactionDeserializer{},
	})
}

func DeserializeTransaction(buf *bytes.Reader) (transaction.Transaction, error) {
	return transaction.DeserializeTransactionFromReader(buf)
}

func findTransactionBoundaries(data []byte, txCount uint32) ([]int, error) {
	boundaries := make([]int, txCount+1)
	offset := 0
	for i := uint32(0); i < txCount; i++ {
		if offset+108 > len(data) {
			return nil, fmt.Errorf("tx %d header out of bounds", i)
		}
		sigLen := binary.LittleEndian.Uint32(data[offset+32+32+32+8:])
		if sigLen > 1024 {
			return nil, fmt.Errorf("signature too large at tx %d: %d", i, sigLen)
		}
		txSize := 108 + int(sigLen)
		if offset+txSize > len(data) {
			return nil, fmt.Errorf("tx %d size mismatch: need %d, have %d", i, txSize, len(data)-offset)
		}
		offset += txSize
		boundaries[i+1] = offset
	}
	return boundaries, nil
}
