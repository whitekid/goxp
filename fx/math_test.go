package fx

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMin(t *testing.T) {
	type args struct {
		collection []int
	}
	tests := [...]struct {
		name string
		args args
		want int
	}{
		{"", args{[]int{1, 2, 3}}, 1},
		{"", args{[]int{2, 1, 3}}, 1},
		{"", args{[]int{1}}, 1},
		{"", args{[]int{3, 2, 1}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, Min(tt.args.collection))
		})
	}
}

func TestMax(t *testing.T) {
	type args struct {
		collection []int
	}
	tests := [...]struct {
		name string
		args args
		want int
	}{
		{"", args{[]int{1, 2, 3}}, 3},
		{"", args{[]int{1, 3, 2}}, 3},
		{"", args{[]int{1}}, 1},
		{"", args{[]int{3, 2, 1}}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, Max(tt.args.collection))
		})
	}
}

func TestSumBy(t *testing.T) {
	require.Equal(t, 6, Sum([]int{1, 2, 3}))
}

func TestSum(t *testing.T) {
	require.Equal(t, 6,
		SumBy([]string{"1", "2", "3"}, func(s string) int {
			v, _ := strconv.Atoi(s)
			return int(v)
		}))
}

func TestScale(t *testing.T) {
	require.Equal(t, []int{2, 4, 6, 8}, Scale([]int{1, 2, 3, 4}, 2))
}
