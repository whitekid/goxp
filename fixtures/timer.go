package fixtures

import (
	"fmt"
	"time"

	log "github.com/whitekid/go-utils/logging"
)

// Timer log execution time
func Timer(format string, args ...interface{}) Teardown {
	start := time.Now()

	return func() {
		span := time.Now().Sub(start)

		log.Debugf("%s takes %s", span, fmt.Sprintf(format, args...))
	}
}
