package goxp

import "runtime"

// Filename return callers filename
// eg /opt/src/goxp/misc.go
func Filename() string {
	_, filename, _, _ := runtime.Caller(1)
	return filename
}
