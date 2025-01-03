package fx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFx(t *testing.T) {
	s := Of(0, 1, 2, 3, 4)
	for i, v := range s.All() {
		require.Equal(t, i, v)
	}

	i := 0
	for v := range s {
		require.Equal(t, i, v)
		i++
	}
}

func TestFilter(t *testing.T) {
	s := Of(0, 1, 2, 3, 4, 5)
	require.Equal(t, []int{0, 2, 4}, s.Filter(func(e int) bool { return e%2 == 0 }).Collect())
}

func TestMap(t *testing.T) {
	s := Of(0, 1, 2)
	require.Equal(t, []int{0, 2, 4}, Map(s, func(i int) int { return i * 2 }).Collect())
}
