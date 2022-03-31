package goxp

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShuffle(t *testing.T) {
	type args struct {
		ary []rune
	}
	tests := [...]struct {
		name string
		args args
	}{
		{"default", args{[]rune("ABCDEF")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Shuffle(tt.args.ary)
			require.NotEqual(t, tt.args.ary, got)

			sorted := make([]rune, len(got))
			copy(sorted, got)
			sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

			require.Equal(t, tt.args.ary, sorted)
		})
	}
}
