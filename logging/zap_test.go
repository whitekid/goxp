package logging

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestNamed(t *testing.T) {
	type args struct {
		loggerName string
		message    string
	}
	tests := [...]struct {
		name string
		args args
	}{
		{"named", args{"logger", "hello"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, logs := observer.New(zapcore.DebugLevel)
			oldLogger := rootLogger
			defer func() {
				rootLogger = oldLogger
			}()
			rootLogger = zap.New(core)

			named := Named(tt.args.loggerName)
			named.Info(tt.args.message)

			all := logs.All()
			require.Equal(t, 1, len(all))
			e := all[0]
			require.Equal(t, tt.args.message, e.Message)
			require.Equal(t, tt.args.loggerName, e.LoggerName)
		})
	}
}
