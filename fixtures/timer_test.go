package fixtures

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTimer(t *testing.T) {
	func() {
		defer Timer("hello world")
	}()
	require.Fail(t, "@@@")
}
