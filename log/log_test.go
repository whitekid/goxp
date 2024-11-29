package log

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestDefaultLoggerHasArgs(t *testing.T) {
	type fields struct {
		logFn func(string, ...any)
	}
	type args struct {
		format string
		args   []any
	}
	tests := [...]struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"nil args", fields{Infof}, args{"Hello", nil}, "Hello"},
		{"has args", fields{Infof}, args{"Hello %s", []any{"world"}}, "Hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, logs := observer.New(zapcore.DebugLevel)

			oldLogger := defaultLogger
			defaultLogger = &zapLogger{
				SugaredLogger: zap.New(core).Sugar(),
			}

			defer func() {
				defaultLogger = oldLogger
			}()

			tt.fields.logFn(tt.args.format, tt.args.args...)

			all := logs.All()
			require.Equal(t, 1, len(all))
			require.Equal(t, tt.want, all[0].Message)
		})
	}
}

func TestLevel(t *testing.T) {
	p, _ := zap.NewProduction()
	p.Sugar().Infof("hello world")

	Info("hello")
	Debug("hello debug")

	defaultLogger.SetLevel(zap.DebugLevel)
	Info("hello")
	Debug("hello debug")

	SetLevel(zap.InfoLevel)

	{
		logger := Named("named")
		logger.Info("named info")
		logger.Debug("named debug")

		Info("named info")
		Debug("named debug")
	}
}
