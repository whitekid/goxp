package slicex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSample(t *testing.T) {
	s := []int{12, 3, 4, 5}
	e := Sample(s)
	require.Contains(t, s, e)
}

func TestSamples(t *testing.T) {
	s := []int{12, 3, 4, 5}

	type args struct {
		size int
	}
	tests := [...]struct {
		name string
		args args
	}{
		{`valid`, args{3}},
		{`valid`, args{-1}},
		{`valid`, args{4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Samples(s, tt.args.size)
			
			if tt.args.size < 1 {
				require.Nil(t, got)
				return
			}
			if tt.args.size >= len(s) {
				require.Equal(t, got, s)
				return
			}

			require.Equal(t, len(got), tt.args.size)
			for _, e := range got {
				require.Contains(t, s, e)
			}
		})
	}
}

func TestToPtr(t *testing.T) {
	s := []string{"hello", "world"}
	p := ToPtr(s)

	for i := 0; i < len(s); i++ {
		require.Equal(t, s[i], *p[i])
		require.Equal(t, &s[i], p[i])
	}
}

func TestIntersect(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := []int{1, 2, 3}

	require.Equal(t, []int{4, 5}, Intersect(s1, s2))
}

func TestUniq(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 6, 3, 4, 5, 3, 2, 1}
	require.Equal(t, []int{1, 2, 3, 4, 5, 6}, Uniq(s))
}

func TestFlatten(t *testing.T) {
	s := [][]int{{1, 2, 3}, {4, 5, 6}}
	require.Equal(t, []int{1, 2, 3, 4, 5, 6}, Flatten(s))
}

func TestGroupBy(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5, 6}
	r := GroupBy(s, func(e int) int { return e % 2 })
	require.Equal(t, map[int][]int{0: {0, 2, 4, 6}, 1: {1, 3, 5}}, r)
}

func TestConcat(t *testing.T) {
	s := Of[[]int](0, 1, 2)
	got := s.Concat([]int{3, 4, 5}).Slice()
	require.Equal(t, []int{0, 1, 2, 3, 4, 5}, got)
}
