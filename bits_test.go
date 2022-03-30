package goxp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetBits(t *testing.T) {
	type args struct {
		in    byte
		index int
		want  byte
	}
	tests := [...]struct {
		name string
		args args
	}{
		{"default", args{0x0, 0, 0b00000001}},
		{"default", args{0x0, 1, 0b00000010}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SetBit(tt.args.in, tt.args.index)
			require.Equal(t, got, tt.args.want, "got=%08b, want=%08b", got, tt.args.want)
		})
	}
}

func TestClearBits(t *testing.T) {
	type args struct {
		in    byte
		index int
		want  byte
	}
	tests := [...]struct {
		name string
		args args
	}{
		{"default", args{0b11111111, 0, 0b11111110}},
		{"default", args{0b11111111, 1, 0b11111101}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClearBit(tt.args.in, tt.args.index)
			require.Equal(t, got, tt.args.want, "got=%08b, want=%08b", got, tt.args.want)
		})
	}
}
