package goxp

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
	"runtime"
)

// Filename return callers filename
// eg /opt/src/goxp/misc.go
func Filename() string {
	_, filename, _, _ := runtime.Caller(1)
	return filename
}

func FileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func JsonRecode(dest, src interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}

	if err := json.NewDecoder(&buf).Decode(dest); err != nil {
		return err
	}

	return nil
}

// ReplaceExt returns with new file extensions
func ReplaceExt(filename string, ext string) string {
	return filename[:len(filename)-len(path.Ext(filename))] + ext
}

// EnvExists return true if environment variabes exists
func EnvExists(k string) bool {
	_, exists := os.LookupEnv(k)
	return exists
}
