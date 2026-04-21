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

// TestBlock_Mine_Errors проверяет ошибочные ситуации при майнинге блока.
// Тестируются следующие сценарии:
//   1. "nil block" – передан nil-указатель на блок.
//      Ожидается: ошибка с текстом "block is nil".
//   2. "negative difficulty" – блок имеет отрицательную сложность (difficulty = -1).
//      Ожидается: ошибка с текстом "block difficulty cannot be negative".
//
// Входные параметры:
//   - block: *block.Block – может быть nil или содержать отрицательную сложность.
//   - config: block.MineConfig – конфигурация майнинга (не влияет на эти ошибки).
//
// Выходные значения:
//   - err: error – nil для успеха, иначе ошибка с указанным сообщением.
//
// Проверки:
//   - Наличие ошибки соответствует wantErr.
//   - Текст ошибки соответствует ожидаемому errString.
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
			block:     helpers.CreateTestBlock(0, -1), // сложность = -1
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

// TestBlock_Mine_Success проверяет успешный майнинг блока с использованием мока,
// который всегда возвращает валидный хеш.
//
// Входные данные:
//   - Создаётся тестовый блок с высотой 1 и сложностью 10.
//   - Используется моковый хеш-вычислитель mocks.MockHashCalculator{Valid: true},
//     который всегда сообщает, что хеш соответствует целевой сложности.
//   - Конфигурация майнинга: NumWorkers = 2, Verbose = false, Timeout = 0 (нет таймаута).
//
// Ожидаемый результат:
//   - Метод Mine() возвращает nil (ошибок нет).
//   - Поле Hash блока установлено (не nil).
//   - Поле Nonce заголовка блока изменено (не равно 0).
//
// Выходные значения:
//   - err: nil.
//   - Блок модифицируется: Hash и Header.Nonce заполняются.
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

// TestBlock_Mine_Timeout проверяет, что майнинг прерывается по таймауту.
//
// Входные данные:
//   - Блок со сложностью 20 (высокая, чтобы мок не находил решение).
//   - Моковый хеш-вычислитель с Valid = false (никогда не находит валидный хеш).
//   - Конфигурация: таймаут 10 миллисекунд, 1 воркер.
//
// Ожидаемый результат:
//   - Метод Mine() возвращает ошибку.
//   - Текст ошибки: "mining timeout" или "context deadline exceeded".
//
// Выходные значения:
//   - err: error с сообщением о таймауте.
//   - Блок не изменяется (хеш не устанавливается).
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

// TestBlock_Mine_Cancel проверяет, что майнинг отменяется через контекст.
//
// Входные данные:
//   - Блок со сложностью 10.
//   - Моковый хеш-вычислитель с Valid = false (решение не находится).
//   - Контекст с функцией cancel, вызываемой через 10 мс.
//   - Конфигурация: 2 воркера, без таймаута.
//
// Ожидаемый результат:
//   - Метод Mine() возвращает ошибку.
//   - Текст ошибки: "context canceled" или "mining timeout".
//
// Выходные значения:
//   - err: error с сообщением об отмене контекста.
//   - Блок не изменяется.
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

// TestBlock_Mine_NonceIncrements проверяет, что майнинг корректно перебирает nonce
// и останавливается на том значении, которое даёт валидный хеш.
//
// Входные данные:
//   - Блок со сложностью 1 (легко найти решение).
//   - Используется условный мок conditionalMockHasher, который считает хеш валидным
//     только когда nonce (извлекаемый из последних 8 байт данных) равен 42.
//   - Конфигурация: 1 воркер.
//
// Ожидаемый результат:
//   - Майнинг завершается успешно (ошибка nil).
//   - Nonce блока устанавливается в 42 (искомое значение).
//
// Выходные значения:
//   - err: nil.
//   - blk.Header.Nonce = 42.
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

// conditionalMockHasher – моковый хеш-вычислитель, который считает хеш валидным,
// только если значение nonce (извлекаемое из последних 8 байт данных) равно targetNonce.
//
// Поля:
//   - targetNonce: uint64 – значение nonce, при котором хеш считается валидным.
//
// Метод Hash(data []byte) возвращает:
//   - interfaces.Hash: моковый хеш с полем Valid = true, если извлечённый nonce == targetNonce,
//     иначе Valid = false. BytesHash всегда пустой (32 нулевых байта).
type conditionalMockHasher struct {
	targetNonce uint64
}

// Hash реализует интерфейс interfaces.Hasher.
// Извлекает nonce из последних 8 байт переданного слайса data.
// Если длина data меньше 8, возвращает невалидный хеш.
// В противном случае сравнивает nonce с targetNonce.
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