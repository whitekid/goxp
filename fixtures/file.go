package fixtures

import (
	"io/ioutil"
	"os"
)

// TempDir create temporary directory and remove it
func TempDir(name, pattern string, callbacks ...func(string)) func() {
	dir, _ := ioutil.TempDir(name, pattern)

	for _, callback := range callbacks {
		callback(dir)
	}

	return func() {
		os.RemoveAll(dir)
	}
}
