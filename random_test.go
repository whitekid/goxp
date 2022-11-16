package goxp

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomString(t *testing.T) {
	s := RandomString(10)
	require.Equal(t, 10, len(s))
	require.NotEqual(t, "", s)
}

func TestRandomStringWithCrypto(t *testing.T) {
	s := RandomStringWithCrypto(10)
	require.Equal(t, 10, len(s))
	require.NotEqual(t, "", s)
}

func BenchmarkRandomString(b *testing.B) {
	type args struct {
		fn func(int) string
	}
	tests := [...]struct {
		name string
		args args
	}{
		{"rand", args{RandomString}},
		{"crypto", args{RandomStringWithCrypto}},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tt.args.fn(100)
			}
		})
	}
}

func TestRandomStringWith(t *testing.T) {
	s := RandomStringWith(10, []rune("abcdefg"))
	require.Equal(t, 10, len(s))
	require.NotEqual(t, "", s)
}

func TestRandomByte(t *testing.T) {
	b := RandomByte(10)
	require.Equal(t, 10, len(b))
	require.NotEqual(t, "", hex.EncodeToString(b))
}

func BenchmarkRandomByte(b *testing.B) {
	type args struct {
		fn func(int) []byte
	}
	tests := [...]struct {
		name string
		args args
	}{
		{"rand", args{RandomByte}},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tt.args.fn(100)
			}
		})
	}
}
