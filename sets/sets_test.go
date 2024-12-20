package sets

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	set := New[int]()
	set.Set(1, 2, 3, 3, 4, 5, 6, 1, 1, 2, 3, 2, 3, 5)
	require.Equal(t, []int{1, 2, 3, 4, 5, 6}, set.Slice())
	require.Equal(t, len(set.Slice()), set.Len())
	for _, e := range set.Slice() {
		require.True(t, set.Contains(e))
	}

	slice := set.Slice()
	set.Remove(slice[0])
	require.Equal(t, len(slice)-1, set.Len())
	require.Equal(t, len(slice)-1, len(set.keys))
	require.False(t, set.Contains(slice[0]))
}

func TestRemove(t *testing.T) {
	s := New[int](0, 1, 2, 3)
	s.Remove(2)

	require.Equal(t, []int{0, 1, 3}, slices.Collect(s.Values()))
}

func TestSetNX(t *testing.T) {
	s := New[int](1, 2, 3, 4)
	type args struct {
		e int
	}
	tests := [...]struct {
		name string
		args args
		want bool
	}{
		{`valid`, args{1}, false},
		{`valid`, args{5}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.SetNX(tt.args.e)
			require.Equal(t, tt.want, got, `SetNX() failed: got = %+v, want = %v`, got, tt.want)
		})
	}
}
