// internal/core/block/block_mine.go
// Реализация параллельного майнинга блока с возможностью настройки количества воркеров,
// таймаута и использования интерфейса для хеширования.
package block

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Alex1997377/weave/internal/core/block/interfaces"
	"github.com/Alex1997377/weave/internal/crypto/hash"
)

// MineConfig содержит параметры майнинга.
// Поля:
//
//	NumWorkers – количество параллельных воркеров (0 или ≤0 → runtime.NumCPU()).
//	Verbose    – если true, выводит прогресс в stdout.
//	Timeout    – максимальное время майнинга; 0 означает без таймаута.
//	Hasher     – реализация интерфейса HashCalculator (обязательно, но если nil, будет использован hash.HashCalculatorImpl).
type MineConfig struct {
	NumWorkers int
	Verbose    bool
	Timeout    time.Duration
	Hasher     interfaces.HashCalculator
}

// headerBufferPool – пул для переиспользования буферов заголовка, чтобы снизить нагрузку на GC.
// Буфер изначально создаётся с ёмкостью 512 байт, что достаточно для заголовка (обычно ≤256).
var headerBufferPool = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 0, 512)
		return &buf
	},
}

// Mine запускает процесс майнинга блока.
// Алгоритм:
//  1. Проверяет блок на nil и неотрицательную сложность.
//  2. Если сложность == 0, сразу вычисляет хеш заголовка и завершается.
//  3. Если Hasher не задан, подставляет реализацию по умолчанию.
//  4. Предварительно сериализует заголовок без nonce (SerializeWithoutNonce), чтобы не делать это на каждой итерации.
//  5. Запускает пул воркеров (количество = min(txCount, NumWorkers)), каждый перебирает свой диапазон nonce.
//  6. Как только один из воркеров находит подходящий хеш, останавливает всех.
//  7. Устанавливает найденный хеш и nonce в блок.
//
// Параметры:
//
//	ctx    – контекст для отмены или таймаута.
//	config – настройки майнинга.
//
// Возвращает:
//
//	error – nil при успехе, иначе описание ошибки.
//
// Примерные величины (на CPU AMD Ryzen 3 3250U, difficulty=10):
//
//	workers=1:  ~200 мкс, 22.8 KB, 804 аллокации
//	workers=4:  ~121 мкс, 26.8 KB, 930 аллокаций
//	workers=8:  ~95 мкс,  25.8 KB, 879 аллокаций
//	workers=16: ~177 мкс, 58.7 KB, 2024 аллокации (деградация из-за конкуренции)
func (b *Block) Mine(ctx context.Context, config MineConfig) error {
	// Базовые проверки
	if b == nil {
		return errors.New("block is nil")
	}
	if b.Header.Difficulty < 0 {
		return errors.New("block difficulty cannot be negative")
	}
	if config.Hasher == nil {
		config.Hasher = &hash.HashCalculatorImpl{}
	}
	// Сложность 0 – любой хеш подходит, вычисляем сразу
	if b.Header.Difficulty == 0 {
		h, err := b.CalculateHash()
		if err != nil {
			return err
		}
		b.Hash = h
		return nil
	}
	// Установка количества воркеров по умолчанию
	if config.NumWorkers <= 0 {
		config.NumWorkers = runtime.NumCPU()
	}
	if config.Verbose {
		fmt.Printf("Mining block %d, difficulty %d, workers %d\n",
			b.Header.Index, b.Header.Difficulty, config.NumWorkers)
	}
	// Предсериализация заголовка без nonce (оставляем 8 байт для вставки nonce)
	baseHeader, nonceOffset, err := b.Header.SerializeWithoutNonce()
	if err != nil {
		return fmt.Errorf("failed to serialize header without nonce: %w", err)
	}
	startTime := time.Now()
	var (
		found       atomic.Bool           // флаг, что решение найдено
		winnerNonce atomic.Uint64         // nonce, который дал решение
		hashResult  []byte                // найденный хеш
		wg          sync.WaitGroup        // для ожидания воркеров
		stopCh      = make(chan struct{}) // канал для остановки воркеров
	)
	if config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, config.Timeout)
		defer cancel()
	}
	// Защита от ситуации, когда Hasher не задан (хотя уже подставлен, но оставим)
	if config.Hasher == nil {
		return errors.New("hasher is required")
	}
	// Шаг для каждого воркера – диапазон nonce, чтобы не пересекались
	step := uint64(1 << 20) // ≈ 1 миллион
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
	// Проверка, найдено ли решение
	if !found.Load() {
		if err := ctx.Err(); err != nil {
			return err // context.Canceled или DeadlineExceeded
		}
		return errors.New("mining failed to find valid nonce")
	}
	b.Hash = hashResult
	b.Header.Nonce = winnerNonce.Load()
	if config.Verbose {
		fmt.Printf("Mined! Nonce=%d, hash=%x, time=%v\n",
			b.Header.Nonce, hashResult, time.Since(startTime))
	}
	return nil
}

// workerArgs – аргументы, передаваемые каждому воркеру.
// Инкапсулируют все необходимые данные, чтобы избежать замыканий на большие структуры.
type workerArgs struct {
	b           *Block
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

// mineWorker – горутина, перебирающая nonce в своём диапазоне.
// Работает следующим образом:
//  1. Начинает с nonce = startNonce.
//  2. На каждой итерации:
//     - Проверяет, не найден ли уже победитель (через found) или не остановлен ли процесс (stopCh, ctx).
//     - Берёт буфер из пула, копирует baseHeader, вставляет текущий nonce по смещению.
//     - Вычисляет хеш через Hasher.
//     - Если хеш удовлетворяет сложности, атомарно устанавливает found, сохраняет nonce и хеш, закрывает stopCh.
//  3. Увеличивает nonce на 1 (при переполнении uint64 завершается).
//  4. Возвращает буфер в пул.
//
// Параметры:
//
//	args – указатель на workerArgs.
//	wg   – WaitGroup, сигнализирует о завершении работы.
func mineWorker(args *workerArgs, wg *sync.WaitGroup) {
	defer wg.Done()
	nonce := args.startNonce
	for {
		// Проверка условий остановки
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
		// Получаем буфер из пула
		headerBufPtr := headerBufferPool.Get().(*[]byte)
		headerBuf := *headerBufPtr
		// Очищаем буфер и копируем baseHeader
		headerBuf = append(headerBuf[:0], args.baseHeader...)
		// Расширяем буфер, если не хватает места под nonce (обычно хватает, т.к. 512 байт)
		if len(headerBuf) < args.nonceOffset+8 {
			newBuf := make([]byte, args.nonceOffset+8)
			copy(newBuf, headerBuf)
			headerBuf = newBuf
			*headerBufPtr = headerBuf
		} else {
			headerBuf = headerBuf[:args.nonceOffset+8]
		}
		// Вставляем nonce как little-endian uint64
		headerBuf[args.nonceOffset] = byte(nonce)
		headerBuf[args.nonceOffset+1] = byte(nonce >> 8)
		headerBuf[args.nonceOffset+2] = byte(nonce >> 16)
		headerBuf[args.nonceOffset+3] = byte(nonce >> 24)
		headerBuf[args.nonceOffset+4] = byte(nonce >> 32)
		headerBuf[args.nonceOffset+5] = byte(nonce >> 40)
		headerBuf[args.nonceOffset+6] = byte(nonce >> 48)
		headerBuf[args.nonceOffset+7] = byte(nonce >> 56)
		// Вычисляем хеш
		h := args.config.Hasher.Hash(headerBuf)
		// Возвращаем буфер в пул (очищаем длину, но сохраняем ёмкость)
		*headerBufPtr = headerBuf[:0]
		headerBufferPool.Put(headerBufPtr)
		// Проверка сложности
		if h.IsValidForDifficulty(args.b.Header.Difficulty) {
			if args.found.CompareAndSwap(false, true) {
				args.winnerNonce.Store(nonce)
				*args.hashResult = h.Bytes()
				close(args.stopCh)
			}
			return
		}
		nonce++
		if nonce == 0 { // переполнение uint64
			return
		}
	}
}
