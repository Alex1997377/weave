package block

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/block/interfaces"
)

// Для переиспользования байтовых беферов
var headerBufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, 32)
	},
}

// Содержит параметры майнинга
type MineConfig struct {
	NumWorkers int
	Verbose    bool
	Timeout    time.Duration
	Hasher     interfaces.HashCalculator
}

func (b *Block) Mine(ctx context.Context, config MineConfig) error {
	if b == nil {
		return errors.New("block is nil")
	}

	if b.Header.Difficulty < 0 {
		return errors.New("block difficulty cannot be negative")
	}

	if config.NumWorkers <= 0 {
		config.NumWorkers = runtime.NumCPU()
	}

	if config.Verbose {
		fmt.Printf("Mining block %d, difficulty %d, workers %d\n",
			b.Header.Index, b.Header.Difficulty, config.NumWorkers)
	}

	baseHeader, nonceOffset, err := b.Header.SerializeWithoutNonce()
	if err != nil {
		return fmt.Errorf("failed to serialize header without nonce: %w", err)
	}

	startTime := time.Now()
	var (
		found       atomic.Bool
		winnerNonce atomic.Uint64
		hashResult  []byte
		wg          sync.WaitGroup
		stopCh      = make(chan struct{})
	)

	if config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, config.Timeout)
		defer cancel()
	}

	step := uint64(1 << 20)
	for i := 0; i < config.NumWorkers; i++ {
		wg.Add(1)
		args := &workerArgs{
			b:           b,
			baseHeader:  baseHeader,
			nonceOffset: nonceOffset,
			startNonce:  uint64(i) * step,
			step:        step,
			config:      config,
			found:       &found,
			winnerNonce: &winnerNonce,
			hashResult:  &hashResult,
			stopCh:      stopCh,
			ctx:         ctx,
		}
		go mineWorker(args, &wg)
	}

	wg.Wait()
	close(stopCh)

	if !found.Load() {
		if ctx.Err() == context.DeadlineExceeded {
			return errors.New("mining timeout")
		}
		return errors.New("mining failed to find valid nonce")
	}

	b.Hash = hashResult
	b.Header.Nonce = int(winnerNonce.Load())
	if config.Verbose {
		fmt.Printf("Mined! Nonce=%d, hash=%x, time=%v\n",
			b.Header.Nonce, hashResult, time.Since(startTime))
	}

	return nil
}

// -----------------------
// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// -----------------------

type workerArgs struct {
	b           *block.Block
	baseHeader  []byte
	nonceOffset int
	startNonce  uint64
	step        uint64
	config      MineConfig
	found       *atomic.Bool
	winnerNonce *atomic.Uint64
	hashResult  *[]byte
	stopCh      chan struct{}
	ctx         context.Context
}

func mineWorker(args *workerArgs, wg *sync.WaitGroup) {
	defer wg.Done()

	nonce := args.startNonce
	for {
		select {
		case <-args.stopCh:
			return
		case <-args.ctx.Done():
			return
		default:
		}
		if args.found.Load() {
			return
		}

		headerBufPtr := headerBufferPool.Get().(*[]byte)
		headerBuf := *headerBufPtr
		headerBuf = append(headerBuf[:0], args.baseHeader...)

		if len(headerBuf) < args.nonceOffset+8 {
			newBuf := make([]byte, args.nonceOffset+8)
			copy(newBuf, headerBuf)
			headerBuf = newBuf
			*headerBufPtr = headerBuf
		} else {
			headerBuf = headerBuf[:args.nonceOffset+8]
		}

		headerBuf[args.nonceOffset] = byte(nonce)
		headerBuf[args.nonceOffset+1] = byte(nonce >> 8)
		headerBuf[args.nonceOffset+2] = byte(nonce >> 16)
		headerBuf[args.nonceOffset+3] = byte(nonce >> 24)
		headerBuf[args.nonceOffset+4] = byte(nonce >> 32)
		headerBuf[args.nonceOffset+5] = byte(nonce >> 40)
		headerBuf[args.nonceOffset+6] = byte(nonce >> 48)
		headerBuf[args.nonceOffset+7] = byte(nonce >> 56)

		hash := args.config.Hash(headerBuf)

		*headerButPtr = headerBuf[:0]
		headerBufferPool.Put(headerButPtr)

		if hash.IsValidForDifficulty(args.b.Header.Difficulty) {
			if args.found.CompareAndSwap(false, true) {
				args.winnerNonce.Store(nonce)
				*args.hashResult = hash.Bytes()
				close(args.stopCh)
			}
			return
		}
		nonce++
		if nonce == 0 {
			return
		}
	}
}
