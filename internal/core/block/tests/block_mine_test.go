package tests

import (
	"context"
	"testing"
	"time"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/block/interfaces"
	"github.com/Alex1997377/weave/internal/core/block/tests/helpers"
	"github.com/Alex1997377/weave/internal/core/block/tests/mocks"
)

// Тест на ошибки (без изменений)
func TestBlock_Mine_Errors(t *testing.T) {
	tests := []struct {
		name      string
		block     *block.Block
		config    block.MineConfig
		wantErr   bool
		errString string
	}{
		{
			name:      "nil block",
			block:     nil,
			config:    block.MineConfig{},
			wantErr:   true,
			errString: "block is nil",
		},
		{
			name:      "negative difficulty",
			block:     helpers.CreateTestBlock(0, -1),
			config:    block.MineConfig{},
			wantErr:   true,
			errString: "block difficulty cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.block.Mine(context.Background(), tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errString {
				t.Errorf("Mine() error = %v, want %v", err, tt.errString)
			}
		})
	}
}

// Тест успешного майнинга (с моком, который всегда возвращает валидный хеш)
func TestBlock_Mine_Success(t *testing.T) {
	blk := helpers.CreateTestBlock(1, 10)
	hasher := &mocks.MockHashCalculator{Valid: true}
	config := block.MineConfig{
		NumWorkers: 2,
		Verbose:    false,
		Hasher:     hasher,
	}

	err := blk.Mine(context.Background(), config)
	if err != nil {
		t.Fatalf("Mine() failed: %v", err)
	}
	if blk.Hash == nil {
		t.Error("Hash not set after mining")
	}
	if blk.Header.Nonce == 0 {
		t.Error("Nonce not changed after mining")
	}
}

// Тест таймаута
func TestBlock_Mine_Timeout(t *testing.T) {
	hasher := &mocks.MockHashCalculator{Valid: false}
	blk := helpers.CreateTestBlock(1, 20)
	config := block.MineConfig{
		NumWorkers: 1,
		Verbose:    false,
		Timeout:    10 * time.Millisecond,
		Hasher:     hasher,
	}

	err := blk.Mine(context.Background(), config)
	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}
	if err.Error() != "mining timeout" && err.Error() != "context deadline exceeded" {
		t.Errorf("Expected timeout error, got %v", err)
	}
}

// Тест отмены через контекст
func TestBlock_Mine_Cancel(t *testing.T) {
	hasher := &mocks.MockHashCalculator{Valid: false}
	blk := helpers.CreateTestBlock(1, 10)
	ctx, cancel := context.WithCancel(context.Background())
	config := block.MineConfig{
		NumWorkers: 2,
		Verbose:    false,
		Hasher:     hasher,
	}

	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	err := blk.Mine(ctx, config)
	if err == nil {
		t.Fatal("Expected cancel error, got nil")
	}
	if err.Error() != "context canceled" && err.Error() != "mining timeout" {
		t.Errorf("Expected context error, got %v", err)
	}
}

// Тест на то, что nonce увеличивается и победитель записывается
func TestBlock_Mine_NonceIncrements(t *testing.T) {
	conditionalHasher := &conditionalMockHasher{targetNonce: 42}
	blk := helpers.CreateTestBlock(1, 1)
	config := block.MineConfig{
		NumWorkers: 1,
		Verbose:    false,
		Hasher:     conditionalHasher,
	}

	err := blk.Mine(context.Background(), config)
	if err != nil {
		t.Fatalf("Mine failed: %v", err)
	}
	if blk.Header.Nonce != 42 {
		t.Errorf("Expected nonce 42, got %d", blk.Header.Nonce)
	}
}

// Реализация условного мока
type conditionalMockHasher struct {
	targetNonce uint64
}

func (c *conditionalMockHasher) Hash(data []byte) interfaces.Hash {
	if len(data) < 8 {
		return mocks.MockHash{Valid: false}
	}
	nonce := uint64(data[len(data)-8]) |
		uint64(data[len(data)-7])<<8 |
		uint64(data[len(data)-6])<<16 |
		uint64(data[len(data)-5])<<24 |
		uint64(data[len(data)-4])<<32 |
		uint64(data[len(data)-3])<<40 |
		uint64(data[len(data)-2])<<48 |
		uint64(data[len(data)-1])<<56
	valid := (nonce == c.targetNonce)
	return mocks.MockHash{Valid: valid, BytesHash: make([]byte, 32)}
}
