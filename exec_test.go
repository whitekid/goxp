package goxp

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	type args struct {
		command []string
		shell   bool
	}
	tests := [...]struct {
		name       string
		args       args
		wantErr    bool
		wantOutput string
	}{
		{`valid`, args{[]string{"ls", "-al"}, false}, false, ""},
		{`valid`, args{[]string{"ls", "-al"}, false}, false, "exec_test.go"},
		{`valid`, args{[]string{"ls", "-al"}, true}, false, ""},
		{`valid`, args{[]string{"ls", "-al"}, true}, false, "exec_test.go"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			exc := Exec(tt.args.command...)
			exc = exc.Shell(tt.args.shell)

			var err error
			var output []byte
			IfThen(tt.wantOutput == "", func() { err = exc.Do(ctx) }, func() { output, err = exc.Output(ctx) })

			require.Truef(t, (err != nil) == tt.wantErr, `Execute.Do() failed: error = %+v, wantErr = %v`, err, tt.wantErr)
			if tt.wantErr {
				return
			}

			if tt.wantOutput != "" {
				require.Contains(t, string(output), tt.wantOutput)
			}
		})
	}
}

func TestExecPipeIn(t *testing.T) {
	type args struct {
		cmd   []string
		stdin string
	}
	tests := [...]struct {
		name       string
		args       args
		wantErr    bool
		wantOutput string
	}{
		{`valid`, args{cmd: []string{"wc", "-c"}}, false, "0"},
		{`valid`, args{cmd: []string{"wc", "-c"}, stdin: "Hello world"}, false, "11"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			exc := Exec(tt.args.cmd...)
			if tt.args.stdin != "" {
				exc = exc.Pipe(func(wc io.Writer) {
					wc.Write([]byte(tt.args.stdin))
				}, nil, nil)
			}

			output, err := exc.Output(ctx)
			require.Truef(t, (err != nil) == tt.wantErr, `Executor.Do() failed: error = %+v, wantErr = %v`, err, tt.wantErr)
			if tt.wantErr {
				return
			}
			require.Equal(t, tt.wantOutput, strings.TrimSpace(string(output)))
		})
	}
}

func TestExecPipeOut(t *testing.T) {
	type args struct {
		cmd []string
	}
	tests := [...]struct {
		name       string
		args       args
		wantErr    bool
		wantOutput string
	}{
		{`valid`, args{cmd: []string{"echo", "-n", "hello world"}}, false, "hello world"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			exc := Exec(tt.args.cmd...)
			var output []byte
			exc = exc.Pipe(nil, func(rc io.Reader) {
				out, err := io.ReadAll(rc)

				require.NoError(t, err)
				output = out
			}, nil)

			err := exc.Do(ctx)
			require.Truef(t, (err != nil) == tt.wantErr, `Executor.Do() failed: error = %+v, wantErr = %v`, err, tt.wantErr)
			if tt.wantErr {
				return
			}
			require.Equal(t, tt.wantOutput, strings.TrimSpace(string(output)))
		})
	}
}
