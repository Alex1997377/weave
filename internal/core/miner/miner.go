package miner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/chain"
	"github.com/Alex1997377/weave/internal/core/transaction"
)

type Mempool interface {
	GetPending() []transaction.Transaction
	Remove(txs []transaction.Transaction)
}

type Config struct {
	Difficulty int
	Mempool    Mempool
	Chain      chain.Chain
	Interval   time.Duration
}

type Miner struct {
	cfg    Config
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewMiner(cfg Config) *Miner {
	ctx, cancel := context.WithCancel(context.Background())
	return &Miner{
		cfg:    cfg,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (m *Miner) StartMine() {
	m.wg.Add(1)
	go m.loop()
}

func (m *Miner) StopMine() {
	m.cancel()
	m.wg.Wait()
}

func (m *Miner) loop() {
	defer m.wg.Done()
	for {
		select {
		case <-m.ctx.Done():
			return
		default:
			txs := m.cfg.Mempool.GetPending()

			if len(txs) == 0 {
				time.Sleep(m.cfg.Interval)
				continue
			}

			last, err := m.cfg.Chain.GetLastBlock()
			if err != nil {
				time.Sleep(m.cfg.Interval)
				continue
			}

			newBlock, err := block.NewBlock(
				txs,
				last.Hash,
				last.Header.Index+1,
				m.cfg.Difficulty,
			)

			if err != nil {
				time.Sleep(m.cfg.Interval)
				continue
			}

			if err := m.cfg.Chain.AddBlock(newBlock); err != nil {
				// логируем ошибку, но не прерываем цикл
				fmt.Printf("failed to add block: %v\n", err)
				continue
			}

			// Удаляем транзакции из мемпула
			m.cfg.Mempool.Remove(txs)
		}
	}
}
