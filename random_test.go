package goxp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomString(t *testing.T) {
	type args struct {
		randomString func(size int, source []rune) string
		source       []rune
	}
	tests := [...]struct {
		name string
		args args
	}{
		{`rand`, args{randomStringWithRand, randomChars}},
		{`crypto`, args{randomStringWithCrypto, randomChars}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRandomString(t, tt.args.randomString, 1000, tt.args.source)
		})
	}
}

func testRandomString(t *testing.T, randomString func(int, []rune) string, size int, source []rune) {
	got := randomString(size, source)
	for _, c := range got {
		require.Contains(t, source, c)
	}

	require.NotEqual(t, got, randomString(size, source))
}

func BenchmarkRandomString(b *testing.B) {
	tests := [...]struct {
		name         string
		randomString func(int, []rune) string
	}{
		{"rand", randomStringWithRand},
		{"crypto", randomStringWithCrypto},
	}
	for _, bb := range tests {
		b.Run(bb.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bb.randomString(100, randomChars)
			}
		})
	}
}

func FuzzRandomString(f *testing.F) {
	f.Add(100)
	f.Fuzz(func(t *testing.T, size int) {
		testRandomString(t, randomStringWithRand, size, randomChars)
		testRandomString(t, randomStringWithCrypto, size, randomChars)
	})
}

func TestRandomByte(t *testing.T) {
	type args struct {
		randomByte func(size int) []byte
	}
	tests := [...]struct {
		name string
		args args
	}{
		{"rand", args{randomByteWithRand}},
		{"crypto", args{randomByteWithCrypto}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRandomByte(t, tt.args.randomByte, 1000)
		})
	}
}

func testRandomByte(t *testing.T, randomByte func(size int) []byte, size int) {
	got := randomByte(size)
	require.Equal(t, size, len(got))
	for _, c := range got {
		require.Contains(t, got, c)
	}
}

func BenchmarkRandomByte(b *testing.B) {
	tests := [...]struct {
		name       string
		randomByte func(int) []byte
	}{
		{"rand", randomByteWithRand},
		{"crypto", randomByteWithCrypto},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tt.randomByte(100)
			}
		})
	}
}

func FuzzRandomByte(f *testing.F) {
	f.Add(100)
	f.Fuzz(func(t *testing.T, size int) {
		testRandomByte(t, randomByteWithRand, size)
		testRandomByte(t, randomByteWithCrypto, size)
	})
}
