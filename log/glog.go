package log

import (
	"io"
	"log"

	"go.uber.org/zap"
)

type logWritter struct {
	logger Interface
	out    func(args ...interface{})
}

func (l *logWritter) Write(p []byte) (int, error) {
	l.logger.Error(string(p))
	return len(p), nil
}

func newLogWriter(level Level) io.Writer {
	w := &logWritter{
		logger: New(zap.AddCallerSkip(-1)),
	}

	w.out = map[Level]func(args ...interface{}){
		DebugLevel:  w.logger.Debug,
		InfoLevel:   w.logger.Info,
		WarnLevel:   w.logger.Warn,
		ErrorLevel:  w.logger.Error,
		DPanicLevel: w.logger.Fatal,
		PanicLevel:  w.logger.Fatal,
		FatalLevel:  w.logger.Fatal,
	}[level]

	return w
}

func NewGoLogger(level Level) *log.Logger {
	return log.New(newLogWriter(level), "", 0)
}
