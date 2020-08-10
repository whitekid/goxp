package utils

import "testing"

func BenchmarkRandomString(b *testing.B) {
	type args struct {
		fn func(int) string
	}
	tests := [...]struct {
		name string
		args args
	}{
		{"", args{RandomString}},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tt.args.fn(100)
			}
		})
	}
}
