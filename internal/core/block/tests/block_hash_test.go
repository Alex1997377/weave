package tests

import (
	"bytes"
	"testing"

	"github.com/Alex1997377/weave/internal/core/block"
	"github.com/Alex1997377/weave/internal/core/header"
)

func createTestBlock(hash []byte) *block.Block {
	return &block.Block{
		Header:      header.Header{},
		Transaction: nil,
		Hash:        hash,
		Size:        0,
	}
}

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
			name:    "narmal block",
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
