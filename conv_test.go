package goxp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAtoiDef(t *testing.T) {
	type args struct {
		a   string
		def int
	}
	tests := [...]struct {
		name string
		args args
		want int
	}{
		{`valid`, args{"1234", 0}, 1234},
		{`valid`, args{"1234a", 12}, 12},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AtoiDef(tt.args.a, tt.args.def)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseBoolDef(t *testing.T) {
	type args struct {
		s   string
		def bool
	}
	tests := [...]struct {
		name string
		args args
		want bool
	}{
		{`valid`, args{"true", true}, true},
		{`valid`, args{"false", true}, false},
		{`valid`, args{"x-true", true}, true},
		{`valid`, args{"x-false", false}, false},
		{`valid`, args{"0", false}, false},
		{`valid`, args{"1", false}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseBoolDef(tt.args.s, tt.args.def)
			require.Equal(t, tt.want, got)
		})
	}
}
