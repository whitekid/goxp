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
