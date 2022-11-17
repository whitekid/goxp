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
//
//	   func doSomething(){
//		      defer Timer("doSomething()")()
//
// .   .      bla.... bla...
//
//	}
func Timer(format string, args ...interface{}) func() {
	logTimerOnce.Do(func() {
		logTimer = log.New(zap.AddCallerSkip(1))
	})

	t := time.Now()

	return func() {
		logTimer.Debugf("time takes %s: %s", time.Since(t), fmt.Sprintf(format, args...))
	}
}
