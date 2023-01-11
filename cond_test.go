package goxp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIfThen(t *testing.T) {
	want := 0

	trueF := func() { want++ }
	falseF := func() { want-- }

	type args struct {
		cond   bool
		trueF  func()
		falseF []func()
	}
	tests := [...]struct {
		name string
		args args
		want int
	}{
		{`valid`, args{true, trueF, []func(){falseF}}, 1},
		{`valid`, args{true, trueF, []func(){}}, 1},
		{`valid`, args{false, trueF, []func(){falseF}}, -1},
		{`valid`, args{false, trueF, []func(){falseF, falseF}}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want = 0

			IfThen(tt.args.cond, tt.args.trueF, tt.args.falseF...)
			require.Equal(t, tt.want, want)
		})
	}
}

func TestTernary(t *testing.T) {
	type args struct {
		value int
	}
	tests := [...]struct {
		name string
		args args
		want string
	}{
		{"even", args{10}, "even"},
		{"odd", args{11}, "odd"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Ternary(tt.args.value%2 == 0, "even", "odd")
			require.Equal(t, tt.want, got)
		})
	}
}

func TestIf(t *testing.T) {
	require.Equal(t, "true", If(func() bool { return true }, "true").Else("false"))
	require.Equal(t, "false", If(func() bool { return false }, "true").Else("false"))
}
