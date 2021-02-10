package fixtures

import "testing"

func TestTempDir(t *testing.T) {
	defer TempDir("", "world")()
}
