package fixtures

import (
	"io/ioutil"
	"os"
)

// TempDir create temporary directory and remove it
//
// Deprecated: use t.TempDir()
func TempDir(name, pattern string, callbacks ...func(string)) func() {
	dir, err := ioutil.TempDir(name, pattern)
	if err != nil {
		panic(err)
	}

	for _, callback := range callbacks {
		callback(dir)
	}

	return func() { os.RemoveAll(dir) }
}

// TempFile tempfile
func TempFile(dir, pattern string, callbacks ...func(string)) func() {
	f, err := ioutil.TempFile(dir, pattern)
	if err != nil {
		panic(err)
	}

	for _, callback := range callbacks {
		callback(f.Name())
	}

	return func() { os.Remove(f.Name()) }
}
