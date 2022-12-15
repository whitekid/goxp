package fixtures

import (
	"os"
)

// TempFile tempfile
func TempFile(dir, pattern string, callbacks ...func(string)) func() {
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		panic(err)
	}

	for _, callback := range callbacks {
		callback(f.Name())
	}

	return func() { os.Remove(f.Name()) }
}
