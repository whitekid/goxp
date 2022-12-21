package fixtures

import (
	"fmt"
	"sync"
	"time"

	"github.com/whitekid/goxp/log"
)

var (
	timerLogger     log.Interface
	timerLoggerOnce sync.Once
)

// Timer log execution time
func Timer(format string, args ...interface{}) Teardown {
	timerLoggerOnce.Do(func() { timerLogger = log.New() })
	start := time.Now()

	var once sync.Once

	return func() {
		once.Do(func() {
			span := time.Since(start)
			timerLogger.Debugf("%s takes %s", span, fmt.Sprintf(format, args...))
		})
	}
}
