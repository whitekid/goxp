package slicex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToPtr(t *testing.T) {
	s := []string{"hello", "world"}
	p := ToPtr(s)

	for i := 0; i < len(s); i++ {
		require.Equal(t, s[i], *p[i])
		require.Equal(t, &s[i], p[i])
	}
}
