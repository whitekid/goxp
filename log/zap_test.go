package log

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestEnv(t *testing.T) {
	type args struct {
		envs map[string]string
	}
	tests := [...]struct {
		name      string
		args      args
		wantLevel zapcore.Level
	}{
		{"debug", args{map[string]string{"LOG_LEVEL": "DEBUG"}}, zapcore.DebugLevel},
		{"info", args{map[string]string{"LOG_LEVEL": "INFO"}}, zapcore.InfoLevel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.args.envs {
				old, exists := os.LookupEnv(k)
				os.Setenv(k, v)
				defer func(k string) {
					if exists {
						os.Setenv(k, old)
					} else {
						os.Unsetenv(k)
					}
				}(k)
			}

			_, level := newLogger()
			require.Equal(t, tt.wantLevel, level.Level())
		})
	}
}

func TestNamed(t *testing.T) {
	type args struct {
		loggerName string
		message    string
	}
	tests := [...]struct {
		name string
		args args
	}{
		{"named", args{"name", "hello"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, logs := observer.New(zapcore.DebugLevel)
			_ = core

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
