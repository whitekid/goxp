package goxp

import (
	"testing"
	"time"
)

func foo() {
	defer Timer("foo()")()

	time.Sleep(500 * time.Millisecond)
}

func TestTimer(t *testing.T) {
	foo()
}
