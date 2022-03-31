package slug

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp"
	"github.com/whitekid/goxp/log"
)

func TestShortner(t *testing.T) {
	encoding := goxp.Shuffle([]byte(urlEncoding))
	shortner := NewShortner(string(encoding))

	for i := 1; i < 100000; i++ {
		s := shortner.Encode(int64(i))

		got, err := shortner.Decode(s)
		require.NoError(t, err)
		require.Equal(t, int64(i), got)
		if i == 99999 {
			log.Infof("%d => %s", i, s)
		}
	}
}
