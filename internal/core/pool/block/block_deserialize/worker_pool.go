package blockdeserialize

import (
	"bytes"
	"errors"
	"sync"

	"github.com/Alex1997377/weave/internal/core/block/interfaces"
)

type WorkerPool struct {
	Tasks          chan *TxTask
	Results        chan *TxResult
	Wg             sync.WaitGroup
	Quit           chan struct{}
	TxDeserializer interfaces.TransactionDeserializer
}

func NewWorkerPool(numWorkers int, txDeserializer interfaces.TransactionDeserializer) *WorkerPool {
	if txDeserializer == nil {
		panic("deserializer cannot be nil")
	}

	wp := &WorkerPool{
		Tasks:          make(chan *TxTask, 1000),
		Results:        make(chan *TxResult, 1000),
		Wg:             sync.WaitGroup{},
		Quit:           make(chan struct{}),
		TxDeserializer: txDeserializer,
	}

	for i := 0; i < numWorkers; i++ {
		wp.Wg.Add(1)
		go wp.worker()
	}

	return wp
}

func (wp *WorkerPool) worker() {
	defer wp.Wg.Done()

	for {
		select {
		case task, ok := <-wp.Tasks:
			if !ok {
				return
			}
			if task == nil {
				continue
			}
			if task.Data == nil {
				wp.Results <- &TxResult{
					Index: task.Index,
					Err:   errors.New("task data is nil"),
				}
				continue
			}

			tx, err := wp.TxDeserializer.DeserializeTransaction(bytes.NewReader(task.Data))

			wp.Results <- &TxResult{
				Index: task.Index,
				Tx:    tx,
				Err:   err,
			}
		case <-wp.Quit:
			return

		}
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.Quit)
	wp.Wg.Wait()
	close(wp.Tasks)
	close(wp.Results)
}
