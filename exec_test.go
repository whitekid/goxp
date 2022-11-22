package goxp

import (
	"context"
	"testing"

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
			IfThen(tt.args.shell, func() { exc = exc.Shell() }, func() { exc = exc.NoShell() })

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
