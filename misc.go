package goxp

import (
	"bytes"
	"encoding/json"
	"os"
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

func JsonRedecode(dest, src interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}

	if err := json.NewDecoder(&buf).Decode(dest); err != nil {
		return err
	}

	return nil
}

func IfThen(condition bool, thenF func(), falseF ...func()) {
	if condition {
		thenF()
	}

	if len(falseF) > 0 {
		falseF[0]()
	}
}
