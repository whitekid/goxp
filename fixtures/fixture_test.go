package fixtures

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFixture(t *testing.T) {
	defer Env("HELLO", "WORLD")()
	require.Equal(t, "WORLD", os.Getenv("HELLO"))
}
