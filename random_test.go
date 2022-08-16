package goxp

import (
	"testing"
)

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
