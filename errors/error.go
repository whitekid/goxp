package errors

import (
	"errors"
	"os"
	"strconv"
)

// standard errors package wrapper

var ErrUnsupported = errors.ErrUnsupported

// Global configuration for stack trace behavior
var (
	// EnableStackTrace can be disabled for performance in production
	EnableStackTrace = true
	
	// MaxStackDepth controls the maximum number of stack frames to capture
	MaxStackDepth = 50
)

func init() {
	// Allow disabling stack traces via environment variable
	if disabled := os.Getenv("GOXP_ERRORS_DISABLE_STACK"); disabled != "" {
		if val, err := strconv.ParseBool(disabled); err == nil {
			EnableStackTrace = !val
		}
	}
	
	// Allow configuring max stack depth via environment variable
	if depth := os.Getenv("GOXP_ERRORS_MAX_STACK_DEPTH"); depth != "" {
		if val, err := strconv.Atoi(depth); err == nil && val > 0 && val <= 200 {
			MaxStackDepth = val
		}
	}
}

func New(message string) error      { return wrap(nil, message, 3) }
func Is(err, target error) bool     { return errors.Is(err, target) }
func As(err error, target any) bool { return errors.As(err, target) }
func Join(errs ...error) error      { return errors.Join(errs...) }
func Unwrap(err error) error        { return errors.Unwrap(err) }

// NewWithoutStack creates an error without stack trace for high-performance scenarios
func NewWithoutStack(message string) error {
	return errors.New(message)
}

// WrapWithoutStack wraps an error without adding stack trace
func WrapWithoutStack(err error, message string) error {
	if err == nil {
		return NewWithoutStack(message)
	}
	return errors.New(message + ": " + err.Error())
}
