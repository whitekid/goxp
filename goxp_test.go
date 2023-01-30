package goxp

import (
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	// check for goroutine
	goleak.VerifyTestMain(m)
}
