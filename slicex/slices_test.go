package slicex

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSample(t *testing.T) {
	s := []int{12, 3, 4, 5}
	slices.Contains(s, Sample(s))
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
