package goxp

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/constraints"
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

func testParseIntDef[T constraints.Integer](t *testing.T, s string, defaultValue, minValue, maxValue, want T) {
	got := ParseIntDef(s, defaultValue, minValue, maxValue)
	require.Equal(t, want, got)
}

func TestParseIntDef(t *testing.T) {
	type args struct {
		s   string
		def int
		min int
		max int
	}
	tests := [...]struct {
		name string
		args args
		want int
	}{
		{`valid`, args{"1", 5, 0, 10}, 1},
		{`valid`, args{"100", 5, 0, 10}, 10},
		{`valid`, args{"-10", 5, 0, 10}, 0},
		{`valid`, args{"-xx", 5, 0, 10}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParseIntDef(t, tt.args.s, tt.args.def, tt.args.min, tt.args.max, tt.want)
		})
	}
}

func FuzzParseIntDef(f *testing.F) {
	f.Add("1", 5, 0, 10)
	f.Fuzz(func(t *testing.T, s string, defValue, minValue, maxValue int) {
		want := ParseIntDef(s, defValue, minValue, maxValue)
		testParseIntDef(t, s, defValue, minValue, maxValue, want)
	})
}
