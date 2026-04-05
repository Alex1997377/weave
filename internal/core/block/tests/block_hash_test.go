package tests

import (
	"bytes"
	"testing"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/header"
)

// createTestBlock создаёт тестовый блок с заданным хешем.
// Поля Header и Transaction остаются пустыми/нулевыми, Size = 0.
//
// Параметры:
//   - hash: слайс байт (ожидается 32 байта), который будет присвоен полю Hash блока.
//
// Возвращаемое значение:
//   - *block.Block: указатель на блок с заполненным полем Hash.
func createTestBlock(hash []byte) *block.Block {
	return &block.Block{
		Header:      header.Header{},
		Transaction: nil,
		Hash:        hash,
		Size:        0,
	}
}

// TestBlock_HashString тестирует метод HashString() у блока.
// Метод HashString должен возвращать шестнадцатеричное представление хеша блока.
//
// Сценарии тестирования:
//  1. "nil block" – передаётся nil-указатель на блок.
//     Ожидается: ошибка (wantErr = true), строка пуста.
//  2. "normal block" – передаётся корректный блок с хешем из 32 байт (0xAB...).
//     Ожидается: ошибки нет, возвращается непустая строка.
//
// Входные данные тестов:
//   - block: *block.Block или nil.
//
// Ожидаемый результат:
//   - got: string – шестнадцатеричная строка (например, "abab...").
//   - err: error – nil для корректного блока, не nil для nil-блока.
//
// Проверки:
//   - Наличие ошибки соответствует ожиданию (wantErr).
//   - Для успешного случая строка не пуста.
func TestBlock_HashString(t *testing.T) {
	tests := []struct {
		name    string
		block   *block.Block
		wantErr bool
	}{
		{
			name:    "nil block",
			block:   nil,
			wantErr: true,
		},
		{
			name:    "normal block",
			block:   createTestBlock(bytes.Repeat([]byte{0xAB}, 32)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.block.HashString()

			if (err != nil) != tt.wantErr {
				t.Errorf("HashString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Error("HashString() returned empty string")
			}
		})
	}
}

// TestBlock_ShortHash тестирует метод ShortHash() у блока.
// Метод ShortHash должен возвращать первые 8 символов (4 байта) шестнадцатеричного хеша,
// что удобно для отображения в UI или логах.
//
// Сценарии тестирования:
//  1. "nil block" – передаётся nil-указатель на блок.
//     Ожидается: ошибка (wantErr = true), строка пуста.
//  2. "normal block" – передаётся корректный блок с хешем из 32 байт (0xAB...).
//     Ожидается: ошибки нет, возвращается непустая строка (длиной 8 символов).
//
// Входные данные тестов:
//   - block: *block.Block или nil.
//
// Ожидаемый результат:
//   - got: string – первые 8 символов шестнадцатеричного представления хеша.
//   - err: error – nil для корректного блока, не nil для nil-блока.
//
// Проверки:
//   - Наличие ошибки соответствует ожиданию (wantErr).
//   - Для успешного случая строка не пуста.
func TestBlock_ShortHash(t *testing.T) {
	tests := []struct {
		name    string
		block   *block.Block
		wantErr bool
	}{
		{
			name:    "nil block",
			block:   nil,
			wantErr: true,
		},
		{
			name:    "normal block",
			block:   createTestBlock(bytes.Repeat([]byte{0xAB}, 32)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.block.ShortHash()
			if (err != nil) != tt.wantErr {
				t.Errorf("ShortHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Error("ShortHash() returned empty string")
			}
		})
	}
}

// TestBlock_FormatHash тестирует метод FormatHash(prefix string) у блока.
// Метод FormatHash возвращает шестнадцатеричное представление хеша с заданным префиксом.
//
// Сценарии тестирования:
//  1. "nil block" – передаётся nil-указатель на блок.
//     Параметры: prefix = "0x".
//     Ожидается: ошибка (wantErr = true), строка пуста.
//  2. "normal block" – передаётся корректный блок с хешем из 32 байт (0xAB...).
//     Параметры: prefix = "hash:".
//     Ожидается: ошибки нет, возвращается непустая строка вида "hash:abab...".
//
// Входные данные тестов:
//   - block: *block.Block или nil.
//   - args.prefix: строка, которая будет добавлена перед шестнадцатеричным хешем.
//
// Ожидаемый результат:
//   - got: string – конкатенация prefix + шестнадцатеричная строка хеша.
//   - err: error – nil для корректного блока, не nil для nil-блока.
//
// Проверки:
//   - Наличие ошибки соответствует ожиданию (wantErr).
//   - Для успешного случая строка не пуста.
func TestBlock_FormatHash(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name    string
		block   *block.Block
		args    args
		wantErr bool
	}{
		{
			name:    "nil block",
			block:   nil,
			args:    args{prefix: "0x"},
			wantErr: true,
		},
		{
			name:    "normal block",
			block:   createTestBlock(bytes.Repeat([]byte{0xAB}, 32)),
			args:    args{prefix: "hash:"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.block.FormatHash(tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Error("FormatHash() returned empty string")
			}
		})
	}
}
