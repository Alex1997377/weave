// Package block предоставляет функции для десериализации блоков из бинарного формата.
// Основная функция DeserializeBlockWithparallelPooled использует параллельный пул воркеров
// для ускоренной обработки транзакций. Ожидается, что входные данные соответствуют формату:
// [заголовок (32 байта)] + [txCount (4 байта)] + [транзакции подряд] + [хеш (32)] + [размер (4)].
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
	"github.com/Alex1997377/weave/pkg/utils"
)

// Константы, определяющие максимальное количество транзакций в блоке (MaxTransactions),
// размер хеша в байтах (HashSize) и минимальный размер блока (minBlockSize).
const (
	MaxTransactions = 10000           // максимальное количество транзакций в блоке (защита от DoS)
	HashSize        = 32              // размер хеша в байтах (SHA-256)
	minBlockSize    = 32 + 4 + 32 + 4 // минимальный размер блока: заголовок(32) + txCount(4) + хеш(32) + размер(4)
)

// DeserializeOptions содержит зависимости для десериализации: десериализатор заголовка и транзакций.
// Позволяет подменять реализации (например, для тестирования).
type DeserializeOptions struct {
	Header interfaces.HeaderDeserializer      // десериализатор заголовка
	Tx     interfaces.TransactionDeserializer // десериализатор транзакций
}

// DeserializeBlockWithparallelPooled преобразует срез байт в структуру Block,
// используя параллельную обработку транзакций (worker pool).
//
// Параметры:
//
//	data - бинарные данные блока (ожидаемый формат описан выше)
//	opts - десериализаторы заголовка и транзакций
//
// Возвращает:
//
//	*Block - указатель на восстановленный блок (при успехе)
//	error  - ошибка валидации, нехватки данных, неверного формата и т.д.
//
// Алгоритм:
//  1. Проверка минимального размера (72 байта)
//  2. Десериализация заголовка
//  3. Чтение количества транзакций
//  4. Проверка лимита (MaxTransactions)
//  5. Вычисление границ транзакций (findTransactionBoundaries)
//  6. Параллельная десериализация транзакций (worker pool)
//  7. Чтение хеша и размера
//  8. Проверка на лишние данные
//
// Примерные величины (на CPU AMD Ryzen 3 3250U):
//   - txCount=0:      ~550 нс,   320 B,   6 аллокаций
//   - txCount=1:      ~24 мкс,  17.6 KB, 24 аллокации
//   - txCount=10:     ~42 мкс,  21.6 KB, 99 аллокаций
//   - txCount=100:    ~140 мкс, 61.5 KB, 819 аллокаций
//   - txCount=1000:   ~1.05 мс, 458 KB,  8019 аллокаций
//   - txCount=10000:  ~8.8 мс,  4.4 MB,  80021 аллокаций
func DeserializeBlockWithparallelPooled(data []byte, opts DeserializeOptions) (*Block, error) {
	// 1. Проверка минимального размера (заголовок + txCount + хеш + размер)
	if len(data) < minBlockSize {
		return nil, fmt.Errorf("data too short for block")
	}

	buf := bytes.NewReader(data)
	block := &Block{}

	// 2. Десериализация заголовка (через переданный десериализатор)
	header, err := opts.Header.DeserializeHeader(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize header: %w", err)
	}
	block.Header = *header

	// 3. Чтение количества транзакций (4 байта, little-endian)
	var txCount uint32
	if err := binary.Read(buf, binary.LittleEndian, &txCount); err != nil {
		return nil, fmt.Errorf("failed to read transaction count: %w", err)
	}

	// 4. Проверка лимита (защита от DoS)
	if txCount > MaxTransactions {
		return nil, fmt.Errorf("transaction count too high: %d (max: %d)", txCount, MaxTransactions)
	}

	// Смещение, где начинаются данные транзакций (после заголовка 32 и txCount 4)
	txDataStart := 32 + 4

	// 5. Вычисление границ каждой транзакции в срезе data[txDataStart:]
	txBoundaries, err := findTransactionBoundaries(data[txDataStart:], txCount)
	if err != nil {
		return nil, fmt.Errorf("failed to find tx boundaries: %w", err)
	}

	// Выделяем срез для транзакций
	block.Transaction = make([]transaction.Transaction, txCount)

	// 6. Параллельная десериализация транзакций (если они есть)
	if txCount > 0 {
		// Количество воркеров = min(txCount, число ядер CPU)
		numWorkers := utils.Min(int(txCount), runtime.NumCPU())
		wp := blockdeserialize.NewWorkerPool(numWorkers, opts.Tx)

		var resultsWg sync.WaitGroup
		resultsWg.Add(1)

		// Горутна-сборщик результатов из канала wp.Results
		go func() {
			defer resultsWg.Done()
			for result := range wp.Results {
				// Сохраняем транзакцию, только если нет ошибки и она не nil
				if result.Err == nil && result.Tx != nil {
					block.Transaction[result.Index] = result.Tx
				}
			}
		}()

		// Отправляем задачи для каждой транзакции
		for i := uint32(0); i < txCount; i++ {
			task := blockdeserialize.TaskPool.Get().(*blockdeserialize.TxTask)
			task.Index = i
			task.Data = data[txDataStart+txBoundaries[i] : txDataStart+txBoundaries[i+1]]
			task.Result = wp.Results
			wp.Tasks <- task
		}

		// Закрываем канал задач, ждём завершения воркеров, закрываем канал результатов,
		// ждём завершения сборщика.
		close(wp.Tasks)
		wp.Wg.Wait()
		close(wp.Results)
		resultsWg.Wait()

		// Проверяем, что все транзакции были успешно десериализованы
		for i, tx := range block.Transaction {
			if tx == nil {
				return nil, fmt.Errorf("tx %d missing", i)
			}
		}
	}

	// 7. Чтение хеша и размера блока (после всех транзакций)
	footerStart := txDataStart + txBoundaries[txCount]
	if len(data) < footerStart+HashSize+4 {
		return nil, errors.New("insufficient data for footer")
	}

	block.Hash = make([]byte, HashSize)
	copy(block.Hash, data[footerStart:footerStart+HashSize])
	block.Size = binary.LittleEndian.Uint32(data[footerStart+HashSize:])

	// Размер блока не может быть нулевым
	if block.Size == 0 {
		return nil, errors.New("invalid block size")
	}

	// 8. Проверка, что после размера нет лишних байт
	if len(data) > footerStart+HashSize+4 {
		return nil, fmt.Errorf("extra data after block deserialization: %d bytes remaining", len(data)-(footerStart+HashSize+4))
	}

	return block, nil
}

// DeserializeBlock упрощённая версия, использующая реальные десериализаторы по умолчанию.
// Предназначена для production-кода.
func DeserializeBlock(data []byte) (*Block, error) {
	return DeserializeBlockWithparallelPooled(data, DeserializeOptions{
		Header: interfaces.RealHeaderDeserializer{},
		Tx:     interfaces.RealTransactionDeserializer{},
	})
}

// DeserializeTransaction обёртка для вызова transaction.DeserializeTransactionFromReader.
// Оставлена для обратной совместимости.
func DeserializeTransaction(buf *bytes.Reader) (transaction.Transaction, error) {
	return transaction.DeserializeTransactionFromReader(buf)
}

// findTransactionBoundaries вычисляет границы (смещения) каждой транзакции в срезе байт,
// содержащем последовательные сериализованные транзакции.
//
// Параметры:
//
//	data    - срез байт, начинающийся с первой транзакции
//	txCount - количество транзакций
//
// Возвращает:
//
//	[]int - слайс длины txCount+1, где boundaries[i] - начало i-й транзакции,
//	        boundaries[i+1] - начало следующей (или конец последней)
//	error - ошибка, если данные слишком короткие, подпись слишком большая,
//	        или размер транзакции не укладывается в данные
//
// Формат одной транзакции (ожидаемый):
//
//	Sender (32) + Recipient (32) + ID (32) + Amount (8) + SigLen (4) + Signature (sigLen)
//	Минимальный размер заголовка без подписи = 108 байт.
//	Полный размер = 108 + sigLen.
//
// Пример:
//
//	Для sigLen = 64 (типовая подпись) txSize = 172 байта.
//	При txCount=10000 общий размер транзакций = 1 720 000 байт (~1.64 МБ),
//	boundaries будет занимать ~80 КБ памяти.
func findTransactionBoundaries(data []byte, txCount uint32) ([]int, error) {
	boundaries := make([]int, txCount+1)
	offset := 0

	for i := uint32(0); i < txCount; i++ {
		// Проверяем, что осталось хотя бы 108 байт (минимальный заголовок)
		if offset+108 > len(data) {
			return nil, fmt.Errorf("tx %d header out of bounds", i)
		}

		// Читаем длину подписи (sigLen) по смещению offset + 32+32+32+8
		sigLen := binary.LittleEndian.Uint32(data[offset+32+32+32+8:])
		// Защита от слишком больших подписей (DoS)
		if sigLen > 1024 {
			return nil, fmt.Errorf("signature too large at tx %d: %d", i, sigLen)
		}

		// Полный размер транзакции = 108 + sigLen
		txSize := 108 + int(sigLen)
		if offset+txSize > len(data) {
			return nil, fmt.Errorf("tx %d size mismatch: need %d, have %d", i, txSize, len(data)-offset)
		}

		offset += txSize
		boundaries[i+1] = offset
	}
	return boundaries, nil
}
