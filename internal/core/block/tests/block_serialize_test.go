package tests

import (
	"testing"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/block/tests/helpers"
)

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

func TestBlock_Serialize_Nil(t *testing.T) {
	var b *block.Block
	_, err := b.Serialize()
	if err == nil {
		t.Error("expected error for nil block")
	}
}

func TestBlock_CalculateHash_Nil(t *testing.T) {
	var b *block.Block
	_, err := b.CalculateHash()
	if err == nil {
		t.Error("expected error for nil block")
	}
}

func TestBlock_CalculateSize_Nil(t *testing.T) {
	var b *block.Block
	_, err := b.CalculateSize()
	if err == nil {
		t.Error("expected error for nil block")
	}
}
