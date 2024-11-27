package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func newError(msg string) error {
	return New("hello")
}

func TestNew(t *testing.T) {
	err := newError("hello")
	got := fmt.Sprintf("%+v", err)
	require.Contains(t, got, "newError")
	t.Logf("err = %s", got)
}
