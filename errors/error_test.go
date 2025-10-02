package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func newError(msg string) error {
	return New(msg)
}

func TestNew(t *testing.T) {
	err := newError("hello")
	got := fmt.Sprintf("%+v", err)

	// Basic error message check
	require.Contains(t, got, "hello")

	// Stack trace should contain the package path
	require.Contains(t, got, "github.com/whitekid/goxp/errors")

	t.Logf("err = %s", got)
}
