package goxp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type myError struct{}

func (err *myError) Error() string { return "myerror" }

func TestErrorAs(t *testing.T) {
	e := &myError{}

	ee, ok := ErrorAs[*myError](e)
	require.True(t, ok)
	require.Equal(t, e, ee)
}
