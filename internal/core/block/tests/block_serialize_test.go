package tests

import (
	"testing"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/block/tests/helpers"
)

// TestBlock_Serialize проверяет сериализацию блока в байтовый поток.
// Метод Serialize должен преобразовывать блок в последовательность байт для
// передачи по сети или хранения на диске.
//
// Сценарий тестирования:
//   - Создаётся тестовый блок через helpers.CreateTestBlockForSerialize().
//   - Вызывается b.Serialize().
//
// Ожидаемый результат:
//   - Ошибка отсутствует.
//   - Полученные данные не пусты (len > 0).
//
// Входные данные:
//   - b: *block.Block – корректный блок, созданный вспомогательной функцией.
//
// Выходные значения:
//   - data: []byte – сериализованное представление блока.
//   - err: nil.
func TestBlock_Serialize(t *testing.T) {
	b := helpers.CreateTestBlockForSerialize()
	data, err := b.Serialize()
	if err != nil {
		t.Fatalf("Serialize error: %v", err)
	}
	if len(data) == 0 {
		t.Error("empty data")
	}
}

// TestBlock_CalculateHash проверяет вычисление хеша блока.
// Метод CalculateHash должен возвращать 32-байтовый хеш (например, SHA-256)
// на основе заголовка и транзакций блока.
//
// Сценарий тестирования:
//   - Используется тестовый блок из helpers.CreateTestBlockForSerialize().
//   - Вызывается b.CalculateHash().
//
// Ожидаемый результат:
//   - Ошибка отсутствует.
//   - Длина хеша равна 32 байтам.
//
// Входные данные:
//   - b: *block.Block – корректный блок.
//
// Выходные значения:
//   - hash: []byte – хеш блока (32 байта).
//   - err: nil.
func TestBlock_CalculateHash(t *testing.T) {
	b := helpers.CreateTestBlockForSerialize()
	hash, err := b.CalculateHash()
	if err != nil {
		t.Fatalf("CalculateHash error: %v", err)
	}
	if len(hash) != 32 {
		t.Errorf("hash length %d, want 32", len(hash))
	}
}

// TestBlock_CalculateSize проверяет вычисление размера блока в байтах.
// Метод CalculateSize должен возвращать точный размер, который получится
// при сериализации блока.
//
// Сценарий тестирования:
//   - Создаётся тестовый блок.
//   - Вычисляется размер через b.CalculateSize().
//   - Выполняется фактическая сериализация блока.
//   - Сравнивается вычисленный размер с длиной сериализованных данных.
//
// Ожидаемый результат:
//   - Ошибка отсутствует.
//   - Вычисленный размер не равен нулю.
//   - Вычисленный размер равен len(data) от Serialize().
//
// Входные данные:
//   - b: *block.Block – корректный блок.
//
// Выходные значения:
//   - size: uint32 – размер блока в байтах.
//   - err: nil.
func TestBlock_CalculateSize(t *testing.T) {
	b := helpers.CreateTestBlockForSerialize()
	size, err := b.CalculateSize()
	if err != nil {
		t.Fatalf("CalculateSize error: %v", err)
	}
	if size == 0 {
		t.Error("size is zero")
	}
	// сравнение с реальной сериализацией
	data, err := b.Serialize()
	if err != nil {
		t.Fatalf("Serialize error: %v", err)
	}
	if uint32(len(data)) != size {
		t.Errorf("size mismatch: calc=%d, serialized=%d", size, len(data))
	}
}

// TestBlock_Serialize_Nil проверяет поведение метода Serialize при вызове
// на nil-указателе блока.
//
// Сценарий тестирования:
//   - Объявляется переменная b типа *block.Block со значением nil.
//   - Вызывается b.Serialize().
//
// Ожидаемый результат:
//   - Возвращается ошибка (не nil).
//
// Выходные значения:
//   - data: nil (или пустой слайс).
//   - err: error с сообщением о nil-блоке.
func TestBlock_Serialize_Nil(t *testing.T) {
	var b *block.Block
	_, err := b.Serialize()
	if err == nil {
		t.Error("expected error for nil block")
	}
}

// TestBlock_CalculateHash_Nil проверяет поведение метода CalculateHash
// при вызове на nil-указателе блока.
//
// Сценарий тестирования:
//   - Переменная b равна nil.
//   - Вызывается b.CalculateHash().
//
// Ожидаемый результат:
//   - Возвращается ошибка.
//
// Выходные значения:
//   - hash: nil.
//   - err: error (не nil).
func TestBlock_CalculateHash_Nil(t *testing.T) {
	var b *block.Block
	_, err := b.CalculateHash()
	if err == nil {
		t.Error("expected error for nil block")
	}
}

// TestBlock_CalculateSize_Nil проверяет поведение метода CalculateSize
// при вызове на nil-указателе блока.
//
// Сценарий тестирования:
//   - b = nil.
//   - Вызывается b.CalculateSize().
//
// Ожидаемый результат:
//   - Возвращается ошибка.
//
// Выходные значения:
//   - size: 0.
//   - err: error (не nil).
func TestBlock_CalculateSize_Nil(t *testing.T) {
	var b *block.Block
	_, err := b.CalculateSize()
	if err == nil {
		t.Error("expected error for nil block")
	}
}
