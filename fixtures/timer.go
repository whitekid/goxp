package fixtures

import (
	"fmt"
	"sync"
	"time"

	"github.com/whitekid/go-utils/log"
)

var (
	timerLogger     log.Interface
	timerLoggerOnce sync.Once
)

// Timer log execution time
func Timer(format string, args ...interface{}) Teardown {
	start := time.Now()

	return func() {
		span := time.Now().Sub(start)

		timerLoggerOnce.Do(func() { timerLogger = log.New() })

		timerLogger.Debugf("%s takes %s", span, fmt.Sprintf(format, args...))
	}
}
