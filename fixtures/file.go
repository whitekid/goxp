package fixtures

import (
	"os"
	"sync"

	"github.com/whitekid/goxp"
)

// TempFile tempfile
func TempFile(dir, pattern string, callbacks ...func(string)) func() {
	f, err := os.CreateTemp(dir, pattern)
	goxp.Must(err)

	for _, callback := range callbacks {
		callback(f.Name())
	}

	var once sync.Once

	return func() { once.Do(func() { os.Remove(f.Name()) }) }
}
