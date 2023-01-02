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

func TestFileExists(t *testing.T) {
	type args struct {
		filename string
	}
	tests := [...]struct {
		name string
		args args
		want bool
	}{
		{"not exists", args{"not-exists"}, false},
		{"not exists", args{Filename()}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FileExists(tt.args.filename)
			require.Equalf(t, tt.want, got, "want: %v, got: %v, file:%v", tt.want, got, tt.args.filename)
		})
	}
}

func TestJsonRecode(t *testing.T) {
	type HelloStruct struct {
		Message string `json:"message"`
	}
	var hello HelloStruct
	var helloMap = map[string]string{
		"message": "world",
	}

	require.NoError(t, JsonRecode(&hello, helloMap))
	require.Equal(t, &HelloStruct{Message: "world"}, &hello)
}

func TestReplaceExt(t *testing.T) {
	filename := "hello.mp3"
	got := ReplaceExt(filename, ".webm")
	require.Equal(t, "hello.webm", got)
}
