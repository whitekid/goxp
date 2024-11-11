package gen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntN(t *testing.T) {
	next := IntN(10)
	for i := 0; i < 10; i++ {
		v, ok := next()
		require.True(t, ok)
		require.Equal(t, i, v)
	}
	v, ok := next()
	require.False(t, ok)
	require.Equal(t, 0, v)
}

func TestCycle(t *testing.T) {
	next := Cycle([]int{})
	v, ok := next()
	require.False(t, ok)
	require.Equal(t, 0, v)
}

func TestSample(t *testing.T) {
	next := Sample([]int{0, 1, 2, 3, 4, 5})
	for range 100000 {
		v, ok := next()
		require.True(t, ok)
		require.LessOrEqual(t, v, 5)
	}
}
