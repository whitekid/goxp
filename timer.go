package goxp

import (
	"fmt"
	"sync"
	"time"

	"github.com/whitekid/goxp/log"
	"go.uber.org/zap"
)

var (
	logTimer     log.Interface
	logTimerOnce sync.Once
)

// Timer check running time
// Usage:
//	defer Timer("check it")()
func Timer(format string, args ...interface{}) func() {
	logTimerOnce.Do(func() {
		log.New(zap.AddCallerSkip(1))
	})

	t := time.Now()

	return func() {
		logTimer.Debugf("time takes %s: %s", time.Since(t), fmt.Sprintf(format, args...))
	}
}
