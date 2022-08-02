package goxp

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilename(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	require.Equal(t, filename, Filename())
}
