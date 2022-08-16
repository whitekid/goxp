package slug

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp/fx"
)

func TestShortner(t *testing.T) {
	encoding := fx.Shuffle([]byte(EncodeURL))
	shortner := NewShortner(string(encoding))

	for i := 1; i < 100000; i++ {
		s := shortner.Encode(int64(i))

		got, err := shortner.Decode(s)
		require.NoError(t, err)
		require.Equal(t, int64(i), got)
	}
}
