package fixtures

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFixture(t *testing.T) {
	teardown := Env("HELLO", "WORLD")
	require.Equal(t, "WORLD", os.Getenv("HELLO"))

	teardown()
	require.Equal(t, "", os.Getenv("HELLO"))
}
