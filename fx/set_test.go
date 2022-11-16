package fx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	set := NewSet[int]()
	set.Append(1, 2, 3, 3, 4, 5, 6, 1, 1, 2, 3, 2, 3, 5)
	require.Equal(t, []int{1, 2, 3, 4, 5, 6}, set.Slice())
	require.Equal(t, len(set.Slice()), set.Len())
	for _, e := range set.Slice() {
		require.True(t, set.Has(e))
	}
}