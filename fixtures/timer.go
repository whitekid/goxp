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

	cleared := false

	return func() {
		if cleared {
			return
		}

		span := time.Now().Sub(start)
		timerLogger.Debugf("%s takes %s", span, fmt.Sprintf(format, args...))

		cleared = true
	}
}
