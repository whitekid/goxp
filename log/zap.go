// Package log is zap based logger
package log

import (
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
)

var (
	rootLogger     *zap.Logger
	rootLoggerOnce sync.Once
)

func root() *zap.Logger {
	rootLoggerOnce.Do(func() {
		level := strings.ToLower(os.Getenv("LOG_LEVEL"))
		switch level {
		case "debug":
			rootLogger, _ = zap.NewDevelopment()
		default:
			rootLogger, _ = zap.NewProduction()
		}
	})

	return rootLogger
}

// New create new logger
func New() *zap.SugaredLogger {
	return root().Sugar()
}

// Named create new named logger
func Named(name string) *zap.SugaredLogger {
	return New().Named(name)
}

// WithOptions create logger with option
func WithOptions(opts ...zap.Option) *zap.SugaredLogger {
	return root().WithOptions(opts...).Sugar()
}
