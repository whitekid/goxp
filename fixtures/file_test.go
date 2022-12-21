package fixtures

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp"
)

func TestTempFile(t *testing.T) {
	var name string
	teardown := TempFile(t.TempDir(), "test-*.tmp", func(n string) { name = n })
	defer teardown()

	require.True(t, goxp.FileExists(name))

	teardown()
	require.False(t, goxp.FileExists(name))
}
