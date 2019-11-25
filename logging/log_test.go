package logging

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestDefaultLoggerHasArgs(t *testing.T) {
	type fields struct {
		logFn func(string, ...interface{})
	}
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"nil args", fields{Infof}, args{"Hello", nil}, "Hello"},
		{"has args", fields{Infof}, args{"Hello %s", []interface{}{"world"}}, "Hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, logs := observer.New(zapcore.DebugLevel)
			oldLogger := defaultLogger
			defer func() {
				defaultLogger = oldLogger
			}()
			defaultLogger = zap.New(core).Sugar()

			tt.fields.logFn(tt.args.format, tt.args.args...)

			all := logs.All()
			require.Equal(t, 1, len(all))
			require.Equal(t, tt.want, all[0].Message)
		})
	}
}
