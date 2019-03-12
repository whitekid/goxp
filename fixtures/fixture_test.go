package fixtures

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixture(t *testing.T) {
	defer Env("HELLO", "WORLD").Teardown()
	assert.Equal(t, "WORLD", os.Getenv("HELLO"))
}
