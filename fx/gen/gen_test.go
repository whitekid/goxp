package gen

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp/iterx"
)

func TestSlice(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5}
	next := Slice(s)

	for i := 0; i < len(s); i++ {
		v, ok := next()
		require.True(t, ok)
		require.Equal(t, i, v)
	}

	_, ok := next()
	require.False(t, ok)
}

func TestNext(t *testing.T) {
	genTo10 := func() func() (int, bool) {
		i := 0
		return func() (int, bool) {
			defer func() { i++ }()
			return i, i < 10
		}
	}

	for i, v := range iterx.All(Next(genTo10).Seq()) {
		require.Equal(t, i, v)
	}
}

func TestMap(t *testing.T) {
	m := map[int]int{
		1: 10,
		2: 20,
	}

	next := Map(m)
	for k, v, ok := next(); ok; k, v, ok = next() {
		require.Equal(t, m[k], v)
	}

	_, _, ok := next()
	require.False(t, ok)
}
