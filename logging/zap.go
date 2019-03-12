// Package logging is zap based logger
//
// usage:
//
// var log = logging.New()
package logging

import "go.uber.org/zap"

var (
	rootLogger *zap.Logger
)

func root() *zap.Logger {
	if rootLogger == nil {
		rootLogger, _ = zap.NewDevelopment()
	}

	return rootLogger
}

// New create new logger
func New() *zap.SugaredLogger {
	return root().Sugar()
}

// Named create new named logger
func Named(name string) *zap.SugaredLogger {
	return New().Named("[" + name + "]")
}
