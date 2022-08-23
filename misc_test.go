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

func TestJsonRedecode(t *testing.T) {
	type HelloStruct struct {
		Message string `json:"message"`
	}
	var hello HelloStruct
	var helloMap = map[string]string{
		"message": "world",
	}

	require.NoError(t, JsonRedecode(&hello, helloMap))
	require.Equal(t, &HelloStruct{Message: "world"}, &hello)
}
