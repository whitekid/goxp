package log

import "go.uber.org/zap/zapcore"

// Interface log.interface
type Interface interface {
	SetLevel(zapcore.Level)

	Debug(args ...any)
	Debugf(fmt string, args ...any)
	Info(args ...any)
	Infof(fmt string, args ...any)
	Warn(args ...any)
	Warnf(fmt string, args ...any)
	Error(args ...any)
	Errorf(fmt string, args ...any)
	Fatal(args ...any)
	Fatalf(fmt string, args ...any)
}

// Level is a logging priority
type Level = zapcore.Level

// LogLevel aliases
const (
	DebugLevel  = zapcore.DebugLevel
	InfoLevel   = zapcore.InfoLevel
	WarnLevel   = zapcore.WarnLevel
	ErrorLevel  = zapcore.ErrorLevel
	DPanicLevel = zapcore.DPanicLevel
	PanicLevel  = zapcore.PanicLevel
	FatalLevel  = zapcore.FatalLevel
)

// Utility functions
func SetLevel(level Level)              { defaultLogger.SetLevel(level) }
func Info(args ...any)                  { defaultLogger.Info(args...) }
func Infof(format string, args ...any)  { defaultLogger.Infof(format, args...) }
func Debug(args ...any)                 { defaultLogger.Debug(args...) }
func Debugf(format string, args ...any) { defaultLogger.Debugf(format, args...) }
func Warn(args ...any)                  { defaultLogger.Warn(args...) }
func Warnf(format string, args ...any)  { defaultLogger.Warnf(format, args...) }
func Error(args ...any)                 { defaultLogger.Error(args...) }
func Errorf(format string, args ...any) { defaultLogger.Errorf(format, args...) }
func Fatal(args ...any)                 { defaultLogger.Fatal(args...) }
func Fatalf(format string, args ...any) { defaultLogger.Fatalf(format, args...) }
