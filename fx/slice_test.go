package fx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlice(t *testing.T) {
	s := S([]string{"1", "2", "3"})
	r := []string{}
	s.Each(func(_ int, e string) { r = append(r, e) })
	require.Equal(t, s.Slice(), r)
}
