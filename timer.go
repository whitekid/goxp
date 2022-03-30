package goxp

import (
	"fmt"
	"time"

	"github.com/whitekid/goxp/log"
	"go.uber.org/zap"
)

var (
	logTimer = log.WithOptions(zap.AddCallerSkip(1))
)

// Timer check running time
// Usage:
//	defer Timer("check it")()
func Timer(format string, args ...interface{}) func() {
	t := time.Now()

	return func() {
		logTimer.Debugf("time takes %s: %s", time.Since(t), fmt.Sprintf(format, args...))
	}
}
