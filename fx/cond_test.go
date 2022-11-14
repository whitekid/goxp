package fx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
