package fixtures

import (
	"testing"
)

func TestTimer(t *testing.T) {
	defer Timer("hello world")()
}
