// Package logging is zap based logger
package logging

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newEncoder() zapcore.Encoder {
	level := strings.ToLower(os.Getenv("LOG_LEVEL"))

	switch level {
	case "debug":
		encoderCfg := zap.NewDevelopmentEncoderConfig()
		return zapcore.NewConsoleEncoder(encoderCfg)
	default:
		encoderCfg := zap.NewProductionEncoderConfig()
		return zapcore.NewJSONEncoder(encoderCfg)
	}
}

func newLogger() (*zap.Logger, zap.AtomicLevel) {
	encoder := newEncoder()
	level := zap.NewAtomicLevel()

	return zap.New(zapcore.NewCore(
		encoder,
		zapcore.Lock(os.Stdout),
		level,
	),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	), level
}

type zapLogger struct {
	*zap.SugaredLogger
	level zap.AtomicLevel
}

func (l *zapLogger) SetLevel(level zapcore.Level) {
	l.level.SetLevel(level)
}

var (
	defaultLogger *zapLogger
)

func init() {
	logger, level := newLogger()
	zap.ReplaceGlobals(logger)
	defaultLogger = &zapLogger{
		SugaredLogger: logger.Sugar(),
		level:         level,
	}
}

// New create new logger
func New() Interface {
	logger, level := newLogger()
	return &zapLogger{
		SugaredLogger: logger.Sugar(),
		level:         level,
	}
}

// Named create new named logger
func Named(name string) Interface {
	logger, level := newLogger()
	return &zapLogger{
		SugaredLogger: logger.Sugar().Named(name),
		level:         level,
	}
}

// WithOptions create logger with option
func WithOptions(opts ...zap.Option) Interface {
	logger, level := newLogger()
	return &zapLogger{
		SugaredLogger: logger.WithOptions(opts...).Sugar(),
		level:         level,
	}
}
