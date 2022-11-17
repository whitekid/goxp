package goxp

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	func() {
		defer Timer("foo()")()

		time.Sleep(500 * time.Millisecond)
	}()
}
